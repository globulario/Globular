// @awareness namespace=globular.platform
// @awareness component=platform_gateway.admin_services_authority_pin
// @awareness file_role=architectural_pin_test_for_gateway_fetchClusterNodeMap_typed_rpc
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package admin

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestServices_FetchClusterNodeMap_NoDirectEtcd pins the refactor
// that moved the gateway off the raw /globular/clustercontroller/state
// etcd read onto cluster_controller.ListNodes.
//
// That key is owned by cluster_controller (its persistence layer),
// so the gateway reading raw etcd violated
// invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.
//
// This test fails if the fetchClusterNodeMap function body
// reintroduces a direct etcd Get/Put/Delete/Watch against
// /globular/* — or a config.GetEtcdClient() call (which strongly
// implies an impending etcd read).
func TestServices_FetchClusterNodeMap_NoDirectEtcd(t *testing.T) {
	body, err := os.ReadFile("services.go")
	if err != nil {
		t.Fatalf("read services.go: %v", err)
	}

	const fnHeader = "func fetchClusterNodeMap(ctx context.Context) map[string]string {"
	startIdx := strings.Index(string(body), fnHeader)
	if startIdx < 0 {
		t.Fatalf("fetchClusterNodeMap function header not found — has the signature changed? " +
			"Update this pin so the regression guard still tracks the function body.")
	}
	endIdx := findMatchingBraceServices(string(body), startIdx+len(fnHeader)-1)
	if endIdx <= startIdx {
		t.Fatalf("could not locate end of fetchClusterNodeMap — pin cannot scope to body")
	}
	fnBody := body[startIdx:endIdx]

	re := regexp.MustCompile(`\.(Get|Put|Delete|Watch)\(\s*[^,)]+,\s*[^,)]*"/globular/`)
	if loc := re.FindIndex(fnBody); loc != nil {
		match := re.FindSubmatch(fnBody)
		t.Errorf("CRITICAL fetchClusterNodeMap contains a direct etcd %s against /globular/* "+
			"(near byte offset %d in the function) — violates "+
			"invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. "+
			"Route through controllerclient.Client.ListNodes; the controller's "+
			"NodeRecord.Identity carries the same Hostname + Ips fields.",
			string(match[1]), loc[0])
	}
	if strings.Contains(string(fnBody), "config_.GetEtcdClient(") {
		t.Errorf("CRITICAL fetchClusterNodeMap calls config_.GetEtcdClient() — " +
			"that path was removed in the typed-RPC refactor. The function MUST consume " +
			"the cluster_controller's typed ListNodes RPC, not its persistence layer directly.")
	}
}

// findMatchingBraceServices returns the index of the closing '}'
// that balances the '{' at openIdx within src. Returns -1 if not
// found.
func findMatchingBraceServices(src string, openIdx int) int {
	depth := 0
	for i := openIdx; i < len(src); i++ {
		switch src[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
