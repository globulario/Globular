package config

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewGetServiceConfig returns a handler that serves /config/<service>.
func NewGetServiceConfig(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceName := strings.TrimPrefix(r.URL.Path, "/config/")
		if idx := strings.Index(serviceName, "/"); idx >= 0 {
			serviceName = serviceName[:idx]
		}
		serviceName = strings.TrimSpace(serviceName)
		if serviceName == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}

		svc, err := p.ServiceConfig(serviceName)
		if err != nil || svc == nil {
			http.Error(w, "no service found with name or id "+serviceName, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(svc)
	})
}
