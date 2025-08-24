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

type Server struct {
	xdsserver server_v3.Server
}

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



// RunServer starts an xDS server at the given port.
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