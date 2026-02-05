package sds

import (
	"os"
	"testing"
)

func TestBuildInternalServerSecret(t *testing.T) {
	tmpDir := t.TempDir()

	certPath := tmpDir + "/cert.pem"
	keyPath := tmpDir + "/key.pem"

	// Write test cert and key
	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	paths := CertPaths{
		CertFile: certPath,
		KeyFile:  keyPath,
	}

	secret, err := BuildInternalServerSecret(paths)
	if err != nil {
		t.Fatalf("BuildInternalServerSecret failed: %v", err)
	}

	if secret.Name != InternalServerCert {
		t.Errorf("expected name %s, got %s", InternalServerCert, secret.Name)
	}

	tlsCert := secret.GetTlsCertificate()
	if tlsCert == nil {
		t.Fatal("expected TlsCertificate, got nil")
	}
}

func TestBuildInternalServerSecret_MissingCert(t *testing.T) {
	paths := CertPaths{
		CertFile: "/nonexistent/cert.pem",
		KeyFile:  "/nonexistent/key.pem",
	}

	_, err := BuildInternalServerSecret(paths)
	if err == nil {
		t.Error("expected error for missing cert file")
	}
}

func TestBuildInternalServerSecret_EmptyPaths(t *testing.T) {
	paths := CertPaths{
		CertFile: "",
		KeyFile:  "",
	}

	_, err := BuildInternalServerSecret(paths)
	if err == nil {
		t.Error("expected error for empty paths")
	}
}

func TestBuildInternalServerSecret_MissingKey(t *testing.T) {
	tmpDir := t.TempDir()

	certPath := tmpDir + "/cert.pem"
	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}

	paths := CertPaths{
		CertFile: certPath,
		KeyFile:  "/nonexistent/key.pem",
	}

	_, err := BuildInternalServerSecret(paths)
	if err == nil {
		t.Error("expected error for missing key file")
	}
}

func TestBuildInternalCASecret(t *testing.T) {
	tmpDir := t.TempDir()

	caPath := tmpDir + "/ca.pem"
	if err := os.WriteFile(caPath, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	secret, err := BuildInternalCASecret(caPath)
	if err != nil {
		t.Fatalf("BuildInternalCASecret failed: %v", err)
	}

	if secret.Name != InternalCABundle {
		t.Errorf("expected name %s, got %s", InternalCABundle, secret.Name)
	}

	validationCtx := secret.GetValidationContext()
	if validationCtx == nil {
		t.Fatal("expected ValidationContext, got nil")
	}
}

func TestBuildInternalCASecret_MissingFile(t *testing.T) {
	_, err := BuildInternalCASecret("/nonexistent/ca.pem")
	if err == nil {
		t.Error("expected error for missing CA file")
	}
}

func TestBuildInternalCASecret_EmptyPath(t *testing.T) {
	_, err := BuildInternalCASecret("")
	if err == nil {
		t.Error("expected error for empty CA path")
	}
}

func TestBuildPublicServerSecret(t *testing.T) {
	tmpDir := t.TempDir()

	certPath := tmpDir + "/public-cert.pem"
	keyPath := tmpDir + "/public-key.pem"

	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	paths := CertPaths{
		CertFile: certPath,
		KeyFile:  keyPath,
	}

	secret, err := BuildPublicServerSecret(paths)
	if err != nil {
		t.Fatalf("BuildPublicServerSecret failed: %v", err)
	}

	if secret.Name != PublicServerCert {
		t.Errorf("expected name %s, got %s", PublicServerCert, secret.Name)
	}

	tlsCert := secret.GetTlsCertificate()
	if tlsCert == nil {
		t.Fatal("expected TlsCertificate, got nil")
	}
}

func TestBuildPublicServerSecret_MissingFiles(t *testing.T) {
	paths := CertPaths{
		CertFile: "/nonexistent/cert.pem",
		KeyFile:  "/nonexistent/key.pem",
	}

	_, err := BuildPublicServerSecret(paths)
	if err == nil {
		t.Error("expected error for missing files")
	}
}

func TestBuildAllSecrets(t *testing.T) {
	tmpDir := t.TempDir()

	// Internal cert and key
	internalCertPath := tmpDir + "/internal-cert.pem"
	internalKeyPath := tmpDir + "/internal-key.pem"
	if err := os.WriteFile(internalCertPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write internal cert: %v", err)
	}
	if err := os.WriteFile(internalKeyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write internal key: %v", err)
	}

	// CA file
	caPath := tmpDir + "/ca.pem"
	if err := os.WriteFile(caPath, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	// Public cert and key
	publicCertPath := tmpDir + "/public-cert.pem"
	publicKeyPath := tmpDir + "/public-key.pem"
	if err := os.WriteFile(publicCertPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write public cert: %v", err)
	}
	if err := os.WriteFile(publicKeyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write public key: %v", err)
	}

	internalPaths := CertPaths{
		CertFile: internalCertPath,
		KeyFile:  internalKeyPath,
	}

	publicPaths := CertPaths{
		CertFile: publicCertPath,
		KeyFile:  publicKeyPath,
	}

	secrets, err := BuildAllSecrets(internalPaths, publicPaths, caPath)
	if err != nil {
		t.Fatalf("BuildAllSecrets failed: %v", err)
	}

	// Should have 3 secrets: internal-server-cert, internal-ca-bundle, public-server-cert
	if len(secrets) != 3 {
		t.Errorf("expected 3 secrets, got %d", len(secrets))
	}

	// Check internal server cert
	if _, ok := secrets[InternalServerCert]; !ok {
		t.Errorf("missing %s", InternalServerCert)
	}

	// Check internal CA bundle
	if _, ok := secrets[InternalCABundle]; !ok {
		t.Errorf("missing %s", InternalCABundle)
	}

	// Check public server cert
	if _, ok := secrets[PublicServerCert]; !ok {
		t.Errorf("missing %s", PublicServerCert)
	}
}

func TestBuildAllSecrets_InternalOnly(t *testing.T) {
	tmpDir := t.TempDir()

	internalCertPath := tmpDir + "/internal-cert.pem"
	internalKeyPath := tmpDir + "/internal-key.pem"
	caPath := tmpDir + "/ca.pem"

	if err := os.WriteFile(internalCertPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(internalKeyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	if err := os.WriteFile(caPath, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	internalPaths := CertPaths{
		CertFile: internalCertPath,
		KeyFile:  internalKeyPath,
	}

	// No public cert
	publicPaths := CertPaths{}

	secrets, err := BuildAllSecrets(internalPaths, publicPaths, caPath)
	if err != nil {
		t.Fatalf("BuildAllSecrets failed: %v", err)
	}

	// Should have 2 secrets: internal-server-cert, internal-ca-bundle
	if len(secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(secrets))
	}

	if _, ok := secrets[InternalServerCert]; !ok {
		t.Error("missing internal-server-cert")
	}

	if _, ok := secrets[InternalCABundle]; !ok {
		t.Error("missing internal-ca-bundle")
	}

	// Public cert should not be present
	if _, ok := secrets[PublicServerCert]; ok {
		t.Error("public-server-cert should not be present")
	}
}

func TestBuildAllSecrets_NoSecrets(t *testing.T) {
	internalPaths := CertPaths{}
	publicPaths := CertPaths{}

	_, err := BuildAllSecrets(internalPaths, publicPaths, "")
	if err == nil {
		t.Error("expected error when no secrets can be built")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file
	testFile := tmpDir + "/test.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	// Test existing file
	if !fileExists(testFile) {
		t.Error("fileExists returned false for existing file")
	}

	// Test non-existent file
	if fileExists("/nonexistent/file.txt") {
		t.Error("fileExists returned true for non-existent file")
	}

	// Test directory (should return false)
	if fileExists(tmpDir) {
		t.Error("fileExists returned true for directory")
	}
}
