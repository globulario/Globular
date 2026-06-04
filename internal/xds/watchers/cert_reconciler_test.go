package watchers

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCertReconcilerLifecycle tests basic lifecycle operations.
// v1.2.179: signature change — NewCertReconciler no longer takes an
// etcdClient; cert change detection moved to filesystem watch.
func TestCertReconcilerLifecycle(t *testing.T) {
	reconciler := NewCertReconciler(nil, "", "", "", "", "")

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

	// Create reconciler — internal-cert paths empty, ACME paths set.
	reconciler := NewCertReconciler(nil, "", "", "", certPath, keyPath)
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
		t.Log("ACME certificate change detected via filesystem watch")
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for ACME certificate change event")
	}
}

// TestCertReconcilerInternalFileWatch tests filesystem watching for
// the internal cluster-CA-issued cert/key/CA files. This is the
// v1.2.179 replacement for the prior etcd-watch test.
func TestCertReconcilerInternalFileWatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cert-reconciler-internal-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	certPath := filepath.Join(tmpDir, "service.crt")
	keyPath := filepath.Join(tmpDir, "service.key")
	caPath := filepath.Join(tmpDir, "ca.crt")

	if err := os.WriteFile(certPath, []byte("initial cert content"), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte("initial key content"), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	if err := os.WriteFile(caPath, []byte("initial ca content"), 0644); err != nil {
		t.Fatalf("write ca: %v", err)
	}

	reconciler := NewCertReconciler(nil, certPath, keyPath, caPath, "", "")
	if err := reconciler.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer reconciler.Stop()

	// Wait for initialization.
	time.Sleep(500 * time.Millisecond)

	// Rotate the cert content.
	if err := os.WriteFile(certPath, []byte("rotated cert content"), 0644); err != nil {
		t.Fatalf("write rotated cert: %v", err)
	}

	select {
	case <-reconciler.CertChangedChan():
		t.Log("internal certificate change detected via filesystem watch")
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for internal certificate change event")
	}
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

	reconciler := NewCertReconciler(nil, "", "", "", certPath, keyPath)
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

	t.Logf("Certificate reconciler initialized with hashes: cert=%s, key=%s",
		certHash[:8], keyHash[:8])
}
