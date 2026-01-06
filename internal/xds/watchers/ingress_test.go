package watchers

import (
	"context"
	"strings"
	"testing"

	cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/globulario/Globular/internal/xds/builder"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type fakeEtcdGetter struct {
	responses map[string][]*mvccpb.KeyValue
}

func (f *fakeEtcdGetter) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	if kvs, ok := f.responses[key]; ok {
		return &clientv3.GetResponse{Count: int64(len(kvs)), Kvs: kvs}, nil
	}
	var agg []*mvccpb.KeyValue
	for storedKey, kvs := range f.responses {
		if strings.HasPrefix(storedKey, key) {
			agg = append(agg, kvs...)
		}
	}
	if len(agg) > 0 {
		return &clientv3.GetResponse{Count: int64(len(agg)), Kvs: agg}, nil
	}
	return &clientv3.GetResponse{}, nil
}

func kv(key, value string) *mvccpb.KeyValue {
	return &mvccpb.KeyValue{
		Key:   []byte(key),
		Value: []byte(value),
	}
}

func TestParseEtcdIngress(t *testing.T) {
	getter := &fakeEtcdGetter{
		responses: map[string][]*mvccpb.KeyValue{
			etcdIngressPrefix + "/enabled":           {kv(etcdIngressPrefix+"/enabled", "true")},
			etcdIngressPrefix + "/listener_host":     {kv(etcdIngressPrefix+"/listener_host", "0.0.0.0")},
			etcdIngressPrefix + "/https_port":        {kv(etcdIngressPrefix+"/https_port", "8443")},
			etcdIngressPrefix + "/http_port":         {kv(etcdIngressPrefix+"/http_port", "80")},
			etcdIngressPrefix + "/tls/cert_file":     {kv(etcdIngressPrefix+"/tls/cert_file", "/etc/foo.crt")},
			etcdIngressPrefix + "/tls/key_file":      {kv(etcdIngressPrefix+"/tls/key_file", "/etc/foo.key")},
			etcdIngressPrefix + "/tls/issuer_file":   {kv(etcdIngressPrefix+"/tls/issuer_file", "/etc/ca.crt")},
			etcdIngressPrefix + "/redirect_to_https": {kv(etcdIngressPrefix+"/redirect_to_https", "true")},
			etcdRoutesPrefix + "route1/prefix":       {kv(etcdRoutesPrefix+"route1/prefix", "/")},
			etcdRoutesPrefix + "route1/cluster":      {kv(etcdRoutesPrefix+"route1/cluster", "gateway_http")},
			etcdRoutesPrefix + "route1/domains":      {kv(etcdRoutesPrefix+"route1/domains", "example.com,foo.example.com")},
			etcdRoutesPrefix + "route1/authority":    {kv(etcdRoutesPrefix+"route1/authority", "example.com")},
			etcdRoutesPrefix + "route1/host_rewrite": {kv(etcdRoutesPrefix+"route1/host_rewrite", "internal.host")},
			etcdClustersPrefix + "gateway_http/endpoints/node1": {
				kv(etcdClustersPrefix+"gateway_http/endpoints/node1", `{"host":"10.0.0.10","port":8080}`),
			},
			etcdClustersPrefix + "gateway_http/endpoints/node2": {
				kv(etcdClustersPrefix+"gateway_http/endpoints/node2", `{"host":"10.0.0.11","port":8080}`),
			},
		},
	}

	spec, err := parseEtcdIngress(context.Background(), getter)
	if err != nil {
		t.Fatalf("parseEtcdIngress: %v", err)
	}
	if spec == nil {
		t.Fatalf("expected ingress spec")
	}
	if spec.Listener.Port != 8443 {
		t.Fatalf("expected https port 8443, got %d", spec.Listener.Port)
	}
	if len(spec.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(spec.Routes))
	}
	if spec.Routes[0].Domains == nil || len(spec.Routes[0].Domains) != 2 {
		t.Fatalf("unexpected domains %v", spec.Routes[0].Domains)
	}
	if len(spec.Clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(spec.Clusters))
	}
	if spec.HTTPPort != 80 {
		t.Fatalf("expected http port 80, got %d", spec.HTTPPort)
	}
	if !spec.EnableHTTPRedirect {
		t.Fatalf("expected redirect enabled")
	}
	if !spec.RedirectConfigured {
		t.Fatalf("expected redirect configuration to be marked as set")
	}
}

func TestIngressFromFallback(t *testing.T) {
	w := &Watcher{}
	fallback := &FallbackConfig{
		Enabled: true,
		Ingress: &FallbackIngress{
			ListenerHost: "0.0.0.0",
			HTTPSPort:    8443,
			TLS: TLSPaths{
				CertFile:   "/etc/certs/fullchain.pem",
				KeyFile:    "/etc/certs/privkey.pem",
				IssuerFile: "/etc/certs/ca.pem",
			},
			Routes: []FallbackRoute{
				{Prefix: "/", Cluster: "gateway_http", Domains: "example.com"},
			},
		},
		Clusters: []FallbackCluster{
			{
				Name: "gateway_http",
				Endpoints: []FallbackEndpoint{
					{Host: "10.0.0.10", Port: 8080},
				},
			},
		},
	}
	spec := w.ingressFromFallback(fallback)
	if spec == nil {
		t.Fatalf("expected fallback spec")
	}
	if len(spec.Clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(spec.Clusters))
	}
	if spec.Listener.CertFile == "" {
		t.Fatalf("expected TLS paths, got %v", spec.Listener)
	}
}

func TestGatewayHTTPClusterUsesGatewayPort(t *testing.T) {
	input := builder.Input{
		NodeID: "test-gateway",
		Clusters: []builder.Cluster{
			{
				Name: "gateway_http",
				Endpoints: []builder.Endpoint{
					{Host: "127.0.0.1", Port: 8080},
				},
			},
		},
		GatewayPort:     8080,
		IngressHTTPPort: 80,
	}
	snap, err := builder.BuildSnapshot(input, "version-1")
	if err != nil {
		t.Fatalf("BuildSnapshot: %v", err)
	}
	clusters := snap.GetResources(resource_v3.ClusterType)
	var gwCluster *cluster_v3.Cluster
	for _, res := range clusters {
		c, ok := res.(*cluster_v3.Cluster)
		if !ok {
			continue
		}
		if c.Name == "gateway_http" {
			gwCluster = c
			break
		}
	}
	if gwCluster == nil {
		t.Fatalf("missing gateway_http cluster")
	}
	endpoints := gwCluster.LoadAssignment.GetEndpoints()
	if len(endpoints) == 0 {
		t.Fatalf("gateway_http missing endpoints")
	}
	lbEndpoints := endpoints[0].LbEndpoints
	if len(lbEndpoints) == 0 {
		t.Fatalf("gateway_http endpoint missing lb entry")
	}
	socket := lbEndpoints[0].GetEndpoint().GetAddress().GetSocketAddress()
	if socket == nil {
		t.Fatalf("gateway_http lb endpoint missing socket address")
	}
	if socket.GetPortValue() != 8080 {
		t.Fatalf("gateway_http endpoint port = %d, want 8080", socket.GetPortValue())
	}
}
