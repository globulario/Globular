// @awareness namespace=globular.platform
// @awareness component=platform_xds.watcher_pki_authority_pin
// @awareness file_role=architectural_pin_test_for_xds_pki_bundles_read_removal
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestWatcher_NoCheckCertificateGeneration pins the v1.2.178 removal
// of watcher.checkCertificateGeneration. That function polled
// /globular/pki/bundles/{domain} directly from etcd as a fallback
// when w.certReconciler was nil. The fallback violated
// invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage
// (PKI is owned by the security/PKI service) and was a no-op in
// practice because both it and the reconciler share the same
// w.etcdClient prerequisite.
//
// This test fails if a future contributor reintroduces:
//   - the checkCertificateGeneration function name in watcher.go, or
//   - any direct etcd Get / Put / Delete / Watch against
//     /globular/pki/* in watcher.go.
//
// (cert_reconciler.go is intentionally NOT pinned here — that file
// still legitimately watches PKI for cluster-internal CA rotation
// events; it's the next ratchet target. When that ratchet lands,
// extend this pin to cover cert_reconciler.go too.)
func TestWatcher_NoCheckCertificateGeneration(t *testing.T) {
	body, err := os.ReadFile("watcher.go")
	if err != nil {
		t.Fatalf("read watcher.go: %v", err)
	}

	if strings.Contains(string(body), "func (w *Watcher) checkCertificateGeneration(") {
		t.Errorf("CRITICAL watcher.go declares func checkCertificateGeneration — " +
			"removed in v1.2.178 because it scanned /globular/pki/bundles/* directly. " +
			"Reintroducing this function reopens the bypass vector against the " +
			"security/PKI service's owned prefix. If a real fallback is needed, " +
			"add a typed RPC on the owner. Anchored by " +
			"invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.")
	}

	re := regexp.MustCompile(`\.(Get|Put|Delete|Watch)\(\s*[^,)]+,\s*[^,)]*"/globular/pki/`)
	if loc := re.FindIndex(body); loc != nil {
		match := re.FindSubmatch(body)
		t.Errorf("CRITICAL watcher.go contains a direct etcd %s against /globular/pki/* "+
			"(near byte offset %d) — violates invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. "+
			"PKI bundles are owned by the security/PKI service; route reads through a typed RPC, "+
			"not raw etcd.",
			string(match[1]), loc[0])
	}
}
