package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// ── Job record ──────────────────────────────────────────────────────────────

// UpgradeJobRecord is a persistent record of an upgrade operation.
type UpgradeJobRecord struct {
	OperationID string              `json:"operation_id"`
	StartedAt   int64               `json:"started_at"`  // unix ms
	FinishedAt  int64               `json:"finished_at"` // unix ms, 0 if still running
	Status      string              `json:"status"`      // pending | running | success | failed | rolled_back
	Services    []UpgradeJobService `json:"services"`
	Error       string              `json:"error,omitempty"`
	IssuedBy    string              `json:"issued_by"`
}

// UpgradeJobService tracks a single service within a job record.
type UpgradeJobService struct {
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
}

// ── File-based store ────────────────────────────────────────────────────────

// JobStore persists upgrade job records as individual JSON files.
type JobStore struct {
	mu  sync.Mutex
	dir string // e.g. /var/lib/globular/data/upgrade_jobs
}

// NewJobStore creates a store backed by the given directory.
func NewJobStore(dataDir string) *JobStore {
	dir := filepath.Join(dataDir, "upgrade_jobs")
	return &JobStore{dir: dir}
}

// Save writes or updates a job record to disk.
func (s *JobStore) Save(rec UpgradeJobRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(s.dir, rec.OperationID+".json")
	return os.WriteFile(path, data, 0o600)
}

// Get reads a single job record by operation ID.
func (s *JobStore) Get(operationID string) (*UpgradeJobRecord, error) {
	path := filepath.Join(s.dir, operationID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rec UpgradeJobRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

// List returns all job records, sorted by StartedAt descending (newest first).
// Limits to maxResults; 0 means no limit.
func (s *JobStore) List(maxResults int) ([]UpgradeJobRecord, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var records []UpgradeJobRecord
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			continue
		}
		var rec UpgradeJobRecord
		if err := json.Unmarshal(data, &rec); err != nil {
			continue
		}
		records = append(records, rec)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].StartedAt > records[j].StartedAt
	})

	if maxResults > 0 && len(records) > maxResults {
		records = records[:maxResults]
	}
	return records, nil
}

// ── History handler ─────────────────────────────────────────────────────────

// UpgradeHistoryResponse is the JSON envelope for GET /admin/upgrades/history.
type UpgradeHistoryResponse struct {
	Jobs []UpgradeJobRecord `json:"jobs"`
}

// NewUpgradeHistoryHandler returns a GET handler for /admin/upgrades/history.
// Query params: ?limit=N (default 50)
func NewUpgradeHistoryHandler(store *JobStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		limit := 50
		if v := r.URL.Query().Get("limit"); v != "" {
			if n := parseInt(v); n > 0 {
				limit = n
			}
		}

		jobs, err := store.List(limit)
		if err != nil {
			writeUpgradeJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "read history: " + err.Error(),
			})
			return
		}
		if jobs == nil {
			jobs = []UpgradeJobRecord{}
		}

		writeUpgradeJSON(w, http.StatusOK, UpgradeHistoryResponse{Jobs: jobs})
	})
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// ── Background job finalizer ────────────────────────────────────────────────

// StartJobFinalizer launches a goroutine that periodically checks
// incomplete jobs and updates their final status from the node-agent.
func StartJobFinalizer(store *JobStore, agent NodeAgentProvider, interval time.Duration) {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			finalizeJobs(store, agent)
		}
	}()
}

func newTimeoutCtx(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

func finalizeJobs(store *JobStore, agent NodeAgentProvider) {
	jobs, err := store.List(0)
	if err != nil {
		return
	}
	for _, job := range jobs {
		if job.Status == "success" || job.Status == "failed" || job.Status == "rolled_back" {
			continue // already terminal
		}
		// Query node-agent for current status.
		ctx, cancel := newTimeoutCtx(10 * time.Second)
		status, err := agent.GetPlanStatus(ctx, job.OperationID)
		cancel()
		if err != nil {
			continue
		}

		newStatus := planStateToString(status.GetState())
		if newStatus == job.Status {
			continue // no change
		}

		job.Status = newStatus
		job.Error = status.GetErrorMessage()

		isTerminal := newStatus == "success" || newStatus == "failed" || newStatus == "rolled_back"
		if isTerminal && job.FinishedAt == 0 {
			job.FinishedAt = time.Now().UnixMilli()
		}

		_ = store.Save(job)
	}
}
