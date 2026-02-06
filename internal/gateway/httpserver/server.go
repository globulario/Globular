package httpserver

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"log/slog"

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
		// v1 Conformance: Use stable TLS directory (security violation INV-1.5)
		// REMOVED: domain-based directory selection
		// Certificate storage MUST NOT depend on domain configuration
		// Using stable path ensures certs accessible across domain changes
		credDir := filepath.Join(config_.GetConfigDir(), "tls")
		sup.TLS = &server.TLSFiles{
			CertFile: filepath.Join(credDir, "fullchain.pem"),
			KeyFile:  filepath.Join(credDir, "server.key"),
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
