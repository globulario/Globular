package watchers

import (
	"reflect"
	"testing"

	"github.com/globulario/Globular/internal/xds/builder"
)

// xds_emits_routes_for_hosted_services: each declared HostedServices entry must
// produce a prefix route to the host service's cluster, so a co-hosted gRPC
// service (behavioral_memory on the ai_memory process) is reachable through the
// gateway instead of falling through to the HTML catch-all.
func TestHostedServiceRoutes(t *testing.T) {
	const svc = "ai_memory.AiMemoryService"
	const cluster = "ai_memory_AiMemoryService_cluster"

	got := hostedServiceRoutes(map[string]interface{}{
		"HostedServices": []interface{}{"behavioral_memory.BehavioralMemoryService", svc, ""},
	}, svc, cluster)

	want := []builder.Route{{Prefix: "/behavioral_memory.BehavioralMemoryService/", Cluster: cluster}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("hostedServiceRoutes = %+v, want %+v (self-name and empty must be skipped)", got, want)
	}

	// No HostedServices => no extra routes.
	if r := hostedServiceRoutes(map[string]interface{}{"Name": svc}, svc, cluster); len(r) != 0 {
		t.Fatalf("expected no routes when HostedServices absent, got %+v", r)
	}
}

// hostedServiceNames must handle the etcd-JSON shape ([]interface{}), the
// local-config shape ([]string), and absence — so a co-hosted gRPC service
// (e.g. behavioral_memory on the ai_memory process) gets a gateway route.
func TestHostedServiceNames(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]interface{}
		want []string
	}{
		{
			name: "etcd json []interface{}",
			in:   map[string]interface{}{"HostedServices": []interface{}{"behavioral_memory.BehavioralMemoryService", " ", ""}},
			want: []string{"behavioral_memory.BehavioralMemoryService"},
		},
		{
			name: "local config []string",
			in:   map[string]interface{}{"HostedServices": []string{"behavioral_memory.BehavioralMemoryService"}},
			want: []string{"behavioral_memory.BehavioralMemoryService"},
		},
		{name: "absent", in: map[string]interface{}{"Name": "ai_memory.AiMemoryService"}, want: nil},
		{name: "nil value", in: map[string]interface{}{"HostedServices": nil}, want: nil},
		{name: "nil map", in: nil, want: nil},
		{name: "wrong type", in: map[string]interface{}{"HostedServices": "not-a-list"}, want: nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := hostedServiceNames(tc.in)
			if len(got) == 0 && len(tc.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("hostedServiceNames(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
