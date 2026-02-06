package watchers

import (
	"os"
	"strings"
	"testing"

	"github.com/globulario/Globular/internal/xds/builder"
	"github.com/globulario/Globular/internal/xds/secrets"
)

// TestSetupSDSSecrets_NoCertificates tests that SDS is disabled when no certificates are configured
func TestSetupSDSSecrets_NoCertificates(t *testing.T) {
	w := &Watcher{}
	listener := builder.Listener{
		CertFile: "",
		KeyFile:  "",
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enableSDS {
		t.Error("SDS should be disabled when no certificates are configured")
	}
	if sdsSecrets != nil {
		t.Errorf("expected nil secrets, got: %v", sdsSecrets)
	}
	t.Log("✓ SDS correctly disabled when no certificates configured")
}

// TestSetupSDSSecrets_WithCertificates tests that SDS is enabled with valid certificate configuration
func TestSetupSDSSecrets_WithCertificates(t *testing.T) {
	w := &Watcher{}
	listener := builder.Listener{
		CertFile:   "/path/to/cert.pem",
		KeyFile:    "/path/to/key.pem",
		IssuerFile: "/path/to/ca.pem",
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enableSDS {
		t.Error("SDS should be enabled when certificates are configured")
	}
	if len(sdsSecrets) != 2 {
		t.Fatalf("expected 2 secrets (server cert + CA bundle), got: %d", len(sdsSecrets))
	}

	// Verify internal server cert
	if sdsSecrets[0].Name != secrets.InternalServerCert {
		t.Errorf("expected first secret name %s, got: %s", secrets.InternalServerCert, sdsSecrets[0].Name)
	}
	if sdsSecrets[0].CertPath != listener.CertFile {
		t.Errorf("expected cert path %s, got: %s", listener.CertFile, sdsSecrets[0].CertPath)
	}
	if sdsSecrets[0].KeyPath != listener.KeyFile {
		t.Errorf("expected key path %s, got: %s", listener.KeyFile, sdsSecrets[0].KeyPath)
	}

	// Verify internal CA bundle
	if sdsSecrets[1].Name != secrets.InternalCABundle {
		t.Errorf("expected second secret name %s, got: %s", secrets.InternalCABundle, sdsSecrets[1].Name)
	}
	if sdsSecrets[1].CAPath != listener.IssuerFile {
		t.Errorf("expected CA path %s, got: %s", listener.IssuerFile, sdsSecrets[1].CAPath)
	}

	t.Log("✓ SDS correctly enabled with internal server cert and CA bundle")
}

// TestSetupSDSSecrets_WithACMECertificate tests that ACME certificates are detected and added
func TestSetupSDSSecrets_WithACMECertificate(t *testing.T) {
	w := &Watcher{}
	listener := builder.Listener{
		CertFile: "/path/to/acme/fullchain.pem",
		KeyFile:  "/path/to/acme/privkey.pem",
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enableSDS {
		t.Error("SDS should be enabled for ACME certificates")
	}
	if len(sdsSecrets) != 2 {
		t.Fatalf("expected 2 secrets (internal cert + public ACME cert), got: %d", len(sdsSecrets))
	}

	// Verify internal server cert
	if sdsSecrets[0].Name != secrets.InternalServerCert {
		t.Errorf("expected first secret name %s, got: %s", secrets.InternalServerCert, sdsSecrets[0].Name)
	}

	// Verify public ingress cert (ACME)
	if sdsSecrets[1].Name != secrets.PublicIngressCert {
		t.Errorf("expected second secret name %s, got: %s", secrets.PublicIngressCert, sdsSecrets[1].Name)
	}
	if sdsSecrets[1].CertPath != listener.CertFile {
		t.Errorf("expected ACME cert path %s, got: %s", listener.CertFile, sdsSecrets[1].CertPath)
	}

	// Verify ACME paths stored in watcher
	if w.acmeCertPath != listener.CertFile {
		t.Errorf("expected watcher.acmeCertPath %s, got: %s", listener.CertFile, w.acmeCertPath)
	}
	if w.acmeKeyPath != listener.KeyFile {
		t.Errorf("expected watcher.acmeKeyPath %s, got: %s", listener.KeyFile, w.acmeKeyPath)
	}

	t.Log("✓ ACME certificates correctly detected and configured")
}

// TestSetupSDSSecrets_SecurityViolation tests that insecure xDS is rejected when SDS is enabled
func TestSetupSDSSecrets_SecurityViolation(t *testing.T) {
	// Set insecure mode
	os.Setenv("GLOBULAR_XDS_INSECURE", "1")
	defer os.Unsetenv("GLOBULAR_XDS_INSECURE")

	w := &Watcher{}
	listener := builder.Listener{
		CertFile: "/path/to/cert.pem",
		KeyFile:  "/path/to/key.pem",
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err == nil {
		t.Fatal("expected security violation error, got nil")
	}
	if !strings.Contains(err.Error(), "security violation") {
		t.Errorf("expected 'security violation' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "plaintext") {
		t.Errorf("expected 'plaintext' in error message, got: %v", err)
	}
	if enableSDS {
		t.Error("SDS should not be enabled when security validation fails")
	}
	if sdsSecrets != nil {
		t.Errorf("expected nil secrets on error, got: %v", sdsSecrets)
	}

	t.Log("✓ Security violation correctly detected for insecure xDS + SDS")
}

// TestSetupSDSSecrets_OnlyCertFile tests that both cert and key are required
func TestSetupSDSSecrets_OnlyCertFile(t *testing.T) {
	w := &Watcher{}
	listener := builder.Listener{
		CertFile: "/path/to/cert.pem",
		KeyFile:  "", // Missing key file
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enableSDS {
		t.Error("SDS should be disabled when key file is missing")
	}
	if sdsSecrets != nil {
		t.Errorf("expected nil secrets, got: %v", sdsSecrets)
	}

	t.Log("✓ SDS correctly disabled when key file is missing")
}

// TestSetupSDSSecrets_WithoutIssuerFile tests that CA bundle is optional
func TestSetupSDSSecrets_WithoutIssuerFile(t *testing.T) {
	w := &Watcher{}
	listener := builder.Listener{
		CertFile:   "/path/to/cert.pem",
		KeyFile:    "/path/to/key.pem",
		IssuerFile: "", // No CA bundle
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enableSDS {
		t.Error("SDS should be enabled even without CA bundle")
	}
	if len(sdsSecrets) != 1 {
		t.Fatalf("expected 1 secret (server cert only), got: %d", len(sdsSecrets))
	}

	// Verify only internal server cert is present
	if sdsSecrets[0].Name != secrets.InternalServerCert {
		t.Errorf("expected secret name %s, got: %s", secrets.InternalServerCert, sdsSecrets[0].Name)
	}

	t.Log("✓ SDS correctly configured without CA bundle")
}
