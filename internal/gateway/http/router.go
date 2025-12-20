package http

import (
	"log/slog"
	"net/http"

	middleware "github.com/globulario/Globular/internal/gateway/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config carries public HTTP knobs the router needs.
type Config struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	RateRPS        int // <=0 disables throttling
	RateBurst      int // <=0 disables throttling
}

// NewRouter creates the outer mux with the middleware chain, and lets you
// inject a rootHandler that will serve "/" (e.g., your static file server).
// More-specific handlers like /healthz and /metrics are registered on the
// inner 'base' mux so they still win over "/".
func NewRouter(logger *slog.Logger, cfg Config, rootHandler http.Handler) *http.ServeMux {
	base := http.NewServeMux()

	// health & metrics (specific paths)
	base.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	base.Handle("/metrics", promhttp.Handler())

	// If provided, this serves "/" (and subpaths) via your static handler.
	// Handlers above (like /healthz) remain higher-priority.
	if rootHandler != nil {
		base.Handle("/", rootHandler)
	}

	// middleware chain (order matters)
	middlewares := []func(http.Handler) http.Handler{
		middleware.Recoverer(logger),
		middleware.SecurityHeaders,
		middleware.CORS(cfg.AllowedOrigins, cfg.AllowedMethods, cfg.AllowedHeaders),
		middleware.Logger(logger),
	}
	if cfg.RateRPS > 0 && cfg.RateBurst > 0 {
		middlewares = append(middlewares, middleware.RateLimiter(
			middleware.NewLimiterStore(float64(cfg.RateRPS), cfg.RateBurst),
		))
	}
	chain := middleware.Compose(middlewares...)

	// Outer mux owns "/" exactly once.
	root := http.NewServeMux()
	root.Handle("/", chain.Then(base))
	return root
}
