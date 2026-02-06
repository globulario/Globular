package watchers

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// TestCertReconcilerLifecycle tests basic lifecycle operations.
func TestCertReconcilerLifecycle(t *testing.T) {
	reconciler := NewCertReconciler(nil, nil, "test.domain", "", "")

	if err := reconciler.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Check channels are available
	if reconciler.CertChangedChan() == nil {
		t.Error("CertChangedChan should not be nil")
	}
	if reconciler.ACMEChangedChan() == nil {
		t.Error("ACMEChangedChan should not be nil")
	}

	reconciler.Stop()
}

// TestCertReconcilerACMEFileWatch tests filesystem watching for ACME certificates.
func TestCertReconcilerACMEFileWatch(t *testing.T) {
	// Create temp directory for test certificates
	tmpDir, err := os.MkdirTemp("", "cert-reconciler-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	certPath := filepath.Join(tmpDir, "fullchain.pem")
	keyPath := filepath.Join(tmpDir, "privkey.pem")

	// Write initial certificates
	testCert := `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU1lfMA0GCSqGSIb3DQEBCwUAMA0xCzAJBgNVBAYTAlVT
-----END CERTIFICATE-----`
	testKey := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg1234567890
-----END PRIVATE KEY-----`

	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	// Create reconciler
	reconciler := NewCertReconciler(nil, nil, "", certPath, keyPath)
	if err := reconciler.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer reconciler.Stop()

	// Wait for initialization
	time.Sleep(500 * time.Millisecond)

	// Modify certificate - should trigger change event
	modifiedCert := testCert + "\n# Modified"
	if err := os.WriteFile(certPath, []byte(modifiedCert), 0644); err != nil {
		t.Fatalf("write modified cert: %v", err)
	}

	// Wait for change event (with timeout)
	select {
	case <-reconciler.ACMEChangedChan():
		t.Log("✓ ACME certificate change detected via filesystem watch")
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for ACME certificate change event")
	}
}

// TestCertReconcilerEtcdWatch tests etcd watching (requires mock or skip).
func TestCertReconcilerEtcdWatch(t *testing.T) {
	// This test requires a real etcd instance or mock
	// For now, just test that the reconciler can be created with etcd client
	t.Skip("Requires etcd instance - manual test only")

	// In a real test environment with etcd:
	// 1. Create etcd client
	// 2. Write initial certificate generation
	// 3. Start reconciler
	// 4. Update certificate generation in etcd
	// 5. Verify change event received
}

// TestCertReconcilerInitialization tests state initialization.
func TestCertReconcilerInitialization(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cert-reconciler-init-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")

	// Write test files
	if err := os.WriteFile(certPath, []byte("test cert content"), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte("test key content"), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	reconciler := NewCertReconciler(nil, nil, "", certPath, keyPath)
	if err := reconciler.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer reconciler.Stop()

	// Wait for initialization
	time.Sleep(500 * time.Millisecond)

	// Check that hashes were initialized
	reconciler.mu.RLock()
	certHash := reconciler.lastACMECertHash
	keyHash := reconciler.lastACMEKeyHash
	reconciler.mu.RUnlock()

	if certHash == "" {
		t.Error("certificate hash should be initialized")
	}
	if keyHash == "" {
		t.Error("key hash should be initialized")
	}

	t.Logf("✓ Certificate reconciler initialized with hashes: cert=%s, key=%s",
		certHash[:8], keyHash[:8])
}

// Helper to create a mock etcd client (if needed for tests).
func createMockEtcdClient(t *testing.T) *clientv3.Client {
	t.Skip("Mock etcd client not implemented")
	return nil
}

// Helper to write certificate bundle to etcd.
func writeCertBundle(ctx context.Context, client *clientv3.Client, domain string, generation uint64) error {
	key := "/globular/pki/bundles/" + domain
	value, err := json.Marshal(map[string]interface{}{
		"generation": generation,
	})
	if err != nil {
		return err
	}
	_, err = client.Put(ctx, key, string(value))
	return err
}
