package watchers

import (
	"testing"

	"github.com/globulario/Globular/internal/xds/builder"
)

// TestHashInputVersion_ContentAddressed verifies that the xDS snapshot version
// is a stable hash of the input content and excludes the Version field itself.
// This is the guard against the pre-fix behaviour where every poll-tick rebuild
// got a fresh time.Now().UnixNano() version, churning Envoy and dropping
// in-flight mesh RPCs (event circuit-breaker flapping).
func TestHashInputVersion_ContentAddressed(t *testing.T) {
	base := builder.Input{
		NodeID:      "node-1",
		GatewayPort: 8443,
		Clusters:    []builder.Cluster{{Name: "a"}, {Name: "b"}},
	}

	// Identical content -> identical version (stable across rebuilds).
	v1 := hashInputVersion(base)
	v2 := hashInputVersion(base)
	if v1 == "" {
		t.Fatal("hashInputVersion returned empty version")
	}
	if v1 != v2 {
		t.Fatalf("identical input produced different versions: %q vs %q", v1, v2)
	}

	// The Version field must not influence the hash — otherwise the version is
	// self-referential and can never be stable.
	withVer := base
	withVer.Version = "arbitrary-previous-version"
	if got := hashInputVersion(withVer); got != v1 {
		t.Fatalf("Version field leaked into hash: %q vs %q", got, v1)
	}

	// Different content -> different version (real changes must still push).
	changed := base
	changed.GatewayPort = 9443
	if hashInputVersion(changed) == v1 {
		t.Fatal("changed content (GatewayPort) produced the same version — real config change would be suppressed")
	}

	changedClusters := base
	changedClusters.Clusters = []builder.Cluster{{Name: "a"}, {Name: "c"}}
	if hashInputVersion(changedClusters) == v1 {
		t.Fatal("changed cluster set produced the same version — real config change would be suppressed")
	}
}
