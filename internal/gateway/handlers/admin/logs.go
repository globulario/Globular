package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// ── Allowlisted systemd units ────────────────────────────────────────────────

var allowedUnits = map[string]bool{
	"globular-node-agent.service":         true,
	"globular-cluster-controller.service": true,
	"globular-etcd.service":               true,
	"etcd.service":                        true,
	"envoy.service":                       true,
	"globular-envoy.service":              true,
	"globular-xds.service":                true,
	"prometheus.service":                  true,
	"globular-prometheus.service":         true,
	"globular-dns.service":                true,
	"globular-scylladb.service":           true,
	"scylla-server.service":               true,
	"globular-minio.service":              true,
	"minio.service":                       true,
}

// ── Journal reader interface ────────────────────────────────────────────────

// JournalReader reads systemd journal entries for a given unit.
// The implementation lives outside internal/gateway (in internal/journal)
// to satisfy the gateway no-exec lint rule.
type JournalReader interface {
	// ReadUnit returns the last N lines from the unit's journal.
	ReadUnit(ctx context.Context, unit string, lines int, sinceSec int) JournalResult
}

// JournalResult holds the output of a journal read.
type JournalResult struct {
	Unit      string
	Lines     []string
	Truncated bool
	Error     string
}

// ── JSON response type ──────────────────────────────────────────────────────

// ServiceLogsResponse is the response for GET /admin/service/logs.
type ServiceLogsResponse struct {
	Unit      string   `json:"unit"`
	Lines     []string `json:"lines"`
	Truncated bool     `json:"truncated"`
	Timestamp int64    `json:"timestamp"`
	Error     string   `json:"error,omitempty"`
}

// ── Constants ───────────────────────────────────────────────────────────────

const (
	defaultLogLines = 100
	maxLogLines     = 500
	defaultSinceSec = 3600 // 1 hour
)

// ── Handler ─────────────────────────────────────────────────────────────────

// NewLogsHandler returns a handler for GET /admin/service/logs.
//
// Query parameters:
//   - unit   (required) systemd unit name — must be in the allowlist
//   - lines  (optional) number of lines, default 100, max 500
//   - since  (optional) seconds to look back, default 3600
func NewLogsHandler(reader JournalReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		unit := r.URL.Query().Get("unit")
		if unit == "" {
			writeLogsError(w, "", "unit parameter is required")
			return
		}

		if !allowedUnits[unit] {
			writeLogsError(w, unit, "unit not in allowlist")
			return
		}

		lines := defaultLogLines
		if v := r.URL.Query().Get("lines"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				lines = n
			}
		}
		if lines < 1 {
			lines = 1
		}
		if lines > maxLogLines {
			lines = maxLogLines
		}

		sinceSec := defaultSinceSec
		if v := r.URL.Query().Get("since"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				sinceSec = n
			}
		}

		result := reader.ReadUnit(r.Context(), unit, lines, sinceSec)

		resp := ServiceLogsResponse{
			Unit:      result.Unit,
			Lines:     result.Lines,
			Truncated: result.Truncated,
			Timestamp: time.Now().UnixMilli(),
			Error:     result.Error,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}

func writeLogsError(w http.ResponseWriter, unit, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp := ServiceLogsResponse{
		Unit:      unit,
		Error:     errMsg,
		Timestamp: time.Now().UnixMilli(),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}
