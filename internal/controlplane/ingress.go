// internal/controlplane/ingress.go
package controlplane

import (
	"time"

	route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cors_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/types/known/durationpb"
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

// MakeRoutes builds a single RouteConfiguration that contains multiple
// routes (one per service) for use by a shared ingress listener.
func MakeRoutes(routeName string, items []IngressRoute) *route_v3.RouteConfiguration {
	var rs []*route_v3.Route
	for _, it := range items {
		r := &route_v3.Route{
			Match: &route_v3.RouteMatch{
				PathSpecifier: &route_v3.RouteMatch_Prefix{Prefix: it.Prefix},
			},
			Action: &route_v3.Route_Route{
				Route: &route_v3.RouteAction{
					ClusterSpecifier: &route_v3.RouteAction_Cluster{Cluster: it.Cluster},
					Timeout:          durationpb.New(time.Duration(0)),
				},
			},
		}
		if it.HostRewrite != "" {
			r.GetRoute().HostRewriteSpecifier = &route_v3.RouteAction_HostRewriteLiteral{
				HostRewriteLiteral: it.HostRewrite,
			}
		}
		rs = append(rs, r)
	}

	return &route_v3.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route_v3.VirtualHost{{
			Name:    "ingress",
			Domains: []string{"*"},
			Routes:  rs,
			TypedPerFilterConfig: map[string]*any.Any{
				"envoy.filters.http.cors": toAny(defaultCORSPolicy()),
			},
		}},
	}
}

func defaultCORSPolicy() *cors_v3.CorsPolicy {
	return &cors_v3.CorsPolicy{
		AllowOriginStringMatch: []*matcher_v3.StringMatcher{
			{MatchPattern: &matcher_v3.StringMatcher_Prefix{Prefix: "*"}},
		},
		AllowMethods:  "GET, PUT, DELETE, POST, OPTIONS",
		AllowHeaders:  "keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout, domain, address, token, application, path, routing",
		MaxAge:        "1728000",
		ExposeHeaders: "custom-header-1,grpc-status,grpc-message",
	}
}
