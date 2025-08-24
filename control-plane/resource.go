package controlplane

import (
	"time"
	cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cors_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	grpc_web_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_web/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	structpb_v3 "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

/**
 * makeCluster creates a new cluster from a name, LB policy, and a list of hosts.
 * The cluster is configured with a DNS resolver to allow Envoy to resolve the
 * cluster IP address from the DNS name.
 */
func MakeCluster(clusterName, certFilePath, keyFilePath, caFilePath string, endPoints []EndPoint) *cluster_v3.Cluster {

	// Create the cluster
	c := &cluster_v3.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster_v3.Cluster_Type{Type: cluster_v3.Cluster_STRICT_DNS},
		LbPolicy:             cluster_v3.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName, endPoints),
		DnsLookupFamily:      cluster_v3.Cluster_V4_ONLY,
		Http2ProtocolOptions: &core_v3.Http2ProtocolOptions{},
	}

	// In case of TLS, we need to set the transport socket
	if len(keyFilePath) > 0 {
		c.TransportSocket = makeUpstreamTls(certFilePath, keyFilePath, caFilePath)
	}

	return c
}

func toAny(msg protoreflect.ProtoMessage) *any.Any {
	anyMsg, err := anypb.New(msg)
	if err != nil {
		panic(err)
	}
	return anyMsg
}

/**
 * make TLS creates a TLS transport socket config for upstream connections.
 * This config is intended for use with SDS.
 */
func makeUpstreamTls(certFilePath, keyFilePath, caFilePath string) *core_v3.TransportSocket {
	return &core_v3.TransportSocket{
		Name: "envoy.transport_sockets.tls",
		ConfigType: &core_v3.TransportSocket_TypedConfig{
			TypedConfig: toAny(&tls_v3.UpstreamTlsContext{
				CommonTlsContext: &tls_v3.CommonTlsContext{
					TlsCertificates: []*tls_v3.TlsCertificate{
						{
							CertificateChain: &core_v3.DataSource{
								Specifier: &core_v3.DataSource_Filename{
									Filename: certFilePath,
								},
							},
							PrivateKey: &core_v3.DataSource{
								Specifier: &core_v3.DataSource_Filename{
									Filename: keyFilePath,
								},
							},
						},
					},
					ValidationContextType: &tls_v3.CommonTlsContext_ValidationContext{
						ValidationContext: &tls_v3.CertificateValidationContext{
							TrustedCa: &core_v3.DataSource{
								Specifier: &core_v3.DataSource_Filename{
									Filename: caFilePath,
								},
							},
						},
					},
					AlpnProtocols: []string{"h2", "http/1.1"},
					TlsParams:     &tls_v3.TlsParameters{TlsMinimumProtocolVersion: tls_v3.TlsParameters_TLSv1_3, TlsMaximumProtocolVersion: tls_v3.TlsParameters_TLSv1_3},
				},
			}),
		},
	}
}

/**
 * make TLS creates a TLS transport socket config for downstream connections.
 * This config is intended for use with SDS. It is not necessary to specify
 * the downstream TLS context if the downstream client is not configured for TLS.
 * In this case, Envoy will use plaintext downstream connections.
 */
func makeDownstreamTls(certFilePath, keyFilePath, caFilePath string) *core_v3.TransportSocket {
	return &core_v3.TransportSocket{
		Name: "envoy.transport_sockets.tls",
		ConfigType: &core_v3.TransportSocket_TypedConfig{
			TypedConfig: toAny(&tls_v3.DownstreamTlsContext{
				CommonTlsContext: &tls_v3.CommonTlsContext{
					TlsCertificates: []*tls_v3.TlsCertificate{
						{
							CertificateChain: &core_v3.DataSource{
								Specifier: &core_v3.DataSource_Filename{
									Filename: certFilePath,
								},
							},
							PrivateKey: &core_v3.DataSource{
								Specifier: &core_v3.DataSource_Filename{
									Filename: keyFilePath,
								},
							},
						},
					},
					ValidationContextType: &tls_v3.CommonTlsContext_ValidationContext{
						ValidationContext: &tls_v3.CertificateValidationContext{
							TrustedCa: &core_v3.DataSource{
								Specifier: &core_v3.DataSource_Filename{
									Filename: caFilePath,
								},
							},
						},
					},
					AlpnProtocols: []string{"h2", "http/1.1"},
					TlsParams:     &tls_v3.TlsParameters{TlsMinimumProtocolVersion: tls_v3.TlsParameters_TLSv1_3, TlsMaximumProtocolVersion: tls_v3.TlsParameters_TLSv1_3},
				},
			}),
		},
	}
}

// makeEndpoint constructs an Envoy ClusterLoadAssignment for the given cluster name and endpoints.
// It creates a list of LbEndpoint objects, each representing an endpoint with its address, port, and priority metadata.
// The resulting ClusterLoadAssignment can be used to configure Envoy's load balancing for the specified cluster.
//
// Parameters:
//   - clusterName: The name of the cluster for which the load assignment is created.
//   - endPoints: A slice of EndPoint structs, each containing host, port, and priority information.
//
// Returns:
//   - A pointer to an endpoint.ClusterLoadAssignment populated with the provided endpoints.
func makeEndpoint(clusterName string, endPoints []EndPoint) *endpoint_v3.ClusterLoadAssignment {
	var lbEndpoints []*endpoint_v3.LbEndpoint

	for _, endPoint := range endPoints {

		lbEndpoint := &endpoint_v3.LbEndpoint{
			HostIdentifier: &endpoint_v3.LbEndpoint_Endpoint{
				Endpoint: &endpoint_v3.Endpoint{
					Address: &core_v3.Address{
						Address: &core_v3.Address_SocketAddress{
							SocketAddress: &core_v3.SocketAddress{
								Protocol: core_v3.SocketAddress_TCP,
								Address:  endPoint.Host,
								PortSpecifier: &core_v3.SocketAddress_PortValue{
									PortValue: endPoint.Port,
								},
							},
						},
					},
				},
			},
			Metadata: &core_v3.Metadata{
				FilterMetadata: map[string]*structpb_v3.Struct{
					"envoy.lb": {
						Fields: map[string]*structpb_v3.Value{
							"priority": {
								Kind: &structpb_v3.Value_NumberValue{
									NumberValue: float64(endPoint.Priority),
								},
							},
						},
					},
				},
			},
		}

		lbEndpoints = append(lbEndpoints, lbEndpoint)
	}

	return &endpoint_v3.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint_v3.LocalityLbEndpoints{{
			LbEndpoints: lbEndpoints,
		}},
	}
}


// MakeRoute creates and returns a RouteConfiguration for Envoy with the specified route name,
// cluster name, and host. The configuration includes a virtual host with a single route that matches
// all paths ("/") and forwards requests to the specified cluster, rewriting the host header.
// It also sets an infinite timeout for requests and configures CORS policy to allow all origins,
// specific HTTP methods, headers, and exposes certain headers.
func MakeRoute(routeName string, clusterName, host string) *route_v3.RouteConfiguration {
	return &route_v3.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route_v3.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*route_v3.Route{{

				Match: &route_v3.RouteMatch{
					PathSpecifier: &route_v3.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &route_v3.Route_Route{
					Route: &route_v3.RouteAction{
						ClusterSpecifier: &route_v3.RouteAction_Cluster{
							Cluster: clusterName,
						},
						HostRewriteSpecifier: &route_v3.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: host,
						},
						Timeout: ptypes.DurationProto(time.Duration(0)), // Infinite timeout
					},
				},
			}},
			TypedPerFilterConfig: map[string]*any.Any{
				"envoy.filters.http.cors": toAny(&cors_v3.CorsPolicy{
					AllowOriginStringMatch: []*matcher_v3.StringMatcher{
						{
							MatchPattern: &matcher_v3.StringMatcher_Prefix{
								Prefix: "*",
							},
						},
					},
					AllowMethods:  "GET, PUT, DELETE, POST, OPTIONS",
					AllowHeaders:  "keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout, domain, address, token, application, path, routing",
					MaxAge:        "1728000",
					ExposeHeaders: "custom-header-1,grpc-status,grpc-message",
				}),
			},
		}},
	}
}



// MakeHTTPListener creates and configures an Envoy HTTP listener with the specified parameters.
// It sets up HTTP filters for gRPC-Web, CORS, and routing, and optionally configures TLS if certificate and key file paths are provided.
//
// Parameters:
//   listenerHost   - The host address to bind the listener.
//   listenerPort   - The port number to bind the listener.
//   listenerName   - The name of the listener.
//   clusterName    - The name of the cluster (not directly used in this function).
//   routeName      - The name of the route configuration.
//   certFilePath   - Path to the TLS certificate file (optional).
//   keyFilePath    - Path to the TLS key file (optional).
//   caFilePath     - Path to the CA certificate file (optional).
//
// Returns:
//   *listener.Listener - A pointer to the configured Envoy listener.
func MakeHTTPListener(listenerHost string, listenerPort uint32, listenerName, clusterName, routeName, certFilePath, keyFilePath, caFilePath string) *listener_v3.Listener {
	// HTTP filter configuration
	manager := &hcm_v3.HttpConnectionManager{
		CodecType:  hcm_v3.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm_v3.HttpConnectionManager_Rds{
			Rds: &hcm_v3.Rds{
				ConfigSource:    makeConfigSource(), // TODO: ---> xds_cluster
				RouteConfigName: routeName,
			},
		},
		// Had necessary filters...
		HttpFilters: []*hcm_v3.HttpFilter{
			{
				Name: "envoy.filters.http.grpc_web",
				ConfigType: &hcm_v3.HttpFilter_TypedConfig{
					TypedConfig: toAny(&grpc_web_v3.GrpcWeb{}),
				},
			},
			{
				Name: "envoy.filters.http.cors",
				ConfigType: &hcm_v3.HttpFilter_TypedConfig{
					TypedConfig: toAny(&cors_v3.Cors{}),
				},
			},
			{
				Name: "envoy.filters.http.router",
				ConfigType: &hcm_v3.HttpFilter_TypedConfig{
					TypedConfig: toAny(&router.Router{}),
				},
			},
		},
	}

	// Add GrpcWeb filter if needed
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	l := &listener_v3.Listener{
		Name: listenerName,
		Address: &core_v3.Address{
			Address: &core_v3.Address_SocketAddress{
				SocketAddress: &core_v3.SocketAddress{
					Protocol: core_v3.SocketAddress_TCP,
					Address:  listenerHost,
					PortSpecifier: &core_v3.SocketAddress_PortValue{
						PortValue: listenerPort,
					},
				},
			},
		},
		FilterChains: []*listener_v3.FilterChain{{
			Filters: []*listener_v3.Filter{{
				Name: "envoy.filters.network.http_connection_manager",
				ConfigType: &listener_v3.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}

	// I case of TLS, we need to set the transport socket
	if len(keyFilePath) > 0 {
		l.FilterChains[0].TransportSocket = makeDownstreamTls(certFilePath, keyFilePath, caFilePath)
	}

	return l
}

// !!!! here the clusterName can be the one in envoy.yml...
// ---> xds_cluster
func makeConfigSource() *core_v3.ConfigSource {
	source := &core_v3.ConfigSource{}
	source.ResourceApiVersion = resource_v3.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core_v3.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core_v3.ApiConfigSource{
			TransportApiVersion:       resource_v3.DefaultAPIVersion,
			ApiType:                   core_v3.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core_v3.GrpcService{{
				TargetSpecifier: &core_v3.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core_v3.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}
