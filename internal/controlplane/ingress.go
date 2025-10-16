// internal/controlplane/ingress.go
package controlplane

import (
	route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cors_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// IngressRoute describes one route entry for the shared ingress listener.
type IngressRoute struct {
	// Prefix is the path prefix to match (for gRPC use "/pkg.Service/").
	Prefix string
	// Cluster is the name of the upstream cluster to route to.
	Cluster string
	// HostRewrite (optional) sets :authority when forwarding upstream.
	HostRewrite string
}

// MakeRoutes builds the shared ingress RouteConfiguration with per-route prefixes.
// IMPORTANT: Enables CORS (with credentials) for browser callers.
func MakeRoutes(routeName string, rs []IngressRoute) *route_v3.RouteConfiguration {
	// EXACT origins only when AllowCredentials=true
	allowedOrigins := []*matcher_v3.StringMatcher{
		{MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: "http://localhost:5173"}},
		{MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: "https://globule-ryzen.globular.io"}},
		// add any other UIs here…
	}

	// Build the list of service routes
	var routes []*route_v3.Route
	for _, r := range rs {
		routes = append(routes, &route_v3.Route{
			Match: &route_v3.RouteMatch{
				PathSpecifier: &route_v3.RouteMatch_Prefix{Prefix: r.Prefix},
			},
			Action: &route_v3.Route_Route{
				Route: &route_v3.RouteAction{
					ClusterSpecifier: &route_v3.RouteAction_Cluster{Cluster: r.Cluster},
					Timeout:          durationpb.New(0), // no per-route timeout
				},
			},
		})
	}

	// Add a final catch-all route to your “globular_https” (optional, if you already add it elsewhere)
	// routes = append(routes, &route_v3.Route{ ... })
	vhost := &route_v3.VirtualHost{
		Name:    "ingress_vhost",
		Domains: []string{"*"},
		Routes:  routes,
		TypedPerFilterConfig: map[string]*anypb.Any{
			"envoy.filters.http.cors": toAny(&cors_v3.CorsPolicy{
				AllowOriginStringMatch: allowedOrigins,
				AllowCredentials:       wrapperspb.Bool(true),
				AllowMethods:           "GET, PUT, DELETE, POST, OPTIONS",
				AllowHeaders:           "keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout,domain,address,token,application,path,routing,authorization",
				ExposeHeaders:          "custom-header-1,grpc-status,grpc-message",
				MaxAge:                 "1728000",
			}),
		},
	}

	return &route_v3.RouteConfiguration{
		Name:         routeName,
		VirtualHosts: []*route_v3.VirtualHost{vhost},
	}
}
