package watchers

import (
	"context"
	"testing"

	"github.com/globulario/Globular/internal/controlplane"
	"github.com/globulario/Globular/internal/xds/builder"
)

func TestBuildIngressModeInput_Basic(t *testing.T) {
	w := &Watcher{}
	ing := &IngressSpec{
		Listener: builder.Listener{Host: "", Port: 443},
		Clusters: []builder.Cluster{{Name: "ingress"}},
		Routes:   []builder.Route{{Prefix: "/", Cluster: "ingress"}},
		HTTPPort: 80,
	}

	svcClusters := []builder.Cluster{{Name: "svc"}}
	svcRoutes := []builder.Route{{Prefix: "/svc", Cluster: "svc"}}

	in, err := w.buildIngressModeInput(context.Background(), nil, svcClusters, svcRoutes, ing)
	if err != nil {
		t.Fatalf("buildIngressModeInput error: %v", err)
	}

	if in.Listener.Host != "0.0.0.0" {
		t.Fatalf("listener host defaulted to %q", in.Listener.Host)
	}
	if in.Listener.RouteName != defaultRouteName {
		t.Fatalf("route name defaulted to %q", in.Listener.RouteName)
	}
	if findClusterByName(in.Clusters, "svc") == nil || findClusterByName(in.Clusters, "ingress") == nil {
		t.Fatalf("expected svc and ingress clusters to be present, got %v", in.Clusters)
	}
	if in.IngressHTTPPort != 80 {
		t.Fatalf("IngressHTTPPort expected 80, got %d", in.IngressHTTPPort)
	}
}

func TestBuildLegacyModeInput_Basic(t *testing.T) {
	w := &Watcher{protocol: "https"}
	svcClusters := []builder.Cluster{{Name: "svc"}}
	svcRoutes := []builder.Route{{Prefix: "/svc", Cluster: "svc"}}

	in, err := w.buildLegacyModeInput(context.Background(), nil, svcClusters, svcRoutes)
	if err != nil {
		t.Fatalf("buildLegacyModeInput error: %v", err)
	}

	if len(in.Clusters) == 0 {
		t.Fatal("expected at least one cluster from legacy resources")
	}
	if len(in.Routes) == 0 {
		t.Fatal("expected at least one route from legacy resources")
	}
}

func TestSetupSDSSecrets_Stable(t *testing.T) {
	w := &Watcher{}
	in := &builder.Input{
		Listener: builder.Listener{CertFile: "/tmp/fullchain.pem", KeyFile: "/tmp/privkey.pem"},
	}

	if err := w.setupSDSSecrets(in); err != nil {
		t.Fatalf("setupSDSSecrets error: %v", err)
	}
	firstLen := len(in.SDSSecrets)
	if !in.EnableSDS || firstLen == 0 {
		t.Fatalf("expected SDS enabled with secrets, got enabled=%v len=%d", in.EnableSDS, firstLen)
	}

	// Second call should not grow secrets unexpectedly
	if err := w.setupSDSSecrets(in); err != nil {
		t.Fatalf("second setupSDSSecrets error: %v", err)
	}
	if len(in.SDSSecrets) != firstLen {
		t.Fatalf("expected secrets length stable, got %d -> %d", firstLen, len(in.SDSSecrets))
	}
}

func TestBuildDynamicInput_SelectsCorrectMode(t *testing.T) {
	ctx := context.Background()
	calledIngress := false

	w := &Watcher{
		nodeID: "node-1",
		serviceResourcesFn: func(context.Context, *XDSConfig) ([]builder.Cluster, []builder.Route, error) {
			return []builder.Cluster{{Name: "svc"}}, []builder.Route{{Prefix: "/svc", Cluster: "svc"}}, nil
		},
		ingressSpecFn: func(context.Context, *XDSConfig) (*IngressSpec, error) {
			if calledIngress {
				return &IngressSpec{
					Listener: builder.Listener{Port: controlplane.DefaultIngressPort("0.0.0.0")},
					Clusters: []builder.Cluster{{Name: "ingress"}},
					Routes:   []builder.Route{{Prefix: "/", Cluster: "ingress"}},
					HTTPPort: 80,
				}, nil
			}
			return nil, nil
		},
	}

	// Legacy path first (ingressSpec nil)
	input, _, err := w.buildDynamicInput(ctx, nil)
	if err != nil {
		t.Fatalf("buildDynamicInput legacy error: %v", err)
	}
	if findClusterByName(input.Clusters, "ingress") != nil {
		t.Fatalf("unexpected ingress cluster in legacy mode")
	}

	// Ingress path
	calledIngress = true
	input, _, err = w.buildDynamicInput(ctx, nil)
	if err != nil {
		t.Fatalf("buildDynamicInput ingress error: %v", err)
	}
	if findClusterByName(input.Clusters, "ingress") == nil {
		t.Fatalf("expected ingress cluster in ingress mode")
	}
}
