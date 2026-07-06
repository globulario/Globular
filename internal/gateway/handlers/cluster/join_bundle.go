package cluster

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var joinBundleInstallerCandidates = []string{
	"/usr/lib/globular/bin/globular-installer",
	"/var/lib/globular/bin/globular-installer",
	"/usr/local/bin/globular-installer",
}

var (
	joinBundlePackagesDir  = "/var/lib/globular/packages"
	joinBundleWorkflowsDir = "/var/lib/globular/workflows"
	joinBundleReleaseIndex = "/var/lib/globular/release-index.json"
)

// bundleModTime is a fixed timestamp stamped into every tar header so the bundle
// is byte-reproducible: the tarball and its separately-fetched .sha256 are
// assembled from independent HTTP requests and must produce identical bytes.
var bundleModTime = time.Unix(0, 0).UTC()

// NewJoinBundleHandler serves a controller-local bootstrap bundle for Day-1
// joins when the requested platform release tarball is not published on GitHub.
//
// The bundle is generated from the founder node's installed state instead of a
// repository checkout, so it works for local/dev builds too. The join script
// only needs globular-installer, release-index.json, packages, and workflow
// definitions.
//
// Two routes are served under /join/bundle/:
//   - <tarball>          → the assembled gzip tarball
//   - <tarball>.sha256   → the sha256 of that tarball (sha256sum format)
//
// The join script fetches both and verifies the checksum, mirroring the GitHub
// release path. NOTE: the checksum guards against corruption/truncation. It is
// NOT, on its own, MITM protection when the bundle is fetched over `curl -k`
// (an attacker serving the tarball serves a matching checksum too). Pinning the
// cluster CA for this transport is a tracked follow-up; see the join trust
// bootstrap. release-index.json is MANDATORY — it is the platform-release BOM
// truth object, and a joining node must not come up without it.
func NewJoinBundleHandler(platformVersion string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		version := strings.TrimSpace(activePlatformVersion(platformVersion))
		if version == "" || version == "0.0.0-dev" || version == "0.0.0-detect" {
			http.Error(w, "active platform version unavailable", http.StatusServiceUnavailable)
			return
		}

		wantName := fmt.Sprintf("globular-%s-linux-amd64.tar.gz", version)
		if !strings.HasPrefix(r.URL.Path, "/join/bundle/") {
			http.NotFound(w, r)
			return
		}
		gotName := strings.TrimPrefix(r.URL.Path, "/join/bundle/")
		wantChecksum := false
		switch gotName {
		case wantName:
			wantChecksum = false
		case wantName + ".sha256":
			wantChecksum = true
		default:
			http.NotFound(w, r)
			return
		}

		bundleRoot := fmt.Sprintf("globular-%s-linux-amd64", version)
		data, statusCode, err := assembleJoinBundle(bundleRoot)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		sum := sha256.Sum256(data)
		sumHex := hex.EncodeToString(sum[:])

		if wantChecksum {
			// sha256sum-compatible line: "<hex>  <filename>".
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "%s  %s\n", sumHex, wantName)
			return
		}

		w.Header().Set("Content-Type", "application/gzip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", wantName))
		w.Header().Set("X-Content-SHA256", sumHex)
		_, _ = w.Write(data)
	})
}

// assembleJoinBundle builds the Day-1 fallback tarball into a buffer so the same
// bytes back both the tarball response and its .sha256. Assembly is deterministic
// (zero gzip mtime, fixed tar mtimes, zeroed uid/gid, lexical walk order) so the
// separately-fetched checksum always matches the served bytes. On any missing
// input it returns an HTTP status and error for the handler to surface.
func assembleJoinBundle(bundleRoot string) ([]byte, int, error) {
	installerPath := ""
	for _, candidate := range joinBundleInstallerCandidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			installerPath = candidate
			break
		}
	}
	if installerPath == "" {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("globular-installer not found on controller")
	}

	packagesDir := joinBundlePackagesDir
	workflowsDir := joinBundleWorkflowsDir
	releaseIndexPath := joinBundleReleaseIndex

	if info, err := os.Stat(packagesDir); err != nil || !info.IsDir() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("packages directory unavailable on controller")
	}
	if info, err := os.Stat(workflowsDir); err != nil || !info.IsDir() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("workflows directory unavailable on controller")
	}
	// release-index.json is the platform-release BOM truth object — REQUIRED.
	// Refuse to serve a BOM-less bundle rather than let a node join blind.
	if info, err := os.Stat(releaseIndexPath); err != nil || info.IsDir() {
		return nil, http.StatusServiceUnavailable,
			fmt.Errorf("release-index.json unavailable on controller; refusing to serve a BOM-less join bundle")
	}

	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	if err := writeTarFile(tw, installerPath, bundleRoot+"/globular-installer", 0o755); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("bundle installer: %w", err)
	}
	if err := writeTarFile(tw, releaseIndexPath, bundleRoot+"/release-index.json", 0o644); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("bundle release-index: %w", err)
	}
	if err := writeTarTree(tw, packagesDir, bundleRoot+"/packages"); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("bundle packages: %w", err)
	}
	if err := writeTarTree(tw, workflowsDir, bundleRoot+"/workflows"); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("bundle workflows: %w", err)
	}
	if err := tw.Close(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("bundle tar close: %w", err)
	}
	if err := gzw.Close(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("bundle gzip close: %w", err)
	}
	return buf.Bytes(), http.StatusOK, nil
}

func writeTarTree(tw *tar.Writer, srcDir, dstRoot string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == srcDir {
			return nil
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(filepath.Join(dstRoot, rel))
		if info.IsDir() {
			hdr := &tar.Header{
				Name:     name + "/",
				Mode:     0o755,
				Typeflag: tar.TypeDir,
				ModTime:  bundleModTime,
			}
			return tw.WriteHeader(hdr)
		}
		return writeTarFile(tw, path, name, info.Mode().Perm())
	})
}

func writeTarFile(tw *tar.Writer, srcPath, dstName string, perm os.FileMode) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	hdr.Name = filepath.ToSlash(dstName)
	hdr.Mode = int64(perm)
	hdr.Uid = 0
	hdr.Gid = 0
	hdr.Uname = ""
	hdr.Gname = ""
	// Fixed mtime for byte-reproducibility (see bundleModTime).
	hdr.ModTime = bundleModTime
	hdr.AccessTime = time.Time{}
	hdr.ChangeTime = time.Time{}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(tw, f)
	return err
}
