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
	MetricsServices      http.Handler // GET /admin/metrics/services
	MetricsStorage       http.Handler // GET /admin/metrics/storage
	MetricsEnvoy         http.Handler // GET /admin/metrics/envoy
	ServiceLogs          http.Handler // GET /admin/service/logs
	CertificatesOverview http.Handler // GET /admin/certificates
	CertificatesCluster  http.Handler // GET /admin/certificates/cluster
	RenewPublic          http.Handler // POST /admin/certificates/renew-public
	RegenerateInternal   http.Handler // POST /admin/certificates/regenerate-internal
	UpgradesStatus       http.Handler // GET /admin/upgrades/status
	UpgradesPlan         http.Handler // POST /admin/upgrades/plan
	UpgradesApply        http.Handler // POST /admin/upgrades/apply
	UpgradesJobStatus    http.Handler // GET /admin/upgrades/jobs
	UpgradesHistory      http.Handler // GET /admin/upgrades/history
	InstalledPackages    http.Handler // GET /admin/packages
	RepoSearch           http.Handler // GET /admin/repository/search
	RepoManifest         http.Handler // GET /admin/repository/manifest
	RepoVersions         http.Handler // GET /admin/repository/versions
	RepoDelete           http.Handler // DELETE /admin/repository/artifact
	StateAlignment       http.Handler // GET /admin/state-alignment
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
	if d.CertificatesOverview != nil {
		mux.Handle("/admin/certificates", d.CertificatesOverview)
	}
	if d.CertificatesCluster != nil {
		mux.Handle("/admin/certificates/cluster", d.CertificatesCluster)
	}
	if d.RenewPublic != nil {
		mux.Handle("/admin/certificates/renew-public", d.RenewPublic)
	}
	if d.RegenerateInternal != nil {
		mux.Handle("/admin/certificates/regenerate-internal", d.RegenerateInternal)
	}
	if d.UpgradesStatus != nil {
		mux.Handle("/admin/upgrades/status", d.UpgradesStatus)
	}
	if d.UpgradesPlan != nil {
		mux.Handle("/admin/upgrades/plan", d.UpgradesPlan)
	}
	if d.UpgradesApply != nil {
		mux.Handle("/admin/upgrades/apply", d.UpgradesApply)
	}
	if d.UpgradesJobStatus != nil {
		mux.Handle("/admin/upgrades/jobs", d.UpgradesJobStatus)
	}
	if d.UpgradesHistory != nil {
		mux.Handle("/admin/upgrades/history", d.UpgradesHistory)
	}
	if d.InstalledPackages != nil {
		mux.Handle("/admin/packages", d.InstalledPackages)
	}
	if d.RepoSearch != nil {
		mux.Handle("/admin/repository/search", d.RepoSearch)
	}
	if d.RepoManifest != nil {
		mux.Handle("/admin/repository/manifest", d.RepoManifest)
	}
	if d.RepoVersions != nil {
		mux.Handle("/admin/repository/versions", d.RepoVersions)
	}
	if d.RepoDelete != nil {
		mux.Handle("/admin/repository/artifact", d.RepoDelete)
	}
	if d.StateAlignment != nil {
		mux.Handle("/admin/state-alignment", d.StateAlignment)
	}
}
