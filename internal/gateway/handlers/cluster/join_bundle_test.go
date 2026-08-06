package cluster

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestJoinBundleHandler_ServesBootstrapBundle(t *testing.T) {
	tmp := t.TempDir()

	installerPath := filepath.Join(tmp, "globular-installer")
	if err := os.WriteFile(installerPath, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write installer: %v", err)
	}
	packagesDir := filepath.Join(tmp, "packages")
	if err := os.MkdirAll(packagesDir, 0o755); err != nil {
		t.Fatalf("mkdir packages: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packagesDir, "node-agent_1.2.3_linux_amd64.tgz"), []byte("pkg"), 0o644); err != nil {
		t.Fatalf("write package: %v", err)
	}
	workflowsDir := filepath.Join(tmp, "workflows")
	if err := os.MkdirAll(workflowsDir, 0o755); err != nil {
		t.Fatalf("mkdir workflows: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workflowsDir, "node.join.yaml"), []byte("kind: workflow\n"), 0o644); err != nil {
		t.Fatalf("write workflow: %v", err)
	}
	releaseIndexPath := filepath.Join(tmp, "release-index.json")
	if err := os.WriteFile(releaseIndexPath, []byte(`{"platform_release":"1.2.3"}`), 0o644); err != nil {
		t.Fatalf("write release-index: %v", err)
	}

	oldCandidates := joinBundleInstallerCandidates
	oldPackagesDir := joinBundlePackagesDir
	oldWorkflowsDir := joinBundleWorkflowsDir
	oldReleaseIndex := joinBundleReleaseIndex
	oldReleasesDir := joinBundleReleasesDir
	t.Cleanup(func() {
		joinBundleInstallerCandidates = oldCandidates
		joinBundlePackagesDir = oldPackagesDir
		joinBundleWorkflowsDir = oldWorkflowsDir
		joinBundleReleaseIndex = oldReleaseIndex
		joinBundleReleasesDir = oldReleasesDir
	})
	joinBundleInstallerCandidates = []string{installerPath}
	joinBundlePackagesDir = packagesDir
	joinBundleWorkflowsDir = workflowsDir
	joinBundleReleaseIndex = releaseIndexPath
	joinBundleReleasesDir = filepath.Join(tmp, "releases")

	version := activePlatformVersion("1.2.3")
	name := "globular-" + version + "-linux-amd64.tar.gz"
	req := httptest.NewRequest(http.MethodGet, "/join/bundle/"+name, nil)
	rr := httptest.NewRecorder()

	NewJoinBundleHandler("1.2.3").ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d body=%s", rr.Code, rr.Body.String())
	}

	gzr, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	want := map[string]bool{
		"globular-" + version + "-linux-amd64/globular-installer":                        false,
		"globular-" + version + "-linux-amd64/packages/node-agent_1.2.3_linux_amd64.tgz": false,
		"globular-" + version + "-linux-amd64/workflows/node.join.yaml":                  false,
		"globular-" + version + "-linux-amd64/release-index.json":                        false,
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("tar next: %v", err)
		}
		if _, ok := want[hdr.Name]; ok {
			want[hdr.Name] = true
		}
	}
	for path, seen := range want {
		if !seen {
			t.Errorf("bundle missing %s", path)
		}
	}
}

// bundleTestFixture wires the package-level source paths to a temp dir with all
// required inputs present, and returns the version + tarball name. Pass
// includeReleaseIndex=false to exercise the BOM-required guard.
func bundleTestFixture(t *testing.T, includeReleaseIndex bool) (version, name string) {
	t.Helper()
	tmp := t.TempDir()

	installerPath := filepath.Join(tmp, "globular-installer")
	if err := os.WriteFile(installerPath, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write installer: %v", err)
	}
	packagesDir := filepath.Join(tmp, "packages")
	if err := os.MkdirAll(packagesDir, 0o755); err != nil {
		t.Fatalf("mkdir packages: %v", err)
	}
	if err := os.WriteFile(filepath.Join(packagesDir, "node-agent_1.2.3_linux_amd64.tgz"), []byte("pkg"), 0o644); err != nil {
		t.Fatalf("write package: %v", err)
	}
	workflowsDir := filepath.Join(tmp, "workflows")
	if err := os.MkdirAll(workflowsDir, 0o755); err != nil {
		t.Fatalf("mkdir workflows: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workflowsDir, "node.join.yaml"), []byte("kind: workflow\n"), 0o644); err != nil {
		t.Fatalf("write workflow: %v", err)
	}
	releaseIndexPath := filepath.Join(tmp, "release-index.json")
	if includeReleaseIndex {
		if err := os.WriteFile(releaseIndexPath, []byte(`{"platform_release":"1.2.3"}`), 0o644); err != nil {
			t.Fatalf("write release-index: %v", err)
		}
	}

	oldCandidates := joinBundleInstallerCandidates
	oldPackagesDir := joinBundlePackagesDir
	oldWorkflowsDir := joinBundleWorkflowsDir
	oldReleaseIndex := joinBundleReleaseIndex
	oldReleasesDir := joinBundleReleasesDir
	t.Cleanup(func() {
		joinBundleInstallerCandidates = oldCandidates
		joinBundlePackagesDir = oldPackagesDir
		joinBundleWorkflowsDir = oldWorkflowsDir
		joinBundleReleaseIndex = oldReleaseIndex
		joinBundleReleasesDir = oldReleasesDir
	})
	joinBundleInstallerCandidates = []string{installerPath}
	joinBundlePackagesDir = packagesDir
	joinBundleWorkflowsDir = workflowsDir
	joinBundleReleaseIndex = releaseIndexPath
	joinBundleReleasesDir = filepath.Join(tmp, "releases")

	version = activePlatformVersion("1.2.3")
	name = "globular-" + version + "-linux-amd64.tar.gz"
	return version, name
}

func serveBundle(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	NewJoinBundleHandler("1.2.3").ServeHTTP(rr, req)
	return rr
}

// TestJoinBundleHandler_HeadReportsSizeWithoutBody pins the contract the Day-1
// join relies on at phase 4.1: the bundle route answers HEAD with a usable
// Content-Length and no body.
//
// Regression. The handler used to reject every non-GET method with 405. The
// join script probes this endpoint with `curl -sfL -I` to announce the download
// size, and it runs under `set -euo pipefail` — so the 405 became curl exit 22,
// pipefail promoted it past the exit-0 awk, and the assignment aborted the join
// at phase 4. Every joiner in the 5-node simulation died there on v1.2.295,
// with the ~750 MB download the probe was only describing never attempted.
// The script no longer lets the probe be fatal; this test pins the other half,
// so HEAD keeps working and the size stays reportable.
func TestJoinBundleHandler_HeadReportsSizeWithoutBody(t *testing.T) {
	_, name := bundleTestFixture(t, true)

	getRR := serveBundle(t, "/join/bundle/"+name)
	if getRR.Code != http.StatusOK {
		t.Fatalf("GET: unexpected status %d body=%s", getRR.Code, getRR.Body.String())
	}

	req := httptest.NewRequest(http.MethodHead, "/join/bundle/"+name, nil)
	headRR := httptest.NewRecorder()
	NewJoinBundleHandler("1.2.3").ServeHTTP(headRR, req)

	if headRR.Code != http.StatusOK {
		t.Fatalf("HEAD must be allowed (join phase 4.1 probes it), got %d body=%s",
			headRR.Code, headRR.Body.String())
	}

	// The probe greps Content-Length out of the headers; without it the size
	// is unreportable and the feature that motivated the probe is dead.
	gotLen := headRR.Header().Get("Content-Length")
	wantLen := strconv.Itoa(getRR.Body.Len())
	if gotLen != wantLen {
		t.Errorf("HEAD Content-Length = %q, want %q (must match the GET body size)",
			gotLen, wantLen)
	}

	if headRR.Body.Len() != 0 {
		t.Errorf("HEAD must not return a body, got %d bytes", headRR.Body.Len())
	}
}

// TestJoinBundleHandler_RejectsMutatingMethods keeps the widened method guard
// honest: allowing HEAD must not have opened the endpoint to anything else.
func TestJoinBundleHandler_RejectsMutatingMethods(t *testing.T) {
	_, name := bundleTestFixture(t, true)

	for _, method := range []string{
		http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
	} {
		req := httptest.NewRequest(method, "/join/bundle/"+name, nil)
		rr := httptest.NewRecorder()
		NewJoinBundleHandler("1.2.3").ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", method, rr.Code)
		}
	}
}

// TestJoinBundleHandler_ChecksumMatchesTarball verifies the .sha256 route returns
// the sha256 of the exact bytes the tarball route serves (deterministic assembly).
func TestJoinBundleHandler_ChecksumMatchesTarball(t *testing.T) {
	_, name := bundleTestFixture(t, true)

	tarRR := serveBundle(t, "/join/bundle/"+name)
	if tarRR.Code != http.StatusOK {
		t.Fatalf("tarball: unexpected status %d body=%s", tarRR.Code, tarRR.Body.String())
	}
	sumRR := serveBundle(t, "/join/bundle/"+name+".sha256")
	if sumRR.Code != http.StatusOK {
		t.Fatalf("sha256: unexpected status %d body=%s", sumRR.Code, sumRR.Body.String())
	}

	actual := sha256.Sum256(tarRR.Body.Bytes())
	wantLine := hex.EncodeToString(actual[:]) + "  " + name + "\n"
	if got := sumRR.Body.String(); got != wantLine {
		t.Fatalf(".sha256 mismatch:\n got=%q\nwant=%q", got, wantLine)
	}
}

// TestJoinBundleHandler_RequiresReleaseIndex verifies the handler refuses to serve
// a BOM-less bundle: release-index.json is mandatory.
func TestJoinBundleHandler_RequiresReleaseIndex(t *testing.T) {
	_, name := bundleTestFixture(t, false)

	rr := serveBundle(t, "/join/bundle/"+name)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when release-index.json is absent, got %d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "release-index.json") {
		t.Errorf("error body should name the missing release-index.json, got %q", rr.Body.String())
	}
}

// TestJoinBundleHandler_ServedBytesMatchStoredArtifact pins the provenance
// contract (docs/design/package-identity-single-authority.md §4): the served
// join bundle is byte-identical to the artifact in the release-addressed
// store, and the .sha256 route serves the stored sidecar.
func TestJoinBundleHandler_ServedBytesMatchStoredArtifact(t *testing.T) {
	version, name := bundleTestFixture(t, true)

	rr := serveBundle(t, "/join/bundle/"+name)
	if rr.Code != http.StatusOK {
		t.Fatalf("tarball: unexpected status %d body=%s", rr.Code, rr.Body.String())
	}

	storedPath := filepath.Join(joinBundleReleasesDir, version, name)
	stored, err := os.ReadFile(storedPath)
	if err != nil {
		t.Fatalf("stored artifact missing after first request: %v", err)
	}
	if string(stored) != rr.Body.String() {
		t.Fatal("served bundle bytes differ from the stored release artifact")
	}

	sumStored, err := os.ReadFile(storedPath + ".sha256")
	if err != nil {
		t.Fatalf("stored checksum sidecar missing: %v", err)
	}
	sumRR := serveBundle(t, "/join/bundle/"+name+".sha256")
	if sumRR.Code != http.StatusOK || sumRR.Body.String() != string(sumStored) {
		t.Fatalf(".sha256 route must serve the stored sidecar: got %q want %q", sumRR.Body.String(), string(sumStored))
	}
}

// TestJoinBundleHandler_NoPerRequestReassembly pins that the bundle is
// materialized ONCE: mutating the live packages dir after the first request
// must NOT change the served bytes — the stored release artifact is the
// authority, not the controller's current disk state.
func TestJoinBundleHandler_NoPerRequestReassembly(t *testing.T) {
	_, name := bundleTestFixture(t, true)

	first := serveBundle(t, "/join/bundle/"+name)
	if first.Code != http.StatusOK {
		t.Fatalf("first request: unexpected status %d body=%s", first.Code, first.Body.String())
	}

	// Drift the live source dir — a per-request assembler would pick this up.
	if err := os.WriteFile(filepath.Join(joinBundlePackagesDir, "drift_9.9.9_linux_amd64.tgz"), []byte("drift"), 0o644); err != nil {
		t.Fatalf("write drift package: %v", err)
	}

	second := serveBundle(t, "/join/bundle/"+name)
	if second.Code != http.StatusOK {
		t.Fatalf("second request: unexpected status %d body=%s", second.Code, second.Body.String())
	}
	if first.Body.String() != second.Body.String() {
		t.Fatal("served bundle changed between requests — per-request reassembly detected; the stored artifact must be served")
	}
}

// TestJoinBundleHandler_PrePopulatedStoreWins pins the Day-0 lane: when the
// release store already holds an artifact for the active version (e.g. the
// original release tarball stored at install), the handler serves those exact
// bytes and never assembles from live disk.
func TestJoinBundleHandler_PrePopulatedStoreWins(t *testing.T) {
	version, name := bundleTestFixture(t, true)

	pre := []byte("original-release-artifact-bytes")
	dir := filepath.Join(joinBundleReleasesDir, version)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir releases: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), pre, 0o644); err != nil {
		t.Fatalf("pre-populate store: %v", err)
	}

	rr := serveBundle(t, "/join/bundle/"+name)
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status %d body=%s", rr.Code, rr.Body.String())
	}
	if rr.Body.String() != string(pre) {
		t.Fatal("handler did not serve the pre-populated stored artifact byte-for-byte")
	}
	// The checksum sidecar must pin the STORED bytes.
	sum := sha256.Sum256(pre)
	wantLine := hex.EncodeToString(sum[:]) + "  " + name + "\n"
	sumRR := serveBundle(t, "/join/bundle/"+name+".sha256")
	if sumRR.Body.String() != wantLine {
		t.Fatalf("checksum must pin stored bytes: got %q want %q", sumRR.Body.String(), wantLine)
	}
}
