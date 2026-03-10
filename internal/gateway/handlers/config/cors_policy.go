package config

import (
	"encoding/json"
	"net/http"
	"strings"

	coreConfig "github.com/globulario/Globular/internal/config"
	goconfig "github.com/globulario/services/golang/config"
)

// Type aliases so handler code reads cleanly without qualifying every use.
type CorsPolicy = coreConfig.CorsPolicy

// Re-export shared functions for local use.
var (
	DefaultGatewayCorsPolicy = coreConfig.DefaultGatewayCorsPolicy
	DefaultServiceCorsPolicy = coreConfig.DefaultServiceCorsPolicy
	EffectivePolicy          = coreConfig.EffectivePolicy
)

// ── Interfaces ─────────────────────────────────────────────────────────────

// GatewayCorsPolicyProvider reads the gateway CORS policy.
type GatewayCorsPolicyProvider interface {
	GetCorsPolicy() *CorsPolicy
}

// GatewayCorsPolicySaver persists the gateway CORS policy.
type GatewayCorsPolicySaver interface {
	SetCorsPolicy(p *CorsPolicy) error
	Validate(token string) error
}

// ── Gateway CORS policy endpoints ──────────────────────────────────────────

// NewGetGatewayCorsPolicy returns GET /api/cors-policy.
func NewGetGatewayCorsPolicy(p GatewayCorsPolicyProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		policy := p.GetCorsPolicy()
		if policy == nil {
			policy = DefaultGatewayCorsPolicy()
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(policy)
	})
}

// NewSetGatewayCorsPolicy returns POST /api/cors-policy.
func NewSetGatewayCorsPolicy(s GatewayCorsPolicySaver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := r.Header.Get("token")
		if token == "" {
			token = r.Header.Get("Authorization")
		}
		if token == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}
		if err := s.Validate(token); err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var policy CorsPolicy
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate
		if warnings := validateCorsPolicy(&policy); len(warnings) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"saved":    true,
				"warnings": warnings,
			})
			_ = s.SetCorsPolicy(&policy)
			return
		}

		if err := s.SetCorsPolicy(&policy); err != nil {
			http.Error(w, "failed to save: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// ── Per-service CORS policy endpoints ──────────────────────────────────────

// NewGetServiceCorsPolicy returns GET /api/service-cors-policy?id=...
func NewGetServiceCorsPolicy(p GatewayCorsPolicyProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query parameter required", http.StatusBadRequest)
			return
		}

		cfg, err := goconfig.GetServiceConfigurationById(id)
		if err != nil || cfg == nil {
			http.Error(w, "service not found: "+id, http.StatusNotFound)
			return
		}

		svcPolicy := ExtractServiceCorsPolicy(cfg)
		gwPolicy := p.GetCorsPolicy()
		if gwPolicy == nil {
			gwPolicy = DefaultGatewayCorsPolicy()
		}
		effective := EffectivePolicy(gwPolicy, svcPolicy)

		resp := struct {
			Service   *CorsPolicy `json:"service"`
			Effective *CorsPolicy `json:"effective"`
		}{
			Service:   svcPolicy,
			Effective: effective,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}

// NewSetServiceCorsPolicy returns POST /api/service-cors-policy?id=...
func NewSetServiceCorsPolicy(s GatewayCorsPolicySaver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := r.Header.Get("token")
		if token == "" {
			token = r.Header.Get("Authorization")
		}
		if token == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}
		if err := s.Validate(token); err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query parameter required", http.StatusBadRequest)
			return
		}

		var policy CorsPolicy
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Ensure mode is valid for service policy
		switch policy.Mode {
		case "inherit", "override", "disabled":
		default:
			policy.Mode = "inherit"
		}

		cfg, err := goconfig.GetServiceConfigurationById(id)
		if err != nil || cfg == nil {
			http.Error(w, "service not found: "+id, http.StatusNotFound)
			return
		}

		// Persist as structured JSON under "Cors" key
		policyMap, _ := toMap(&policy)
		cfg["Cors"] = policyMap

		// Also update legacy fields for backward compatibility
		if policy.Mode == "override" {
			cfg["AllowAllOrigins"] = policy.AllowAllOrigins
			if policy.AllowAllOrigins {
				cfg["AllowedOrigins"] = "*"
			} else if len(policy.AllowedOrigins) > 0 {
				cfg["AllowedOrigins"] = strings.Join(policy.AllowedOrigins, ",")
			}
		}

		if err := goconfig.SaveServiceConfiguration(cfg); err != nil {
			http.Error(w, "failed to save: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// ── All services CORS summary ──────────────────────────────────────────────

// NewGetAllServicesCorsPolicy returns GET /api/services-cors-policy.
func NewGetAllServicesCorsPolicy(p GatewayCorsPolicyProvider, sp ServiceCorsProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cfgs, err := sp.AllServiceConfigs()
		if err != nil {
			http.Error(w, "failed to read services: "+err.Error(), http.StatusInternalServerError)
			return
		}

		gwPolicy := p.GetCorsPolicy()
		if gwPolicy == nil {
			gwPolicy = DefaultGatewayCorsPolicy()
		}

		type svcEntry struct {
			ID        string      `json:"id"`
			Name      string      `json:"name"`
			Service   *CorsPolicy `json:"service"`
			Effective *CorsPolicy `json:"effective"`
		}

		var entries []svcEntry
		for _, c := range cfgs {
			id, _ := c["Id"].(string)
			if id == "" {
				continue
			}
			name, _ := c["Name"].(string)
			svcPolicy := ExtractServiceCorsPolicy(c)
			effective := EffectivePolicy(gwPolicy, svcPolicy)
			entries = append(entries, svcEntry{
				ID:        id,
				Name:      name,
				Service:   svcPolicy,
				Effective: effective,
			})
		}
		if entries == nil {
			entries = []svcEntry{}
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(entries)
	})
}

// ── CORS Diagnostics ───────────────────────────────────────────────────────

// CorsDiagResult is the JSON response from the diagnostics endpoint.
type CorsDiagResult struct {
	Origin           string      `json:"origin"`
	ServiceID        string      `json:"service_id,omitempty"`
	Allowed          bool        `json:"allowed"`
	EffectivePolicy  *CorsPolicy `json:"effective_policy"`
	EnforcementLayer string      `json:"enforcement_layer"` // "gateway", "envoy", "service", "disabled"
	Warnings         []string    `json:"warnings,omitempty"`
	CurlExample      string      `json:"curl_example"`
}

// NewCorsDiagnostics returns GET /api/cors-diagnostics?origin=...&service=...
// It simulates a preflight check and reports which layer enforces CORS,
// the effective policy, any warnings, and a sample curl command.
func NewCorsDiagnostics(p GatewayCorsPolicyProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		origin := strings.TrimSpace(r.URL.Query().Get("origin"))
		serviceID := strings.TrimSpace(r.URL.Query().Get("service"))
		method := strings.TrimSpace(r.URL.Query().Get("method"))
		if method == "" {
			method = "POST"
		}

		gwPolicy := p.GetCorsPolicy()
		if gwPolicy == nil {
			gwPolicy = DefaultGatewayCorsPolicy()
		}

		var effective *CorsPolicy
		var layer string
		var svcPolicy *CorsPolicy

		if serviceID != "" {
			cfg, err := goconfig.GetServiceConfigurationById(serviceID)
			if err != nil || cfg == nil {
				http.Error(w, "service not found: "+serviceID, http.StatusNotFound)
				return
			}
			svcPolicy = ExtractServiceCorsPolicy(cfg)
			effective = EffectivePolicy(gwPolicy, svcPolicy)

			switch {
			case svcPolicy.Mode == "disabled":
				layer = "disabled"
			case svcPolicy.Mode == "override":
				layer = "service"
			default:
				// inherit — gateway/envoy is authoritative
				layer = "gateway+envoy"
			}
		} else {
			effective = gwPolicy
			layer = "gateway+envoy"
		}

		// Check if origin is allowed
		allowed := false
		if !effective.Enabled {
			allowed = false
		} else if effective.AllowAllOrigins {
			allowed = true
		} else if origin != "" {
			for _, o := range effective.AllowedOrigins {
				if o == origin {
					allowed = true
					break
				}
			}
		} else {
			// No origin header = same-origin request, always allowed
			allowed = true
		}

		// Build warnings
		var warnings []string
		if origin != "" && !allowed {
			warnings = append(warnings, "Origin '"+origin+"' is NOT in the allowed list. The browser will block the request.")
		}
		warnings = append(warnings, validateCorsPolicy(effective)...)

		if svcPolicy != nil && svcPolicy.Mode == "override" {
			// Check if override duplicates gateway exactly
			if policiesEqual(svcPolicy, gwPolicy) {
				warnings = append(warnings, "Service override is identical to gateway policy — consider using 'inherit' instead.")
			}
		}

		if !effective.Enabled {
			warnings = append(warnings, "CORS is disabled — no Access-Control-* headers will be sent. Browsers will block cross-origin requests.")
		}

		// Check method is allowed
		if method != "" && effective.Enabled {
			methodAllowed := false
			for _, m := range effective.AllowedMethods {
				if strings.EqualFold(m, method) {
					methodAllowed = true
					break
				}
			}
			if !methodAllowed {
				warnings = append(warnings, "Method '"+method+"' is NOT in allowed_methods.")
			}
		}

		// Build curl example
		host := r.Host
		if host == "" {
			host = "localhost"
		}
		scheme := "https"
		if r.TLS == nil {
			scheme = "http"
		}
		target := scheme + "://" + host + "/"
		if serviceID != "" {
			target = scheme + "://" + host + "/api/cors-diagnostics"
		}
		if origin == "" {
			origin = "https://example.com"
		}
		curl := "curl -s -D - -o /dev/null \\\n" +
			"  -X OPTIONS \\\n" +
			"  -H 'Origin: " + origin + "' \\\n" +
			"  -H 'Access-Control-Request-Method: " + method + "' \\\n" +
			"  -H 'Access-Control-Request-Headers: content-type,authorization' \\\n" +
			"  '" + target + "'"

		result := CorsDiagResult{
			Origin:           origin,
			ServiceID:        serviceID,
			Allowed:          allowed,
			EffectivePolicy:  effective,
			EnforcementLayer: layer,
			Warnings:         warnings,
			CurlExample:      curl,
		}
		if result.Warnings == nil {
			result.Warnings = []string{}
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
	})
}

// policiesEqual returns true if two policies have identical values.
func policiesEqual(a, b *CorsPolicy) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Enabled != b.Enabled || a.AllowAllOrigins != b.AllowAllOrigins ||
		a.AllowCredentials != b.AllowCredentials || a.AllowPrivateNetwork != b.AllowPrivateNetwork ||
		a.GrpcWebEnabled != b.GrpcWebEnabled || a.MaxAgeSeconds != b.MaxAgeSeconds {
		return false
	}
	if !slicesEqual(a.AllowedOrigins, b.AllowedOrigins) ||
		!slicesEqual(a.AllowedMethods, b.AllowedMethods) ||
		!slicesEqual(a.AllowedHeaders, b.AllowedHeaders) ||
		!slicesEqual(a.ExposedHeaders, b.ExposedHeaders) {
		return false
	}
	return true
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ── Helpers ────────────────────────────────────────────────────────────────

// ExtractServiceCorsPolicy reads the structured Cors field from a service
// config map, falling back to legacy fields if not present.
func ExtractServiceCorsPolicy(cfg map[string]any) *CorsPolicy {
	// Try structured "Cors" field first
	if raw, ok := cfg["Cors"]; ok {
		if corsMap, ok := raw.(map[string]any); ok {
			p := &CorsPolicy{}
			if b, err := json.Marshal(corsMap); err == nil {
				if err := json.Unmarshal(b, p); err == nil && p.Mode != "" {
					return p
				}
			}
		}
	}

	// Fall back to legacy fields
	allowAll, _ := cfg["AllowAllOrigins"].(bool)
	originsStr, _ := cfg["AllowedOrigins"].(string)

	p := DefaultServiceCorsPolicy()

	if allowAll || originsStr == "*" {
		p.AllowAllOrigins = true
		p.Mode = "override"
	} else if originsStr != "" {
		p.AllowAllOrigins = false
		p.AllowedOrigins = splitTrimmed(originsStr, ",")
		p.Mode = "override"
	}

	return p
}

func validateCorsPolicy(p *CorsPolicy) []string {
	var warnings []string
	if p.AllowAllOrigins && p.AllowCredentials {
		warnings = append(warnings, "allow_all_origins with allow_credentials is invalid per CORS spec; browsers will reject. Use specific origins instead.")
	}
	if p.GrpcWebEnabled {
		needHeaders := []string{"x-grpc-web", "grpc-timeout"}
		for _, h := range needHeaders {
			found := false
			for _, ah := range p.AllowedHeaders {
				if ah == h {
					found = true
					break
				}
			}
			if !found {
				warnings = append(warnings, "grpc_web_enabled but '"+h+"' missing from allowed_headers")
			}
		}
		needExposed := []string{"grpc-status", "grpc-message"}
		for _, h := range needExposed {
			found := false
			for _, eh := range p.ExposedHeaders {
				if eh == h {
					found = true
					break
				}
			}
			if !found {
				warnings = append(warnings, "grpc_web_enabled but '"+h+"' missing from exposed_headers")
			}
		}
	}
	return warnings
}

func toMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}

func splitTrimmed(s, sep string) []string {
	var result []string
	for _, part := range strings.Split(s, sep) {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
