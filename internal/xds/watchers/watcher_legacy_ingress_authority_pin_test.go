// @awareness namespace=globular.platform
// @awareness component=platform_xds.watcher_legacy_ingress_authority_pin
// @awareness file_role=architectural_pin_test_for_xds_self_owned_legacy_etcd_ingress
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestParseEtcdIngress_SelfOwnedPrefixOnly pins the one remaining direct etcd
// read path in xDS: the legacy bootstrap parser for /globular/xds/v1/*.
//
// This is allowed because xDS owns that prefix and the parser is read-only. The
// test exists so future edits cannot quietly widen that exception onto foreign
// owner prefixes such as /globular/services/*, /globular/domains/*, or
// /globular/pki/*.
func TestParseEtcdIngress_SelfOwnedPrefixOnly(t *testing.T) {
	body, err := os.ReadFile("etcd_ingress.go")
	if err != nil {
		t.Fatalf("read etcd_ingress.go: %v", err)
	}
	src := string(body)

	if !strings.Contains(src, `etcdBasePrefix     = "/globular/xds/v1"`) {
		t.Fatalf("expected etcd_ingress.go to pin its owned prefix at /globular/xds/v1")
	}
	if strings.Contains(src, ".Put(") || strings.Contains(src, ".Delete(") || strings.Contains(src, ".Watch(") {
		t.Errorf("CRITICAL etcd_ingress.go must stay read-only; found Put/Delete/Watch usage")
	}

	foreignPrefixRE := regexp.MustCompile(`"/globular/(services|domains|pki|acme|routing|clustercontroller)/`)
	if loc := foreignPrefixRE.FindStringIndex(src); loc != nil {
		t.Errorf("CRITICAL etcd_ingress.go references a foreign-owned etcd prefix near byte offset %d; "+
			"the legacy exception is only for xDS-owned /globular/xds/v1/*", loc[0])
	}
}

// TestWatcher_IngressFromEtcd_DelegatesToOwnedParser pins that watcher.go keeps
// the legacy etcd path bounded to the owned ingress parser instead of growing
// new direct reads inline.
func TestWatcher_IngressFromEtcd_DelegatesToOwnedParser(t *testing.T) {
	body, err := os.ReadFile("watcher.go")
	if err != nil {
		t.Fatalf("read watcher.go: %v", err)
	}

	const fnHeader = "func (w *Watcher) ingressFromEtcd(ctx context.Context, cfg *XDSConfig) (*IngressSpec, error) {"
	startIdx := strings.Index(string(body), fnHeader)
	if startIdx < 0 {
		t.Fatalf("ingressFromEtcd function header not found — update this pin if the signature changed")
	}
	endIdx := findMatchingBraceWatcher(string(body), startIdx+len(fnHeader)-1)
	if endIdx <= startIdx {
		t.Fatalf("could not locate end of ingressFromEtcd — pin cannot scope to body")
	}
	fnBody := string(body[startIdx:endIdx])

	if !strings.Contains(fnBody, "parseEtcdIngress(") {
		t.Fatalf("expected ingressFromEtcd to delegate parsing to parseEtcdIngress")
	}

	foreignPrefixRE := regexp.MustCompile(`"/globular/(services|domains|pki|acme|routing|clustercontroller)/`)
	if loc := foreignPrefixRE.FindStringIndex(fnBody); loc != nil {
		t.Errorf("CRITICAL ingressFromEtcd references a foreign-owned etcd prefix near byte offset %d; "+
			"legacy bootstrap must stay scoped to the xDS-owned parser only", loc[0])
	}
}
