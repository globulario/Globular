package config

import "net/http"

// Deps lists handlers to mount (all optional).
type Deps struct {
	GetConfig             http.Handler
	GetServiceConfig      http.Handler
	SaveConfig            http.Handler
	SaveServiceConfig     http.Handler
	GetServicePermissions http.Handler
	GetCACertificate      http.Handler
	SignCACertificate     http.Handler
	GetSANConf            http.Handler
	DescribeService       http.Handler
	GetServicesCors       http.Handler // legacy
	SetServiceCors        http.Handler // legacy

	// Structured CORS policy (PR1)
	GetGatewayCorsPolicy     http.Handler
	SetGatewayCorsPolicy     http.Handler
	GetServiceCorsPolicy     http.Handler
	SetServiceCorsPolicy     http.Handler
	GetAllServicesCorsPolicy http.Handler
	CorsDiagnostics          http.Handler
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
	if d.SaveServiceConfig != nil {
		mux.Handle("/api/save-service-config", d.SaveServiceConfig)
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

	// Structured CORS policy endpoints (PR1)
	if d.GetGatewayCorsPolicy != nil {
		mux.Handle("/api/cors-policy", d.GetGatewayCorsPolicy)
	}
	if d.SetGatewayCorsPolicy != nil {
		mux.Handle("/api/set-cors-policy", d.SetGatewayCorsPolicy)
	}
	if d.GetServiceCorsPolicy != nil {
		mux.Handle("/api/service-cors-policy", d.GetServiceCorsPolicy)
	}
	if d.SetServiceCorsPolicy != nil {
		mux.Handle("/api/set-service-cors-policy", d.SetServiceCorsPolicy)
	}
	if d.GetAllServicesCorsPolicy != nil {
		mux.Handle("/api/services-cors-policy", d.GetAllServicesCorsPolicy)
	}
	if d.CorsDiagnostics != nil {
		mux.Handle("/api/cors-diagnostics", d.CorsDiagnostics)
	}
}
