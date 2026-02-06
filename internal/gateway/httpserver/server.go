package httpserver

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/globulario/Globular/internal/config"
	globpkg "github.com/globulario/Globular/internal/globule"
	"github.com/globulario/Globular/internal/server"
	config_ "github.com/globulario/services/golang/config"
)

// Server wraps the existing HTTP supervisor with gateway-specific TLS wiring.
type Server struct {
	sup *server.Supervisor
}

// New constructs the HTTP/HTTPS supervisor configured from Globule.
func New(logger *slog.Logger, globule *globpkg.Globule, httpAddr, httpsAddr string) *Server {
	sup := &server.Supervisor{
		Logger:            logger,
		HTTPAddr:          httpAddr,
		HTTPSAddr:         httpsAddr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if httpsAddr != "" {
		// v1 Conformance: Use stable TLS directory (see config.CertPaths)
		certPaths := config.NewCertPaths(config_.GetConfigDir())
		sup.TLS = &server.TLSFiles{
			CertFile: certPaths.InternalServerCert(),
			KeyFile:  certPaths.InternalServerKey(),
		}
	}

	return &Server{sup: sup}
}

// Start launches the HTTP/HTTPS servers with the supplied handler.
func (s *Server) Start(handler http.Handler) error {
	return s.sup.Start(handler)
}

// Stop gracefully shuts down the HTTP/HTTPS listeners.
func (s *Server) Stop(ctx context.Context) error {
	return s.sup.Stop(ctx)
}

// Ready exposes the supervisor readiness channel.
func (s *Server) Ready() <-chan struct{} {
	return s.sup.Ready
}
