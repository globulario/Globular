package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/globulario/Globular/internal/controlplane"
	"github.com/globulario/Globular/internal/xds/server"
	"github.com/globulario/Globular/internal/xds/watchers"
	globconfig "github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
	"gopkg.in/yaml.v3"
)

const (
	defaultGRPCAddr           = "127.0.0.1:18000"
	defaultXDSConfig          = "/var/lib/globular/xds/config.json"
	defaultEnvoyBootstrapPath = "/run/globular/envoy/envoy-bootstrap.json"
	defaultServiceConfigPath  = "/var/lib/globular/xds/xds.yaml"
	xdsServiceID              = "xds.XdsService"
)

func main() {
	var (
		grpcAddr           = flag.String("grpc_addr", defaultGRPCAddr, "address for the xDS gRPC server")
		xdsConfigPath      = flag.String("xds_config", defaultXDSConfig, "path to the xDS config JSON")
		serviceConfigPath  = flag.String("config", defaultServiceConfigPath, "path to a globular-xds YAML service config")
		nodeID             = flag.String("node_id", "globular-xds", "node id for xDS snapshots")
		bootstrapPath      = flag.String("envoy_bootstrap", defaultEnvoyBootstrapPath, "path to write Envoy bootstrap JSON")
		bootstrapHost      = flag.String("bootstrap_host", "127.0.0.1", "host advertised in envoy bootstrap")
		bootstrapPort      = flag.Int("bootstrap_port", 18000, "port advertised in envoy bootstrap")
		bootstrapCluster   = flag.String("bootstrap_cluster", "globular-cluster", "cluster name used in bootstrap")
		bootstrapAdminPort = flag.Int("bootstrap_admin", 9901, "admin port in bootstrap")
		bootstrapMaxConn   = flag.Uint64("bootstrap_max_downstream_connections", 0, "max active downstream connections (0 disables)")
		downstreamTLSMode  = flag.String("downstream_tls_mode", "optional", "downstream TLS mode (disabled|optional|required)")
		watchInterval      = flag.Duration("watch_interval", 5*time.Second, "static config polling interval")
		debugSignals       = flag.Bool("debug_signals", false, "log signals received while running")
		describeFlag       = flag.Bool("describe", false, "print xds metadata as JSON and exit")
		controllerAddr     = flag.String("controller_addr", "", "cluster controller address for DNS-based routing (optional, enables cluster mode)")
	)
	flag.Parse()

	// Allow controller address to be set via environment variable (PR4)
	if *controllerAddr == "" {
		if envAddr := os.Getenv("CLUSTER_CONTROLLER_ADDR"); envAddr != "" {
			*controllerAddr = envAddr
		}
	}

	xdsCfg, cfgErr := loadXDSServiceConfig(*serviceConfigPath)
	var missingConfigErr error
	if cfgErr != nil {
		if errors.Is(cfgErr, os.ErrNotExist) {
			missingConfigErr = cfgErr
			xdsCfg = nil
		} else {
			fmt.Fprintf(os.Stderr, "load service config %q: %v\n", *serviceConfigPath, cfgErr)
			os.Exit(1)
		}
	}
	portCfg, portCfgErr := loadServicePortConfig(xdsServiceID)

	logLevel := slog.LevelInfo
	if xdsCfg != nil {
		if lvlStr := strings.TrimSpace(xdsCfg.Logging.Level); lvlStr != "" {
			if parsed, err := parseLogLevel(lvlStr); err == nil {
				logLevel = parsed
			} else {
				fmt.Fprintf(os.Stderr, "invalid logging level %q: %v; falling back to info\n", lvlStr, err)
			}
		}
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	if xdsCfg != nil {
		logger.Info("loaded service config", "path", *serviceConfigPath)
	} else if missingConfigErr != nil {
		logger.Info("service config missing; using defaults", "path", *serviceConfigPath, "err", missingConfigErr)
	} else {
		logger.Info("running without service config; using defaults", "path", *serviceConfigPath)
	}
	if portCfgErr != nil {
		logger.Warn("read service port config", "service", xdsServiceID, "err", portCfgErr)
	}

	if xdsCfg != nil {
		if sd := strings.TrimSpace(xdsCfg.Runtime.StateDir); sd != "" {
			sd = filepath.Clean(sd)
			if err := Utility.CreateDirIfNotExist(sd); err != nil {
				logger.Error("ensure runtime state dir", "err", err, "path", sd)
				os.Exit(1)
			}
			if err := os.Setenv("GLOBULAR_STATE_DIR", sd); err != nil {
				logger.Error("export runtime state dir", "err", err, "path", sd)
				os.Exit(1)
			}
			logger.Info("runtime state dir applied", "path", sd)
		}
	}

	opts := controlplane.BootstrapOptions{
		NodeID:    *nodeID,
		Cluster:   *bootstrapCluster,
		XDSHost:   *bootstrapHost,
		XDSPort:   *bootstrapPort,
		AdminPort: *bootstrapAdminPort,
	}
	if *bootstrapMaxConn > 0 {
		opts.MaxActiveDownstreamConns = *bootstrapMaxConn
	}
	writeBootstrap(logger, opts, *bootstrapPath)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	xdsServer := server.New(logger, ctx)
	errCh := make(chan error, 2)
	grpcListenAddr := strings.TrimSpace(*grpcAddr)
	if xdsCfg != nil {
		if cfgAddr := strings.TrimSpace(xdsCfg.GRPC.Address); cfgAddr != "" {
			grpcListenAddr = cfgAddr
		}
	}
	if grpcListenAddr == "" {
		grpcListenAddr = defaultGRPCAddr
	}
	grpcListenAddr = applyServicePort(grpcListenAddr, portCfg)
	if grpcListenAddr == "" {
		grpcListenAddr = defaultGRPCAddr
	}
	if p := portFromAddress(grpcListenAddr); p > 0 {
		opts.XDSPort = p
	}

	if *describeFlag {
		if err := emitXDSDescribe(grpcListenAddr, portCfg); err != nil {
			logger.Error("describe failed", "err", err)
			os.Exit(1)
		}
		return
	}

	go func() {
		if err := xdsServer.Serve(ctx, grpcListenAddr); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		}
	}()

	watcher := watchers.New(logger, xdsServer, *xdsConfigPath, *nodeID, *watchInterval, watchers.ParseDownstreamTLSMode(*downstreamTLSMode), *controllerAddr)
	go func() {
		if err := watcher.Run(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			errCh <- err
			return
		}
		errCh <- fmt.Errorf("watcher exited unexpectedly")
	}()

	if *debugSignals {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
		go func() {
			for sig := range sigCh {
				logger.Info("signal received", "signal", sig)
			}
		}()
		defer signal.Stop(sigCh)
	}

	go func() {
		<-ctx.Done()
		logger.Info("root context canceled", "err", ctx.Err())
	}()

	select {
	case err := <-errCh:
		logger.Error("fatal xDS component exited", "err", err)
		stop()
		os.Exit(1)
	case <-ctx.Done():
	}
}

func writeBootstrap(logger *slog.Logger, opts controlplane.BootstrapOptions, path string) {
	if path == "" {
		path = defaultEnvoyBootstrapPath
	}
	if err := controlplane.WriteBootstrap(path, opts); err != nil {
		logger.Warn("write bootstrap", "err", err, "path", path)
		return
	}
	logger.Info("wrote envoy bootstrap", "path", path, "xds", fmt.Sprintf("%s:%d", opts.XDSHost, opts.XDSPort))
}

type xdsServiceConfig struct {
	GRPC struct {
		Address string `yaml:"address"`
	} `yaml:"grpc"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	Runtime struct {
		StateDir string `yaml:"state_dir"`
	} `yaml:"runtime"`
}

func loadXDSServiceConfig(path string) (*xdsServiceConfig, error) {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return nil, fmt.Errorf("service config path is empty")
	}
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}
	var cfg xdsServiceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", cleanPath, err)
	}
	return &cfg, nil
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown logging level %q", value)
	}
}

type servicePortConfig struct {
	Id      string `json:"Id"`
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

func loadServicePortConfig(serviceID string) (*servicePortConfig, error) {
	if serviceID == "" {
		return nil, fmt.Errorf("service id is empty")
	}
	path := filepath.Join(globconfig.GetServicesConfigDir(), serviceID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var cfg servicePortConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if cfg.Port == 0 {
		cfg.Port = portFromAddress(cfg.Address)
	}
	if cfg.Id == "" {
		cfg.Id = serviceID
	}
	return &cfg, nil
}

func applyServicePort(addr string, cfg *servicePortConfig) string {
	addr = strings.TrimSpace(addr)
	if cfg == nil {
		return addr
	}
	port := cfg.Port
	if port == 0 {
		port = portFromAddress(cfg.Address)
	}
	host := hostFromAddress(addr, "127.0.0.1")
	if strings.TrimSpace(cfg.Address) != "" {
		host = hostFromAddress(cfg.Address, host)
	}
	if port > 0 && host != "" {
		return fmt.Sprintf("%s:%d", host, port)
	}
	return addr
}

func emitXDSDescribe(addr string, cfg *servicePortConfig) error {
	addr = strings.TrimSpace(addr)
	port := portFromAddress(addr)
	outAddr := addr

	if cfg != nil {
		if cfg.Port > 0 {
			port = cfg.Port
		} else if p := portFromAddress(cfg.Address); p > 0 {
			port = p
		}
		if strings.TrimSpace(cfg.Address) != "" {
			outAddr = strings.TrimSpace(cfg.Address)
		}
	}

	if port <= 0 {
		return fmt.Errorf("xds port unavailable for describe")
	}
	if outAddr == "" {
		outAddr = fmt.Sprintf("localhost:%d", port)
	} else if portFromAddress(outAddr) == 0 {
		host := hostFromAddress(outAddr, "localhost")
		outAddr = fmt.Sprintf("%s:%d", host, port)
	}

	payload := servicePortConfig{
		Id:      xdsServiceID,
		Address: outAddr,
		Port:    port,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func hostFromAddress(addr, fallback string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return fallback
	}
	if strings.HasPrefix(addr, ":") {
		return fallback
	}
	if h, _, err := net.SplitHostPort(addr); err == nil {
		if h != "" {
			return h
		}
	}
	parts := strings.Split(addr, ":")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	if addr != "" && !strings.Contains(addr, ":") {
		return addr
	}
	return fallback
}

func portFromAddress(addr string) int {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return 0
	}
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	if _, portStr, err := net.SplitHostPort(addr); err == nil {
		if p, err := strconv.Atoi(portStr); err == nil && p > 0 {
			return p
		}
	}
	if idx := strings.LastIndex(addr, ":"); idx > 0 {
		if p, err := strconv.Atoi(addr[idx+1:]); err == nil && p > 0 {
			return p
		}
	}
	if p, err := strconv.Atoi(addr); err == nil && p > 0 {
		return p
	}
	return 0
}
