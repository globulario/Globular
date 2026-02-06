package controlplane

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
)

var (
	testCert string
	testKey  string
	testCA   string
)

func init() {
	cert, key, ca := generateTestCertKey()
	testCert = string(cert)
	testKey = string(key)
	testCA = string(ca)
}

func TestMakeSecret_TLSCertificate(t *testing.T) {
	tmpDir := t.TempDir()

	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(certFile, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyFile, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	secret, err := MakeSecret("test-cert", certFile, keyFile, "")
	if err != nil {
		t.Fatalf("MakeSecret failed: %v", err)
	}

	if secret.Name != "test-cert" {
		t.Errorf("expected name test-cert, got %s", secret.Name)
	}

	tlsCert := secret.GetTlsCertificate()
	if tlsCert == nil {
		t.Fatal("expected TlsCertificate, got nil")
	}

	certChain := tlsCert.GetCertificateChain().GetInlineBytes()
	if string(certChain) != testCert {
		t.Error("certificate chain mismatch")
	}

	privKey := tlsCert.GetPrivateKey().GetInlineBytes()
	if string(privKey) != testKey {
		t.Error("private key mismatch")
	}
}

func TestMakeSecret_ValidationContext(t *testing.T) {
	tmpDir := t.TempDir()

	caFile := filepath.Join(tmpDir, "ca.pem")
	if err := os.WriteFile(caFile, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	secret, err := MakeSecret("test-ca", "", "", caFile)
	if err != nil {
		t.Fatalf("MakeSecret failed: %v", err)
	}

	if secret.Name != "test-ca" {
		t.Errorf("expected name test-ca, got %s", secret.Name)
	}

	validationCtx := secret.GetValidationContext()
	if validationCtx == nil {
		t.Fatal("expected ValidationContext, got nil")
	}

	ca := validationCtx.GetTrustedCa().GetInlineBytes()
	if string(ca) != testCA {
		t.Error("CA certificate mismatch")
	}
}

func TestMakeSecret_EmptyName(t *testing.T) {
	_, err := MakeSecret("", "cert.pem", "key.pem", "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestMakeSecret_NoPaths(t *testing.T) {
	_, err := MakeSecret("test", "", "", "")
	if err == nil {
		t.Error("expected error when no paths provided")
	}
}

func TestMakeSecret_MissingFiles(t *testing.T) {
	_, err := MakeSecret("test", "/nonexistent/cert.pem", "/nonexistent/key.pem", "")
	if err == nil {
		t.Error("expected error for missing files")
	}
}

func TestMakeTLSCertificate(t *testing.T) {
	tmpDir := t.TempDir()

	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(certFile, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyFile, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	tlsCert, err := MakeTLSCertificate(certFile, keyFile)
	if err != nil {
		t.Fatalf("MakeTLSCertificate failed: %v", err)
	}

	if tlsCert.GetCertificateChain() == nil {
		t.Error("certificate chain is nil")
	}
	if tlsCert.GetPrivateKey() == nil {
		t.Error("private key is nil")
	}
}

func TestMakeCertificateValidationContext(t *testing.T) {
	tmpDir := t.TempDir()

	caFile := filepath.Join(tmpDir, "ca.pem")
	if err := os.WriteFile(caFile, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	validationCtx, err := MakeCertificateValidationContext(caFile)
	if err != nil {
		t.Fatalf("MakeCertificateValidationContext failed: %v", err)
	}

	if validationCtx.GetTrustedCa() == nil {
		t.Error("trusted CA is nil")
	}
}

func TestHashSecret(t *testing.T) {
	tmpDir := t.TempDir()

	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(certFile, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyFile, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	secret1, _ := MakeSecret("test", certFile, keyFile, "")
	secret2, _ := MakeSecret("test", certFile, keyFile, "")

	hash1, err := HashSecret(secret1)
	if err != nil {
		t.Fatalf("HashSecret failed: %v", err)
	}

	hash2, err := HashSecret(secret2)
	if err != nil {
		t.Fatalf("HashSecret failed: %v", err)
	}

	// Same content should produce same hash
	if hash1 != hash2 {
		t.Error("identical secrets produced different hashes")
	}

	// Hash should be hex string
	if len(hash1) != 64 { // SHA256 = 32 bytes = 64 hex chars
		t.Errorf("expected 64-char hex hash, got %d chars", len(hash1))
	}
}

func TestHashSecret_DifferentCerts(t *testing.T) {
	tmpDir := t.TempDir()

	cert1File := filepath.Join(tmpDir, "cert1.pem")
	cert2File := filepath.Join(tmpDir, "cert2.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(cert1File, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert1: %v", err)
	}
	// Slightly different cert (add a space)
	if err := os.WriteFile(cert2File, []byte(testCert+" "), 0644); err != nil {
		t.Fatalf("write cert2: %v", err)
	}
	if err := os.WriteFile(keyFile, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	secret1, _ := MakeSecret("test", cert1File, keyFile, "")
	secret2, _ := MakeSecret("test", cert2File, keyFile, "")

	hash1, _ := HashSecret(secret1)
	hash2, _ := HashSecret(secret2)

	// Different content should produce different hashes
	if hash1 == hash2 {
		t.Error("different secrets produced same hash")
	}
}

func TestSecretVersion(t *testing.T) {
	tmpDir := t.TempDir()

	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(certFile, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyFile, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	secret, _ := MakeSecret("test", certFile, keyFile, "")

	version, err := SecretVersion(secret)
	if err != nil {
		t.Fatalf("SecretVersion failed: %v", err)
	}

	// Version should start with "sds-"
	if len(version) < 5 || version[:4] != "sds-" {
		t.Errorf("expected version to start with 'sds-', got %s", version)
	}

	// Version should be reasonably short (sds- + 16 hex chars = 20 total)
	if len(version) != 20 {
		t.Errorf("expected version length 20, got %d", len(version))
	}
}

func TestValidateCertKeyPair(t *testing.T) {
	tests := []struct {
		name    string
		cert    []byte
		key     []byte
		wantErr bool
	}{
		{
			name:    "valid PEM",
			cert:    []byte(testCert),
			key:     []byte(testKey),
			wantErr: false,
		},
		{
			name:    "empty cert",
			cert:    []byte{},
			key:     []byte(testKey),
			wantErr: true,
		},
		{
			name:    "empty key",
			cert:    []byte(testCert),
			key:     []byte{},
			wantErr: true,
		},
		{
			name:    "not PEM format cert",
			cert:    []byte("not a certificate"),
			key:     []byte(testKey),
			wantErr: true,
		},
		{
			name:    "not PEM format key",
			cert:    []byte(testCert),
			key:     []byte("not a key"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCertKeyPair(tt.cert, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCertKeyPair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashSecret_Nil(t *testing.T) {
	var nilSecret *tls_v3.Secret
	_, err := HashSecret(nilSecret)
	if err == nil {
		t.Error("expected error for nil secret")
	}
}

// generateTestCertKey creates a self-signed RSA certificate and returns cert, key, and CA PEM.
// For test purposes the CA is the same self-signed cert.
func generateTestCertKey() (certPEM, keyPEM, caPEM []byte) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	serial, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	tpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "test.local"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, _ := x509.CreateCertificate(rand.Reader, tpl, tpl, &priv.PublicKey, priv)

	certBuf := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyBuf := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certBuf, keyBuf, certBuf
}
