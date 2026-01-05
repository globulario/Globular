package controlplane

import (
	"fmt"
	"testing"

	route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
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
	if len(routes) != 2 {
		t.Fatalf("expected 2 route configs, got %d", len(routes))
	}
	if _, ok := routes["ingress_routes"]; !ok {
		t.Fatalf("missing ingress_routes config")
	}
	if _, ok := routes["ingress_routes_http_redirect"]; !ok {
		t.Fatalf("missing redirect config")
	}

	listeners := snap.GetResources(resource_v3.ListenerType)
	if len(listeners) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(listeners))
	}
	if _, ok := listeners["ingress_listener"]; !ok {
		t.Fatalf("missing ingress_listener")
	}
	if _, ok := listeners["ingress_listener_http"]; !ok {
		t.Fatalf("missing ingress_listener_http")
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

	rc := MakeRoutes("test", routes)
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
