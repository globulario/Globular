// internal/controlplane/ingress.go
package controlplane

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"

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
	// Domains optionally restricts the route to these virtual host domains.
	Domains []string
	// Authority enables an exact :authority match within the route.
	Authority string
}

// MakeRoutes builds the shared ingress RouteConfiguration with per-route prefixes.
// IMPORTANT: Enables CORS (with credentials) for browser callers.
func MakeRoutes(routeName string, rs []IngressRoute) *route_v3.RouteConfiguration {
	allowedOrigins := []*matcher_v3.StringMatcher{
		{MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: "http://localhost:5173"}},
		{MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: "https://globule-ryzen.globular.io"}},
	}

	routeGroups := map[string][]*route_v3.Route{}
	domainsByKey := map[string][]string{}
	for _, r := range rs {
		domains, key := normalizeDomains(r.Domains)
		if _, ok := domainsByKey[key]; !ok {
			domainsByKey[key] = domains
		}
		match := &route_v3.RouteMatch{
			PathSpecifier: &route_v3.RouteMatch_Prefix{Prefix: r.Prefix},
		}
		if r.Authority != "" {
			match.Headers = []*route_v3.HeaderMatcher{
				{
					Name: ":authority",
					HeaderMatchSpecifier: &route_v3.HeaderMatcher_ExactMatch{
						ExactMatch: r.Authority,
					},
				},
			}
		}
		action := &route_v3.RouteAction{
			ClusterSpecifier: &route_v3.RouteAction_Cluster{Cluster: r.Cluster},
			Timeout:          durationpb.New(0),
		}
		if r.HostRewrite != "" {
			action.HostRewriteSpecifier = &route_v3.RouteAction_HostRewriteLiteral{
				HostRewriteLiteral: r.HostRewrite,
			}
		}
		routeGroups[key] = append(routeGroups[key], &route_v3.Route{
			Match: match,
			Action: &route_v3.Route_Route{
				Route: action,
			},
		})
	}

	vhosts := make([]*route_v3.VirtualHost, 0, len(routeGroups))
	for key, routes := range routeGroups {
		domains := domainsByKey[key]
		if len(domains) == 0 {
			domains = []string{"*"}
		}
		name := fmt.Sprintf("ingress_vhost_%s", hashKey(key))
		vhosts = append(vhosts, &route_v3.VirtualHost{
			Name:    name,
			Domains: domains,
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
		})
	}

	return &route_v3.RouteConfiguration{
		Name:         routeName,
		VirtualHosts: vhosts,
	}
}

func normalizeDomains(in []string) ([]string, string) {
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}
	for _, raw := range in {
		d := strings.TrimSpace(raw)
		if d == "" {
			continue
		}
		if _, ok := seen[d]; ok {
			continue
		}
		seen[d] = struct{}{}
		out = append(out, d)
	}
	sort.Strings(out)
	if len(out) == 0 {
		return nil, "*"
	}
	return out, strings.Join(out, ";")
}

func hashKey(key string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return fmt.Sprintf("%x", h.Sum32())
}
