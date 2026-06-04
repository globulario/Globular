// @awareness namespace=globular.platform
// @awareness component=platform_xds.watcher_acme_authority_pin
// @awareness file_role=architectural_pin_test_for_xds_acme_certs_etcd_removal
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"os"
	"regexp"
	"testing"
)

// TestWatcher_NoDirectEtcdAgainstACMECerts pins the v1.2.180 refactor
// that removed xDS's etcd Watch on /globular/acme/certs/* — that
// prefix is owned by the domain service
// (services/golang/domain/reconciler.go::publishCertToEtcd), and node_agent
// already syncs etcd → local disk under /var/lib/globular/pki/acme/.
// xDS now detects ACME rotation via a periodic hash-walk of that local
// directory (watchACMECertDirectory).
//
// This test fails if a future contributor reintroduces:
//   - a direct etcd Get / Put / Delete / Watch against
//     /globular/acme/* in watcher.go.
//
// The clientv3 import is NOT pinned for removal in watcher.go: other
// open-violation paths (loadExternalDomains, buildServiceResourcesFromEtcd)
// still legitimately use it pending their own ratchets.
func TestWatcher_NoDirectEtcdAgainstACMECerts(t *testing.T) {
	body, err := os.ReadFile("watcher.go")
	if err != nil {
		t.Fatalf("read watcher.go: %v", err)
	}

	re := regexp.MustCompile(`\.(Get|Put|Delete|Watch)\(\s*[^,)]+,\s*[^,)]*"/globular/acme/`)
	if loc := re.FindIndex(body); loc != nil {
		match := re.FindSubmatch(body)
		t.Errorf("CRITICAL watcher.go contains a direct etcd %s against /globular/acme/* "+
			"(near byte offset %d) — violates invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage. "+
			"ACME certs are owned by the domain service; node_agent.WatchACMECerts already syncs "+
			"etcd → local disk under /var/lib/globular/pki/acme/. Detect rotation via the v1.2.180 "+
			"watchACMECertDirectory pattern (periodic hash-walk).",
			string(match[1]), loc[0])
	}
}
