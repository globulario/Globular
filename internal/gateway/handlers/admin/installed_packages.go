package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/globulario/services/golang/installed_state"
	node_agentpb "github.com/globulario/services/golang/node_agent/node_agentpb"
	"google.golang.org/protobuf/encoding/protojson"
)

// InstalledPackagesProvider is the interface for fetching installed packages.
// Default implementation reads from etcd via the installed_state package.
type InstalledPackagesProvider interface {
	ListInstalledPackages(ctx context.Context, nodeID, kind, name string) ([]*node_agentpb.InstalledPackage, error)
}

// etcdInstalledPackagesProvider reads directly from etcd.
type etcdInstalledPackagesProvider struct{}

func (etcdInstalledPackagesProvider) ListInstalledPackages(ctx context.Context, nodeID, kind, name string) ([]*node_agentpb.InstalledPackage, error) {
	if nodeID != "" {
		pkgs, err := installed_state.ListInstalledPackages(ctx, nodeID, kind)
		if err != nil {
			return nil, err
		}
		// Apply name filter for per-node queries (ListInstalledPackages doesn't filter by name).
		if name != "" {
			nameLower := strings.ToLower(name)
			filtered := pkgs[:0]
			for _, p := range pkgs {
				if strings.ToLower(p.GetName()) == nameLower {
					filtered = append(filtered, p)
				}
			}
			return filtered, nil
		}
		return pkgs, nil
	}
	return installed_state.ListAllNodes(ctx, kind, name)
}

// NewInstalledPackagesHandler returns an HTTP handler that serves
// GET /admin/packages?node_id=&kind=&name=
func NewInstalledPackagesHandler() http.Handler {
	return newInstalledPackagesHandler(etcdInstalledPackagesProvider{})
}

func newInstalledPackagesHandler(prov InstalledPackagesProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		nodeID := strings.TrimSpace(r.URL.Query().Get("node_id"))
		kind := strings.TrimSpace(r.URL.Query().Get("kind"))
		name := strings.TrimSpace(r.URL.Query().Get("name"))

		pkgs, err := prov.ListInstalledPackages(r.Context(), nodeID, kind, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// When no kind filter is specified, exclude COMMAND packages —
		// they are CLI tools (claude, etcdctl, mc, etc.), not service
		// instances, and showing them in the "Service Instances" UI
		// is confusing. Callers that explicitly want commands pass
		// ?kind=COMMAND.
		if kind == "" {
			filtered := pkgs[:0]
			for _, pkg := range pkgs {
				if strings.EqualFold(pkg.GetKind(), "COMMAND") {
					continue
				}
				filtered = append(filtered, pkg)
			}
			pkgs = filtered
		}

		// Marshal each package with protojson for consistent field naming,
		// then wrap in a JSON array.
		items := make([]json.RawMessage, 0, len(pkgs))
		for _, pkg := range pkgs {
			data, err := protojson.Marshal(pkg)
			if err != nil {
				continue
			}
			items = append(items, json.RawMessage(data))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"packages": items,
			"count":    len(items),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
