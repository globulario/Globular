package watchers

import (
	"reflect"
	"testing"
)

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
