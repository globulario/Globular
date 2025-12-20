package builder

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cache_v3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/globulario/Globular/internal/controlplane"
)

// Endpoint describes a single upstream endpoint for a cluster.
type Endpoint struct {
	Host     string `json:"host"`
	Port     uint32 `json:"port"`
	Priority uint32 `json:"priority,omitempty"`
}

// Cluster describes an Envoy cluster and its TLS settings.
type Cluster struct {
	Name       string     `json:"name"`
	Endpoints  []Endpoint `json:"endpoints"`
	ServerCert string     `json:"server_cert,omitempty"`
	KeyFile    string     `json:"key_file,omitempty"`
	CAFile     string     `json:"ca_file,omitempty"`
	SNI        string     `json:"sni,omitempty"`
}

// Listener defines the shared ingress listener.
type Listener struct {
	Name       string `json:"listener_name"`
	RouteName  string `json:"route_name"`
	Host       string `json:"host"`
	Port       uint32 `json:"port"`
	CertFile   string `json:"cert_file,omitempty"`
	KeyFile    string `json:"key_file,omitempty"`
	IssuerFile string `json:"issuer_file,omitempty"`
}

// Route maps a prefix to a cluster.
type Route struct {
	Prefix      string   `json:"prefix"`
	Cluster     string   `json:"cluster"`
	HostRewrite string   `json:"host_rewrite,omitempty"`
	Domains     []string `json:"domains,omitempty"`
	Authority   string   `json:"authority,omitempty"`
}

// Input is the data required to build an xDS snapshot.
type Input struct {
	NodeID   string    `json:"node_id"`
	Version  string    `json:"version,omitempty"`
	Nodes    []string  `json:"nodes,omitempty"` // unused but reserved for future use
	Listener Listener  `json:"listener"`
	Routes   []Route   `json:"routes"`
	Clusters []Cluster `json:"clusters"`
}

// BuildSnapshot returns a go-control-plane snapshot based on the configured input.
func BuildSnapshot(input Input, version string) (*cache_v3.Snapshot, error) {
	if strings.TrimSpace(input.NodeID) == "" {
		return nil, fmt.Errorf("node_id is required")
	}
	if version == "" {
		version = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	resources := make(map[string][]types.Resource)
	added := map[string]struct{}{}

	for _, cluster := range input.Clusters {
		name := strings.TrimSpace(cluster.Name)
		if name == "" || len(cluster.Endpoints) == 0 {
			continue
		}
		if _, ok := added[name]; ok {
			continue
		}
		added[name] = struct{}{}

		c := controlplane.MakeCluster(
			name,
			strings.TrimSpace(cluster.ServerCert),
			strings.TrimSpace(cluster.KeyFile),
			strings.TrimSpace(cluster.CAFile),
			strings.TrimSpace(cluster.SNI),
			toControlplaneEndpoints(cluster.Endpoints),
		)
		resources[resource_v3.ClusterType] = append(resources[resource_v3.ClusterType], c)
	}

	if input.Listener.Port > 0 && input.Listener.Name != "" && input.Listener.RouteName != "" {
		var ingressRoutes []controlplane.IngressRoute
		for _, route := range input.Routes {
			prefix := strings.TrimSpace(route.Prefix)
			if prefix == "" {
				continue
			}
			ingressRoutes = append(ingressRoutes, controlplane.IngressRoute{
				Prefix:      prefix,
				Cluster:     strings.TrimSpace(route.Cluster),
				HostRewrite: strings.TrimSpace(route.HostRewrite),
				Domains:     trimValues(route.Domains),
				Authority:   strings.TrimSpace(route.Authority),
			})
		}
		if len(ingressRoutes) > 0 {
			resources[resource_v3.RouteType] = append(resources[resource_v3.RouteType], controlplane.MakeRoutes(input.Listener.RouteName, ingressRoutes))

			host := strings.TrimSpace(input.Listener.Host)
			if host == "" {
				host = "0.0.0.0"
			}

			listener := controlplane.MakeHTTPListener(
				host,
				input.Listener.Port,
				input.Listener.Name,
				input.Listener.RouteName,
				strings.TrimSpace(input.Listener.CertFile),
				strings.TrimSpace(input.Listener.KeyFile),
				strings.TrimSpace(input.Listener.IssuerFile),
			)
			resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], listener)
		}
	}

	return cache_v3.NewSnapshot(version, resources)
}

func toControlplaneEndpoints(ends []Endpoint) []controlplane.EndPoint {
	out := make([]controlplane.EndPoint, 0, len(ends))
	for _, ep := range ends {
		host := strings.TrimSpace(ep.Host)
		if host == "" || ep.Port == 0 {
			continue
		}
		out = append(out, controlplane.EndPoint{
			Host:     host,
			Port:     ep.Port,
			Priority: ep.Priority,
		})
	}
	return out
}

func trimValues(in []string) []string {
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}
	for _, v := range in {
		s := strings.TrimSpace(v)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
