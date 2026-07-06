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
	t.Cleanup(func() {
		joinBundleInstallerCandidates = oldCandidates
		joinBundlePackagesDir = oldPackagesDir
		joinBundleWorkflowsDir = oldWorkflowsDir
		joinBundleReleaseIndex = oldReleaseIndex
	})
	joinBundleInstallerCandidates = []string{installerPath}
	joinBundlePackagesDir = packagesDir
	joinBundleWorkflowsDir = workflowsDir
	joinBundleReleaseIndex = releaseIndexPath

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
	t.Cleanup(func() {
		joinBundleInstallerCandidates = oldCandidates
		joinBundlePackagesDir = oldPackagesDir
		joinBundleWorkflowsDir = oldWorkflowsDir
		joinBundleReleaseIndex = oldReleaseIndex
	})
	joinBundleInstallerCandidates = []string{installerPath}
	joinBundlePackagesDir = packagesDir
	joinBundleWorkflowsDir = workflowsDir
	joinBundleReleaseIndex = releaseIndexPath

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
