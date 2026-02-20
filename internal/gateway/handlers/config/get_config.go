package config

import (
	"encoding/json"
	"net/http"
)

// Provider is the minimal surface the handler needs.
// You can adapt it to your current config_ / globule code in main.go.
type Provider interface {
	Address() (string, error)
	MyIP() string

	// local snapshot
	LocalConfig() map[string]any
	ServiceConfig(idOrName string) (map[string]any, error)
	RootDir() string
	DataDir() string
	ConfigDir() string
	WebRootDir() string
	PublicDirs() []string
}

// NewGetConfig returns an http.Handler with the core logic ONLY.
// Redirect + preflight are handled by middleware (see #2).
func NewGetConfig(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Local config
		conf := p.LocalConfig()

		// Enrich with helpful path info (kept from your legacy handler)
		conf["Root"] = p.RootDir()
		conf["DataPath"] = p.DataDir()
		conf["ConfigPath"] = p.ConfigDir()
		conf["WebRoot"] = p.WebRootDir()
		conf["Public"] = p.PublicDirs()
		conf["OAuth2ClientSecret"] = "********" // mask

		// If ?id=serviceId is present, return that service’s config.
		if serviceID := r.URL.Query().Get("id"); serviceID != "" {
			svc, err := p.ServiceConfig(serviceID)
			if err != nil || svc == nil {
				http.Error(w, "no service found with name or id "+serviceID, http.StatusBadRequest)
				return
			}
			conf = svc
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(conf)
	})
}

func toInt(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return n
		}
		n = n*10 + int(c-'0')
	}
	return n
}
