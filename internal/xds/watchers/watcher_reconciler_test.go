package watchers

import (
	"log/slog"
	"testing"
)

// TestCertReconcilerCreatesWhenACMEConfigured asserts that the
// reconciler is created as soon as any cert source (internal or
// ACME) is available. v1.2.179: the prior wait-for-cluster-domain
// branch is gone — initializeCertReconciler no longer needs a
// cluster domain because cert change detection moved off the
// /globular/pki/bundles/{domain} etcd key onto filesystem paths
// that are domain-independent (service.crt / service.key / ca.crt
// live under the same well-known directory regardless of cluster
// domain).
func TestCertReconcilerCreatesWhenACMEConfigured(t *testing.T) {
	w := &Watcher{
		logger:         slog.Default(),
		controllerAddr: "controller:443",
		acmeCertPath:   "/tmp/fullchain.pem",
		acmeKeyPath:    "/tmp/privkey.pem",
	}

	w.initializeCertReconciler()
	if w.certReconciler == nil {
		t.Fatalf("expected certReconciler to be created immediately when ACME paths are configured")
	}
	if w.certInitPending {
		t.Fatalf("certInitPending should be false after initialization")
	}
}
