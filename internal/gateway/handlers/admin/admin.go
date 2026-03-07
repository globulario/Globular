// Package admin exposes /admin/metrics/* HTTP endpoints that return
// fully-derived service health and storage models so the frontend
// becomes a pure renderer with no client-side health computation.
package admin

import "net/http"

// AdminProvider is the minimal surface the admin handlers need from the
// gateway / globule layer.
type AdminProvider interface {
	AllServiceConfigs() ([]map[string]any, error)
	PublicDirs() []string
	DataDir() string
	StateDir() string
	Hostname() string
	IP() string
}

// Deps lists handlers to mount (all optional).
type Deps struct {
	MetricsServices http.Handler // GET /admin/metrics/services
	MetricsStorage  http.Handler // GET /admin/metrics/storage
	MetricsEnvoy    http.Handler // GET /admin/metrics/envoy
	ServiceLogs     http.Handler // GET /admin/service/logs
}

// Mount registers only the endpoints provided.
func Mount(mux *http.ServeMux, d Deps) {
	if d.MetricsServices != nil {
		mux.Handle("/admin/metrics/services", d.MetricsServices)
	}
	if d.MetricsStorage != nil {
		mux.Handle("/admin/metrics/storage", d.MetricsStorage)
	}
	if d.MetricsEnvoy != nil {
		mux.Handle("/admin/metrics/envoy", d.MetricsEnvoy)
	}
	if d.ServiceLogs != nil {
		mux.Handle("/admin/service/logs", d.ServiceLogs)
	}
}
