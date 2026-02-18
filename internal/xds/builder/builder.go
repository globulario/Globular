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
	"github.com/globulario/Globular/internal/xds/secrets"
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
	Timeout     int      `json:"timeout,omitempty"` // 0 = disabled, >0 = seconds
}

// ExternalDomain represents an external FQDN with TLS certificate for SNI routing (PR3c).
type ExternalDomain struct {
	FQDN          string `json:"fqdn"`           // Fully-qualified domain name
	CertFile      string `json:"cert_file"`      // Path to fullchain.pem
	KeyFile       string `json:"key_file"`       // Path to privkey.pem
	TargetCluster string `json:"target_cluster"` // Backend cluster to route to (e.g., "gateway_http")
}

// Input is the data required to build an xDS snapshot.
type Input struct {
	NodeID             string           `json:"node_id"`
	Version            string           `json:"version,omitempty"`
	Listener           Listener         `json:"listener"`
	Routes             []Route          `json:"routes"`
	Clusters           []Cluster        `json:"clusters"`
	IngressHTTPPort    uint32           `json:"ingress_http_port"`
	EnableHTTPRedirect bool             `json:"enable_http_redirect"`
	GatewayPort        uint32           `json:"gateway_port"`
	EnableSDS          bool             `json:"enable_sds,omitempty"`       // Use SDS for TLS certificates
	SDSSecrets         []Secret         `json:"sds_secrets,omitempty"`      // Secrets to include in snapshot
	ExternalDomains    []ExternalDomain `json:"external_domains,omitempty"` // External domains for SNI routing (PR3c)
}

// Secret represents a TLS secret for SDS.
type Secret struct {
	Name     string `json:"name"`
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
	CAPath   string `json:"ca_path,omitempty"`
}

// buildExternalDomainVirtualHosts creates VirtualHost configurations for external domains (PR3c).
// Each domain gets a VirtualHost that routes all traffic (prefix "/") to the target cluster.
func buildExternalDomainVirtualHosts(domains []ExternalDomain) []controlplane.IngressRoute {
	var routes []controlplane.IngressRoute
	for _, domain := range domains {
		routes = append(routes, controlplane.IngressRoute{
			Prefix:  "/",
			Cluster: domain.TargetCluster,
			Domains: []string{domain.FQDN},
		})
	}
	return routes
}

// buildExternalDomainSecrets creates SDS secrets for external domain certificates (PR3c).
// Secret names follow the pattern: ext-cert/<fqdn>
func buildExternalDomainSecrets(domains []ExternalDomain) ([]Secret, error) {
	var secrets []Secret
	for _, domain := range domains {
		// Verify certificate files exist before creating secret
		if !fileExists(domain.CertFile) || !fileExists(domain.KeyFile) {
			return nil, fmt.Errorf("external domain %s: certificate files not found (cert=%s, key=%s)",
				domain.FQDN, domain.CertFile, domain.KeyFile)
		}

		secrets = append(secrets, Secret{
			Name:     fmt.Sprintf("ext-cert/%s", domain.FQDN),
			CertPath: domain.CertFile,
			KeyPath:  domain.KeyFile,
			// CAPath is optional - ACME certificates include the full chain in CertPath
		})
	}
	return secrets, nil
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

	// DEBUG: Log function entry
	fmt.Fprintf(os.Stderr, "[DEBUG] === BuildSnapshot ENTERED === version=%s, len(ExternalDomains)=%d, len(Routes)=%d\n",
		version, len(input.ExternalDomains), len(input.Routes))

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
			caSecretName := ""
			if cluster.CAFile != "" {
				caSecretName = secrets.InternalCABundle
			}
			clientCertSecretName := ""
			if cluster.ServerCert != "" && cluster.KeyFile != "" {
				clientCertSecretName = secrets.InternalClientCert
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

	// Add external domain VirtualHosts (PR3c)
	// These go first to ensure external FQDNs match before internal routes
	fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Before routes - len(ExternalDomains)=%d, len(Routes)=%d\n", len(input.ExternalDomains), len(input.Routes))
	if len(input.ExternalDomains) > 0 {
		extRoutes := buildExternalDomainVirtualHosts(input.ExternalDomains)
		fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Created %d external routes from %d domains\n", len(extRoutes), len(input.ExternalDomains))
		for i, d := range input.ExternalDomains {
			fmt.Fprintf(os.Stderr, "[DEBUG]   Domain %d: FQDN=%s, TargetCluster=%s\n", i, d.FQDN, d.TargetCluster)
		}
		ingressRoutes = append(ingressRoutes, extRoutes...)
	}

	// Add internal service routes
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
			Timeout:     route.Timeout,
		})
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Total ingressRoutes=%d\n", len(ingressRoutes))
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
				// Choose secret name based on certificate type (ACME vs internal)
				serverCertSecretName := secrets.InternalServerCert
				fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Listener creation - EnableSDS=true, len(ExternalDomains)=%d\n", len(input.ExternalDomains))
				// When external domains are configured, use public ingress cert for default fallback
				if len(input.ExternalDomains) > 0 || strings.Contains(certFile, "fullchain.pem") {
					// ACME certificate (Let's Encrypt) for public ingress
					serverCertSecretName = secrets.PublicIngressCert
				}
				caSecretName := ""
				if issuerFile != "" {
					caSecretName = secrets.InternalCABundle
				}

				// Build SNI filter chains for external domains
				var extraChains []*listener_v3.FilterChain
				if len(input.ExternalDomains) > 0 {
					for _, d := range input.ExternalDomains {
						// Include both apex domain and wildcard subdomain in SNI match
						// This allows matching both "example.com" and "*.example.com"
						serverNames := []string{
							d.FQDN,                      // apex domain (e.g., "globular.app")
							fmt.Sprintf("*.%s", d.FQDN), // wildcard (e.g., "*.globular.app")
						}
						fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Creating SNI filter chain for domain=%s with serverNames=%v\n", d.FQDN, serverNames)
						fc, err := controlplane.MakeSNIHTTPFilterChainWithSDS(
							routeName,
							serverNames,
							fmt.Sprintf("ext-cert/%s", d.FQDN),
							caSecretName,
						)
						if err != nil {
							return nil, fmt.Errorf("failed to create SNI filter chain for %s: %w", d.FQDN, err)
						}
						extraChains = append(extraChains, fc)
					}

					// Use listener with SNI chains + default fallback
					listener = controlplane.MakeHTTPListenerWithSDSFilterChains(
						host,
						httpsPort,
						listenerName,
						routeName,
						serverCertSecretName,
						caSecretName,
						extraChains,
					)
				} else {
					// No external domains - use simple listener with default cert
					listener = controlplane.MakeHTTPListenerWithSDS(
						host,
						httpsPort,
						listenerName,
						routeName,
						serverCertSecretName,
						caSecretName,
					)
				}
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
	fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: EnableSDS=%v, SDSSecrets count=%d\n", input.EnableSDS, len(input.SDSSecrets))
	if input.EnableSDS && len(input.SDSSecrets) > 0 {
		fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Processing %d SDS secrets\n", len(input.SDSSecrets))
		for i, secret := range input.SDSSecrets {
			name := strings.TrimSpace(secret.Name)
			if name == "" {
				fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Skipping secret %d (empty name)\n", i)
				continue
			}

			certPath := strings.TrimSpace(secret.CertPath)
			keyPath := strings.TrimSpace(secret.KeyPath)
			caPath := strings.TrimSpace(secret.CAPath)

			fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Creating secret %d: name=%s, cert=%s, key=%s, ca=%s\n",
				i, name, certPath, keyPath, caPath)

			// Build secret resource
			s, err := controlplane.MakeSecret(name, certPath, keyPath, caPath)
			if err != nil {
				// Log but don't fail the snapshot - secret may be optional
				fmt.Fprintf(os.Stderr, "[ERROR] BuildSnapshot: Failed to build secret %s: %v\n", name, err)
				continue
			}

			fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Successfully created secret %s\n", name)
			resources[resource_v3.SecretType] = append(resources[resource_v3.SecretType], s)
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Added %d secrets to snapshot\n", len(resources[resource_v3.SecretType]))
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: SDS disabled or no secrets to add\n")
	}

	// Debug: Log all resource types before creating snapshot
	fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Creating snapshot with resources:\n")
	for resType, resList := range resources {
		fmt.Fprintf(os.Stderr, "[DEBUG]   - %s: %d resources\n", resType, len(resList))
	}

	snapshot, err := cache_v3.NewSnapshot(version, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Snapshot created, version=%s\n", version)

	// Debug: Verify secrets are in the snapshot
	secretResources := snapshot.GetResources(resource_v3.SecretType)
	fmt.Fprintf(os.Stderr, "[DEBUG] BuildSnapshot: Snapshot contains %d secrets after NewSnapshot\n", len(secretResources))

	return snapshot, nil
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
