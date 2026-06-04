// @awareness namespace=globular.platform
// @awareness component=platform_xds.cert_reconciler_authority_pin
// @awareness file_role=architectural_pin_test_for_cert_reconciler_filesystem_delivery
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestCertReconciler_NoDirectEtcd pins the v1.2.179 refactor of
// cert_reconciler.go. Before v1.2.179 the reconciler watched the etcd
// key /globular/pki/bundles/{domain} via clientv3 (Violations 6+7 in
// the xDS audit). The owner of /globular/pki/* is node_agent's
// internal/certs package
// (services/golang/node_agent/node_agent_server/internal/certs/etcd_kv.go::PutBundle),
// so xDS reading raw etcd violated
// invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.
//
// The refactor switches to filesystem-based change detection on the
// local cert/key/CA files that node_agent already writes during
// cert issuance. No new RPC, no new dial path.
//
// This test fails loudly if a future contributor:
//   - reintroduces the clientv3 import in cert_reconciler.go, or
//   - reintroduces a direct etcd Get / Put / Delete / Watch against
//     /globular/* in this file.
func TestCertReconciler_NoDirectEtcd(t *testing.T) {
	body, err := os.ReadFile("cert_reconciler.go")
	if err != nil {
		t.Fatalf("read cert_reconciler.go: %v", err)
	}

	if strings.Contains(string(body), `clientv3 "go.etcd.io/etcd/client/v3"`) ||
		strings.Contains(string(body), `"go.etcd.io/etcd/client/v3"`) {
		t.Errorf("CRITICAL cert_reconciler.go imports go.etcd.io/etcd/client/v3 — " +
			"violates invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. " +
			"The reconciler MUST detect cert changes via fsnotify on local files " +
			"(node_agent owns /globular/pki/bundles/*). Reintroducing the etcd " +
			"client re-opens the bypass vector closed in v1.2.179.")
	}

	re := regexp.MustCompile(`\.(Get|Put|Delete|Watch)\(\s*[^,)]+,\s*[^,)]*"/globular/`)
	if loc := re.FindIndex(body); loc != nil {
		match := re.FindSubmatch(body)
		t.Errorf("CRITICAL cert_reconciler.go contains a direct etcd %s against /globular/* "+
			"(near byte offset %d) — violates invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. "+
			"PKI bundles are owned by node_agent (internal/certs); xDS detects rotation via "+
			"fsnotify on the local cert/key/CA files (v1.2.179 pattern).",
			string(match[1]), loc[0])
	}
}
