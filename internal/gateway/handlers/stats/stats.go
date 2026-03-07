// Package stats exposes a /stats HTTP endpoint that returns real-time
// host metrics (CPU, memory, disk, network) plus Go runtime stats.
package stats

import (
	"encoding/json"
	"net/http"
	"time"
)

// ── JSON response types ─────────────────────────────────────────────────────

type StatsResponse struct {
	Hostname  string    `json:"hostname"`
	UptimeSec float64   `json:"uptimeSec"`
	CPU       CPUStats  `json:"cpu"`
	Memory    MemStats  `json:"memory"`
	Disk      DiskStats `json:"disk"`
	Network   NetStats  `json:"network"`
	Go        GoStats   `json:"go"`
}

type CPUStats struct {
	Count    int       `json:"count"`
	UsagePct float64   `json:"usagePct"`
	PerCore  []float64 `json:"perCore"`
}

type MemStats struct {
	TotalBytes uint64  `json:"totalBytes"`
	UsedBytes  uint64  `json:"usedBytes"`
	UsedPct    float64 `json:"usedPct"`
}

type DiskStats struct {
	TotalBytes uint64  `json:"totalBytes"`
	UsedBytes  uint64  `json:"usedBytes"`
	FreePct    float64 `json:"freePct"`
	Path       string  `json:"path"`
}

type NetStats struct {
	RxBytes uint64 `json:"rxBytes"`
	TxBytes uint64 `json:"txBytes"`
}

type GoStats struct {
	Goroutines int    `json:"goroutines"`
	HeapAlloc  uint64 `json:"heapAllocBytes"`
	GCPauseNs  uint64 `json:"gcPauseNs"`
	NumGC      uint32 `json:"numGC"`
}

// ── TimeProvider ────────────────────────────────────────────────────────────

// TimeProvider exposes the process start time (satisfied by *Globule).
type TimeProvider interface {
	StartTime() time.Time
}

// ── Handler ─────────────────────────────────────────────────────────────────

// NewStatsHandler returns a GET-only handler that collects host metrics and
// writes them as JSON.
func NewStatsHandler(tp TimeProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		resp := collect(tp)
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}

// ── Mount ───────────────────────────────────────────────────────────────────

// Deps holds the pre-wrapped handlers for route registration.
type Deps struct {
	Stats http.Handler
}

// Mount registers the /stats route on the mux.
func Mount(mux *http.ServeMux, d Deps) {
	if d.Stats != nil {
		mux.Handle("/stats", d.Stats)
	}
}
