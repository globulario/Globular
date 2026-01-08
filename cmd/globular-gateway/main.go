package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	gatewayconfig "github.com/globulario/Globular/internal/config"
	gatewayhandlers "github.com/globulario/Globular/internal/gateway/handlers"
	gatewayhttp "github.com/globulario/Globular/internal/gateway/httpserver"
	globpkg "github.com/globulario/Globular/internal/globule"
	globconfig "github.com/globulario/services/golang/config"
)

var (
	maxUpload           = flag.Int64("max-upload", 2<<30, "max upload size in bytes")
	rateRPS             = flag.Int("rate-rps", 50, "max requests/sec per client; <=0 disables throttling")
	rateBurst           = flag.Int("rate-burst", 200, "max burst per client; <=0 disables throttling")
	modeFlag            = flag.String("mode", "direct", "routing mode (direct|mesh)")
	envoyHttpAddr       = flag.String("envoy_http_addr", "127.0.0.1:8080", "HTTP address of the Envoy ingress for mesh mode")
	requireTLSBoot      = flag.Bool("require-tls-bootstrap", false, "fail startup if TLS bootstrap fails or TLS assets are unavailable when HTTPS is configured")
	httpPortOverride    = flag.String("http", "", "override HTTP port when config file is used")
	httpsPortOverride   = flag.String("https", "", "override HTTPS port when config file is used")
	gatewayConfigPath   = flag.String("config", "", "path to gateway config file (JSON)")
	printDefaultGateway = flag.Bool("print-default-config", false, "print default gateway config and exit")
)

func main() {
	flag.Parse()
	if *printDefaultGateway {
		data, err := json.MarshalIndent(gatewayconfig.DefaultGatewayConfig(), "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal default config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	finalCfg := gatewayconfig.DefaultGatewayConfig()
	gatewayCfgPath := strings.TrimSpace(*gatewayConfigPath)
	if gatewayCfgPath != "" {
		loadedCfg, err := gatewayconfig.LoadGatewayConfig(gatewayCfgPath)
		if err != nil {
			logger.Error("load gateway config", "path", gatewayCfgPath, "err", err)
			os.Exit(1)
		}
		finalCfg = loadedCfg
		logger.Info("loaded gateway config", "path", gatewayCfgPath)
	}

	finalCfg.Mode = strings.TrimSpace(*modeFlag)
	finalCfg.EnvoyHTTPAddr = strings.TrimSpace(*envoyHttpAddr)
	finalCfg.MaxUpload = *maxUpload
	finalCfg.RateRPS = *rateRPS
	finalCfg.RateBurst = *rateBurst

	requireTLSBootFlag := *requireTLSBoot

	if httpPort, err := parsePortOverride(*httpPortOverride); err != nil {
		logger.Error("parse http override", "err", err, "value", *httpPortOverride)
		os.Exit(1)
	} else if httpPort > 0 {
		finalCfg.HTTPPort = httpPort
	}

	if httpsPort, err := parsePortOverride(*httpsPortOverride); err != nil {
		logger.Error("parse https override", "err", err, "value", *httpsPortOverride)
		os.Exit(1)
	} else if httpsPort > 0 {
		finalCfg.HTTPSPort = httpsPort
	}

	if err := finalCfg.Validate(); err != nil {
		logger.Error("validate gateway config", "err", err)
		os.Exit(1)
	}

	bootstrapped, err := globconfig.EnsureLocalConfig()
	if err != nil {
		logger.Error("ensure local config failed", "err", err)
		os.Exit(1)
	}
	if bootstrapped {
		logger.Info("bootstrapped local config")
	}

	globule := globpkg.New(logger)
	applyGatewayConfigToGlobule(logger, globule, finalCfg)

	mode := strings.ToLower(strings.TrimSpace(finalCfg.Mode))
	switch mode {
	case "direct":
		globule.UseEnvoy = false
	case "mesh":
		globule.UseEnvoy = true
		logger.Info("running in mesh mode; expecting Envoy + xDS externally")
	default:
		logger.Error("invalid mode", "mode", finalCfg.Mode)
		os.Exit(1)
	}

	var httpAddr, httpsAddr string
	protocol := strings.ToLower(strings.TrimSpace(globule.Protocol))
	switch protocol {
	case "https":
		err := globule.BootstrapTLSAndDNS(context.Background())
		if err != nil {
			if requireTLSBootFlag {
				logger.Error("tls bootstrap failed", "err", err)
				os.Exit(1)
			}
			logger.Warn("tls bootstrap warning", "err", err)
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
		MaxUpload:      finalCfg.MaxUpload,
		RateRPS:        finalCfg.RateRPS,
		RateBurst:      finalCfg.RateBurst,
		ControllerAddr: globpkg.ControllerAddress(),
		EnvoyHTTPAddr:  strings.TrimSpace(finalCfg.EnvoyHTTPAddr),
		Mode:           mode,
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
	domainDir := tlsDomain(g)
	candidates := []string{
		filepath.Join(globconfig.GetRuntimeTLSDir(), domainDir, "fullchain.pem"),
		filepath.Join(globconfig.GetConfigDir(), "tls", domainDir, "fullchain.pem"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return true
		}
	}
	return false
}

func tlsDomain(g *globpkg.Globule) string {
	if g.Domain != "" {
		return g.Domain
	}
	return "localhost"
}

func parsePortOverride(raw string) (int, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, nil
	}
	port, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: %w", raw, err)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535")
	}
	return port, nil
}

func applyGatewayConfigToGlobule(logger *slog.Logger, g *globpkg.Globule, cfg gatewayconfig.GatewayConfig) {
	if cfg.Domain != "" {
		logger.Warn("ignoring gateway config Domain; cluster networking is controller-managed", "value", cfg.Domain)
	}
	if cfg.Protocol != "" {
		logger.Warn("ignoring gateway config Protocol; cluster networking is controller-managed", "value", cfg.Protocol)
	}
	if cfg.HTTPPort > 0 {
		g.PortHTTP = cfg.HTTPPort
	}
	if cfg.HTTPSPort > 0 {
		g.PortHTTPS = cfg.HTTPSPort
	}
}
