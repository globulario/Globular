package config

import (
	"encoding/json"
	"net/http"
)

// ServicePermissions abstracts how we obtain a service's permissions by ID.
type ServicePermissions interface {
	LoadPermissions(serviceID string) ([]any, error)
}

// NewGetServicePermissions implements GET /getServicePermissions?id=<serviceID>
// - Redirect + CORS preflight are handled by middleware outside.
// - On success: 201 Created with the JSON array of permissions.
// - If the provider returns nil, we respond with an empty array [].
func NewGetServicePermissions(sp ServicePermissions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceID := r.URL.Query().Get("id")

		perms, err := sp.LoadPermissions(serviceID)
		if err != nil {
			http.Error(w, "fail to get service permissions with error "+err.Error(), http.StatusBadRequest)
			return
		}
		if perms == nil {
			perms = []any{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(perms)
	})
}
