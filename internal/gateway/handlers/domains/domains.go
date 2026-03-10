// Package domains exposes /api/domains/* HTTP endpoints for managing
// external domain specs and DNS provider configurations.
package domains

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/globulario/services/golang/dnsprovider"
	"github.com/globulario/services/golang/domain"
)

// ── Provider interfaces ─────────────────────────────────────────────────────

// StoreProvider gives handlers access to the domain store.
type StoreProvider interface {
	DomainStore() domain.DomainStore
}

// ── Deps / Mount ────────────────────────────────────────────────────────────

// Deps lists handlers to mount (all optional).
type Deps struct {
	ListProviders  http.Handler // GET  /api/domains/providers
	GetProvider    http.Handler // GET  /api/domains/providers?name=<ref>
	PutProvider    http.Handler // POST /api/domains/providers
	DeleteProvider http.Handler // DELETE /api/domains/providers?name=<ref>

	ListDomains  http.Handler // GET    /api/domains/specs
	GetDomain    http.Handler // GET    /api/domains/specs?fqdn=<fqdn>
	PutDomain    http.Handler // POST   /api/domains/specs
	DeleteDomain http.Handler // DELETE /api/domains/specs?fqdn=<fqdn>
}

// Mount registers only the endpoints provided.
func Mount(mux *http.ServeMux, d Deps) {
	if d.ListProviders != nil || d.GetProvider != nil || d.PutProvider != nil || d.DeleteProvider != nil {
		mux.Handle("/api/domains/providers", methodRouter(map[string]http.Handler{
			http.MethodGet:    firstNonNil(d.ListProviders, d.GetProvider),
			http.MethodPost:   d.PutProvider,
			http.MethodDelete: d.DeleteProvider,
		}))
	}
	if d.ListDomains != nil || d.GetDomain != nil || d.PutDomain != nil || d.DeleteDomain != nil {
		mux.Handle("/api/domains/specs", methodRouter(map[string]http.Handler{
			http.MethodGet:    firstNonNil(d.ListDomains, d.GetDomain),
			http.MethodPost:   d.PutDomain,
			http.MethodDelete: d.DeleteDomain,
		}))
	}
}

// ── Provider handlers ───────────────────────────────────────────────────────

// NewListProviders returns a GET handler listing all DNS provider configs (credentials masked).
func NewListProviders(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If "name" query param is present, delegate to single-get
		if name := r.URL.Query().Get("name"); name != "" {
			NewGetProvider(prov).ServeHTTP(w, r)
			return
		}

		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		// Use named listing if available (includes etcd key as name)
		type namedResult struct {
			Name        string            `json:"name"`
			Type        string            `json:"type"`
			Zone        string            `json:"zone"`
			Credentials map[string]string `json:"credentials,omitempty"`
			DefaultTTL  int               `json:"default_ttl,omitempty"`
		}

		if etcdStore, ok := store.(*domain.EtcdDomainStore); ok {
			named, err := etcdStore.ListNamedProviderConfigs(r.Context())
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("list providers: %v", err)))
				return
			}
			results := make([]namedResult, len(named))
			for i, n := range named {
				masked := domain.MaskCredentials(n.Config)
				results[i] = namedResult{
					Name:        n.Name,
					Type:        masked.Type,
					Zone:        masked.Zone,
					Credentials: masked.Credentials,
					DefaultTTL:  masked.DefaultTTL,
				}
			}
			writeJSON(w, http.StatusOK, results)
			return
		}

		// Fallback: no names available
		configs, err := store.ListProviderConfigs(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("list providers: %v", err)))
			return
		}
		masked := make([]*dnsprovider.Config, len(configs))
		for i, cfg := range configs {
			masked[i] = domain.MaskCredentials(cfg)
		}
		writeJSON(w, http.StatusOK, masked)
	}
}

// NewGetProvider returns a GET handler for a single provider config (credentials masked).
func NewGetProvider(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			writeJSON(w, http.StatusBadRequest, errorResp("missing 'name' query parameter"))
			return
		}

		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		cfg, _, err := store.GetProviderConfig(r.Context(), name)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSON(w, http.StatusNotFound, errorResp(fmt.Sprintf("provider %q not found", name)))
				return
			}
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("get provider: %v", err)))
			return
		}

		writeJSON(w, http.StatusOK, domain.MaskCredentials(cfg))
	}
}

// NewPutProvider returns a POST handler to create/update a DNS provider config.
// If a credential field has a masked value (****...), the existing value is preserved.
func NewPutProvider(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		// Decode into raw map first to handle timeout (time.Duration
		// doesn't unmarshal from JSON strings like "30s").
		var raw map[string]any
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResp(fmt.Sprintf("invalid JSON: %v", err)))
			return
		}

		var cfg dnsprovider.Config
		cfg.Type, _ = raw["type"].(string)
		cfg.Zone, _ = raw["zone"].(string)
		if ttl, ok := raw["default_ttl"].(float64); ok {
			cfg.DefaultTTL = int(ttl)
		}
		if creds, ok := raw["credentials"].(map[string]any); ok {
			cfg.Credentials = make(map[string]string, len(creds))
			for k, v := range creds {
				if s, ok := v.(string); ok {
					cfg.Credentials[k] = s
				}
			}
		}
		// Parse timeout from string or number
		switch t := raw["timeout"].(type) {
		case string:
			cfg.Timeout, _ = time.ParseDuration(t)
		case float64:
			cfg.Timeout = time.Duration(int64(t))
		}

		if cfg.Type == "" {
			writeJSON(w, http.StatusBadRequest, errorResp("provider type is required"))
			return
		}

		// Derive name from the config — use Type + Zone as a key if no explicit name
		name := cfg.Type
		if cfg.Zone != "" {
			name = cfg.Type + "-" + strings.ReplaceAll(cfg.Zone, ".", "-")
		}

		// Merge masked credentials with existing values
		existing, _, err := store.GetProviderConfig(r.Context(), name)
		if err == nil && existing != nil {
			for k, v := range cfg.Credentials {
				if domain.IsMaskedValue(v) {
					if orig, ok := existing.Credentials[k]; ok {
						cfg.Credentials[k] = orig
					}
				}
			}
		}

		if err := store.PutProviderConfig(r.Context(), name, &cfg); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("save provider: %v", err)))
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"ok":   true,
			"name": name,
		})
	}
}

// NewDeleteProvider returns a DELETE handler to remove a DNS provider config.
func NewDeleteProvider(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			writeJSON(w, http.StatusBadRequest, errorResp("missing 'name' query parameter"))
			return
		}

		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		if err := store.DeleteProviderConfig(r.Context(), name); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("delete provider: %v", err)))
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	}
}

// ── Domain spec handlers ────────────────────────────────────────────────────

// NewListDomains returns a GET handler listing all domain specs with their status.
func NewListDomains(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If "fqdn" query param is present, delegate to single-get
		if fqdn := r.URL.Query().Get("fqdn"); fqdn != "" {
			NewGetDomain(prov).ServeHTTP(w, r)
			return
		}

		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		specs, err := store.ListSpecs(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("list domains: %v", err)))
			return
		}

		// Enrich each spec with its status
		type specWithStatus struct {
			*domain.ExternalDomainSpec
			Status *domain.ExternalDomainStatus `json:"status,omitempty"`
		}

		results := make([]specWithStatus, 0, len(specs))
		for _, spec := range specs {
			item := specWithStatus{ExternalDomainSpec: spec}
			if status, _, err := store.GetStatus(r.Context(), spec.FQDN); err == nil && status != nil {
				item.Status = status
			}
			results = append(results, item)
		}

		writeJSON(w, http.StatusOK, results)
	}
}

// NewGetDomain returns a GET handler for a single domain spec + status.
func NewGetDomain(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fqdn := r.URL.Query().Get("fqdn")
		if fqdn == "" {
			writeJSON(w, http.StatusBadRequest, errorResp("missing 'fqdn' query parameter"))
			return
		}

		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		spec, _, err := store.GetSpec(r.Context(), fqdn)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSON(w, http.StatusNotFound, errorResp(fmt.Sprintf("domain %q not found", fqdn)))
				return
			}
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("get domain: %v", err)))
			return
		}

		result := map[string]any{
			"spec": spec,
		}
		if status, _, err := store.GetStatus(r.Context(), fqdn); err == nil && status != nil {
			result["status"] = status
		}

		writeJSON(w, http.StatusOK, result)
	}
}

// NewPutDomain returns a POST handler to create/update an external domain spec.
func NewPutDomain(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		var spec domain.ExternalDomainSpec
		if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResp(fmt.Sprintf("invalid JSON: %v", err)))
			return
		}

		if err := spec.Validate(); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResp(fmt.Sprintf("validation failed: %v", err)))
			return
		}

		if err := store.PutSpec(r.Context(), &spec); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("save domain: %v", err)))
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"ok":   true,
			"fqdn": spec.FQDN,
		})
	}
}

// NewDeleteDomain returns a DELETE handler to remove an external domain spec.
func NewDeleteDomain(prov StoreProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fqdn := r.URL.Query().Get("fqdn")
		if fqdn == "" {
			writeJSON(w, http.StatusBadRequest, errorResp("missing 'fqdn' query parameter"))
			return
		}

		store := prov.DomainStore()
		if store == nil {
			writeJSON(w, http.StatusServiceUnavailable, errorResp("domain store not available"))
			return
		}

		if err := store.DeleteSpec(r.Context(), fqdn); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResp(fmt.Sprintf("delete domain: %v", err)))
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	}
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func errorResp(msg string) map[string]any {
	return map[string]any{"error": msg}
}

func methodRouter(handlers map[string]http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, ok := handlers[r.Method]
		if !ok || h == nil {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func firstNonNil(handlers ...http.Handler) http.Handler {
	for _, h := range handlers {
		if h != nil {
			return h
		}
	}
	return nil
}
