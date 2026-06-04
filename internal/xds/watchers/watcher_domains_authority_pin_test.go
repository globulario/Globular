// @awareness namespace=globular.platform
// @awareness component=platform_xds.watcher_domains_authority_pin
// @awareness file_role=architectural_pin_test_for_xds_external_domains_typed_rpc
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestWatcher_LoadExternalDomains_NoDirectEtcd pins the v1.2.181/182
// refactor that moved xDS off the raw /globular/domains/v1/* + per-
// FQDN /status etcd reads onto cluster_controller.ListExternalDomains.
//
// That prefix is owned by the controller's embedded domain reconciler
// (services/golang/domain/reconciler.go), so xDS reading raw etcd
// violated invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.
//
// This test fails if the loadExternalDomains function body
// reintroduces a direct etcd Get/Put/Delete/Watch against
// /globular/domains/* — even if the file as a whole retains
// clientv3 for other paths (buildServiceResourcesFromEtcd is still
// pending its own ratchet).
func TestWatcher_LoadExternalDomains_NoDirectEtcd(t *testing.T) {
	body, err := os.ReadFile("watcher.go")
	if err != nil {
		t.Fatalf("read watcher.go: %v", err)
	}

	const fnHeader = "func (w *Watcher) loadExternalDomains(ctx context.Context) ([]ExternalDomainRuntime, error) {"
	startIdx := strings.Index(string(body), fnHeader)
	if startIdx < 0 {
		t.Fatalf("loadExternalDomains function header not found — has the signature changed? " +
			"Update this pin so the regression guard still tracks the function body.")
	}
	endIdx := findMatchingBraceWatcher(string(body), startIdx+len(fnHeader)-1)
	if endIdx <= startIdx {
		t.Fatalf("could not locate end of loadExternalDomains — pin cannot scope to body")
	}
	fnBody := body[startIdx:endIdx]

	re := regexp.MustCompile(`\.(Get|Put|Delete|Watch)\(\s*[^,)]+,\s*[^,)]*"/globular/`)
	if loc := re.FindIndex(fnBody); loc != nil {
		match := re.FindSubmatch(fnBody)
		t.Errorf("CRITICAL watcher.go::loadExternalDomains contains a direct etcd %s against /globular/* "+
			"(near byte offset %d in the function) — violates "+
			"invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. "+
			"Route through controllerclient.Client.ListExternalDomains (the v1.2.181/182 "+
			"pattern). The controller's handler reads via the owner's typed EtcdDomainStore; "+
			"xDS must not rebuild the read in a consumer.",
			string(match[1]), loc[0])
	}
}

// findMatchingBraceWatcher returns the index of the closing '}' that
// balances the '{' at openIdx within src. Returns -1 if not found.
// Simple brace counter; does not handle braces inside string literals —
// fine for the function bodies this pin scans.
func findMatchingBraceWatcher(src string, openIdx int) int {
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
