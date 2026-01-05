package controlplane

import (
	"encoding/json"
	"testing"
)

func TestMarshalBootstrapStaticCluster(t *testing.T) {
	opt := BootstrapOptions{
		XDSHost: "127.0.0.1",
		XDSPort: 18000,
	}
	data, err := MarshalBootstrap(opt)
	if err != nil {
		t.Fatalf("MarshalBootstrap: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}

	static, ok := doc["static_resources"].(map[string]any)
	if !ok {
		t.Fatalf("static_resources missing")
	}
	clusters, ok := static["clusters"].([]any)
	if !ok {
		t.Fatalf("clusters missing")
	}
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(clusters))
	}
	cluster, ok := clusters[0].(map[string]any)
	if !ok {
		t.Fatalf("cluster entry malformed")
	}

	if got := cluster["type"]; got != "STATIC" {
		t.Fatalf("expected STATIC cluster type, got %v", got)
	}
	if timeout, ok := cluster["connect_timeout"].(string); !ok || timeout == "" {
		t.Fatalf("connect_timeout missing or invalid: %v", cluster["connect_timeout"])
	}
	loadAssignment, ok := cluster["load_assignment"].(map[string]any)
	if !ok {
		t.Fatalf("load_assignment missing")
	}
	if loadAssignment["cluster_name"] != "xds_cluster" {
		t.Fatalf("unexpected cluster_name: %v", loadAssignment["cluster_name"])
	}
	endpoints, ok := loadAssignment["endpoints"].([]any)
	if !ok || len(endpoints) == 0 {
		t.Fatalf("endpoints missing")
	}
	first, ok := endpoints[0].(map[string]any)
	if !ok {
		t.Fatalf("endpoint entry malformed")
	}
	lb, ok := first["lb_endpoints"].([]any)
	if !ok || len(lb) == 0 {
		t.Fatalf("lb_endpoints missing")
	}
	endpoint, ok := lb[0].(map[string]any)
	if !ok {
		t.Fatalf("lb endpoint malformed")
	}
	endpointObj, ok := endpoint["endpoint"].(map[string]any)
	if !ok {
		t.Fatalf("endpoint.endpoint missing or invalid")
	}
	addrObj, ok := endpointObj["address"].(map[string]any)
	if !ok {
		t.Fatalf("endpoint.address missing or invalid")
	}
	socket, ok := addrObj["socket_address"].(map[string]any)
	if !ok {
		t.Fatalf("socket_address missing or invalid")
	}
	if socket["address"] != "127.0.0.1" {
		t.Fatalf("unexpected address: %v", socket["address"])
	}
	if socket["port_value"] != float64(18000) {
		t.Fatalf("unexpected port: %v", socket["port_value"])
	}

	typed, ok := cluster["typed_extension_protocol_options"].(map[string]any)
	if !ok {
		t.Fatalf("typed_extension_protocol_options missing")
	}
	hp, ok := typed["envoy.extensions.upstreams.http.v3.HttpProtocolOptions"].(map[string]any)
	if !ok {
		t.Fatalf("HttpProtocolOptions missing or invalid")
	}
	explicit, ok := hp["explicit_http_config"].(map[string]any)
	if !ok {
		t.Fatalf("explicit_http_config missing or invalid")
	}
	if _, ok := explicit["http2_protocol_options"]; !ok {
		t.Fatalf("http2_protocol_options missing")
	}
}
