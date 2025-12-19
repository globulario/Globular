package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	gatewayhandlers "github.com/globulario/Globular/internal/gateway/handlers"
	gatewayhttp "github.com/globulario/Globular/internal/gateway/httpserver"
	globpkg "github.com/globulario/Globular/internal/globule"
	"github.com/globulario/services/golang/config"
)

var (
	maxUpload      = flag.Int64("max-upload", 2<<30, "max upload size in bytes")
	rateRPS        = flag.Int("rate-rps", 50, "max requests/sec per client; <=0 disables throttling")
	rateBurst      = flag.Int("rate-burst", 200, "max burst per client; <=0 disables throttling")
	modeFlag       = flag.String("mode", "direct", "routing mode (direct|mesh)")
	envoyHttpAddr  = flag.String("envoy_http_addr", "127.0.0.1:8080", "HTTP address of the Envoy ingress for mesh mode")
	requireDNSBoot = flag.Bool("require-dns-bootstrap", false, "fail startup if DNS bootstrap cannot contact DNS service")
)

func main() {
	_ = flag.String("http", "", "ignored: HTTP port is taken from Globule config")
	_ = flag.String("https", "", "ignored: HTTPS port is taken from Globule config")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	globule := globpkg.New(logger)
	globule.SkipLocalDNS = true

	switch mode := strings.ToLower(strings.TrimSpace(*modeFlag)); mode {
	case "direct":
		globule.UseEnvoy = false
	case "mesh":
		globule.UseEnvoy = true
		logger.Info("running in mesh mode; expecting Envoy + xDS externally")
	default:
		logger.Error("invalid mode", "mode", *modeFlag)
		os.Exit(1)
	}

	var httpAddr, httpsAddr string
	protocol := strings.ToLower(strings.TrimSpace(globule.Protocol))
	switch protocol {
	case "https":
		err := globule.BootstrapTLSAndDNS(context.Background())
		if err != nil {
			if *requireDNSBoot {
				logger.Error("tls/dns bootstrap failed", "err", err)
				os.Exit(1)
			}
			logger.Warn("tls/dns bootstrap warning", "err", err)
			if hasExistingTLSCert(globule) {
				logger.Info("existing TLS certificate found; continuing with HTTPS")
				httpsAddr = fmt.Sprintf(":%d", globule.PortHTTPS)
				logger.Info("starting HTTPS (from globule config)", "addr", httpsAddr, "domain", globule.Domain)
			} else {
				httpAddr = fmt.Sprintf(":%d", globule.PortHTTP)
				logger.Warn("falling back to HTTP; TLS certificate not available", "addr", httpAddr)
				logger.Info("starting HTTP (from globule config)", "addr", httpAddr, "domain", globule.Domain)
			}
		} else {
			httpsAddr = fmt.Sprintf(":%d", globule.PortHTTPS)
			logger.Info("starting HTTPS (from globule config)", "addr", httpsAddr, "domain", globule.Domain)
		}
	default:
		if err := globule.InitFS(); err != nil {
			logger.Error("bootstrap failed", "err", err)
			os.Exit(1)
		}
		httpAddr = fmt.Sprintf(":%d", globule.PortHTTP)
		logger.Info("starting HTTP (from globule config)", "addr", httpAddr, "domain", globule.Domain)
	}

	limiterDisabled := *rateRPS <= 0 || *rateBurst <= 0
	if limiterDisabled {
		logger.Warn("http rate limiter disabled", "rateRPS", *rateRPS, "rateBurst", *rateBurst)
	}

	handlerSet := gatewayhandlers.New(globule, gatewayhandlers.HandlerConfig{
		MaxUpload:      *maxUpload,
		RateRPS:        *rateRPS,
		RateBurst:      *rateBurst,
		NodeAgentAddr:  globpkg.NodeAgentAddress(),
		ControllerAddr: globpkg.ControllerAddress(),
		EnvoyHTTPAddr:  strings.TrimSpace(*envoyHttpAddr),
		Mode:           strings.ToLower(strings.TrimSpace(*modeFlag)),
	})
	mux := handlerSet.Router(logger)

	httpServer := gatewayhttp.New(logger, globule, httpAddr, httpsAddr)
	if err := httpServer.Start(mux); err != nil {
		logger.Error("start failed", "err", err)
		os.Exit(1)
	}

	<-httpServer.Ready()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shCtx, shCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shCancel()
	if err := httpServer.Stop(shCtx); err != nil {
		logger.Error("HTTP/S shutdown error", "err", err)
	}
}

func hasExistingTLSCert(g *globpkg.Globule) bool {
	fullchain := filepath.Join(config.GetConfigDir(), "tls", g.LocalDomain(), "fullchain.pem")
	if _, err := os.Stat(fullchain); err == nil {
		return true
	}
	return false
}
