package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	clusterservice_v3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice_v3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice_v3 "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice_v3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice_v3 "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	cache_v3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	server_v3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	test_v3 "github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

// XDSServer hosts Envoy ADS/xDS services backed by an in-memory snapshot cache.
type XDSServer struct {
	logger *slog.Logger
	cache  cache_v3.SnapshotCache
	server server_v3.Server
}

// New builds a managed xDS server instance.
func New(logger *slog.Logger, ctx context.Context) *XDSServer {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	cb := &test_v3.Callbacks{}
	cache := cache_v3.NewSnapshotCache(false, cache_v3.IDHash{}, nil)
	server := server_v3.NewServer(ctx, cache, cb)
	return &XDSServer{
		logger: logger,
		cache:  cache,
		server: server,
	}
}

// Serve starts the gRPC xDS server and returns when Serve exits.
func (s *XDSServer) Serve(ctx context.Context, addr string) error {
	if strings.TrimSpace(addr) == "" {
		return fmt.Errorf("grpc address is required")
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	}
	grpcServer := grpc.NewServer(opts...)
	registerServices(grpcServer, s.server)

	s.logger.Info("xDS server listening", "addr", addr)
	go func() {
		<-ctx.Done()
		s.logger.Info("shutting down xDS server")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		if ctx.Err() == nil {
			return fmt.Errorf("serve: %w", err)
		}
	}
	return nil
}

// SetSnapshot replaces snapshot data for the provided node ID.
func (s *XDSServer) SetSnapshot(nodeID string, snap *cache_v3.Snapshot) error {
	if strings.TrimSpace(nodeID) == "" {
		return fmt.Errorf("node id required")
	}
	if snap == nil {
		return fmt.Errorf("snapshot is nil")
	}
	if err := snap.Consistent(); err != nil {
		return fmt.Errorf("snapshot consistency: %w", err)
	}
	return s.cache.SetSnapshot(context.Background(), nodeID, snap)
}

func registerServices(grpcServer *grpc.Server, server server_v3.Server) {
	discoverygrpc_v3.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice_v3.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice_v3.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice_v3.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice_v3.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secretservice_v3.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtimeservice_v3.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
}
