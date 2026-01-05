// controlplane/bootstrap.go
package controlplane

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type BootstrapOptions struct {
	NodeID    string
	Cluster   string
	XDSHost   string
	XDSPort   int
	AdminPort int

	// set a global cap on active downstream connections (0 = omit)
	MaxActiveDownstreamConns uint64
}

// MarshalBootstrap builds the Envoy bootstrap JSON bytes without writing to disk.
func MarshalBootstrap(opt BootstrapOptions) ([]byte, error) {
	if opt.NodeID == "" {
		opt.NodeID = "globular-xds"
	}
	if opt.Cluster == "" {
		opt.Cluster = "globular-cluster"
	}
	if opt.XDSHost == "" {
		opt.XDSHost = "127.0.0.1"
	}
	if opt.XDSPort == 0 {
		opt.XDSPort = 18000
	}
	if opt.AdminPort == 0 {
		opt.AdminPort = 9901
	}

	type socketAddr struct {
		Address   string `json:"address"`
		PortValue int    `json:"port_value"`
	}
	type address struct {
		SocketAddress socketAddr `json:"socket_address"`
	}

	doc := map[string]any{
		"node": map[string]any{
			"cluster": opt.Cluster,
			"id":      opt.NodeID,
		},

		"dynamic_resources": map[string]any{
			"ads_config": map[string]any{
				"api_type":              "GRPC",
				"transport_api_version": "V3",
				"grpc_services": []any{
					map[string]any{
						"envoy_grpc": map[string]any{
							"cluster_name": "xds_cluster",
						},
					},
				},
			},
			"cds_config": map[string]any{
				"resource_api_version": "V3",
				"ads":                  map[string]any{},
			},
			"lds_config": map[string]any{
				"resource_api_version": "V3",
				"ads":                  map[string]any{},
			},
		},

		"static_resources": map[string]any{
			"clusters": []any{
				map[string]any{
					"type":            clusterTypeForHost(opt.XDSHost),
					"connect_timeout": "1s",
					"typed_extension_protocol_options": map[string]any{
						"envoy.extensions.upstreams.http.v3.HttpProtocolOptions": map[string]any{
							"@type": "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions",
							"explicit_http_config": map[string]any{
								"http2_protocol_options": map[string]any{},
							},
						},
					},
					"name": "xds_cluster",
					"load_assignment": map[string]any{
						"cluster_name": "xds_cluster",
						"endpoints": []any{
							map[string]any{
								"lb_endpoints": []any{
									map[string]any{
										"endpoint": map[string]any{
											"address": address{
												SocketAddress: socketAddr{
													Address:   opt.XDSHost,
													PortValue: opt.XDSPort,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"admin": map[string]any{
			"address": address{SocketAddress: socketAddr{
				Address:   "127.0.0.1",
				PortValue: opt.AdminPort,
			}},
			"access_log": []any{
				map[string]any{
					"name": "envoy.access_loggers.stdout",
					"typed_config": map[string]any{
						"@type": "type.googleapis.com/envoy.extensions.access_loggers.stream.v3.StdoutAccessLog",
					},
				},
			},
		},
	}

	if opt.MaxActiveDownstreamConns > 0 {
		doc["overload_manager"] = map[string]any{
			"resource_monitors": []any{
				map[string]any{
					"name": "envoy.resource_monitors.global_downstream_max_connections",
					"typed_config": map[string]any{
						"@type":                             "type.googleapis.com/envoy.extensions.resource_monitors.downstream_connections.v3.DownstreamConnectionsConfig",
						"max_active_downstream_connections": opt.MaxActiveDownstreamConns,
					},
				},
			},
		}
	}

	return json.Marshal(doc)
}

// WriteBootstrap writes the generated Envoy bootstrap to disk.
func WriteBootstrap(path string, opt BootstrapOptions) error {
	data, err := MarshalBootstrap(opt)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	_ = f.Sync()
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func clusterTypeForHost(host string) string {
	if isIP(host) {
		return "STATIC"
	}
	return "STRICT_DNS"
}

func isIP(host string) bool {
	h := strings.TrimSpace(host)
	if h == "" {
		return false
	}
	if parsedHost, _, err := net.SplitHostPort(h); err == nil {
		h = parsedHost
	}
	return net.ParseIP(h) != nil
}
