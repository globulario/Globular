package watchers

import (
	"log/slog"
	"testing"

	clustercontrollerpb "github.com/globulario/services/golang/clustercontroller/clustercontrollerpb"
)

// TestCertReconcilerWaitsForClusterNetwork ensures reconciler creation defers until cluster config exists.
func TestCertReconcilerWaitsForClusterNetwork(t *testing.T) {
	w := &Watcher{
		logger:         slog.Default(),
		controllerAddr: "controller:443",
		acmeCertPath:   "/tmp/fullchain.pem",
		acmeKeyPath:    "/tmp/privkey.pem",
	}

	w.initializeCertReconciler()
	if w.certReconciler != nil {
		t.Fatalf("expected certReconciler to be nil while cluster network missing")
	}
	if !w.certInitPending {
		t.Fatalf("expected certInitPending to be true when cluster not ready")
	}

	w.clusterNetwork = &clustercontrollerpb.ClusterNetwork{Spec: &clustercontrollerpb.ClusterNetworkSpec{ClusterDomain: "example.local"}}
	w.initializeCertReconciler()
	if w.certReconciler == nil {
		t.Fatalf("expected certReconciler to be created after cluster network available")
	}
	if w.certInitPending {
		t.Fatalf("certInitPending should be false after initialization")
	}
}
