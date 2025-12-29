package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/globulario/Globular/internal/controlplane"
	"github.com/globulario/Globular/internal/xds/server"
	"github.com/globulario/Globular/internal/xds/watchers"
)

const (
	defaultGRPCAddr  = "0.0.0.0:18000"
	defaultXDSConfig = "/etc/globular/xds/config.json"
)

func main() {
	var (
		grpcAddr           = flag.String("grpc_addr", defaultGRPCAddr, "address for the xDS gRPC server")
		xdsConfigPath      = flag.String("xds_config", defaultXDSConfig, "path to the xDS config JSON")
		nodeID             = flag.String("node_id", "globular-xds", "node id for xDS snapshots")
		bootstrapPath      = flag.String("envoy_bootstrap", "", "path to write Envoy bootstrap YAML")
		bootstrapHost      = flag.String("bootstrap_host", "127.0.0.1", "host advertised in envoy bootstrap")
		bootstrapPort      = flag.Int("bootstrap_port", 18000, "port advertised in envoy bootstrap")
		bootstrapCluster   = flag.String("bootstrap_cluster", "globular-cluster", "cluster name used in bootstrap")
		bootstrapAdminPort = flag.Int("bootstrap_admin", 9901, "admin port in bootstrap")
		bootstrapMaxConn   = flag.Uint64("bootstrap_max_downstream_connections", 0, "max active downstream connections (0 disables)")
		downstreamTLSMode  = flag.String("downstream_tls_mode", "optional", "downstream TLS mode (disabled|optional|required)")
		watchInterval      = flag.Duration("watch_interval", 5*time.Second, "static config polling interval")
		debugSignals       = flag.Bool("debug_signals", false, "log signals received while running")
	)
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

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
	go func() {
		if err := xdsServer.Serve(ctx, *grpcAddr); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		}
	}()

	watcher := watchers.New(logger, xdsServer, *xdsConfigPath, *nodeID, *watchInterval, watchers.ParseDownstreamTLSMode(*downstreamTLSMode))
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

func writeBootstrap(logger *slog.Logger, opts controlplane.BootstrapOptions, overridePath string) {
	data, err := controlplane.MarshalBootstrap(opts)
	if err != nil {
		logger.Warn("marshal bootstrap", "err", err)
		return
	}

	paths := []string{}
	if overridePath != "" {
		paths = append(paths, overridePath)
	} else {
		paths = bootstrapCandidates()
	}

	for _, p := range paths {
		dir := filepath.Dir(p)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			logger.Warn("bootstrap dir create", "err", err, "dir", dir)
			continue
		}
		if err := os.WriteFile(p, data, 0o644); err != nil {
			logger.Warn("write bootstrap", "err", err, "path", p)
			continue
		}
		logger.Info("wrote bootstrap", "path", p)
		return
	}
	logger.Warn("failed to write bootstrap to any candidate")
}

func bootstrapCandidates() []string {
	return []string{
		filepath.Join("/run/globular", "envoy", "envoy.yml"),
		filepath.Join("/var/lib/globular", "envoy", "envoy.yml"),
	}
}
