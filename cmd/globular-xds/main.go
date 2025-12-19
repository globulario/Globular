package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/globulario/Globular/internal/controlplane"
	"github.com/globulario/Globular/internal/xds/server"
	"github.com/globulario/Globular/internal/xds/watchers"
	"github.com/globulario/services/golang/config"
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
		watchInterval      = flag.Duration("watch_interval", 5*time.Second, "static config polling interval")
	)
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if *bootstrapPath == "" {
		*bootstrapPath = filepath.Join(config.GetConfigDir(), "envoy.yml")
	}
	if err := writeBootstrap(*bootstrapPath, *nodeID, *bootstrapCluster, *bootstrapHost, *bootstrapPort, *bootstrapAdminPort, *bootstrapMaxConn); err != nil {
		logger.Warn("write bootstrap", "err", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	xdsServer := server.New(logger, ctx)
	go func() {
		if err := xdsServer.Serve(ctx, *grpcAddr); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("xDS serve failed", "err", err)
		}
	}()

	watcher := watchers.New(logger, xdsServer, *xdsConfigPath, *nodeID, *watchInterval)
	if err := watcher.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("xDS watcher stopped", "err", err)
	}
}

func writeBootstrap(path, nodeID, cluster, host string, port, admin int, maxConn uint64) error {
	opts := controlplane.BootstrapOptions{
		NodeID:    nodeID,
		Cluster:   cluster,
		XDSHost:   host,
		XDSPort:   port,
		AdminPort: admin,
	}
	if maxConn > 0 {
		opts.MaxActiveDownstreamConns = maxConn
	}
	return controlplane.WriteBootstrap(path, opts)
}
