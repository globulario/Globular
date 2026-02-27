package config

import (
	"encoding/json"
	"net/http"
)

// ServicePatcher fetches and saves individual service configurations.
type ServicePatcher interface {
	GetServiceConfig(idOrName string) (map[string]any, error)
	SaveServiceConfig(cfg map[string]any) error
}

// NewSaveServiceConfig returns POST /api/save-service-config.
//
// The request body must be a JSON object containing at least "Id".
// All other fields are merged on top of the current service config and
// persisted via SaveServiceConfiguration (which correctly routes desired
// fields to etcd /config and runtime fields — State, Process, ProxyProcess —
// to etcd /runtime).
//
// This is a patch operation: fields absent from the body are not changed.
func NewSaveServiceConfig(p ServicePatcher, auth TokenValidator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := r.Header.Get("token")
		if token == "" {
			http.Error(w, "no token in header", http.StatusUnauthorized)
			return
		}
		if err := auth.Validate(token); err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		defer r.Body.Close()
		var patch map[string]any
		if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		id, _ := patch["Id"].(string)
		if id == "" {
			http.Error(w, "missing Id field", http.StatusBadRequest)
			return
		}

		// Fetch the full current config so we only overwrite the supplied fields.
		current, err := p.GetServiceConfig(id)
		if err != nil || current == nil {
			http.Error(w, "service not found: "+id, http.StatusNotFound)
			return
		}

		// Merge patch fields into current (Id is always preserved from current).
		for k, v := range patch {
			current[k] = v
		}

		if err := p.SaveServiceConfig(current); err != nil {
			http.Error(w, "save failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
