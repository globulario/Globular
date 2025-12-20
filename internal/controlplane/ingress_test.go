package controlplane

import (
	"testing"
)

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
