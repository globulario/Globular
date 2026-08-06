package cluster

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Package staging must publish transactionally: a package is either absent from
// the cache or complete in it, never observable half-written at its final name.
//
// The defect this pins, 2026-08-05: Phase 4.3 used
//
//	mv "${TARBALL_DIR}/packages/"* "${STATE_DIR}/packages/"
//
// TARBALL_DIR is /tmp and STATE_DIR is on disk. Verified in the running
// containers: /tmp is tmpfs (dev 129), /var/lib/globular is /dev/nvme0n1p3
// (dev 66307). Across a filesystem boundary mv cannot rename — it degrades to
// copy+unlink, so the destination appears at its final discoverable name
// immediately and fills over seconds, with mtime already set from the archive
// so nothing looks unfinished.
//
// node.join reads that directory during convergence and lost the race on 2 of
// 4 joiners: node-3 on prometheus (62 MB) and node-5 on alertmanager (18 MB),
// both "stage package: gzip: EOF". Both tarballs passed gzip -t minutes later.
// The bytes were never corrupt; the reader was early.
//
// Invariant: repository.cached_artifact_must_not_be_visible_before_verified_commit

// TestPartialPublishIsObservable_NegativeControl proves the failure mechanism is
// real rather than assumed. It reproduces what a cross-filesystem mv does — copy
// straight to the final name — and shows a concurrent consumer can read a
// truncated archive from the published path.
//
// Without this control, the positive test below could pass for the wrong reason
// (e.g. the writer finishing before the reader ever looks).
func TestPartialPublishIsObservable_NegativeControl(t *testing.T) {
	dir := t.TempDir()
	final := filepath.Join(dir, "prometheus_3.5.1_linux_amd64.tgz")

	// A slow producer writing directly to the final path, exactly as
	// copy+unlink does.
	f, err := os.Create(final)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("PARTIAL-ARCHIVE-PREFIX")); err != nil {
		t.Fatal(err)
	}
	// Producer has NOT finished. A consumer scanning the cache right now:
	got, err := os.ReadFile(final)
	if err != nil {
		t.Fatalf("consumer could not read the published path at all: %v", err)
	}
	f.Close()

	if len(got) != len("PARTIAL-ARCHIVE-PREFIX") {
		t.Fatalf("expected to observe the truncated prefix, read %d bytes", len(got))
	}
	// This is the whole bug: the file is at its final, discoverable name and is
	// incomplete. If this control ever stops reproducing, the test below is no
	// longer proving anything.
	t.Logf("negative control reproduced: consumer read %d incomplete bytes from %s",
		len(got), filepath.Base(final))
}

// TestAtomicPublishNeverExposesPartialArtifact runs the directive's scenario:
// start a deliberately slow cache write, concurrently request the package, and
// require that the consumer either does not see it at all or sees it whole —
// never a partial final artifact.
func TestAtomicPublishNeverExposesPartialArtifact(t *testing.T) {
	cache := t.TempDir()
	staging := filepath.Join(cache, ".incoming.test")
	if err := os.MkdirAll(staging, 0o755); err != nil {
		t.Fatal(err)
	}
	final := filepath.Join(cache, "prometheus_3.5.1_linux_amd64.tgz")
	complete := strings.Repeat("COMPLETE-ARCHIVE-BYTES", 4096) // ~90 KB

	tmp := filepath.Join(staging, "prometheus_3.5.1_linux_amd64.tgz")
	f, err := os.Create(tmp)
	if err != nil {
		t.Fatal(err)
	}

	// Slow producer: write in chunks, polling the consumer's view between each.
	// The staging file is inside a dot-directory, which the node-agent's cache
	// scan skips (it requires a .tgz suffix and rejects directories), so nothing
	// here is discoverable.
	chunk := len(complete) / 8
	for i := 0; i < 8; i++ {
		end := (i + 1) * chunk
		if i == 7 {
			end = len(complete)
		}
		if _, err := f.Write([]byte(complete[i*chunk : end])); err != nil {
			t.Fatal(err)
		}

		// The consumer looks at the cache mid-write, every time.
		if data, err := os.ReadFile(final); err == nil {
			if string(data) != complete {
				t.Fatalf("chunk %d: consumer observed a PARTIAL artifact (%d of %d bytes) "+
					"at the final path — publication was not atomic", i, len(data), len(complete))
			}
		} else if !os.IsNotExist(err) {
			t.Fatalf("chunk %d: unexpected consumer error: %v", i, err)
		}
		// os.IsNotExist is the correct pre-publication observation: absent, not
		// half-there. That is what "structurally invisible" means.
	}
	if err := f.Sync(); err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Publish with the atomic rename. Same filesystem, so this really is atomic.
	if err := os.Rename(tmp, final); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(final)
	if err != nil {
		t.Fatalf("after publication the artifact must be readable: %v", err)
	}
	if string(got) != complete {
		t.Fatalf("published artifact is not the complete one: %d of %d bytes",
			len(got), len(complete))
	}
}

// TestJoinScriptPublishesPackagesAtomically pins the shipped script, because the
// mechanism above only protects the cluster if Phase 4.3 actually uses it.
func TestJoinScriptPublishesPackagesAtomically(t *testing.T) {
	script := joinScriptForTest()

	i := strings.Index(script, "[4.3] Staging packages")
	if i < 0 {
		t.Fatal("Phase 4.3 staging block not found — renamed? update this test")
	}
	block := script[i:]
	if j := strings.Index(block, "[4.4]"); j >= 0 {
		block = block[:j]
	}

	// Strip comment lines before checking for the defective form. The block
	// deliberately QUOTES the old `mv` in its explanatory comment — that
	// documentation is the point, and matching it would make this test fire on
	// the very text that explains why it exists.
	var code []string
	for _, line := range strings.Split(block, "\n") {
		if t := strings.TrimSpace(line); t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		code = append(code, line)
	}
	block = strings.Join(code, "\n")

	// The exact defective form must not come back.
	if strings.Contains(block, `mv "${TARBALL_DIR}/packages/"* "${STATE_DIR}/packages/"`) {
		t.Error("Phase 4.3 moves packages directly from TARBALL_DIR (/tmp) into the " +
			"live cache: that crosses a filesystem boundary, so mv is copy+unlink and " +
			"publishes half-written packages at their final names")
	}
	// It must stage inside the destination filesystem first...
	if !strings.Contains(block, "STAGING_DIR=") ||
		!strings.Contains(block, `${STATE_DIR}/packages/.incoming`) {
		t.Error("Phase 4.3 must stage into a hidden directory on the SAME filesystem as " +
			"the cache, so the publishing rename is atomic")
	}
	// ...and publish by rename.
	if !strings.Contains(block, `mv -f "${staged}"`) {
		t.Error("Phase 4.3 must publish each package with an atomic rename")
	}
	// The completion receipt is the BOM-level barrier and must be written LAST:
	// a receipt that lands before the packages it vouches for is worse than none.
	receipt := strings.Index(block, ".staging-complete.json")
	lastRename := strings.LastIndex(block, `mv -f "${staged}"`)
	if receipt < 0 {
		t.Error("Phase 4.3 must write a completion receipt for the staged BOM")
	} else if receipt < lastRename {
		t.Error("the completion receipt is written BEFORE the packages are published; " +
			"it must be the last thing that happens or it vouches for nothing")
	}
}

// TestStagedArtifactsAreInvisibleToTheCacheScan pins the consumer-side half of
// the contract in the form the node-agent actually implements it: its cache scan
// reads the packages directory, skips directories, and requires a .tgz suffix.
// Both staging artifacts must fall outside that filter.
func TestStagedArtifactsAreInvisibleToTheCacheScan(t *testing.T) {
	cache := t.TempDir()
	if err := os.MkdirAll(filepath.Join(cache, ".incoming.4242"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, ".incoming.4242", "half.tgz"), []byte("partial"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, ".staging-complete.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, "workflow_1.2.298_linux_amd64.tgz"), []byte("real"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(cache)
	if err != nil {
		t.Fatal(err)
	}
	var visible []string
	for _, e := range entries {
		// This mirrors node_agent heartbeat.go localPackageCacheProofs.
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tgz") {
			continue
		}
		visible = append(visible, e.Name())
	}

	if len(visible) != 1 || visible[0] != "workflow_1.2.298_linux_amd64.tgz" {
		t.Fatalf("cache scan saw %v; it must see only the published package — the "+
			"staging directory and the receipt must be structurally invisible", visible)
	}
}

// TestStagingBlockIsValidShell keeps the rewritten block from shipping with a
// syntax error, which would fail the join at a point where recovery is manual.
func TestStagingBlockIsValidShell(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}
	cmd := exec.Command("bash", "-n")
	cmd.Stdin = strings.NewReader(joinScriptForTest())
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("join script is not valid bash: %v\n%s", err, out)
	}
}
