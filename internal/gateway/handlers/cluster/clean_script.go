package cluster

import (
	_ "embed"
	"net/http"
)

//go:embed clean-node.sh
var cleanNodeScript []byte

// NewCleanScriptHandler serves the clean-node.sh script that prepares a node
// for a fresh Day-1 join by stopping all Globular/ScyllaDB services and wiping
// their state.
//
// Usage on the target node:
//
//	curl -sfL https://<gateway>:8443/clean -k | sudo bash
//	curl -sfL https://<gateway>:8443/clean -k | sudo bash -s -- --force
//
// The --force flag skips the interactive confirmation prompt. AI agents that
// call this endpoint must pass --force to avoid blocking on stdin.
func NewCleanScriptHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/x-shellscript")
		w.Header().Set("Content-Disposition", `attachment; filename="clean-node.sh"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(cleanNodeScript)
	})
}
