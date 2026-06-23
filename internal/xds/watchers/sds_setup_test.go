package watchers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/globulario/Globular/internal/xds/builder"
	"github.com/globulario/Globular/internal/xds/secrets"
)

// fixtureSDSStateRoot points GetStateRootDir at a temp dir (via GLOBULAR_STATE_DIR)
// and writes the canonical internal PKI material setupSDSSecrets reads, so the
// tests assert the real contract deterministically instead of depending on live
// /var/lib/globular state (which made the secret count 2-vs-3 by whether the
// node happens to have pki/ca.crt). withCA controls whether the optional CA
// bundle (canonical pki/ca.crt) is present. Returns the canonical cert/key/ca
// paths the SUT will emit.
func fixtureSDSStateRoot(t *testing.T, withCA bool) (certPath, keyPath, caPath string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("GLOBULAR_STATE_DIR", root)
	svcDir := filepath.Join(root, "pki", "issued", "services")
	if err := os.MkdirAll(svcDir, 0o755); err != nil {
		t.Fatalf("mkdir services dir: %v", err)
	}
	certPath = filepath.Join(svcDir, "service.crt")
	keyPath = filepath.Join(svcDir, "service.key")
	if err := os.WriteFile(certPath, []byte("cert"), 0o644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte("key"), 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	caPath = filepath.Join(root, "pki", "ca.crt")
	if withCA {
		if err := os.WriteFile(caPath, []byte("ca"), 0o644); err != nil {
			t.Fatalf("write ca: %v", err)
		}
	}
	return certPath, keyPath, caPath
}

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

// TestSetupSDSSecrets_WithCertificates verifies the CURRENT setupSDSSecrets
// contract. AUTHORITY = code: commit 2992c06 ("publish internal-client-cert SDS
// secret for upstream mTLS") and 6070dd4 ("default filter chain must serve
// internal cert, not LE cert") deliberately changed this function to (a) emit
// internal-CA server AND client certs and (b) source them from the canonical PKI
// paths (GetStateRootDir/pki/...), NOT from listener.CertFile. This test predated
// both and was stale (it expected [server, CA] at the listener's path). Aligned
// to the real contract and made hermetic via a fixture state root.
func TestSetupSDSSecrets_WithCertificates(t *testing.T) {
	certPath, keyPath, caPath := fixtureSDSStateRoot(t, true)

	w := &Watcher{}
	// listener cert/key just have to be non-empty to enable SDS; the SUT
	// intentionally ignores them for the internal secrets.
	listener := builder.Listener{
		CertFile:   "/ignored/by/internal/secrets/cert.pem",
		KeyFile:    "/ignored/by/internal/secrets/key.pem",
		IssuerFile: "/ignored/by/internal/secrets/ca.pem",
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enableSDS {
		t.Error("SDS should be enabled when certificates are configured")
	}
	if len(sdsSecrets) != 3 {
		t.Fatalf("expected 3 secrets (internal server + client cert + CA bundle), got: %d", len(sdsSecrets))
	}

	// Internal server cert — canonical PKI path, not the listener's.
	if sdsSecrets[0].Name != secrets.InternalServerCert {
		t.Errorf("expected first secret name %s, got: %s", secrets.InternalServerCert, sdsSecrets[0].Name)
	}
	if sdsSecrets[0].CertPath != certPath {
		t.Errorf("expected canonical cert path %s, got: %s", certPath, sdsSecrets[0].CertPath)
	}
	if sdsSecrets[0].KeyPath != keyPath {
		t.Errorf("expected canonical key path %s, got: %s", keyPath, sdsSecrets[0].KeyPath)
	}

	// Internal client cert (added by 2992c06 for upstream mTLS).
	if sdsSecrets[1].Name != secrets.InternalClientCert {
		t.Errorf("expected second secret name %s, got: %s", secrets.InternalClientCert, sdsSecrets[1].Name)
	}

	// Internal CA bundle — canonical pki/ca.crt, present because the fixture wrote it.
	if sdsSecrets[2].Name != secrets.InternalCABundle {
		t.Errorf("expected third secret name %s, got: %s", secrets.InternalCABundle, sdsSecrets[2].Name)
	}
	if sdsSecrets[2].CAPath != caPath {
		t.Errorf("expected canonical CA path %s, got: %s", caPath, sdsSecrets[2].CAPath)
	}

	t.Log("✓ SDS emits internal server+client certs (canonical PKI) + CA bundle")
}

// TestSetupSDSSecrets_ACMELookingListenerStillUsesInternalPKI: an ACME-looking
// listener cert path must NOT change the internal secrets setupSDSSecrets emits.
// AUTHORITY = code: 6070dd4 ("default filter chain must serve internal cert, not
// LE cert") moved ACME/public-ingress serving OUT of this function — it no longer
// emits a PublicIngressCert secret or records w.acmeCertPath. The old test
// (named _WithACMECertificate) asserted that removed behavior and was stale.
// Asserting the current contract: internal server+client certs from canonical
// PKI, independent of the (ACME-looking) listener path.
func TestSetupSDSSecrets_ACMELookingListenerStillUsesInternalPKI(t *testing.T) {
	certPath, keyPath, _ := fixtureSDSStateRoot(t, false) // no ca.crt → no CA bundle

	w := &Watcher{}
	listener := builder.Listener{
		CertFile: "/var/lib/globular/config/tls/acme/example.com/fullchain.pem",
		KeyFile:  "/var/lib/globular/config/tls/acme/example.com/privkey.pem",
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enableSDS {
		t.Error("SDS should be enabled when cert+key are configured")
	}
	if len(sdsSecrets) != 2 {
		t.Fatalf("expected 2 internal secrets (server + client cert), got: %d", len(sdsSecrets))
	}
	if sdsSecrets[0].Name != secrets.InternalServerCert {
		t.Errorf("expected first secret %s, got: %s", secrets.InternalServerCert, sdsSecrets[0].Name)
	}
	// The internal cert path must be the canonical PKI cert, NOT the ACME listener path.
	if sdsSecrets[0].CertPath != certPath {
		t.Errorf("expected canonical internal cert %s, got the listener/ACME path %s", certPath, sdsSecrets[0].CertPath)
	}
	if sdsSecrets[0].KeyPath != keyPath {
		t.Errorf("expected canonical internal key %s, got: %s", keyPath, sdsSecrets[0].KeyPath)
	}
	if sdsSecrets[1].Name != secrets.InternalClientCert {
		t.Errorf("expected second secret %s, got: %s", secrets.InternalClientCert, sdsSecrets[1].Name)
	}

	t.Log("✓ ACME-looking listener still yields internal-CA secrets (ACME serving lives elsewhere)")
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

// TestSetupSDSSecrets_WithoutCABundle: the CA-bundle secret is optional —
// included only when the canonical pki/ca.crt exists. AUTHORITY = code: the CA
// bundle is gated by fileExists(GetStateRootDir/pki/ca.crt), and the internal
// server+client certs (2992c06) are always emitted. The old test expected a
// single secret and was stale; with no ca.crt fixture, the contract is 2 secrets
// (server + client), no CA bundle. listener.IssuerFile is ignored by the SUT.
func TestSetupSDSSecrets_WithoutCABundle(t *testing.T) {
	certPath, keyPath, _ := fixtureSDSStateRoot(t, false) // no ca.crt written

	w := &Watcher{}
	listener := builder.Listener{
		CertFile:   "/ignored/cert.pem",
		KeyFile:    "/ignored/key.pem",
		IssuerFile: "", // ignored by setupSDSSecrets
	}

	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enableSDS {
		t.Error("SDS should be enabled even without a CA bundle")
	}
	if len(sdsSecrets) != 2 {
		t.Fatalf("expected 2 secrets (internal server + client cert, no CA bundle), got: %d", len(sdsSecrets))
	}
	if sdsSecrets[0].Name != secrets.InternalServerCert || sdsSecrets[0].CertPath != certPath || sdsSecrets[0].KeyPath != keyPath {
		t.Errorf("unexpected server cert secret: %+v", sdsSecrets[0])
	}
	if sdsSecrets[1].Name != secrets.InternalClientCert {
		t.Errorf("expected second secret %s, got: %s", secrets.InternalClientCert, sdsSecrets[1].Name)
	}

	t.Log("✓ SDS emits internal server+client certs with no CA bundle when pki/ca.crt is absent")
}
