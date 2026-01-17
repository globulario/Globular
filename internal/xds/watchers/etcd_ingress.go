package watchers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/globulario/Globular/internal/xds/builder"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	etcdBasePrefix     = "/globular/xds/v1"
	etcdIngressPrefix  = etcdBasePrefix + "/ingress"
	etcdRoutesPrefix   = etcdIngressPrefix + "/routes/"
	etcdClustersPrefix = etcdBasePrefix + "/clusters/"
)

// EtcdGetter is a minimal subset of the etcd client required by the ingress parser.
type EtcdGetter interface {
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
}

// IngressSpec describes the data we need to populate builder.Input when ingress is enabled.
type IngressSpec struct {
	Listener           builder.Listener
	Routes             []builder.Route
	Clusters           []builder.Cluster
	HTTPPort           uint32
	EnableHTTPRedirect bool
	GatewayPort        uint32
	RedirectConfigured bool
	FromEtcd           bool
}

func parseEtcdIngress(ctx context.Context, getter EtcdGetter) (*IngressSpec, error) {
	enabled, err := getStringValue(ctx, getter, etcdIngressPrefix+"/enabled")
	if err != nil {
		return nil, err
	}
	if strings.ToLower(strings.TrimSpace(enabled)) != "true" {
		return nil, nil
	}

	listenerHost, err := getStringValue(ctx, getter, etcdIngressPrefix+"/listener_host")
	if err != nil {
		return nil, err
	}
	httpsPort, err := parseUint32Value(ctx, getter, etcdIngressPrefix+"/https_port")
	if err != nil {
		return nil, err
	}
	httpPort, err := parseUint32Value(ctx, getter, etcdIngressPrefix+"/http_port")
	if err != nil {
		return nil, err
	}

	tlsCert, err := lookupTLSValue(ctx, getter, etcdIngressPrefix+"/tls/cert_chain_path", etcdIngressPrefix+"/tls/cert_file")
	if err != nil {
		return nil, err
	}
	tlsKey, err := lookupTLSValue(ctx, getter, etcdIngressPrefix+"/tls/private_key_path", etcdIngressPrefix+"/tls/key_file")
	if err != nil {
		return nil, err
	}
	tlsIssuer, err := lookupTLSValue(ctx, getter, etcdIngressPrefix+"/tls/ca_path", etcdIngressPrefix+"/tls/issuer_file")
	if err != nil {
		return nil, err
	}
	redirect, redirectSet, err := parseBoolValue(ctx, getter, etcdIngressPrefix+"/redirect_to_https")
	if err != nil {
		return nil, err
	}

	routes, err := collectRoutes(ctx, getter)
	if err != nil {
		return nil, err
	}
	if len(routes) == 0 {
		return nil, nil
	}

	clusters, err := collectClusters(ctx, getter)
	if err != nil {
		return nil, err
	}
	if len(clusters) == 0 {
		return nil, nil
	}

	spec := &IngressSpec{
		Listener: builder.Listener{
			Host:       listenerHost,
			Port:       httpsPort,
			RouteName:  defaultRouteName,
			CertFile:   tlsCert,
			KeyFile:    tlsKey,
			IssuerFile: tlsIssuer,
		},
		HTTPPort:           httpPort,
		EnableHTTPRedirect: redirect,
		RedirectConfigured: redirectSet,
		Routes:             routes,
		Clusters:           clusters,
		FromEtcd:           true,
	}
	return spec, nil
}

func getStringValue(ctx context.Context, getter EtcdGetter, key string) (string, error) {
	resp, err := getter.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", nil
	}
	return strings.TrimSpace(string(resp.Kvs[0].Value)), nil
}

func parseUint32Value(ctx context.Context, getter EtcdGetter, key string) (uint32, error) {
	s, err := getStringValue(ctx, getter, key)
	if err != nil {
		return 0, err
	}
	if s == "" {
		return 0, nil
	}
	v, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func lookupTLSValue(ctx context.Context, getter EtcdGetter, primary, fallback string) (string, error) {
	if val, err := getStringValue(ctx, getter, primary); err != nil {
		return "", err
	} else if val != "" {
		return val, nil
	}
	return getStringValue(ctx, getter, fallback)
}

func parseBoolValue(ctx context.Context, getter EtcdGetter, key string) (bool, bool, error) {
	val, err := getStringValue(ctx, getter, key)
	if err != nil {
		return false, false, err
	}
	if val == "" {
		return false, false, nil
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(val))
	if err != nil {
		return false, false, err
	}
	return parsed, true, nil
}

func collectRoutes(ctx context.Context, getter EtcdGetter) ([]builder.Route, error) {
	resp, err := getter.Get(ctx, etcdRoutesPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	routeFields := map[string]map[string]string{}
	for _, kv := range resp.Kvs {
		keyPath := strings.TrimPrefix(string(kv.Key), etcdRoutesPrefix)
		parts := strings.SplitN(keyPath, "/", 2)
		if len(parts) < 2 {
			continue
		}
		id := parts[0]
		field := parts[1]
		if id == "" || field == "" {
			continue
		}
		if _, ok := routeFields[id]; !ok {
			routeFields[id] = map[string]string{}
		}
		routeFields[id][field] = string(kv.Value)
	}

	routes := make([]builder.Route, 0, len(routeFields))
	for _, fields := range routeFields {
		prefix := strings.TrimSpace(fields["prefix"])
		cluster := strings.TrimSpace(fields["cluster"])
		if prefix == "" || cluster == "" {
			continue
		}
		r := builder.Route{
			Prefix:      prefix,
			Cluster:     cluster,
			Authority:   strings.TrimSpace(fields["authority"]),
			HostRewrite: strings.TrimSpace(fields["host_rewrite"]),
		}
		if domains := strings.TrimSpace(fields["domains"]); domains != "" {
			r.Domains = parseDomains(domains)
		}
		routes = append(routes, r)
	}
	return routes, nil
}

func collectClusters(ctx context.Context, getter EtcdGetter) ([]builder.Cluster, error) {
	resp, err := getter.Get(ctx, etcdClustersPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	clusterEndpoints := map[string][]builder.Endpoint{}
	for _, kv := range resp.Kvs {
		keyPath := strings.TrimPrefix(string(kv.Key), etcdClustersPrefix)
		parts := strings.Split(keyPath, "/")
		if len(parts) < 3 {
			continue
		}
		clusterName := parts[0]
		if parts[1] != "endpoints" || clusterName == "" {
			continue
		}
		var value etcdEndpoint
		if err := json.Unmarshal(kv.Value, &value); err != nil {
			continue
		}
		if value.Host == "" || value.Port == 0 {
			continue
		}
		ep := builder.Endpoint{
			Host:     value.Host,
			Port:     value.Port,
			Priority: value.Priority,
		}
		clusterEndpoints[clusterName] = append(clusterEndpoints[clusterName], ep)
	}
	clusters := make([]builder.Cluster, 0, len(clusterEndpoints))
	for name, endpoints := range clusterEndpoints {
		if len(endpoints) == 0 {
			continue
		}
		clusters = append(clusters, builder.Cluster{
			Name:      name,
			Endpoints: endpoints,
		})
	}
	return clusters, nil
}

func parseDomains(in string) []string {
	parts := strings.Split(in, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			if _, ok := seen[trimmed]; ok {
				continue
			}
			seen[trimmed] = struct{}{}
			out = append(out, trimmed)
		}
	}
	return out
}

type etcdEndpoint struct {
	Host     string `json:"host"`
	Port     uint32 `json:"port"`
	Priority uint32 `json:"priority,omitempty"`
}
