package cluster

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// NewJoinBinHandler serves node-agent and globularcli binaries for joining nodes.
// Only serves files explicitly in the allowlist — no directory traversal.
//
// Routes:
//
//	/join/bin/node_agent_server
//	/join/bin/globularcli
func NewJoinBinHandler(binDir string) http.Handler {
	allowed := map[string]bool{
		"node_agent_server": true,
		"globularcli":       true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract filename from /join/bin/<name>
		name := strings.TrimPrefix(r.URL.Path, "/join/bin/")
		name = filepath.Base(name) // prevent traversal

		if !allowed[name] {
			http.NotFound(w, r)
			return
		}

		path := filepath.Join(binDir, name)
		if _, err := os.Stat(path); err != nil {
			http.Error(w, "binary not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		http.ServeFile(w, r, path)
	})
}
