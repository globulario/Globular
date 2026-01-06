package watchers

import (
	"testing"

	"github.com/globulario/Globular/internal/xds/builder"
)

func TestReadGatewayAddressFromGatewayListenVariants(t *testing.T) {
	cfg := map[string]any{
		"Gateway": map[string]any{
			"Listen": "0.0.0.0:8080",
		},
	}
	host, port := readGatewayAddressFrom(cfg)
	if port != 8080 {
		t.Fatalf("expected port 8080 from Gateway.Listen, got %d", port)
	}
	if host != "0.0.0.0" {
		t.Fatalf("expected host 0.0.0.0 from Gateway.Listen, got %q", host)
	}
}

func TestDefaultGatewayPortUsedWhenListenMissing(t *testing.T) {
	if got := defaultGatewayPort(0); got != 8080 {
		t.Fatalf("expected default gateway port 8080, got %d", got)
	}
	if got := defaultGatewayPort(1234); got != 1234 {
		t.Fatalf("expected provided port 1234 to be preserved, got %d", got)
	}
}

func TestReadGatewayAddressFromConfigIgnoresPortHTTP(t *testing.T) {
	cfg := map[string]any{
		"PortHTTP": 8181,
	}
	host, port := readGatewayAddressFromConfig(cfg)
	if port != 8080 {
		t.Fatalf("expected default gateway port 8080 despite PortHTTP=8181, got %d", port)
	}
	if host != "127.0.0.1" {
		t.Fatalf("expected default host 127.0.0.1, got %q", host)
	}
}

func TestNormalizeGatewayHTTPClusterOverridesLocalEndpoints(t *testing.T) {
	spec := &IngressSpec{
		Clusters: []builder.Cluster{
			{
				Name: "gateway_http",
				Endpoints: []builder.Endpoint{
					{Host: "", Port: 8181},
					{Host: "0.0.0.0", Port: 8181},
					{Host: "external.host", Port: 9090},
				},
			},
		},
	}
	normalizeGatewayHTTPCluster(spec, 8080)
	pts := spec.Clusters[0].Endpoints
	if pts[0].Port != 8080 {
		t.Fatalf("expected endpoint[0] port replaced with 8080, got %d", pts[0].Port)
	}
	if pts[1].Port != 8080 {
		t.Fatalf("expected endpoint[1] port replaced with 8080, got %d", pts[1].Port)
	}
	if pts[2].Port != 9090 {
		t.Fatalf("expected external endpoint port untouched (9090), got %d", pts[2].Port)
	}
}
