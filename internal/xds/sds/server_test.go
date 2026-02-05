package sds

import (
	"context"
	"os"
	"testing"
	"time"

	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/globulario/Globular/internal/controlplane"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Test certificate data (same as secret_test.go)
const testCert = `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU6jSMA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNVBAMMB3Rl
c3QtY2EwHhcNMjQwMTAxMDAwMDAwWhcNMjUwMTAxMDAwMDAwWjASMRAwDgYDVQQD
DAd0ZXN0LWNhMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJhPHhITqQbPklG3ib
SNKcz5RB7aFpVSYKL6vxKLQE6zxMkTx0l1N8FqwL5xQ9l7FZgQmCgIaF0OVc5GmC
Ep8CAwEAATANBgkqhkiG9w0BAQsFAANBAGO6L0Qx9pMd5H2vqQKDyT8HVqKJDxCh
4xP2qQtmR7E7gK7xQ5F2L4L0Q9hVFE9pNqHVXL1pQqJ3xC8VqN4L0pE=
-----END CERTIFICATE-----`

const testKey = `-----BEGIN PRIVATE KEY-----
MIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEA0smE8eEhOpBs+SUb
eJtI0pzPlEHtoWlVJgovq/EotATrPEyRPHSXU3wWrAvnFD2XsVmBCYKAhoXQ5Vzk
aYISnwIDAQABAkA0qOK+oE6EFOkXLdLQaH1PwX9F3xQmxKTY3Q5L7T4LxPSqYc8Z
kK3D8A1HqP5RXJ9fC1qPQxLqYBbL6L9Q5xYhAiEA7VxL5QqH2L7F0Q9L4Q5L7Q8L
5Q6L7Q7L5Q5L7Q4L5QkCIQDk5L7Q9L5Q8L7Q7L5Q6L7Q5L4Q3L2Q1L0QzLyQxLwwJ
AiAL5Q7L5Q6L7Q5L4Q3L2Q1L0QzLyQxLwQvLuQtLsQrCQIgQqQpQoQnQmQlQkQjQ
iQhQgQfQeQdQcQbQaQZQYQXQIhBAkEA5L7Q9L5Q8L7Q7L5Q6L7Q5L4Q3L2Q1L0Qz
LyQxLwQvLuQtLsQrLqQpLoQnLmQlLkQjLiQhLgQfLeQdLcQbLaQ=
-----END PRIVATE KEY-----`

const testCA = `-----BEGIN CERTIFICATE-----
MIIBkjCCATogAwIBAgIJAKHHCgVZU6jTMA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNV
BAMMB3Rlc3QtY2EwHhcNMjQwMTAxMDAwMDAwWhcNMzQwMTAxMDAwMDAwWjASMRAw
DgYDVQQDDAd0ZXN0LWNhMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJhPHhITqQ
bPklG3ibSNKcz5RB7aFpVSYKL6vxKLQE6zxMkTx0l1N8FqwL5xQ9l7FZgQmCgIaF
0OVc5GmCEp8CAwEAAaMQMA4wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAANB
AGO6L0Qx9pMd5H2vqQKDyT8HVqKJDxCh4xP2qQtmR7E7gK7xQ5F2L4L0Q9hVFE9p
NqHVXL1pQqJ3xC8VqN4L0pE=
-----END CERTIFICATE-----`

func TestNewServer(t *testing.T) {
	srv := NewServer()
	if srv == nil {
		t.Fatal("NewServer returned nil")
	}

	if srv.secrets == nil {
		t.Error("secrets map not initialized")
	}

	if srv.cache == nil {
		t.Error("snapshot cache not initialized")
	}

	if srv.version != "v0" {
		t.Errorf("initial version should be v0, got %s", srv.version)
	}

	if len(srv.secrets) != 0 {
		t.Errorf("initial secrets should be empty, got %d", len(srv.secrets))
	}
}

func TestUpdateSecrets(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test certificate files
	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)

	// Build secrets
	secret1, err := controlplane.MakeSecret("test-cert-1", certPath, keyPath, "")
	if err != nil {
		t.Fatalf("MakeSecret failed: %v", err)
	}

	secret2, err := controlplane.MakeSecret("test-ca-1", "", "", caPath)
	if err != nil {
		t.Fatalf("MakeSecret failed: %v", err)
	}

	secrets := map[string]*tls_v3.Secret{
		"test-cert-1": secret1,
		"test-ca-1":   secret2,
	}

	srv := NewServer()

	// Test successful update
	if err := srv.UpdateSecrets(secrets); err != nil {
		t.Fatalf("UpdateSecrets failed: %v", err)
	}

	// Verify secrets stored
	if len(srv.secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(srv.secrets))
	}

	storedSecret1, ok := srv.GetSecret("test-cert-1")
	if !ok {
		t.Error("test-cert-1 not found after update")
	}
	if storedSecret1.Name != "test-cert-1" {
		t.Errorf("expected name test-cert-1, got %s", storedSecret1.Name)
	}

	// Verify version changed
	if srv.GetVersion() == "v0" {
		t.Error("version should have changed from v0")
	}
}

func TestUpdateSecrets_EmptyMap(t *testing.T) {
	srv := NewServer()

	err := srv.UpdateSecrets(map[string]*tls_v3.Secret{})
	if err == nil {
		t.Error("expected error for empty secret map")
	}
}

func TestGetSecret(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath, _ := createTestCertFiles(t, tmpDir)

	secret, err := controlplane.MakeSecret("test-cert", certPath, keyPath, "")
	if err != nil {
		t.Fatalf("MakeSecret failed: %v", err)
	}

	srv := NewServer()
	srv.secrets["test-cert"] = secret
	srv.version = "v1"

	// Test existing secret
	retrieved, ok := srv.GetSecret("test-cert")
	if !ok {
		t.Error("GetSecret returned false for existing secret")
	}
	if retrieved.Name != "test-cert" {
		t.Errorf("expected name test-cert, got %s", retrieved.Name)
	}

	// Test non-existent secret
	_, ok = srv.GetSecret("non-existent")
	if ok {
		t.Error("GetSecret returned true for non-existent secret")
	}
}

func TestGetSecretNames(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)

	secret1, _ := controlplane.MakeSecret("cert-1", certPath, keyPath, "")
	secret2, _ := controlplane.MakeSecret("cert-2", certPath, keyPath, "")
	secret3, _ := controlplane.MakeSecret("ca-1", "", "", caPath)

	srv := NewServer()
	srv.secrets = map[string]*tls_v3.Secret{
		"cert-1": secret1,
		"cert-2": secret2,
		"ca-1":   secret3,
	}

	names := srv.GetSecretNames()
	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}

	// Check all names present
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	if !nameSet["cert-1"] || !nameSet["cert-2"] || !nameSet["ca-1"] {
		t.Errorf("missing expected names, got: %v", names)
	}
}

func TestFetchSecrets(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)

	secret1, _ := controlplane.MakeSecret("test-cert", certPath, keyPath, "")
	secret2, _ := controlplane.MakeSecret("test-ca", "", "", caPath)

	srv := NewServer()
	srv.secrets = map[string]*tls_v3.Secret{
		"test-cert": secret1,
		"test-ca":   secret2,
	}
	srv.version = "v1"

	// Test requesting specific secret
	req := &discovery_v3.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret",
		ResourceNames: []string{"test-cert"},
	}

	resp, err := srv.FetchSecrets(context.Background(), req)
	if err != nil {
		t.Fatalf("FetchSecrets failed: %v", err)
	}

	if len(resp.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(resp.Resources))
	}

	if resp.VersionInfo != "v1" {
		t.Errorf("expected version v1, got %s", resp.VersionInfo)
	}
}

func TestFetchSecrets_AllSecrets(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)

	secret1, _ := controlplane.MakeSecret("test-cert", certPath, keyPath, "")
	secret2, _ := controlplane.MakeSecret("test-ca", "", "", caPath)

	srv := NewServer()
	srv.secrets = map[string]*tls_v3.Secret{
		"test-cert": secret1,
		"test-ca":   secret2,
	}
	srv.version = "v2"

	// Test requesting all secrets (empty ResourceNames)
	req := &discovery_v3.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret",
		ResourceNames: []string{},
	}

	resp, err := srv.FetchSecrets(context.Background(), req)
	if err != nil {
		t.Fatalf("FetchSecrets failed: %v", err)
	}

	if len(resp.Resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(resp.Resources))
	}

	if resp.VersionInfo != "v2" {
		t.Errorf("expected version v2, got %s", resp.VersionInfo)
	}
}

func TestFetchSecrets_NotFound(t *testing.T) {
	srv := NewServer()
	srv.secrets = map[string]*tls_v3.Secret{}
	srv.version = "v1"

	req := &discovery_v3.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret",
		ResourceNames: []string{"non-existent"},
	}

	resp, err := srv.FetchSecrets(context.Background(), req)
	if err != nil {
		t.Fatalf("FetchSecrets failed: %v", err)
	}

	// Should return empty resources list, not error
	if len(resp.Resources) != 0 {
		t.Errorf("expected 0 resources for non-existent secret, got %d", len(resp.Resources))
	}
}

func TestComputeVersion(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath, _ := createTestCertFiles(t, tmpDir)

	secret, _ := controlplane.MakeSecret("test", certPath, keyPath, "")

	srv := NewServer()

	// Test version computation
	version1, err := srv.computeVersion(map[string]*tls_v3.Secret{"test": secret})
	if err != nil {
		t.Fatalf("computeVersion failed: %v", err)
	}

	if len(version1) == 0 {
		t.Error("version should not be empty")
	}

	if len(version1) < 4 || version1[:4] != "sds-" {
		t.Errorf("version should start with 'sds-', got %s", version1)
	}

	// Same secret should produce same version
	version2, err := srv.computeVersion(map[string]*tls_v3.Secret{"test": secret})
	if err != nil {
		t.Fatalf("computeVersion failed: %v", err)
	}

	if version1 != version2 {
		t.Error("identical secrets produced different versions")
	}
}

func TestComputeVersion_EmptyMap(t *testing.T) {
	srv := NewServer()

	version, err := srv.computeVersion(map[string]*tls_v3.Secret{})
	if err != nil {
		t.Fatalf("computeVersion failed: %v", err)
	}

	if version != "v0" {
		t.Errorf("empty map should produce v0, got %s", version)
	}
}

func TestVersionChangesOnUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	certPath1, keyPath1, _ := createTestCertFiles(t, tmpDir)

	secret1, _ := controlplane.MakeSecret("test", certPath1, keyPath1, "")

	srv := NewServer()

	// First update
	err := srv.UpdateSecrets(map[string]*tls_v3.Secret{"test": secret1})
	if err != nil {
		t.Fatalf("UpdateSecrets failed: %v", err)
	}
	version1 := srv.GetVersion()

	// Create different certificate
	certPath2, keyPath2, _ := createDifferentTestCertFiles(t, tmpDir)
	secret2, _ := controlplane.MakeSecret("test", certPath2, keyPath2, "")

	// Second update with different cert
	err = srv.UpdateSecrets(map[string]*tls_v3.Secret{"test": secret2})
	if err != nil {
		t.Fatalf("UpdateSecrets failed: %v", err)
	}
	version2 := srv.GetVersion()

	// Versions should differ
	if version1 == version2 {
		t.Error("version should change when secret content changes")
	}
}

func TestDeltaSecrets_Unimplemented(t *testing.T) {
	srv := NewServer()

	err := srv.DeltaSecrets(nil)
	if err == nil {
		t.Error("DeltaSecrets should return unimplemented error")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("error should be gRPC status")
	}

	if st.Code() != codes.Unimplemented {
		t.Errorf("expected Unimplemented code, got %v", st.Code())
	}
}

func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath, _ := createTestCertFiles(t, tmpDir)

	secret, _ := controlplane.MakeSecret("test", certPath, keyPath, "")
	secrets := map[string]*tls_v3.Secret{"test": secret}

	srv := NewServer()
	if err := srv.UpdateSecrets(secrets); err != nil {
		t.Fatalf("UpdateSecrets failed: %v", err)
	}

	// Concurrent reads should not race
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				srv.GetSecret("test")
				srv.GetVersion()
				srv.GetSecretNames()
			}
			done <- true
		}()
	}

	// Wait for all goroutines with timeout
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("concurrent access test timed out")
		}
	}
}

// Helper function to create test certificate files
func createTestCertFiles(t *testing.T, dir string) (certPath, keyPath, caPath string) {
	t.Helper()

	certPath = dir + "/cert.pem"
	keyPath = dir + "/key.pem"
	caPath = dir + "/ca.pem"

	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	if err := os.WriteFile(caPath, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	return certPath, keyPath, caPath
}

// Helper function to create different test certificate files
func createDifferentTestCertFiles(t *testing.T, dir string) (certPath, keyPath, caPath string) {
	t.Helper()

	certPath = dir + "/cert2.pem"
	keyPath = dir + "/key2.pem"
	caPath = dir + "/ca2.pem"

	// Slightly modified cert (add space at end)
	if err := os.WriteFile(certPath, []byte(testCert+" "), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	if err := os.WriteFile(caPath, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	return certPath, keyPath, caPath
}
