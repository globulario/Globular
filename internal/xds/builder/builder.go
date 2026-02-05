package builder

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
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
	NodeID             string    `json:"node_id"`
	Version            string    `json:"version,omitempty"`
	Listener           Listener  `json:"listener"`
	Routes             []Route   `json:"routes"`
	Clusters           []Cluster `json:"clusters"`
	IngressHTTPPort    uint32    `json:"ingress_http_port"`
	EnableHTTPRedirect bool      `json:"enable_http_redirect"`
	GatewayPort        uint32    `json:"gateway_port"`
	EnableSDS          bool      `json:"enable_sds,omitempty"`  // Use SDS for TLS certificates
	SDSSecrets         []Secret  `json:"sds_secrets,omitempty"` // Secrets to include in snapshot
}

// Secret represents a TLS secret for SDS.
type Secret struct {
	Name     string `json:"name"`
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
	CAPath   string `json:"ca_path,omitempty"`
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

		var c *cluster_v3.Cluster
		if input.EnableSDS {
			// Use SDS for TLS certificates
			// Map file paths to secret names for now (in production, secret names would be explicit)
			caSecretName := ""
			if cluster.CAFile != "" {
				caSecretName = "internal-ca-bundle"
			}
			clientCertSecretName := ""
			if cluster.ServerCert != "" && cluster.KeyFile != "" {
				clientCertSecretName = "internal-client-cert"
			}

			c = controlplane.MakeClusterWithSDS(
				name,
				caSecretName,
				clientCertSecretName,
				strings.TrimSpace(cluster.SNI),
				toControlplaneEndpoints(cluster.Endpoints),
			)
		} else {
			// File-based TLS (legacy)
			c = controlplane.MakeCluster(
				name,
				strings.TrimSpace(cluster.ServerCert),
				strings.TrimSpace(cluster.KeyFile),
				strings.TrimSpace(cluster.CAFile),
				strings.TrimSpace(cluster.SNI),
				toControlplaneEndpoints(cluster.Endpoints),
			)
		}
		resources[resource_v3.ClusterType] = append(resources[resource_v3.ClusterType], c)
	}

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
		routeName := strings.TrimSpace(input.Listener.RouteName)
		if routeName == "" {
			routeName = fmt.Sprintf("ingress_routes_%d", controlplane.DefaultIngressPort(strings.TrimSpace(input.Listener.Host)))
		}
		resources[resource_v3.RouteType] = append(resources[resource_v3.RouteType], controlplane.MakeRoutes(routeName, ingressRoutes))

		host := strings.TrimSpace(input.Listener.Host)
		if host == "" {
			host = "0.0.0.0"
		}

		httpsPort := input.Listener.Port
		if httpsPort == 0 {
			httpsPort = controlplane.DefaultIngressPort(host)
		}
		httpPort := input.IngressHTTPPort
		if httpPort == 0 {
			httpPort = controlplane.DefaultIngressHTTPPort(host)
		}
		gatewayPort := input.GatewayPort

		listenerName := strings.TrimSpace(input.Listener.Name)
		if listenerName == "" {
			listenerName = fmt.Sprintf("ingress_listener_%d", httpsPort)
		}
		certFile := strings.TrimSpace(input.Listener.CertFile)
		keyFile := strings.TrimSpace(input.Listener.KeyFile)
		issuerFile := strings.TrimSpace(input.Listener.IssuerFile)
		tlsEnabled := tlsReady(certFile, keyFile)
		httpAllowed := httpPort > 0 && (gatewayPort == 0 || httpPort != gatewayPort)
		if issuerFile != "" && !fileExists(issuerFile) {
			issuerFile = ""
		}

		if tlsEnabled {
			var listener *listener_v3.Listener
			if input.EnableSDS {
				// Use SDS for TLS certificates
				serverCertSecretName := "internal-server-cert"
				caSecretName := ""
				if issuerFile != "" {
					caSecretName = "internal-ca-bundle"
				}

				listener = controlplane.MakeHTTPListenerWithSDS(
					host,
					httpsPort,
					listenerName,
					routeName,
					serverCertSecretName,
					caSecretName,
				)
			} else {
				// File-based TLS (legacy)
				listener = controlplane.MakeHTTPListener(
					host,
					httpsPort,
					listenerName,
					routeName,
					certFile,
					keyFile,
					issuerFile,
				)
			}

			resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], listener)
			if input.EnableHTTPRedirect && httpAllowed {
				redirectRouteName := fmt.Sprintf("%s_http_redirect_%d", routeName, httpPort)
				redirectRC, err := controlplane.MakeRedirectRoutes(redirectRouteName, httpsPort, true)
				if err != nil {
					return nil, err
				}
				redirectListenerName := fmt.Sprintf("%s_http_%d", listenerName, httpPort)
				redirectListener := controlplane.MakeHTTPListener(host, httpPort, redirectListenerName, redirectRouteName, "", "", "")
				resources[resource_v3.RouteType] = append(resources[resource_v3.RouteType], redirectRC)
				resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], redirectListener)
			}
		} else if httpAllowed {
			redirectListenerName := fmt.Sprintf("%s_http_%d", listenerName, httpPort)
			httpListener := controlplane.MakeHTTPListener(host, httpPort, redirectListenerName, routeName, "", "", "")
			resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], httpListener)
		}
	}

	// Add SDS secrets if enabled
	if input.EnableSDS && len(input.SDSSecrets) > 0 {
		for _, secret := range input.SDSSecrets {
			name := strings.TrimSpace(secret.Name)
			if name == "" {
				continue
			}

			certPath := strings.TrimSpace(secret.CertPath)
			keyPath := strings.TrimSpace(secret.KeyPath)
			caPath := strings.TrimSpace(secret.CAPath)

			// Build secret resource
			s, err := controlplane.MakeSecret(name, certPath, keyPath, caPath)
			if err != nil {
				// Log but don't fail the snapshot - secret may be optional
				fmt.Fprintf(os.Stderr, "warning: failed to build secret %s: %v\n", name, err)
				continue
			}

			resources[resource_v3.SecretType] = append(resources[resource_v3.SecretType], s)
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

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func tlsReady(certPath, keyPath string) bool {
	return fileExists(certPath) && fileExists(keyPath)
}
