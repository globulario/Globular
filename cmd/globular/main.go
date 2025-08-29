package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	cfgHandlers "github.com/globulario/Globular/internal/handlers/config"
	httplib "github.com/globulario/Globular/internal/http"
	middleware "github.com/globulario/Globular/internal/http/middleware"
	"github.com/globulario/Globular/internal/server"
	config_ "github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
)

func main() {
	// Flags (minimal for now — we’ll wire more later)
	var (
		httpAddr  = flag.String("http", ":8080", "HTTP listen address (empty to disable)")
		httpsAddr = flag.String("https", "", "HTTPS listen address (empty to disable)")
	)
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Build router (middleware, metrics, health)
	mux := httplib.NewRouter(logger, httplib.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		RateRPS:        20,
		RateBurst:      80,
	})

	// ---- Mount handlers (new refactor path) ----
	wireConfig(mux)
	// wireFiles(mux)    // next
	// wireMedia(mux)    // next
	// wireAuth(mux)     // next
	// -------------------------------------------

	// Server supervisor (HTTP/HTTPS + graceful shutdown)
	sup := server.Supervisor{
		Logger:            logger,
		HTTPAddr:          *httpAddr,
		HTTPSAddr:         *httpsAddr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if err := sup.Start(mux); err != nil {
		logger.Error("start failed", "err", err)
		os.Exit(1)
	}

	// Graceful stop
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	_ = sup.Stop(context.Background())
}

// ---------------------
// Wiring: /getConfig
// ---------------------

// redirector adapts your existing redirect utilities to middleware.Redirector.
type redirector struct{}

func (redirector) RedirectTo(host string) (bool, *middleware.Target) {
	ok, p := redirectTo(host) // p is *resourcepb.Peer
	if !ok || p == nil {
		return false, nil
	}
	return true, &middleware.Target{
		Hostname:   p.Hostname,
		Domain:     p.Domain,
		Protocol:   p.Protocol,
		PortHTTP:   int(p.PortHttp),
		PortHTTPS:  int(p.PortHttps),
		LocalIP:    p.LocalIpAddress,
		ExternalIP: p.ExternalIpAddress,
		Raw:        p, // keep original available if needed
	}
}

func (redirector) HandleRedirect(to *middleware.Target, w http.ResponseWriter, r *http.Request) {
	// Build upstream host:port from the target
	addr := to.Domain
	scheme := "http"
	if to.Protocol == "https" {
		addr += ":" + Utility.ToString(to.PortHTTPS)
	} else {
		addr += ":" + Utility.ToString(to.PortHTTP)
	}
	// Trim ".localhost" like your legacy function
	addr = strings.ReplaceAll(addr, ".localhost", "")

	u, _ := url.Parse(scheme + "://" + addr)
	proxy := httputil.NewSingleHostReverseProxy(u)

	// Forwarded headers (same as legacy)
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	// If you had a custom error handler already:
	proxy.ErrorHandler = ErrHandle

	proxy.ServeHTTP(w, r)
}

// setHeaders mirrors what your setupResponse(&w, r) did.
// If your router already sets CORS headers, you can no-op here.
func setHeaders(w http.ResponseWriter, r *http.Request) {
	// Example (uncomment if needed):
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
}

// cfgProvider bridges the handler to your existing config_/globule/Utility.
type cfgProvider struct{}

func (cfgProvider) Address() (string, error) { return config_.GetAddress() }
func (cfgProvider) RemoteConfig(host string, port int) (map[string]any, error) {
	return config_.GetRemoteConfig(host, port)
}
func (cfgProvider) MyIP() string                { return Utility.MyIP() }
func (cfgProvider) LocalConfig() map[string]any { return globule.getConfig() }
func (cfgProvider) RootDir() string             { return config_.GetRootDir() }
func (cfgProvider) DataDir() string             { return config_.GetDataDir() }
func (cfgProvider) ConfigDir() string           { return config_.GetConfigDir() }
func (cfgProvider) WebRootDir() string          { return config_.GetWebRootDir() }
func (cfgProvider) PublicDirs() []string        { return config_.GetPublicDirs() }

func wireConfig(mux *http.ServeMux) {
	// Build the pure handler from the provider
	getConfig := cfgHandlers.NewGetConfig(cfgProvider{})

	// Add common cross-cutting behavior once (redirect + preflight)
	getConfig = middleware.WithRedirectAndPreflight(redirector{}, setHeaders)(getConfig)

	// Mount on the router
	cfgHandlers.Mount(mux, cfgHandlers.Deps{
		GetConfig: getConfig,
	})
}
