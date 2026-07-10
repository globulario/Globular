package builder

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

// TestBuildSnapshot_TLSNotReady_SurfacesMissing443 is the regression for the
// Envoy LDS/CDS wedge finding (invariant
// envoy.lds_progress_required_for_http_mesh_readiness). When the ingress cert
// is not materialized (e.g. the ACME symlink is missing), the :443 HTTPS
// listener cannot be built. The snapshot must still build and still carry
// clusters (CDS), but it must NOT contain a :443 listener — and, critically,
// the dropped-443 condition must be SURFACED rather than dropped silently, so
// the wedge is visible at the control-plane source even when the doctor's
// Prometheus rule is not scraping.
func TestBuildSnapshot_TLSNotReady_SurfacesMissing443(t *testing.T) {
	input := Input{
		NodeID: "test-node",
		Listener: Listener{
			Name:      "ingress",
			RouteName: "ingress_routes",
			Host:      "0.0.0.0",
			Port:      443,
			CertFile:  "/nonexistent/fullchain.pem", // TLS not ready
			KeyFile:   "/nonexistent/privkey.pem",
		},
		Routes:          []Route{{Prefix: "/", Cluster: "gateway_http"}},
		Clusters:        []Cluster{{Name: "gateway_http", Endpoints: []Endpoint{{Host: "127.0.0.1", Port: 8080}}}},
		IngressHTTPPort: 80,
	}

	// Capture stderr to assert the wedge condition is surfaced.
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	snapshot, err := BuildSnapshot(input, "v1")
	_ = w.Close()
	os.Stderr = oldStderr
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	logs := buf.String()

	if err != nil {
		t.Fatalf("BuildSnapshot failed: %v", err)
	}
	if snapshot == nil {
		t.Fatal("expected non-nil snapshot")
	}

	// CDS present...
	if len(snapshot.GetResources(resource_v3.ClusterType)) == 0 {
		t.Fatal("expected clusters (CDS) in snapshot")
	}

	// ...but no :443 HTTPS listener (the exact wedge shape).
	for name, res := range snapshot.GetResources(resource_v3.ListenerType) {
		if l, ok := res.(*listener_v3.Listener); ok {
			if p := l.GetAddress().GetSocketAddress().GetPortValue(); p == 443 {
				t.Fatalf("unexpected :443 listener %q built when TLS is not ready", name)
			}
		}
	}

	// The condition must be surfaced (previously silent).
	if !strings.Contains(logs, "omitting :443") || !strings.Contains(logs, "LDS :443=0") {
		t.Fatalf("expected a WARN surfacing the clusters-without-:443 wedge, got stderr: %q", logs)
	}
}
