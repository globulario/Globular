package http

import (
	"github.com/globulario/Globular/internal/http/middleware"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter creates and configures a new HTTP router with middleware and route handlers.
// It sets up health check and metrics endpoints, applies a chain of middleware for
// recovery, security headers, CORS, logging, and rate limiting, and returns the root ServeMux.
// Parameters:
//   - logger: a structured logger for logging middleware and error recovery.
//   - cfg: configuration containing allowed CORS origins, methods, headers, and rate limiting settings.
//
// Returns:
//   - *http.ServeMux: the configured HTTP router ready to be used by an HTTP server.
func NewRouter(logger *slog.Logger, cfg Config) *http.ServeMux {
	base := http.NewServeMux()

	// health & metrics
	base.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	base.Handle("/metrics", promhttp.Handler())

	// middleware chain
	chain := middleware.Compose(
		middleware.Recoverer(logger),
		middleware.SecurityHeaders,
		middleware.CORS(cfg.AllowedOrigins, cfg.AllowedMethods, cfg.AllowedHeaders),
		middleware.Logger(logger),
		middleware.RateLimiter(middleware.NewLimiterStore(cfg.RateRPS, cfg.RateBurst)),
	)

	root := http.NewServeMux()
	root.Handle("/", chain.Then(base))
	return root
}
