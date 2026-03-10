// internal/controlplane/ingress.go
package controlplane

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"
	"time"

	route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cors_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	coreConfig "github.com/globulario/Globular/internal/config"
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
	// Timeout sets the route timeout in seconds (0 = disabled for streaming)
	Timeout int
}

// MakeRoutes builds the shared ingress RouteConfiguration with per-route prefixes.
// IMPORTANT: Enables CORS (with credentials) for browser callers.
//
// corsPolicy (optional) provides the structured CORS configuration from the gateway.
// When nil, falls back to allowedOrigins (legacy []string) and hardcoded defaults.
// When allowedOrigins is also empty, a permissive regex ("https?://.*") is used.
func MakeRoutes(routeName string, rs []IngressRoute, allowedOrigins []string, corsPolicy ...*coreConfig.CorsPolicy) *route_v3.RouteConfiguration {
	// Build the Envoy CORS policy from the structured config or legacy fields.
	envoyCors := buildEnvoyCorsPolicy(allowedOrigins, firstOrNil(corsPolicy))

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
		}

		// Set timeout (0 = disabled for streaming)
		if r.Timeout == 0 {
			action.Timeout = durationpb.New(0)
		} else if r.Timeout > 0 {
			action.Timeout = durationpb.New(time.Duration(r.Timeout) * time.Second)
		}
		// If not set, use Envoy default

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
				"envoy.filters.http.cors": toAny(envoyCors),
			},
		})
	}

	return &route_v3.RouteConfiguration{
		Name:         routeName,
		VirtualHosts: vhosts,
	}
}

// buildEnvoyCorsPolicy converts the structured CorsPolicy (or legacy AllowedOrigins)
// into an Envoy cors_v3.CorsPolicy proto.
func buildEnvoyCorsPolicy(allowedOrigins []string, policy *coreConfig.CorsPolicy) *cors_v3.CorsPolicy {
	if policy != nil && !policy.Enabled {
		// CORS disabled — return a minimal policy with no origin matchers.
		return &cors_v3.CorsPolicy{}
	}

	// Origin matchers
	var originMatchers []*matcher_v3.StringMatcher

	if policy != nil {
		if policy.AllowAllOrigins {
			// Permissive: match any http(s) origin.
			// Envoy reflects the exact Origin back (not "*"), so credentials remain valid.
			originMatchers = []*matcher_v3.StringMatcher{
				{MatchPattern: &matcher_v3.StringMatcher_SafeRegex{
					SafeRegex: &matcher_v3.RegexMatcher{Regex: `https?://.*`},
				}},
			}
		} else {
			for _, o := range policy.AllowedOrigins {
				o = strings.TrimSpace(o)
				if o == "" {
					continue
				}
				originMatchers = append(originMatchers, &matcher_v3.StringMatcher{
					MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: o},
				})
			}
		}
	} else {
		// Legacy path: use allowedOrigins []string
		for _, o := range allowedOrigins {
			o = strings.TrimSpace(o)
			if o == "" {
				continue
			}
			originMatchers = append(originMatchers, &matcher_v3.StringMatcher{
				MatchPattern: &matcher_v3.StringMatcher_Exact{Exact: o},
			})
		}
	}

	if len(originMatchers) == 0 {
		// Default permissive fallback
		originMatchers = []*matcher_v3.StringMatcher{
			{MatchPattern: &matcher_v3.StringMatcher_SafeRegex{
				SafeRegex: &matcher_v3.RegexMatcher{Regex: `https?://.*`},
			}},
		}
	}

	// Methods, headers, exposed headers, max age
	allowMethods := "GET, PUT, DELETE, POST, OPTIONS"
	allowHeaders := "keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout,token,application,routing,authorization"
	exposeHeaders := "custom-header-1,grpc-status,grpc-message,grpc-status-details-bin"
	maxAge := "1728000"
	allowCredentials := true
	allowPrivateNetwork := false

	if policy != nil {
		if len(policy.AllowedMethods) > 0 {
			allowMethods = strings.Join(policy.AllowedMethods, ", ")
		}
		if len(policy.AllowedHeaders) > 0 {
			allowHeaders = strings.Join(policy.AllowedHeaders, ",")
		}
		if len(policy.ExposedHeaders) > 0 {
			exposeHeaders = strings.Join(policy.ExposedHeaders, ",")
		}
		if policy.MaxAgeSeconds > 0 {
			maxAge = strconv.Itoa(policy.MaxAgeSeconds)
		}
		allowCredentials = policy.AllowCredentials
		allowPrivateNetwork = policy.AllowPrivateNetwork
	}

	cp := &cors_v3.CorsPolicy{
		AllowOriginStringMatch: originMatchers,
		AllowCredentials:       wrapperspb.Bool(allowCredentials),
		AllowMethods:           allowMethods,
		AllowHeaders:           allowHeaders,
		ExposeHeaders:          exposeHeaders,
		MaxAge:                 maxAge,
	}

	// Private network access (Envoy 1.29+ supports this via allow_private_network_access)
	if allowPrivateNetwork {
		cp.AllowPrivateNetworkAccess = wrapperspb.Bool(true)
	}

	return cp
}

func firstOrNil(ps []*coreConfig.CorsPolicy) *coreConfig.CorsPolicy {
	if len(ps) > 0 {
		return ps[0]
	}
	return nil
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

// MakeRedirectRoutes builds a simple redirect configuration that sends all HTTP requests to HTTPS.
func MakeRedirectRoutes(routeName string, httpsPort uint32, permanent bool) (types.Resource, error) {
	name := strings.TrimSpace(routeName)
	if name == "" {
		name = "http_redirect"
	}

	redirect := &route_v3.RedirectAction{
		SchemeRewriteSpecifier: &route_v3.RedirectAction_HttpsRedirect{
			HttpsRedirect: true,
		},
	}
	if httpsPort > 0 {
		redirect.PortRedirect = httpsPort
	}
	if permanent {
		redirect.ResponseCode = route_v3.RedirectAction_PERMANENT_REDIRECT
	}

	return &route_v3.RouteConfiguration{
		Name: name,
		VirtualHosts: []*route_v3.VirtualHost{{
			Name:    "redirect_vhost",
			Domains: []string{"*"},
			Routes: []*route_v3.Route{{
				Match: &route_v3.RouteMatch{
					PathSpecifier: &route_v3.RouteMatch_Prefix{Prefix: "/"},
				},
				Action: &route_v3.Route_Redirect{
					Redirect: redirect,
				},
			}},
		}},
	}, nil
}
