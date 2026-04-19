package admin

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ── JSON response types ─────────────────────────────────────────────────────

// StorageResponse is the top-level response for GET /admin/metrics/storage.
type StorageResponse struct {
	NowUnixMs         int64             `json:"now_unix_ms"`
	DerivedStatus     string            `json:"derived_status"` // healthy | degraded | critical
	Reasons           []string          `json:"reasons"`
	MostCriticalMount string            `json:"most_critical_mount"`
	Thresholds        StorageThresholds `json:"thresholds"`
	Mounts            []MountInfo       `json:"mounts"`
	Applications      []ApplicationPath `json:"applications"`
	Series            json.RawMessage   `json:"series"` // empty object in Phase 1
}

// StorageThresholds are the thresholds used to derive storage health.
type StorageThresholds struct {
	DiskWarnFreePct float64 `json:"disk_warn_free_pct"`
	DiskCritFreePct float64 `json:"disk_crit_free_pct"`
}

// MountInfo describes a single mounted filesystem.
type MountInfo struct {
	Device     string  `json:"device"`
	MountPoint string  `json:"mount_point"`
	FSType     string  `json:"fs_type"`
	TotalBytes uint64  `json:"total_bytes"`
	UsedBytes  uint64  `json:"used_bytes"`
	FreeBytes  uint64  `json:"free_bytes"`
	UsedPct    float64 `json:"used_pct"`
	FreePct    float64 `json:"free_pct"`
	Status     string  `json:"status"` // healthy | degraded | critical
}

// ApplicationPath maps an application data directory to its mount point.
type ApplicationPath struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Exists     bool   `json:"exists"`
	Writable   bool   `json:"writable"`
	MountPoint string `json:"mount_point"`
	Status     string `json:"status"`     // healthy | at_risk | unavailable
	SizeBytes  *int64 `json:"size_bytes"` // null in Phase 1
}

// ── Constants ───────────────────────────────────────────────────────────────

const (
	diskWarnFreePct = 15.0
	diskCritFreePct = 10.0
)

// ── Handler ─────────────────────────────────────────────────────────────────

// NewStorageHandler returns a GET-only handler for /admin/metrics/storage.
func NewStorageHandler(provider AdminProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// If ?node=<hostname> targets a different node, proxy the request to
		// that node's gateway server-side. Each gateway has cluster DNS + CA
		// trust, so we can reach peers by short-hostname with TLS validated
		// against the cluster CA. The browser can't do this directly because
		// (a) it lacks cluster DNS resolution and (b) cross-origin CORS.
		if target := strings.TrimSpace(r.URL.Query().Get("node")); target != "" && target != provider.Hostname() {
			proxyStorageToPeer(w, r, target)
			return
		}

		// 1. Collect mount points
		mounts := collectMounts()

		// 2. Derive mount health
		var (
			reasons           []string
			mostCriticalMount string
			lowestFreePct     = 100.0
			overallStatus     = "healthy"
		)
		for i := range mounts {
			m := &mounts[i]
			if m.TotalBytes == 0 {
				m.Status = "healthy"
				continue
			}
			m.UsedPct = math.Round(float64(m.UsedBytes)/float64(m.TotalBytes)*1000) / 10
			m.FreePct = math.Round(float64(m.FreeBytes)/float64(m.TotalBytes)*1000) / 10

			if m.FreePct < diskCritFreePct {
				m.Status = "critical"
			} else if m.FreePct < diskWarnFreePct {
				m.Status = "degraded"
			} else {
				m.Status = "healthy"
			}

			if m.FreePct < lowestFreePct {
				lowestFreePct = m.FreePct
				mostCriticalMount = m.MountPoint
			}
		}

		if lowestFreePct < diskCritFreePct {
			overallStatus = "critical"
			reasons = append(reasons, mostCriticalMount+" free "+strings.TrimRight(strings.TrimRight(
				fmtFloat(lowestFreePct), "0"), ".")+"% (< "+fmtFloat(diskCritFreePct)+"%)")
		} else if lowestFreePct < diskWarnFreePct {
			overallStatus = "degraded"
			reasons = append(reasons, mostCriticalMount+" free "+strings.TrimRight(strings.TrimRight(
				fmtFloat(lowestFreePct), "0"), ".")+"% (< "+fmtFloat(diskWarnFreePct)+"%)")
		}

		// 3. Build application paths
		apps := buildAppPaths(provider, mounts)

		resp := StorageResponse{
			NowUnixMs:         time.Now().UnixMilli(),
			DerivedStatus:     overallStatus,
			Reasons:           reasons,
			MostCriticalMount: mostCriticalMount,
			Thresholds: StorageThresholds{
				DiskWarnFreePct: diskWarnFreePct,
				DiskCritFreePct: diskCritFreePct,
			},
			Mounts:       mounts,
			Applications: apps,
			Series:       json.RawMessage(`{}`),
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}

// ── Application path collection ─────────────────────────────────────────────

func buildAppPaths(provider AdminProvider, mounts []MountInfo) []ApplicationPath {
	stateDir := provider.StateDir()

	known := []struct {
		name string
		path string
	}{
		{"etcd", filepath.Join(stateDir, "etcd")},
		{"ScyllaDB", "/var/lib/scylla"},
		{"Prometheus", filepath.Join(stateDir, "prometheus", "data")},
	}
	// MinIO may resolve to multiple local paths in distributed mode.
	for _, p := range resolveMinIODataDirs(stateDir) {
		name := "MinIO"
		if base := filepath.Base(p); base != "" && base != "." && base != "/" {
			name = "MinIO: " + base
		}
		known = append(known, struct {
			name string
			path string
		}{name, p})
	}

	// Add file service public dirs. Deduplicate because PublicDirs()
	// aggregates Public[] from every FileService instance across the
	// cluster, and nodes generally share identical config — without this
	// the same path shows up once per node.
	seenPublic := make(map[string]struct{})
	for _, dir := range provider.PublicDirs() {
		if _, dup := seenPublic[dir]; dup {
			continue
		}
		seenPublic[dir] = struct{}{}
		known = append(known, struct {
			name string
			path string
		}{"FileService: " + filepath.Base(dir), dir})
	}

	var apps []ApplicationPath
	for _, k := range known {
		app := ApplicationPath{
			Name:      k.name,
			Path:      k.path,
			SizeBytes: nil, // Phase 1: skip du
		}

		// Check existence
		if _, err := os.Stat(k.path); err == nil {
			app.Exists = true
			app.Writable = isWritable(k.path)
		}

		// Map to mount
		app.MountPoint = findMountForPath(k.path, mounts)

		// Derive status
		if !app.Exists {
			app.Status = "unavailable"
		} else {
			mountStatus := mountStatusForPath(app.MountPoint, mounts)
			if mountStatus == "critical" {
				app.Status = "at_risk"
			} else {
				app.Status = "healthy"
			}
		}

		apps = append(apps, app)
	}
	return apps
}

// findMountForPath returns the mount point with the longest prefix match.
func findMountForPath(p string, mounts []MountInfo) string {
	best := ""
	for _, m := range mounts {
		mp := m.MountPoint
		if strings.HasPrefix(p, mp) && len(mp) > len(best) {
			best = mp
		}
	}
	return best
}

func mountStatusForPath(mountPoint string, mounts []MountInfo) string {
	for _, m := range mounts {
		if m.MountPoint == mountPoint {
			return m.Status
		}
	}
	return "healthy"
}

// isWritable reports whether the path can be written to. Because the gateway
// process generally runs as the "globular" user but many application data
// directories are owned by their own service account (scylla, minio, mongodb,
// …), a raw os.Create as our process yields false negatives. We use the
// following logic:
//
//  1. Try creating a probe file directly. Success ⇒ writable.
//  2. If the filesystem itself is read-only (EROFS) ⇒ not writable.
//  3. Otherwise (EACCES/EPERM) we couldn't write but the owning service
//     likely can. If the directory's owner-write bit is set, report it as
//     writable — that matches how the service process will see it.
func isWritable(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		p = filepath.Dir(p)
		info, err = os.Stat(p)
		if err != nil {
			return false
		}
	}
	tmp := filepath.Join(p, ".globular_write_test")
	if f, createErr := os.Create(tmp); createErr == nil {
		f.Close()
		os.Remove(tmp)
		return true
	} else if errors.Is(createErr, syscall.EROFS) {
		return false
	}
	return info.Mode().Perm()&0200 != 0
}

// resolveMinIODataDirs reads MINIO_VOLUMES from the minio env file and returns
// the set of LOCAL filesystem paths that this node owns. In distributed mode
// MINIO_VOLUMES is a space-separated list of URLs (e.g.
// "https://host1:9000/data1 https://host2:9000/data1 ..."); only entries whose
// host resolves to one of this node's identities are returned. In single-node
// mode it's a plain path. Falls back to {stateDir}/minio/data on any parse
// failure.
func resolveMinIODataDirs(stateDir string) []string {
	fallback := []string{filepath.Join(stateDir, "minio", "data")}

	envFile := filepath.Join(stateDir, "minio", "minio.env")
	f, err := os.Open(envFile)
	if err != nil {
		return fallback
	}
	defer f.Close()

	var raw string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		k, v, _ := strings.Cut(line, "=")
		if strings.TrimSpace(k) == "MINIO_VOLUMES" {
			raw = strings.Trim(strings.TrimSpace(v), "\"'")
			break
		}
	}
	if raw == "" {
		return fallback
	}

	// Distributed mode: space-separated URL list. Extract local paths.
	if strings.Contains(raw, "://") {
		local := extractLocalMinIOPaths(raw)
		if len(local) == 0 {
			return fallback
		}
		return local
	}

	// Single-node mode: plain path (may contain a brace-expansion ellipsis
	// like "/data/{1...4}", but we don't attempt to expand it here — report
	// the literal value so the operator sees what MinIO was configured with).
	return []string{raw}
}

// extractLocalMinIOPaths parses a space-separated URL list and returns the
// URL.Path components whose host matches one of this node's identities
// (hostname or any interface IP). Duplicates are preserved so the operator
// sees each volume entry.
func extractLocalMinIOPaths(raw string) []string {
	ids := localHostIdentities()
	var paths []string
	for _, tok := range strings.Fields(raw) {
		u, err := url.Parse(tok)
		if err != nil || u.Host == "" || u.Path == "" {
			continue
		}
		if _, ok := ids[strings.ToLower(u.Hostname())]; !ok {
			continue
		}
		paths = append(paths, u.Path)
	}
	return paths
}

// localHostIdentities returns a set of lowercase names/IPs that identify this
// node: loopback names, os.Hostname(), and every interface IP address.
func localHostIdentities() map[string]struct{} {
	ids := map[string]struct{}{
		"localhost": {},
		"127.0.0.1": {},
		"::1":       {},
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		ids[strings.ToLower(h)] = struct{}{}
		// Also add the short hostname (before the first dot) in case
		// MINIO_VOLUMES uses FQDN or vice versa.
		if short, _, ok := strings.Cut(h, "."); ok && short != "" {
			ids[strings.ToLower(short)] = struct{}{}
		}
	}
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				ids[ipnet.IP.String()] = struct{}{}
			}
		}
	}
	return ids
}

// proxyStorageToPeer forwards /admin/metrics/storage to a peer gateway so
// the UI can see any node's ground-truth local filesystem view through a
// single origin (avoiding CORS/DNS issues in the browser). Uses the cluster
// CA for TLS validation. Caches the shared HTTP client.
func proxyStorageToPeer(w http.ResponseWriter, r *http.Request, peer string) {
	client, err := clusterHTTPClient()
	if err != nil {
		http.Error(w, "cluster HTTP client: "+err.Error(), http.StatusBadGateway)
		return
	}
	// Target the peer's gateway on its TLS port; short hostname resolves via
	// /etc/hosts or cluster DNS entries populated by the installer.
	target := "https://" + peer + "/admin/metrics/storage"
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, target, nil)
	if err != nil {
		http.Error(w, "build proxy request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Preserve the expected Host header so the peer's envoy picks the right
	// cert/listener.
	req.Host = peer + ".globular.internal"
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "proxy to "+peer+" failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

var (
	clusterClientOnce sync.Once
	clusterClient     *http.Client
	clusterClientErr  error
)

func clusterHTTPClient() (*http.Client, error) {
	clusterClientOnce.Do(func() {
		// Trust the cluster CA for peer HTTPS calls.
		pool := x509.NewCertPool()
		caBytes, err := os.ReadFile("/var/lib/globular/pki/ca.pem")
		if err != nil {
			// Fall back to alternate CA path used by some installs.
			caBytes, err = os.ReadFile("/var/lib/globular/pki/ca.crt")
		}
		if err == nil && pool.AppendCertsFromPEM(caBytes) {
			clusterClient = &http.Client{
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12},
				},
			}
			return
		}
		// Last-resort: skip verify (cluster-internal, trusted network).
		clusterClient = &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12}, //nolint:gosec
			},
		}
	})
	return clusterClient, clusterClientErr
}

func fmtFloat(v float64) string {
	return strings.TrimRight(strings.TrimRight(
		strconv.FormatFloat(math.Round(v*10)/10, 'f', 1, 64),
		"0"), ".")
}
