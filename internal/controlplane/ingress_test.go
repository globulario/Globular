package controlplane

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

func TestMakeRedirectRoutes(t *testing.T) {
	res, err := MakeRedirectRoutes("redirect", 8443, true)
	if err != nil {
		t.Fatalf("MakeRedirectRoutes: %v", err)
	}
	rc, ok := res.(*route_v3.RouteConfiguration)
	if !ok {
		t.Fatalf("expected RouteConfiguration, got %T", res)
	}
	if rc.Name != "redirect" {
		t.Fatalf("unexpected route name: %q", rc.Name)
	}
	if len(rc.VirtualHosts) != 1 {
		t.Fatalf("expected 1 virtual host, got %d", len(rc.VirtualHosts))
	}
	vh := rc.VirtualHosts[0]
	if len(vh.Domains) != 1 || vh.Domains[0] != "*" {
		t.Fatalf("unexpected domains: %v", vh.Domains)
	}
	if len(vh.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(vh.Routes))
	}
	redirect := vh.Routes[0].GetRedirect()
	if redirect == nil {
		t.Fatalf("route action missing redirect")
	}
	if !redirect.GetHttpsRedirect() {
		t.Fatalf("expected https redirect")
	}
	if redirect.GetPortRedirect() != 8443 {
		t.Fatalf("unexpected port redirect: %d", redirect.GetPortRedirect())
	}
	if redirect.GetResponseCode() != route_v3.RedirectAction_PERMANENT_REDIRECT {
		t.Fatalf("unexpected response code: %v", redirect.GetResponseCode())
	}
}

func TestAddSnapshotIngressRedirect(t *testing.T) {
	id := fmt.Sprintf("test-ingress-redirect-%s", t.Name())
	defer RemoveSnapshot(id)

	cert := filepath.Join(t.TempDir(), "cert.pem")
	key := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(cert, []byte("cert"), 0o644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(key, []byte("key"), 0o644); err != nil {
		t.Fatalf("write key: %v", err)
	}

	err := AddSnapshot(id, "v-test", []Snapshot{
		{
			ListenerName:       "ingress_listener",
			RouteName:          "ingress_routes",
			ListenerHost:       "0.0.0.0",
			ListenerPort:       0,
			IngressRoutes:      []IngressRoute{{Prefix: "/", Cluster: "files"}},
			CertFilePath:       cert,
			KeyFilePath:        key,
			HTTPPort:           80,
			EnableHTTPRedirect: true,
		},
	})
	if err != nil {
		t.Fatalf("AddSnapshot: %v", err)
	}

	snap, err := GetSnapshot(id)
	if err != nil {
		t.Fatalf("GetSnapshot: %v", err)
	}

	routes := snap.GetResources(resource_v3.RouteType)
	if len(routes) != 2 {
		t.Fatalf("expected 2 route configs, got %d", len(routes))
	}
	if _, ok := routes["ingress_routes"]; !ok {
		t.Fatalf("missing ingress_routes config")
	}
	if _, ok := routes["ingress_routes_http_redirect_80"]; !ok {
		t.Fatalf("missing redirect config")
	}

	listeners := snap.GetResources(resource_v3.ListenerType)
	if len(listeners) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(listeners))
	}
	if _, ok := listeners["ingress_listener"]; !ok {
		t.Fatalf("missing ingress_listener")
	}
	if _, ok := listeners["ingress_listener_http_80"]; !ok {
		t.Fatalf("missing ingress_listener_http_80")
	}
}

func TestAddSnapshotIngressHTTPOnly(t *testing.T) {
	id := fmt.Sprintf("test-ingress-http-only-%s", t.Name())
	defer RemoveSnapshot(id)

	err := AddSnapshot(id, "v-test", []Snapshot{
		{
			ListenerName:  "ingress_listener",
			RouteName:     "ingress_routes",
			ListenerHost:  "0.0.0.0",
			ListenerPort:  0,
			IngressRoutes: []IngressRoute{{Prefix: "/", Cluster: "files"}},
		},
	})
	if err != nil {
		t.Fatalf("AddSnapshot: %v", err)
	}

	snap, err := GetSnapshot(id)
	if err != nil {
		t.Fatalf("GetSnapshot: %v", err)
	}

	routes := snap.GetResources(resource_v3.RouteType)
	if len(routes) != 1 {
		t.Fatalf("expected 1 route config, got %d", len(routes))
	}
	if _, ok := routes["ingress_routes"]; !ok {
		t.Fatalf("missing ingress_routes config")
	}

	listeners := snap.GetResources(resource_v3.ListenerType)
	if len(listeners) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(listeners))
	}
	if _, ok := listeners["ingress_listener_http_80"]; !ok {
		t.Fatalf("missing ingress_listener_http_80")
	}
}

func TestAddSnapshotIngressTLSMissing(t *testing.T) {
	id := fmt.Sprintf("test-ingress-tls-missing-%s", t.Name())
	defer RemoveSnapshot(id)

	err := AddSnapshot(id, "v-test", []Snapshot{
		{
			ListenerName:       "ingress_listener",
			RouteName:          "ingress_routes",
			ListenerHost:       "0.0.0.0",
			ListenerPort:       443,
			HTTPPort:           80,
			EnableHTTPRedirect: true,
			IngressRoutes:      []IngressRoute{{Prefix: "/", Cluster: "files"}},
		},
	})
	if err != nil {
		t.Fatalf("AddSnapshot: %v", err)
	}

	snap, err := GetSnapshot(id)
	if err != nil {
		t.Fatalf("GetSnapshot: %v", err)
	}

	routes := snap.GetResources(resource_v3.RouteType)
	if len(routes) != 1 {
		t.Fatalf("expected 1 route config (no redirect), got %d", len(routes))
	}
	if _, ok := routes["ingress_routes"]; !ok {
		t.Fatalf("missing ingress_routes config")
	}

	listeners := snap.GetResources(resource_v3.ListenerType)
	if len(listeners) != 1 {
		t.Fatalf("expected 1 listener (HTTP only), got %d", len(listeners))
	}
	if _, ok := listeners["ingress_listener_http_80"]; !ok {
		t.Fatalf("missing ingress_listener_http_80")
	}
	if _, ok := listeners["ingress_listener"]; ok {
		t.Fatalf("TLS listener should not exist without certs")
	}
}

func TestDefaultIngressHTTPPortAlways80(t *testing.T) {
	if got := defaultIngressHTTPPort("localhost"); got != 80 {
		t.Fatalf("expected default HTTP port 80, got %d", got)
	}
	if got := defaultIngressHTTPPort("0.0.0.0"); got != 80 {
		t.Fatalf("expected default HTTP port 80, got %d", got)
	}
}

func TestAddSnapshotIngressCollisionDisablesRedirect(t *testing.T) {
	id := fmt.Sprintf("test-ingress-collision-redirect-%s", t.Name())
	defer RemoveSnapshot(id)

	cert := filepath.Join(t.TempDir(), "cert.pem")
	key := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(cert, []byte("cert"), 0o644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(key, []byte("key"), 0o644); err != nil {
		t.Fatalf("write key: %v", err)
	}

	err := AddSnapshot(id, "v-test", []Snapshot{
		{
			ListenerName:       "ingress_listener",
			RouteName:          "ingress_routes",
			ListenerHost:       "0.0.0.0",
			ListenerPort:       443,
			HTTPPort:           8080,
			EnableHTTPRedirect: true,
			IngressRoutes:      []IngressRoute{{Prefix: "/", Cluster: "files"}},
			GatewayPort:        8080,
			CertFilePath:       cert,
			KeyFilePath:        key,
			IssuerFilePath:     "/etc/globular/certs/ca.pem",
		},
	})
	if err != nil {
		t.Fatalf("AddSnapshot: %v", err)
	}

	snap, err := GetSnapshot(id)
	if err != nil {
		t.Fatalf("GetSnapshot: %v", err)
	}

	listeners := snap.GetResources(resource_v3.ListenerType)
	if _, ok := listeners["ingress_listener"]; !ok {
		t.Fatalf("missing TLS listener")
	}
	if _, ok := listeners["ingress_listener_http_8080"]; ok {
		t.Fatalf("unexpected collision listener on 8080")
	}
	if len(listeners) != 1 {
		t.Fatalf("expected only TLS listener, got %d", len(listeners))
	}
}

func TestAddSnapshotIngressMissingGatewayPortStillSkips8080(t *testing.T) {
	id := fmt.Sprintf("test-ingress-collision-missing-gw-%s", t.Name())
	defer RemoveSnapshot(id)

	cert := filepath.Join(t.TempDir(), "cert.pem")
	key := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(cert, []byte("cert"), 0o644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(key, []byte("key"), 0o644); err != nil {
		t.Fatalf("write key: %v", err)
	}

	err := AddSnapshot(id, "v-test", []Snapshot{
		{
			ListenerName:       "ingress_listener",
			RouteName:          "ingress_routes",
			ListenerHost:       "0.0.0.0",
			ListenerPort:       443,
			HTTPPort:           8080,
			EnableHTTPRedirect: true,
			IngressRoutes:      []IngressRoute{{Prefix: "/", Cluster: "files"}},
			CertFilePath:       cert,
			KeyFilePath:        key,
		},
	})
	if err != nil {
		t.Fatalf("AddSnapshot: %v", err)
	}

	snap, err := GetSnapshot(id)
	if err != nil {
		t.Fatalf("GetSnapshot: %v", err)
	}

	listeners := snap.GetResources(resource_v3.ListenerType)
	if _, ok := listeners["ingress_listener_http_8080"]; ok {
		t.Fatalf("redirect listener should not bind 8080 even if gateway_port missing")
	}
	if _, ok := listeners["ingress_listener"]; !ok {
		t.Fatalf("missing TLS listener")
	}
}

func TestAddSnapshotIngressIssuerMissing(t *testing.T) {
	id := fmt.Sprintf("test-ingress-issuer-%s", t.Name())
	defer RemoveSnapshot(id)

	cert := filepath.Join(t.TempDir(), "cert.pem")
	key := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(cert, []byte("cert"), 0o644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(key, []byte("key"), 0o644); err != nil {
		t.Fatalf("write key: %v", err)
	}

	err := AddSnapshot(id, "v-test", []Snapshot{
		{
			ListenerName:       "ingress_listener",
			RouteName:          "ingress_routes",
			ListenerHost:       "0.0.0.0",
			ListenerPort:       443,
			EnableHTTPRedirect: true,
			IngressRoutes:      []IngressRoute{{Prefix: "/", Cluster: "files"}},
			CertFilePath:       cert,
			KeyFilePath:        key,
			IssuerFilePath:     "/etc/globular/certs/missing-ca.pem",
		},
	})
	if err != nil {
		t.Fatalf("AddSnapshot: %v", err)
	}

	snap, err := GetSnapshot(id)
	if err != nil {
		t.Fatalf("GetSnapshot: %v", err)
	}

	listeners := snap.GetResources(resource_v3.ListenerType)
	lnRes, ok := listeners["ingress_listener"]
	if !ok {
		t.Fatalf("missing TLS listener")
	}
	ln, ok := lnRes.(*listener_v3.Listener)
	if !ok {
		t.Fatalf("listener is not Listener: %T", lnRes)
	}
	if len(ln.FilterChains) == 0 {
		t.Fatalf("listener lacks filter chains")
	}
	ts := ln.FilterChains[0].TransportSocket
	if ts == nil {
		t.Fatalf("listener missing transport socket")
	}
	ctx := &tls_v3.DownstreamTlsContext{}
	if err := ts.GetTypedConfig().UnmarshalTo(ctx); err != nil {
		t.Fatalf("unmarshal TLS context: %v", err)
	}
	if ctx.GetCommonTlsContext().GetValidationContextType() != nil {
		t.Fatalf("expected no validation context when CA is missing")
	}
}

func TestMakeRoutesGroupsByDomains(t *testing.T) {
	routes := []IngressRoute{
		{
			Prefix:  "/files",
			Cluster: "files",
			Domains: []string{"files.globular.io"},
		},
		{
			Prefix:  "/media",
			Cluster: "media",
			Domains: []string{"media.globular.io", "static.globular.info"},
		},
		{
			Prefix:  "/api/default",
			Cluster: "api",
		},
	}

	rc := MakeRoutes("test", routes, nil)
	if got, want := len(rc.VirtualHosts), 3; got != want {
		t.Fatalf("expected %d virtual hosts, got %d", want, got)
	}

	expected := map[string][]string{
		"files.globular.io":                      {"files.globular.io"},
		"media.globular.io;static.globular.info": {"media.globular.io", "static.globular.info"},
		"*":                                      {"*"},
	}

	found := make(map[string]struct{})
	for _, vh := range rc.VirtualHosts {
		domains, key := normalizeDomains(vh.Domains)
		if key == "" {
			key = "*"
		}
		expectedDomains, ok := expected[key]
		if !ok {
			t.Fatalf("unexpected virtual host with key %q and domains %v", key, vh.Domains)
		}
		found[key] = struct{}{}

		if len(domains) != len(expectedDomains) {
			t.Fatalf("domains mismatch for key %q: got %v, want %v", key, domains, expectedDomains)
		}
		for i := range domains {
			if domains[i] != expectedDomains[i] {
				t.Fatalf("domains mismatch for key %q: got %v, want %v", key, domains, expectedDomains)
			}
		}
	}

	for key := range expected {
		if _, ok := found[key]; !ok {
			t.Fatalf("expected virtual host for key %q not found", key)
		}
	}
}
