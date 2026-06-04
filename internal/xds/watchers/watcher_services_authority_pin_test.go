// @awareness namespace=globular.platform
// @awareness component=platform_xds.watcher_services_authority_pin
// @awareness file_role=architectural_pin_test_for_xds_service_registry_typed_rpc
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestWatcher_BuildServiceResources_NoDirectEtcd pins the v1.2.182
// refactor that moved xDS off the raw /globular/services/* etcd
// scan onto cluster_controller.ListServices.
//
// That prefix is owned by the cluster_controller (writers
// config.PutInstance / config.PutConfig in
// services/golang/config/etcd_service_config.go), so xDS reading
// raw etcd violated
// invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.
//
// This test fails if the buildServiceResourcesFromEtcd function body
// reintroduces:
//   - a direct etcd Get/Put/Delete/Watch against /globular/services/*, or
//   - a clientv3.New call (re-creating the etcd client locally).
func TestWatcher_BuildServiceResources_NoDirectEtcd(t *testing.T) {
	body, err := os.ReadFile("watcher.go")
	if err != nil {
		t.Fatalf("read watcher.go: %v", err)
	}

	const fnHeader = "func (w *Watcher) buildServiceResourcesFromEtcd(ctx context.Context, cfg *XDSConfig) ([]map[string]any, error) {"
	startIdx := strings.Index(string(body), fnHeader)
	if startIdx < 0 {
		t.Fatalf("buildServiceResourcesFromEtcd function header not found — has the signature changed? " +
			"Update this pin so the regression guard still tracks the function body.")
	}
	endIdx := findMatchingBraceWatcher(string(body), startIdx+len(fnHeader)-1)
	if endIdx <= startIdx {
		t.Fatalf("could not locate end of buildServiceResourcesFromEtcd — pin cannot scope to body")
	}
	fnBody := body[startIdx:endIdx]

	etcdRE := regexp.MustCompile(`\.(Get|Put|Delete|Watch)\(\s*[^,)]+,\s*[^,)]*"/globular/services`)
	if loc := etcdRE.FindIndex(fnBody); loc != nil {
		match := etcdRE.FindSubmatch(fnBody)
		t.Errorf("CRITICAL buildServiceResourcesFromEtcd contains a direct etcd %s against /globular/services* "+
			"(near byte offset %d in the function) — violates "+
			"invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. "+
			"Route through controllerclient.Client.ListServices (the v1.2.182 pattern). "+
			"The controller's handler delegates to config.GetServicesConfigurations; "+
			"xDS must not rebuild the registry parse in a consumer.",
			string(match[1]), loc[0])
	}
	if strings.Contains(string(fnBody), "clientv3.New(") {
		t.Errorf("CRITICAL buildServiceResourcesFromEtcd constructs a clientv3.Client locally — " +
			"the typed RPC path makes this unnecessary. The v1.2.182 refactor removed the explicit " +
			"etcd dial in favour of controllerclient.Client.ListServices.")
	}
}
