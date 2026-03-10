package config

// CorsPolicy is the structured CORS configuration shared by the gateway
// (top-level) and per-service overrides.  It lives in internal/config so
// that both the globule package and the gateway handler packages can import
// it without creating circular dependencies.
type CorsPolicy struct {
	Enabled             bool     `json:"enabled"`
	Mode                string   `json:"mode"` // "gateway" | "inherit" | "override" | "disabled"
	AllowAllOrigins     bool     `json:"allow_all_origins"`
	AllowedOrigins      []string `json:"allowed_origins"`
	AllowCredentials    bool     `json:"allow_credentials"`
	AllowedMethods      []string `json:"allowed_methods"`
	AllowedHeaders      []string `json:"allowed_headers"`
	ExposedHeaders      []string `json:"exposed_headers"`
	MaxAgeSeconds       int      `json:"max_age_seconds"`
	AllowPrivateNetwork bool     `json:"allow_private_network"`
	GrpcWebEnabled      bool     `json:"grpc_web_enabled"`
}

// DefaultGatewayCorsPolicy returns sensible defaults for the gateway-level policy.
func DefaultGatewayCorsPolicy() *CorsPolicy {
	return &CorsPolicy{
		Enabled:          true,
		Mode:             "gateway",
		AllowAllOrigins:  true,
		AllowedOrigins:   []string{},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"application",
			"token",
			"video-path",
			"index-path",
			"routing",
			"x-grpc-web",
			"grpc-timeout",
			"x-user-agent",
		},
		ExposedHeaders: []string{
			"grpc-status",
			"grpc-message",
			"grpc-status-details-bin",
		},
		MaxAgeSeconds:       3600,
		AllowPrivateNetwork: true,
		GrpcWebEnabled:      true,
	}
}

// DefaultServiceCorsPolicy returns the default per-service CORS policy (inherit from gateway).
func DefaultServiceCorsPolicy() *CorsPolicy {
	return &CorsPolicy{
		Enabled:             true,
		Mode:                "inherit",
		AllowAllOrigins:     false,
		AllowedOrigins:      []string{},
		AllowCredentials:    true,
		AllowedMethods:      []string{},
		AllowedHeaders:      []string{},
		ExposedHeaders:      []string{},
		MaxAgeSeconds:       0,
		AllowPrivateNetwork: false,
		GrpcWebEnabled:      true,
	}
}

// EffectivePolicy computes the policy actually enforced for a service,
// given the gateway-level policy and the service's own policy.
func EffectivePolicy(gateway, service *CorsPolicy) *CorsPolicy {
	if service == nil || service.Mode == "inherit" || service.Mode == "" {
		return gateway
	}
	if service.Mode == "disabled" {
		return &CorsPolicy{Enabled: false, Mode: "disabled"}
	}
	// Mode == "override"
	return service
}
