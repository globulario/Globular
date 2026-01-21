package watchers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cache_v3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/globulario/Globular/internal/xds/builder"
	"github.com/globulario/services/golang/config"
)

func writeConfig(t *testing.T, cfgDir, domain, protocol string) {
	t.Helper()
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	payload := map[string]any{
		"Domain":   domain,
		"Protocol": protocol,
	}
	data, _ := json.Marshal(payload)
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), data, 0o644); err != nil {
		t.Fatalf("write config.json: %v", err)
	}
}

func ensureTLSFiles(t *testing.T) {
	t.Helper()
	_, certPath, keyPath, caPath := config.CanonicalTLSPaths(config.GetRuntimeConfigDir())
	for _, p := range []string{certPath, keyPath, caPath} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("mkdir tls dir: %v", err)
		}
		if err := os.WriteFile(p, []byte("dummy"), 0o644); err != nil {
			t.Fatalf("write tls file: %v", err)
		}
	}
}

func containsDomain(domains []string, target string) bool {
	for _, d := range domains {
		trim := strings.TrimSpace(d)
		if trim == "*" {
			return true
		}
		if strings.EqualFold(trim, target) {
			return true
		}
	}
	return false
}

func buildSnapshotFromFallback(t *testing.T, cfg XDSConfig) *builder.Input {
	t.Helper()
	cfg.normalize()
	w := &Watcher{}
	spec := w.ingressFromFallback(cfg.Fallback)
	if spec == nil {
		t.Fatalf("ingress spec is nil")
	}
	w.applyIngressSettings(spec, &cfg)
	return &builder.Input{
		NodeID:             "node-test",
		Listener:           spec.Listener,
		Routes:             spec.Routes,
		Clusters:           spec.Clusters,
		IngressHTTPPort:    spec.HTTPPort,
		EnableHTTPRedirect: spec.EnableHTTPRedirect,
		GatewayPort:        spec.GatewayPort,
	}
}

func TestIngressSnapshotTracksDomainChange(t *testing.T) {
	t.Setenv("GLOBULAR_STATE_DIR", filepath.Join(t.TempDir(), "state"))
	t.Setenv("GLOBULAR_CONFIG_DIR", filepath.Join(t.TempDir(), "etc", "globular"))

	// Initial domain
	writeConfig(t, config.GetConfigDir(), "localhost", "http")
	ensureTLSFiles(t)
	// Reload config cache after write
	if _, err := config.GetLocalConfig(false); err != nil {
		t.Fatalf("load config: %v", err)
	}

	cfg := XDSConfig{
		Fallback: &FallbackConfig{
			Enabled: true,
			Ingress: &FallbackIngress{
				ListenerHost: "0.0.0.0",
				HTTPSPort:    8443,
				HTTPPort:     8080,
				TLS:          TLSPaths{},
				Routes: []FallbackRoute{
					{Prefix: "/", Cluster: "gateway_http"},
				},
			},
			Clusters: []FallbackCluster{
				{
					Name: "gateway_http",
					Endpoints: []FallbackEndpoint{
						{Host: "127.0.0.1", Port: 8080},
					},
				},
			},
		},
	}

	input := buildSnapshotFromFallback(t, cfg)
	snap, err := builder.BuildSnapshot(*input, "v1")
	if err != nil {
		t.Fatalf("build snapshot: %v", err)
	}
	if !snapshotHasDomain(snap, "localhost") {
		t.Fatalf("expected localhost domain in route config")
	}

	// Change domain + protocol and rebuild
	writeConfig(t, config.GetConfigDir(), "globular.io", "https")
	if _, err := config.GetLocalConfig(false); err != nil {
		t.Fatalf("reload config: %v", err)
	}
	input = buildSnapshotFromFallback(t, cfg)
	snap, err = builder.BuildSnapshot(*input, "v2")
	if err != nil {
		t.Fatalf("build snapshot v2: %v", err)
	}
	if !snapshotHasDomain(snap, "globular.io") {
		t.Fatalf("expected globular.io domain in route config after update")
	}
	if !listenerHasCanonicalTLS(snap) {
		t.Fatalf("listener TLS paths not canonical")
	}
}

func snapshotHasDomain(snap *cache_v3.Snapshot, domain string) bool {
	rc := snap.GetResources(resource_v3.RouteType)
	for _, res := range rc {
		if rcfg, ok := res.(*route_v3.RouteConfiguration); ok {
			for _, vh := range rcfg.VirtualHosts {
				if containsDomain(vh.Domains, domain) {
					return true
				}
			}
		}
	}
	return false
}

func listenerHasCanonicalTLS(snap *cache_v3.Snapshot) bool {
	ln := snap.GetResources(resource_v3.ListenerType)
	for _, res := range ln {
		if l, ok := res.(*listener_v3.Listener); ok {
			if strings.Contains(l.String(), "/config/tls/") {
				return true
			}
		}
	}
	return false
}
