package cluster

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

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

// NewJoinWorkflowsHandler serves workflow definition YAML files from etcd
// (/globular/workflows/) or from the local filesystem as fallback.
// All workflow definitions are stored in etcd — MinIO is not used.
func NewJoinWorkflowsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, "/join/workflows/")
		name = strings.Trim(name, "/")

		// GET /join/workflows/ → JSON list of available workflow names.
		if name == "" {
			names := make([]string, 0, len(allowedWorkflows))
			for n := range allowedWorkflows {
				names = append(names, n)
			}
			sort.Strings(names)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(names)
			return
		}

		name = filepath.Base(name) // sanitize path traversal
		if !allowedWorkflows[name] {
			http.NotFound(w, r)
			return
		}

		// Primary: etcd /globular/workflows/<name>
		if cli, err := config.GetEtcdClient(); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			resp, etcdErr := cli.Get(ctx, "/globular/workflows/"+name)
			cancel()
			if etcdErr == nil && len(resp.Kvs) > 0 {
				w.Header().Set("Content-Type", "application/x-yaml")
				w.Header().Set("Content-Disposition", "attachment; filename="+name)
				w.Write(resp.Kvs[0].Value)
				return
			}
		}

		// Fallback: local filesystem (bootstrap window before etcd is seeded).
		for _, path := range []string{
			"/var/lib/globular/workflows/" + name,
			"/usr/lib/globular/workflows/" + name,
		} {
			http.ServeFile(w, r, path)
			return
		}

		http.Error(w, "workflow definition not found", http.StatusNotFound)
	})
}
