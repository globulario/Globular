package controlplane

import (
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cors "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	grpc_web "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_web/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	v31 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes/any"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

/**
 * makeCluster creates a new cluster from a name, LB policy, and a list of hosts.
 * The cluster is configured with a DNS resolver to allow Envoy to resolve the
 * cluster IP address from the DNS name.
 */
func MakeCluster(clusterName, certFilePath, keyFilePath, caFilePath string, endPoints []EndPoint) *cluster.Cluster {

	// Create the cluster
	c := &cluster.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName, endPoints),
		DnsLookupFamily:      cluster.Cluster_V4_ONLY,
		Http2ProtocolOptions: &core.Http2ProtocolOptions{},
	}

	// I case of TLS, we need to set the transport socket
	if len(certFilePath) > 0 {
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
func makeUpstreamTls(certFilePath, keyFilePath, caFilePath string) *core.TransportSocket {
	return &core.TransportSocket{
		Name: "envoy.transport_sockets.tls",
		ConfigType: &core.TransportSocket_TypedConfig{
			TypedConfig: toAny(&tls.UpstreamTlsContext{
				CommonTlsContext: &tls.CommonTlsContext{
					TlsCertificates: []*tls.TlsCertificate{
						{
							CertificateChain: &core.DataSource{
								Specifier: &core.DataSource_Filename{
									Filename: certFilePath,
								},
							},
							PrivateKey: &core.DataSource{
								Specifier: &core.DataSource_Filename{
									Filename: keyFilePath,
								},
							},
						},
					},
					ValidationContextType: &tls.CommonTlsContext_ValidationContext{
						ValidationContext: &tls.CertificateValidationContext{
							TrustedCa: &core.DataSource{
								Specifier: &core.DataSource_Filename{
									Filename: caFilePath,
								},
							},
						},
					},
					AlpnProtocols: []string{"h2", "http/1.1"},
					TlsParams:     &tls.TlsParameters{TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_3, TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3},
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
func makeDownstreamTls(certFilePath, keyFilePath, caFilePath string) *core.TransportSocket {
	return &core.TransportSocket{
		Name: "envoy.transport_sockets.tls",
		ConfigType: &core.TransportSocket_TypedConfig{
			TypedConfig: toAny(&tls.DownstreamTlsContext{
				CommonTlsContext: &tls.CommonTlsContext{
					TlsCertificates: []*tls.TlsCertificate{
						{
							CertificateChain: &core.DataSource{
								Specifier: &core.DataSource_Filename{
									Filename: certFilePath,
								},
							},
							PrivateKey: &core.DataSource{
								Specifier: &core.DataSource_Filename{
									Filename: keyFilePath,
								},
							},
						},
					},
					ValidationContextType: &tls.CommonTlsContext_ValidationContext{
						ValidationContext: &tls.CertificateValidationContext{
							TrustedCa: &core.DataSource{
								Specifier: &core.DataSource_Filename{
									Filename: caFilePath,
								},
							},
						},
					},
					AlpnProtocols: []string{"h2", "http/1.1"},
					TlsParams:     &tls.TlsParameters{TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_3, TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3},
				},
			}),
		},
	}
}

func makeEndpoint(clusterName string, endPoints []EndPoint) *endpoint.ClusterLoadAssignment {

	load_assignment := &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{},
		}},
	}

	for _, endPoint := range endPoints {
		load_assignment.Endpoints[0].LbEndpoints = append(load_assignment.Endpoints[0].LbEndpoints,
			&endpoint.LbEndpoint{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  endPoint.Host,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: endPoint.Port,
									},
								},
							},
						},
					},
				},
				Metadata: &core.Metadata{
					FilterMetadata: map[string]*structpb.Struct{
						"envoy.lb": {
							Fields: map[string]*structpb.Value{
								"priority": {
									Kind: &structpb.Value_NumberValue{
										NumberValue: float64(endPoint.Priority),
									},
								},
							},
						},
					},
				},
			},
		)
	}

	return load_assignment
}

/**
 * makeRoute creates a new route for Envoy to forward HTTP requests to the
 * upstream cluster.
 */
func MakeRoute(routeName string, clusterName, host string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: host,
						},
					},
				},
			}},
			// Add CORS policy to allow all origins.
			Cors: &route.CorsPolicy{
				AllowOriginStringMatch: []*v31.StringMatcher{
					{
						MatchPattern: &v31.StringMatcher_Prefix{
							Prefix: "*",
						},
					},
				},
				AllowMethods:  "GET, PUT, DELETE, POST, OPTIONS",
				AllowHeaders:  "keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout, domain, address, token, application, path",
				MaxAge:        "1728000",
				ExposeHeaders: "custom-header-1,grpc-status,grpc-message",
			},
		}},
	}
}

/**
 * makeHTTPListener creates a new HTTP(s) listener.
 * The listener config references the RDS configuration defined below.
 */
func MakeHTTPListener(listenerHost string, listenerPort uint32, listenerName, clusterName, routeName, certFilePath, keyFilePath, caFilePath string) *listener.Listener {

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(), // TODO: ---> xds_cluster
				RouteConfigName: routeName,
			},
		},
		// Had necessary filters...
		HttpFilters: []*hcm.HttpFilter{
			{
				Name: "envoy.filters.http.grpc_web",
				ConfigType: &hcm.HttpFilter_TypedConfig{
					TypedConfig: toAny(&grpc_web.GrpcWeb{}),
				},
			},
			{
				Name: "envoy.filters.http.cors",
				ConfigType: &hcm.HttpFilter_TypedConfig{
					TypedConfig: toAny(&cors.Cors{}),
				},
			},
			{
				Name: "envoy.filters.http.router",
				ConfigType: &hcm.HttpFilter_TypedConfig{
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

	l := &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  listenerHost,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: listenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: "envoy.filters.network.http_connection_manager",
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}

	// I case of TLS, we need to set the transport socket
	if len(certFilePath) > 0 {
		l.FilterChains[0].TransportSocket = makeDownstreamTls(certFilePath, keyFilePath, caFilePath)
	}

	return l
}

// !!!! here the clusterName can be the one in envoy.yaml...
// ---> xds_cluster
func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}