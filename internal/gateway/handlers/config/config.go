package config

import "net/http"

// Deps lists handlers to mount (all optional).
type Deps struct {
	GetConfig             http.Handler
	GetServiceConfig      http.Handler
	SaveConfig            http.Handler
	GetServicePermissions http.Handler
	GetCACertificate      http.Handler
	SignCACertificate     http.Handler
	GetSANConf            http.Handler
	DescribeService       http.Handler
	GetServicesCors       http.Handler
	SetServiceCors        http.Handler
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
		mux.Handle("/api/get-service-permissions", d.GetServicePermissions)
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
	if d.GetServicesCors != nil {
		mux.Handle("/api/services-cors", d.GetServicesCors)
	}
	if d.SetServiceCors != nil {
		mux.Handle("/api/service-cors", d.SetServiceCors)
	}
}
