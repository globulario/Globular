package cluster

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/globulario/services/golang/config"
)

// allowedWorkflows lists workflow definitions that can be served to joining nodes.
var allowedWorkflows = map[string]bool{
	"node.join.yaml":                    true,
	"node.bootstrap.yaml":               true,
	"node.repair.yaml":                  true,
	"day0.bootstrap.yaml":               true,
	"cluster.reconcile.yaml":            true,
	"release.apply.package.yaml":        true,
	"release.apply.infrastructure.yaml": true,
	"release.remove.package.yaml":       true,
}

// NewJoinWorkflowsHandler serves workflow definition YAML files from MinIO
// (globular-config/workflows/) or from the local filesystem as fallback.
// This allows joining nodes to fetch workflow definitions via the gateway
// without needing direct MinIO access.
func NewJoinWorkflowsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, "/join/workflows/")
		name = filepath.Base(name) // sanitize path traversal

		if !allowedWorkflows[name] {
			http.NotFound(w, r)
			return
		}

		// Try MinIO first (globular-config/workflows/<name>).
		data, err := config.GetClusterConfig("workflows/" + name)
		if err == nil && len(data) > 0 {
			w.Header().Set("Content-Type", "application/x-yaml")
			w.Header().Set("Content-Disposition", "attachment; filename="+name)
			w.Write(data)
			return
		}
		if err != nil {
			log.Printf("join/workflows/%s: MinIO fetch failed: %v — trying local", name, err)
		}

		// Fallback to local filesystem.
		candidates := []string{
			"/var/lib/globular/workflows/" + name,
			"/usr/lib/globular/workflows/" + name,
		}
		for _, path := range candidates {
			http.ServeFile(w, r, path)
			return
		}

		http.Error(w, "workflow definition not found", http.StatusNotFound)
	})
}
