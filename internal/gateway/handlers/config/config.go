package config

import "net/http"

// Deps lists handlers to mount (all optional).
type Deps struct {
	GetConfig             http.Handler
	GetServiceConfig      http.Handler
	SaveConfig            http.Handler
	GetServicePermissions http.Handler // NEW
	GetCACertificate      http.Handler // NEW
	SignCACertificate     http.Handler // NEW
	GetSANConf            http.Handler // NEW
	DescribeService       http.Handler
}

// Mount registers only the endpoints provided.
func Mount(mux *http.ServeMux, d Deps) {
	if d.GetConfig != nil {
		mux.Handle("/api/get-config", d.GetConfig)
		mux.Handle("/config", d.GetConfig) // legacy
	}
	if d.GetServiceConfig != nil {
		mux.Handle("/config/", d.GetServiceConfig)
	}
	if d.SaveConfig != nil {
		mux.Handle("/api/save-config", d.SaveConfig)
	}
	if d.GetServicePermissions != nil {
		mux.Handle("/api/get-service-permissions", d.GetServicePermissions) // NEW
	}
	if d.DescribeService != nil {
		mux.Handle("/api/describe-service", d.DescribeService)
	}

	if d.GetCACertificate != nil {
		mux.Handle("/get_ca_certificate", d.GetCACertificate)
	}
	if d.SignCACertificate != nil {
		mux.Handle("/sign_ca_certificate", d.SignCACertificate)
	}
	if d.GetSANConf != nil {
		mux.Handle("/get_san_conf", d.GetSANConf)
	}
}
