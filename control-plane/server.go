// Package controlplane provides an xDS management server implementation for Envoy,
// leveraging the go-control-plane library. It offers functionality to create, configure,
// and run a gRPC management server with custom keepalive and stream options to support
// high concurrency and robust connection management.
//
// The Server type wraps the go-control-plane server implementation and exposes methods
// for registering xDS services and running the management server. Utility functions are
// provided for direct server registration and startup.
//
// Key Features:
//   - Customizable gRPC keepalive and stream settings for reliability and scalability.
//   - Registration of all major xDS services (ADS, Cluster, Endpoint, Route, Listener, Secret, Runtime).
//   - Simple API for server instantiation and execution.
//
// Example usage:
//
//	cache := ... // create a cache_v3.Cache implementation
//	callbacks := ... // create test_v3.Callbacks
//	srv := controlplane.NewServer(context.Background(), cache, callbacks)
//	go srv.Run(18000)
//
//	// Or using RunServer directly:
//	controlplane.RunServer(srv.xdsserver, 18000)
package controlplane

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

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
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

// Server represents an xDS management server for Envoy, wrapping the go-control-plane server implementation.
type Server struct {
	xdsserver server_v3.Server
}

// NewServer creates and returns a new Server instance using the provided context, cache, and callbacks.
// It wraps the underlying server_v3.Server with additional functionality as defined by the Server type.
//
// Parameters:
//
//	ctx   - The context for controlling cancellation and deadlines.
//	cache - The cache implementation to be used by the server.
//	cb    - The callbacks to handle server events.
//
// Returns:
//
//	A new Server instance.
func NewServer(ctx context.Context, cache cache_v3.Cache, cb *test_v3.Callbacks) Server {
	srv := server_v3.NewServer(ctx, cache, cb)
	return Server{srv}
}

func (s *Server) registerServer(grpcServer *grpc.Server) {
	// register services
	discoverygrpc_v3.RegisterAggregatedDiscoveryServiceServer(grpcServer, s.xdsserver)
	endpointservice_v3.RegisterEndpointDiscoveryServiceServer(grpcServer, s.xdsserver)
	clusterservice_v3.RegisterClusterDiscoveryServiceServer(grpcServer, s.xdsserver)
	routeservice_v3.RegisterRouteDiscoveryServiceServer(grpcServer, s.xdsserver)
	listenerservice_v3.RegisterListenerDiscoveryServiceServer(grpcServer, s.xdsserver)
	secretservice_v3.RegisterSecretDiscoveryServiceServer(grpcServer, s.xdsserver)
	runtimeservice_v3.RegisterRuntimeDiscoveryServiceServer(grpcServer, s.xdsserver)
}

// Run starts the gRPC management server on the specified port with custom keepalive and stream options.
// It configures the server to handle a high number of concurrent streams and sets keepalive parameters
// to ensure connection reliability, especially when requests are multiplexed over a single TCP connection.
// The server is registered and begins listening for incoming connections. If the server fails to start,
// the error is logged and the process is terminated.
//
// Parameters:
//
//	port - The TCP port number on which the server will listen.
//
// The function logs the listening port and any errors encountered during startup or while serving.
func (s *Server) Run(port uint) {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems. Keepalive timeouts based on connection_keepalive parameter
	// https://www.envoyproxy.io/docs/envoy/latest/configuration/overview/examples#dynamic
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s.registerServer(grpcServer)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}

func registerServer(grpcServer *grpc.Server, server server_v3.Server) {
	// register services
	discoverygrpc_v3.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice_v3.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice_v3.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice_v3.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice_v3.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secretservice_v3.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtimeservice_v3.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
}

// RunServer starts a gRPC server for the provided Server implementation on the specified port.
// It configures gRPC server options including maximum concurrent streams and keepalive parameters
// to ensure robust connection management. The function listens on the given TCP port, registers
// the server implementation, and begins serving incoming gRPC requests. If the server fails to
// start or encounters an error, it logs the error and terminates the process.
//
// Parameters:
//
//	srv  - The server_v3.Server implementation to register and serve.
//	port - The TCP port number to listen on for incoming gRPC connections.
func RunServer(srv server_v3.Server, port uint) {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems. Keepalive timeouts based on connection_keepalive parameter
	// https://www.envoyproxy.io/docs/envoy/latest/configuration/overview/examples#dynamic
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	registerServer(grpcServer, srv)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}
