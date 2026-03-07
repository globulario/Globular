package admin

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	dataDir := provider.DataDir()

	known := []struct {
		name string
		path string
	}{
		{"etcd", filepath.Join(dataDir, "etcd")},
		{"ScyllaDB", "/var/lib/scylla"},
		{"Prometheus", filepath.Join(dataDir, "prometheus-data")},
		{"MinIO", filepath.Join(stateDir, "minio", "data")},
	}

	// Add file service public dirs
	for _, dir := range provider.PublicDirs() {
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

// isWritable tests if the path is writable by attempting to create a temp file.
func isWritable(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		// For files, check the parent directory
		p = filepath.Dir(p)
	}
	tmp := filepath.Join(p, ".globular_write_test")
	f, err := os.Create(tmp)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(tmp)
	return true
}

func fmtFloat(v float64) string {
	return strings.TrimRight(strings.TrimRight(
		strconv.FormatFloat(math.Round(v*10)/10, 'f', 1, 64),
		"0"), ".")
}
