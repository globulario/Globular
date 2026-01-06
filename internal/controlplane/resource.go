// Package controlplane provides utilities for constructing Envoy xDS resources,
// including clusters, endpoints, routes, listeners, and TLS transport sockets.
// It facilitates dynamic configuration of Envoy proxies for service discovery,
// load balancing, routing, and secure communication.
//
// Functions:
//
//   - MakeCluster: Creates an Envoy Cluster resource with optional TLS configuration and endpoints.
//   - makeUpstreamTLS: Generates an upstream TLS transport socket configuration for secure connections.
//   - makeDownstreamTLS: Generates a downstream TLS transport socket configuration for secure client connections.
//   - makeEndpoint: Constructs a ClusterLoadAssignment for load balancing across specified endpoints.
//   - MakeRoute: Builds a RouteConfiguration with CORS policy and host header rewriting.
//   - MakeHTTPListener: Configures an HTTP listener with filters for gRPC-Web, CORS, and routing, with optional TLS.
//   - makeConfigSource: Creates a ConfigSource for xDS gRPC API configuration.
//
// The package leverages Envoy's go-control-plane APIs and protobuf types to
// programmatically generate configuration resources for Envoy proxies.
package controlplane

import (
	"net"
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
	http_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes/any"
	structpb_v3 "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Builds: typed_extension_protocol_options["envoy.extensions.upstreams.http.v3.HttpProtocolOptions"]
func h2TypedOption() (map[string]*anypb.Any, error) {
	hp := &http_v3.HttpProtocolOptions{
		// Optional but useful:
		// CommonHttpProtocolOptions:    &corev3.HttpProtocolOptions{ /* idle_timeout, etc. */ },
		// UpstreamHttpProtocolOptions:  &corev3.UpstreamHttpProtocolOptions{ /* auto_sni, etc. */ },

		// *** This is the important part: choose the oneof = explicit_http_config
		UpstreamProtocolOptions: &http_v3.HttpProtocolOptions_ExplicitHttpConfig_{
			ExplicitHttpConfig: &http_v3.HttpProtocolOptions_ExplicitHttpConfig{
				// Inside that msg, choose protocol_config = http2_protocol_options
				ProtocolConfig: &http_v3.
					HttpProtocolOptions_ExplicitHttpConfig_Http2ProtocolOptions{
					Http2ProtocolOptions: &core_v3.Http2ProtocolOptions{},
				},
			},
		},
	}

	any, err := anypb.New(hp)
	if err != nil {
		return nil, err
	}
	return map[string]*anypb.Any{
		"envoy.extensions.upstreams.http.v3.HttpProtocolOptions": any,
	}, nil
}

// MakeCluster function constructs a cluster using the STRICT_DNS discovery type, sets a 5-second connection timeout,
// configures round-robin load balancing, and enables HTTP/2 protocol options. If a key file path is provided,
// it configures the cluster to use TLS for upstream connections by setting the transport socket.
//
// Parameters:
//   - clusterName:      The name of the Envoy cluster.
//   - certFilePath:     Path to the TLS certificate file (optional; required for TLS).
//   - keyFilePath:      Path to the TLS private key file (optional; required for TLS).
//   - caFilePath:       Path to the CA certificate file (optional; required for TLS).
//   - endPoints:        Slice of EndPoint structs specifying the cluster's endpoints.
//
// Returns:
//   - *cluster_v3.Cluster: A pointer to the constructed Envoy Cluster resource.
//
// MakeCluster builds an upstream cluster. If caFilePath != "" OR both cert/key are set,
// Envoy will use TLS (and HTTP/2) to the upstream. Pass SNI = the DNS name on the
// Globular 8181 certificate (e.g., "globular.io" or your host FQDN).
func MakeCluster(
	clusterName, certFilePath, keyFilePath, caFilePath, sni string,
	endPoints []EndPoint,
) *cluster_v3.Cluster {
	typed, _ := h2TypedOption() // HTTP/2 upstream

	// STATIC for IPs, STRICT_DNS for hostnames. (Never EDS here.)
	useStrictDNS := false
	for _, ep := range endPoints {
		if net.ParseIP(ep.Host) == nil {
			useStrictDNS = true
			break
		}
	}
	discoveryType := cluster_v3.Cluster_STATIC
	if useStrictDNS {
		discoveryType = cluster_v3.Cluster_STRICT_DNS
	}

	c := &cluster_v3.Cluster{
		Name:                          clusterName,
		ConnectTimeout:                durationpb.New(5 * time.Second),
		ClusterDiscoveryType:          &cluster_v3.Cluster_Type{Type: discoveryType},
		LbPolicy:                      cluster_v3.Cluster_ROUND_ROBIN,
		LoadAssignment:                makeEndpoint(clusterName, endPoints),
		DnsLookupFamily:               cluster_v3.Cluster_V4_ONLY,
		TypedExtensionProtocolOptions: typed, // HTTP/2 upstream
	}

	// Enable TLS if we have any trust material (CA) or mTLS (cert+key).
	if caFilePath != "" || (certFilePath != "" && keyFilePath != "") {
		c.TransportSocket = makeUpstreamTLS(certFilePath, keyFilePath, caFilePath, sni)
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

// makeUpstreamTLS creates an Envoy TransportSocket configuration for upstream TLS connections.
// It sets up the TLS context using the provided certificate, private key, and CA file paths.
// The function configures ALPN protocols for HTTP/2 and HTTP/1.1, and enforces TLS 1.3 as both
// the minimum and maximum protocol version.
//
// Parameters:
//
//	certFilePath - Path to the TLS certificate file.
//	keyFilePath  - Path to the TLS private key file.
//	caFilePath   - Path to the CA certificate file for validation.
//
// Returns:
//
//	*core_v3.TransportSocket - A pointer to the configured TransportSocket for TLS.
func makeUpstreamTLS(certFilePath, keyFilePath, caFilePath, sni string) *core_v3.TransportSocket {
	utls := &tls_v3.UpstreamTlsContext{
		Sni: sni,
		CommonTlsContext: &tls_v3.CommonTlsContext{
			AlpnProtocols: []string{"h2"}, // gRPC over TLS
		},
	}

	// Validate server (Globular) with CA if provided
	if caFilePath != "" {
		utls.CommonTlsContext.ValidationContextType = &tls_v3.CommonTlsContext_ValidationContext{
			ValidationContext: &tls_v3.CertificateValidationContext{
				TrustedCa: &core_v3.DataSource{
					Specifier: &core_v3.DataSource_Filename{Filename: caFilePath},
				},
			},
		}
	}

	// Optional mTLS: client cert/key from Envoy to Globular
	if certFilePath != "" && keyFilePath != "" {
		utls.CommonTlsContext.TlsCertificates = []*tls_v3.TlsCertificate{{
			CertificateChain: &core_v3.DataSource{
				Specifier: &core_v3.DataSource_Filename{Filename: certFilePath},
			},
			PrivateKey: &core_v3.DataSource{
				Specifier: &core_v3.DataSource_Filename{Filename: keyFilePath},
			},
		}}
	}

	return &core_v3.TransportSocket{
		Name: "envoy.transport_sockets.tls",
		ConfigType: &core_v3.TransportSocket_TypedConfig{
			TypedConfig: toAny(utls),
		},
	}
}

// makeDownstreamTLS creates an Envoy TransportSocket configuration for downstream TLS connections.
// It sets up the TLS context using the provided certificate, private key, and CA file paths.
// The function configures ALPN protocols for HTTP/2 and HTTP/1.1, and enforces TLS 1.3 as both
// the minimum and maximum protocol version.
//
// Parameters:
//
//	certFilePath - Path to the TLS certificate file.
//	keyFilePath  - Path to the TLS private key file.
//	caFilePath   - Path to the CA certificate file for client certificate validation.
//
// Returns:
//
//	*core_v3.TransportSocket - A pointer to the configured TransportSocket for downstream TLS.
func makeDownstreamTLS(certFilePath, keyFilePath, caFilePath string) *core_v3.TransportSocket {
	common := &tls_v3.CommonTlsContext{
		TlsCertificates: []*tls_v3.TlsCertificate{{
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
		}},
		AlpnProtocols: []string{"h2", "http/1.1"},
		TlsParams: &tls_v3.TlsParameters{
			TlsMinimumProtocolVersion: tls_v3.TlsParameters_TLSv1_3,
			TlsMaximumProtocolVersion: tls_v3.TlsParameters_TLSv1_3,
		},
	}

	if caFilePath != "" {
		common.ValidationContextType = &tls_v3.CommonTlsContext_ValidationContext{
			ValidationContext: &tls_v3.CertificateValidationContext{
				TrustedCa: &core_v3.DataSource{
					Specifier: &core_v3.DataSource_Filename{
						Filename: caFilePath,
					},
				},
			},
		}
	}

	return &core_v3.TransportSocket{
		Name: "envoy.transport_sockets.tls",
		ConfigType: &core_v3.TransportSocket_TypedConfig{
			TypedConfig: toAny(&tls_v3.DownstreamTlsContext{
				CommonTlsContext: common,
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

// MakeRoute creates an Envoy RouteConfiguration with a single virtual host and route.
// The route matches all paths ("/") and forwards requests to the specified cluster,
// rewriting the host header to the provided host value. CORS policy is configured to
// allow all origins, specific HTTP methods, and headers. The route has an infinite timeout.
//
// Parameters:
//   - routeName:     Name of the route configuration.
//   - clusterName:   Name of the target cluster for routing.
//   - host:          Host value to rewrite in the request.
//
// Returns:
//   - *route_v3.RouteConfiguration: Pointer to the constructed RouteConfiguration.
func MakeRoute(routeName string, clusterName, host string) *route_v3.RouteConfiguration {
	allowedOrigins := []*matcher_v3.StringMatcher{
		// List each origin that may send credentials
		{MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: "http://localhost:5173"}},
		{MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: "https://globule-ryzen.globular.io"}},
		// add more exact origins here if needed
	}

	return &route_v3.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route_v3.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*route_v3.Route{{
				Match: &route_v3.RouteMatch{
					PathSpecifier: &route_v3.RouteMatch_Prefix{Prefix: "/"},
				},
				Action: &route_v3.Route_Route{
					Route: &route_v3.RouteAction{
						ClusterSpecifier: &route_v3.RouteAction_Cluster{Cluster: clusterName},
						HostRewriteSpecifier: &route_v3.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: host,
						},
						Timeout: durationpb.New(time.Duration(0)), // no per-route timeout
					},
				},
			}},
			TypedPerFilterConfig: map[string]*anypb.Any{
				"envoy.filters.http.cors": toAny(&cors_v3.CorsPolicy{
					AllowOriginStringMatch: allowedOrigins,
					// Must be true when the browser sends credentials (cookies, auth headers)
					AllowCredentials: wrapperspb.Bool(true),

					AllowMethods:  "GET, PUT, DELETE, POST, OPTIONS",
					AllowHeaders:  "keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout,domain,address,token,application,path,routing,authorization",
					ExposeHeaders: "custom-header-1,grpc-status,grpc-message",
					MaxAge:        "1728000",
				}),
			},
		}},
	}
}

// MakeHTTPListener creates and configures an Envoy HTTP listener with the specified parameters.
// It sets up HTTP filters for gRPC-Web, CORS, and routing, and optionally configures TLS if certificate and key file paths are provided.
//
// Parameters:
//
//	listenerHost   - The host address to bind the listener.
//	listenerPort   - The port number to bind the listener.
//	listenerName   - The name of the listener.
//	routeName      - The name of the route configuration.
//	certFilePath   - Path to the TLS certificate file (optional).
//	keyFilePath    - Path to the TLS key file (optional).
//	caFilePath     - Path to the CA certificate file (optional).
//
// Returns:
//
//	*listener.Listener - A pointer to the configured Envoy listener.
func MakeHTTPListener(listenerHost string, listenerPort uint32, listenerName, routeName, certFilePath, keyFilePath, caFilePath string) *listener_v3.Listener {

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
		l.FilterChains[0].TransportSocket = makeDownstreamTLS(certFilePath, keyFilePath, caFilePath)
	}

	return l
}

// !!!! here the clusterName must be the one in envoy.yml...
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

func fileDS(path string) *core_v3.DataSource {
	if path == "" {
		return nil
	}
	return &core_v3.DataSource{Specifier: &core_v3.DataSource_Filename{Filename: path}}
}
