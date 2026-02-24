package config

import (
	"encoding/json"
	"net/http"

	"github.com/globulario/services/golang/config"
)

// ServiceCorsProvider can enumerate all service configs.
type ServiceCorsProvider interface {
	// AllServiceConfigs returns the raw config map for every registered service.
	AllServiceConfigs() ([]map[string]any, error)
}

// ServiceCorsSaver can persist an updated service config.
type ServiceCorsSaver interface {
	Validate(token string) error
	SaveServiceConfig(cfg map[string]any) error
}

// NewGetServicesCors returns GET /api/services-cors.
// Response: JSON array of { id, name, allowAllOrigins, allowedOrigins }.
// No auth required — CORS settings are not sensitive.
func NewGetServicesCors(p ServiceCorsProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfgs, err := p.AllServiceConfigs()
		if err != nil {
			http.Error(w, "failed to read service configurations: "+err.Error(), http.StatusInternalServerError)
			return
		}

		type policy struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			AllowAllOrigins bool   `json:"allowAllOrigins"`
			AllowedOrigins  string `json:"allowedOrigins"`
		}

		var policies []policy
		for _, c := range cfgs {
			id, _ := c["Id"].(string)
			if id == "" {
				continue
			}
			name, _ := c["Name"].(string)
			allowAll, _ := c["AllowAllOrigins"].(bool)
			origins, _ := c["AllowedOrigins"].(string)
			policies = append(policies, policy{
				ID:              id,
				Name:            name,
				AllowAllOrigins: allowAll,
				AllowedOrigins:  origins,
			})
		}

		if policies == nil {
			policies = []policy{}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(policies)
	})
}

// NewSetServiceCors returns POST /api/service-cors.
// Body: { "id": "...", "allowAllOrigins": bool, "allowedOrigins": "..." }
// Requires token header (JWT-only validation, no Scylla needed).
func NewSetServiceCors(s ServiceCorsSaver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := r.Header.Get("token")
		if token == "" {
			http.Error(w, "token header required", http.StatusUnauthorized)
			return
		}
		if err := s.Validate(token); err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var body struct {
			ID              string `json:"id"`
			AllowAllOrigins bool   `json:"allowAllOrigins"`
			AllowedOrigins  string `json:"allowedOrigins"`
		}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if body.ID == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		cfg, err := config.GetServiceConfigurationById(body.ID)
		if err != nil || cfg == nil {
			http.Error(w, "service not found: "+body.ID, http.StatusNotFound)
			return
		}

		cfg["AllowAllOrigins"] = body.AllowAllOrigins
		cfg["AllowedOrigins"] = body.AllowedOrigins

		if err := s.SaveServiceConfig(cfg); err != nil {
			http.Error(w, "failed to save: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
