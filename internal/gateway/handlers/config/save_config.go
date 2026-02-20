package config

import (
	"encoding/json"
	"net/http"
)

// TokenValidator checks an access token.
type TokenValidator interface {
	Validate(token string) error
}

// Saver persists the supplied configuration map.
type Saver interface {
	Save(m map[string]any) error
}

// NewSaveConfig returns POST /saveConfig.
// - OPTIONS & redirect are handled by middleware (wrap this in main.go).
// - Requires a token in header "token" or query "token" (same as legacy).
// - 204 on success.
func NewSaveConfig(s Saver, auth TokenValidator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// v2 Conformance: Token from header ONLY (security violation INV-3.1)
		// REMOVED: Query parameter token fallback
		// Tokens in URLs leak via access logs, proxy logs, browser history
		// Header-only prevents token exposure
		token := r.Header.Get("token")
		if token == "" {
			http.Error(w, "no token in header - query parameter tokens not accepted (use Authorization header)", http.StatusUnauthorized)
			return
		}

		// Validate token
		if err := auth.Validate(token); err != nil {
			http.Error(w, "fail to validate token with error "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Decode payload
		defer r.Body.Close()
		var cfg map[string]any
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			http.Error(w, "fail to decode configuration with error "+err.Error(), http.StatusBadRequest)
			return
		}

		// Persist
		if err := s.Save(cfg); err != nil {
			http.Error(w, "fail to set configuration with error "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
