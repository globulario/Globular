package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

// TLSConfig holds TLS configuration for the xDS server.
// If nil or all fields empty, server runs without TLS (insecure, for localhost only).
type TLSConfig struct {
	// ServerCertPath is the path to the server certificate PEM file
	ServerCertPath string
	// ServerKeyPath is the path to the server private key PEM file
	ServerKeyPath string
	// ClientCAPath is the path to the CA bundle for validating client certificates
	// If set, client certificate authentication is required (mTLS)
	ClientCAPath string
}

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
// If tlsConfig is nil or has empty paths, server runs without TLS (insecure).
func (s *XDSServer) Serve(ctx context.Context, addr string, tlsConfig *TLSConfig) error {
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

	// Add TLS credentials if configured
	if tlsConfig != nil && tlsConfig.ServerCertPath != "" && tlsConfig.ServerKeyPath != "" {
		tlsCreds, err := buildTLSCredentials(tlsConfig)
		if err != nil {
			return fmt.Errorf("build TLS credentials: %w", err)
		}
		opts = append(opts, grpc.Creds(tlsCreds))
		if tlsConfig.ClientCAPath != "" {
			s.logger.Info("xDS server using mTLS", "addr", addr, "client_auth", "required")
		} else {
			s.logger.Info("xDS server using TLS", "addr", addr, "client_auth", "none")
		}
	} else {
		s.logger.Warn("xDS server running without TLS (insecure)", "addr", addr)
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

// buildTLSCredentials creates TLS credentials for the xDS gRPC server.
// If ClientCAPath is set, requires and validates client certificates (mTLS).
func buildTLSCredentials(cfg *TLSConfig) (credentials.TransportCredentials, error) {
	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(cfg.ServerCertPath, cfg.ServerKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load server cert/key: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Require TLS 1.2+
	}

	// If client CA is configured, require and validate client certificates
	if cfg.ClientCAPath != "" {
		caPEM, err := os.ReadFile(cfg.ClientCAPath)
		if err != nil {
			return nil, fmt.Errorf("read client CA: %w", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caPEM) {
			return nil, fmt.Errorf("failed to parse client CA certificate")
		}

		tlsConfig.ClientCAs = certPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return credentials.NewTLS(tlsConfig), nil
}
