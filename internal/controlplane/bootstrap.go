// controlplane/bootstrap.go
package controlplane

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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

// MarshalBootstrap builds the Envoy bootstrap YAML bytes without writing to disk.
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
		Address   string `yaml:"address"`
		PortValue int    `yaml:"port_value"`
	}
	type address struct {
		SocketAddress socketAddr `yaml:"socket_address"`
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
					"type": "STRICT_DNS",
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
				Address:   "0.0.0.0",
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

	return yaml.Marshal(doc)
}

// WriteBootstrap writes the generated Envoy bootstrap to disk.
func WriteBootstrap(path string, opt BootstrapOptions) error {
	b, err := MarshalBootstrap(opt)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
