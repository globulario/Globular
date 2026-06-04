// @awareness namespace=globular.platform
// @awareness component=platform_xds.routing_reconciler_authority_pin
// @awareness file_role=architectural_pin_test_for_routing_reconciler_typed_rpc_routing
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// Architectural pin for the v1.2.177 refactor of
// routing_reconciler.go.
//
// Before v1.2.177 routing_reconciler.go watched the etcd key
// /globular/routing/refresh-generation directly via clientv3.
// That prefix is owned by cluster_controller
// (writeRoutingRefresh in
// services/golang/cluster_controller/cluster_controller_server/server.go),
// so the consumer reading raw etcd violated
// invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.
//
// The refactor routes through the typed
// cluster_controller.GetRoutingRefresh RPC (added in v1.2.177).
// xDS polls every ~2s instead of watching etcd.
//
// This test fails loudly if a future contributor:
//   - reintroduces the clientv3 import in routing_reconciler.go, or
//   - reintroduces a direct etcd Get / Put / Delete / Watch against
//     /globular/routing/* in this file.
//
// Anchored by:
//
//	invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage
//	forbidden_fix:read_owned_etcd_prefix_directly_instead_of_calling_owner_rpc
func TestRoutingReconciler_NoDirectEtcd(t *testing.T) {
	body, err := os.ReadFile("routing_reconciler.go")
	if err != nil {
		t.Fatalf("read routing_reconciler.go: %v", err)
	}

	if strings.Contains(string(body), `clientv3 "go.etcd.io/etcd/client/v3"`) ||
		strings.Contains(string(body), `"go.etcd.io/etcd/client/v3"`) {
		t.Errorf("CRITICAL routing_reconciler.go imports go.etcd.io/etcd/client/v3 — " +
			"violates invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. " +
			"The reconciler MUST poll cluster_controller.GetRoutingRefresh via " +
			"controllerclient.Client, never watch etcd directly. Reintroducing the " +
			"etcd client re-opens the bypass vector closed in v1.2.177.")
	}

	re := regexp.MustCompile(`\.(Get|Put|Delete|Watch)\(\s*[^,)]+,\s*"/globular/`)
	if loc := re.FindIndex(body); loc != nil {
		match := re.FindSubmatch(body)
		t.Errorf("CRITICAL routing_reconciler.go contains a direct etcd %s against /globular/* "+
			"(near byte offset %d) — violates invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. "+
			"The routing-refresh signal is owned by cluster_controller; consume it via "+
			"controllerclient.Client.GetRoutingRefresh (the v1.2.177 pattern).",
			string(match[1]), loc[0])
	}
}
