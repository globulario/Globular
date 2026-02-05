package watchers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/globulario/Globular/internal/controlplane"
	"github.com/globulario/Globular/internal/xds/builder"
	"github.com/globulario/Globular/internal/xds/server"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Test certificate data (self-signed, for testing only)
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

// TestSnapshotContainsSecrets verifies that snapshots include Secret resources when EnableSDS is true.
func TestSnapshotContainsSecrets(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test certificate files
	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")
	caPath := filepath.Join(tmpDir, "ca.pem")

	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	if err := os.WriteFile(caPath, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	// Build input with SDS enabled
	input := builder.Input{
		NodeID:    "test-node",
		EnableSDS: true,
		SDSSecrets: []builder.Secret{
			{
				Name:     "internal-server-cert",
				CertPath: certPath,
				KeyPath:  keyPath,
			},
			{
				Name:   "internal-ca-bundle",
				CAPath: caPath,
			},
		},
		Listener: builder.Listener{
			Name:      "test-listener",
			Port:      443,
			RouteName: "test-routes",
			CertFile:  certPath,
			KeyFile:   keyPath,
		},
		Routes: []builder.Route{
			{Prefix: "/", Cluster: "test-cluster"},
		},
		Clusters: []builder.Cluster{
			{
				Name:      "test-cluster",
				Endpoints: []builder.Endpoint{{Host: "localhost", Port: 8080}},
			},
		},
	}

	// Build snapshot
	snapshot, err := builder.BuildSnapshot(input, "v1")
	if err != nil {
		t.Fatalf("BuildSnapshot failed: %v", err)
	}

	// Verify snapshot contains Secret resources
	secrets := snapshot.GetResources(resource_v3.SecretType)
	if len(secrets) == 0 {
		t.Fatal("snapshot should contain Secret resources when EnableSDS=true")
	}

	if len(secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(secrets))
	}

	// Verify secret names
	secretNames := make(map[string]bool)
	for name := range secrets {
		secretNames[name] = true
	}

	if !secretNames["internal-server-cert"] {
		t.Error("snapshot missing internal-server-cert")
	}
	if !secretNames["internal-ca-bundle"] {
		t.Error("snapshot missing internal-ca-bundle")
	}

	t.Logf("✓ Snapshot contains %d secrets: %v", len(secrets), getSecretNames(secrets))
}

// TestSnapshotWithoutSDS verifies that snapshots don't include secrets when SDS is disabled.
func TestSnapshotWithoutSDS(t *testing.T) {
	tmpDir := t.TempDir()

	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	// Build input with SDS disabled (file-based TLS)
	input := builder.Input{
		NodeID:    "test-node",
		EnableSDS: false, // Explicitly disabled
		Listener: builder.Listener{
			Name:      "test-listener",
			Port:      443,
			RouteName: "test-routes",
			CertFile:  certPath,
			KeyFile:   keyPath,
		},
		Routes: []builder.Route{
			{Prefix: "/", Cluster: "test-cluster"},
		},
		Clusters: []builder.Cluster{
			{
				Name:      "test-cluster",
				Endpoints: []builder.Endpoint{{Host: "localhost", Port: 8080}},
			},
		},
	}

	// Build snapshot
	snapshot, err := builder.BuildSnapshot(input, "v1")
	if err != nil {
		t.Fatalf("BuildSnapshot failed: %v", err)
	}

	// Verify snapshot does NOT contain Secret resources
	secrets := snapshot.GetResources(resource_v3.SecretType)
	if len(secrets) != 0 {
		t.Errorf("snapshot should not contain secrets when EnableSDS=false, got %d", len(secrets))
	}

	t.Log("✓ Snapshot correctly excludes secrets when SDS disabled")
}

// TestListenerUsesSDS verifies that listeners have SDS secret configs when SDS is enabled.
func TestListenerUsesSDS(t *testing.T) {
	tmpDir := t.TempDir()

	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")

	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	// Build input with SDS enabled
	input := builder.Input{
		NodeID:    "test-node",
		EnableSDS: true,
		SDSSecrets: []builder.Secret{
			{
				Name:     "internal-server-cert",
				CertPath: certPath,
				KeyPath:  keyPath,
			},
		},
		Listener: builder.Listener{
			Name:      "test-listener",
			Port:      443,
			RouteName: "test-routes",
			CertFile:  certPath,
			KeyFile:   keyPath,
		},
		Routes: []builder.Route{
			{Prefix: "/", Cluster: "test-cluster"},
		},
		Clusters: []builder.Cluster{
			{
				Name:      "test-cluster",
				Endpoints: []builder.Endpoint{{Host: "localhost", Port: 8080}},
			},
		},
	}

	// Build snapshot
	snapshot, err := builder.BuildSnapshot(input, "v1")
	if err != nil {
		t.Fatalf("BuildSnapshot failed: %v", err)
	}

	// Get listener from snapshot
	listeners := snapshot.GetResources(resource_v3.ListenerType)
	if len(listeners) == 0 {
		t.Fatal("snapshot should contain listeners")
	}

	// Verify listener has SDS config
	for _, resource := range listeners {
		listener, ok := resource.(*listener_v3.Listener)
		if !ok {
			continue
		}

		// Check filter chains for TLS config
		for _, chain := range listener.FilterChains {
			if chain.TransportSocket == nil {
				continue
			}

			// Unmarshal the transport socket config
			var tlsContext tls_v3.DownstreamTlsContext
			if err := chain.TransportSocket.GetTypedConfig().UnmarshalTo(&tlsContext); err != nil {
				t.Logf("failed to unmarshal TLS context: %v", err)
				continue
			}

			// Verify SDS secret configs are present
			if tlsContext.CommonTlsContext == nil {
				t.Error("listener TLS context missing CommonTlsContext")
				continue
			}

			sdsConfigs := tlsContext.CommonTlsContext.TlsCertificateSdsSecretConfigs
			if len(sdsConfigs) == 0 {
				t.Error("listener should have tls_certificate_sds_secret_configs when SDS enabled")
			} else {
				t.Logf("✓ Listener has SDS config: secret_name=%s", sdsConfigs[0].Name)

				// Verify secret name
				if sdsConfigs[0].Name != "internal-server-cert" {
					t.Errorf("expected secret name 'internal-server-cert', got '%s'", sdsConfigs[0].Name)
				}

				// Verify ADS config
				if sdsConfigs[0].SdsConfig == nil {
					t.Error("SDS config missing sds_config")
				} else if sdsConfigs[0].SdsConfig.GetAds() == nil {
					t.Error("SDS config should use ADS")
				}
			}
		}
	}
}

// TestClusterUsesSDS verifies that clusters have SDS validation context when SDS is enabled.
func TestClusterUsesSDS(t *testing.T) {
	tmpDir := t.TempDir()

	certPath := filepath.Join(tmpDir, "cert.pem")
	keyPath := filepath.Join(tmpDir, "key.pem")
	caPath := filepath.Join(tmpDir, "ca.pem")

	if err := os.WriteFile(certPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	if err := os.WriteFile(caPath, []byte(testCA), 0644); err != nil {
		t.Fatalf("write CA: %v", err)
	}

	// Build input with SDS enabled and upstream TLS
	input := builder.Input{
		NodeID:    "test-node",
		EnableSDS: true,
		SDSSecrets: []builder.Secret{
			{
				Name:   "internal-ca-bundle",
				CAPath: caPath,
			},
		},
		Listener: builder.Listener{
			Name:      "test-listener",
			Port:      443,
			RouteName: "test-routes",
		},
		Routes: []builder.Route{
			{Prefix: "/", Cluster: "test-cluster"},
		},
		Clusters: []builder.Cluster{
			{
				Name:       "test-cluster",
				Endpoints:  []builder.Endpoint{{Host: "backend.example.com", Port: 8080}},
				ServerCert: certPath,
				KeyFile:    keyPath,
				CAFile:     caPath,
				SNI:        "backend.example.com",
			},
		},
	}

	// Build snapshot
	snapshot, err := builder.BuildSnapshot(input, "v1")
	if err != nil {
		t.Fatalf("BuildSnapshot failed: %v", err)
	}

	// Get clusters from snapshot
	clusters := snapshot.GetResources(resource_v3.ClusterType)
	if len(clusters) == 0 {
		t.Fatal("snapshot should contain clusters")
	}

	// Verify cluster has SDS validation context
	foundSDSCluster := false
	for _, resource := range clusters {
		cluster, ok := resource.(*cluster_v3.Cluster)
		if !ok {
			continue
		}

		if cluster.TransportSocket == nil {
			continue
		}

		// Unmarshal the transport socket config
		var tlsContext tls_v3.UpstreamTlsContext
		if err := cluster.TransportSocket.GetTypedConfig().UnmarshalTo(&tlsContext); err != nil {
			t.Logf("failed to unmarshal TLS context: %v", err)
			continue
		}

		if tlsContext.CommonTlsContext == nil {
			continue
		}

		// Check for SDS validation context
		sdsValidation := tlsContext.CommonTlsContext.GetValidationContextSdsSecretConfig()
		if sdsValidation != nil {
			foundSDSCluster = true
			t.Logf("✓ Cluster has SDS validation context: secret_name=%s", sdsValidation.Name)

			// Verify secret name
			if sdsValidation.Name != "internal-ca-bundle" {
				t.Errorf("expected secret name 'internal-ca-bundle', got '%s'", sdsValidation.Name)
			}

			// Verify ADS config
			if sdsValidation.SdsConfig == nil {
				t.Error("validation context SDS config missing sds_config")
			} else if sdsValidation.SdsConfig.GetAds() == nil {
				t.Error("validation context SDS config should use ADS")
			}
		}
	}

	if !foundSDSCluster {
		t.Error("no cluster found with SDS validation context config")
	}
}

// TestCertificateRotationTriggersSnapshot verifies that cert generation changes trigger snapshot updates.
func TestCertificateRotationTriggersSnapshot(t *testing.T) {
	// This test requires a real etcd instance - skip if not available
	if os.Getenv("ETCD_ENDPOINTS") == "" {
		t.Skip("Skipping: ETCD_ENDPOINTS not set (integration test)")
	}

	ctx := context.Background()

	// Create etcd client
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{os.Getenv("ETCD_ENDPOINTS")},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("etcd connection failed: %v", err)
	}
	defer etcdClient.Close()

	// Create test server
	xdsServer := server.New(nil, ctx)

	// Create watcher
	watcher := New(nil, xdsServer, "", "test-node", 1*time.Second, DownstreamTLSDisabled, "")
	watcher.SetEtcdClient(etcdClient)

	// Set initial certificate generation in etcd
	domain := "test.globular.internal"
	key := fmt.Sprintf("/globular/pki/bundles/%s", domain)

	bundle1 := map[string]interface{}{
		"generation": 1,
		"updated_ms": time.Now().UnixMilli(),
	}
	data1, _ := json.Marshal(bundle1)
	_, err = etcdClient.Put(ctx, key, string(data1))
	if err != nil {
		t.Fatalf("failed to set initial generation: %v", err)
	}

	// Check generation (should initialize)
	changed1 := watcher.checkCertificateGeneration(ctx)
	if changed1 {
		t.Error("first check should not detect change (initialization)")
	}

	if watcher.lastCertGeneration != 1 {
		t.Errorf("expected lastCertGeneration=1, got %d", watcher.lastCertGeneration)
	}

	// Update certificate generation
	bundle2 := map[string]interface{}{
		"generation": 2,
		"updated_ms": time.Now().UnixMilli(),
	}
	data2, _ := json.Marshal(bundle2)
	_, err = etcdClient.Put(ctx, key, string(data2))
	if err != nil {
		t.Fatalf("failed to update generation: %v", err)
	}

	// Check generation again (should detect change)
	changed2 := watcher.checkCertificateGeneration(ctx)
	if !changed2 {
		t.Error("second check should detect generation change")
	}

	if watcher.lastCertGeneration != 2 {
		t.Errorf("expected lastCertGeneration=2, got %d", watcher.lastCertGeneration)
	}

	t.Logf("✓ Certificate rotation detection works: 1 → 2")

	// Clean up
	_, _ = etcdClient.Delete(ctx, key)
}

// TestSecretContentHash verifies that secret hashes change when content changes.
func TestSecretContentHash(t *testing.T) {
	tmpDir := t.TempDir()

	// Create first certificate
	cert1Path := filepath.Join(tmpDir, "cert1.pem")
	key1Path := filepath.Join(tmpDir, "key1.pem")

	if err := os.WriteFile(cert1Path, []byte(testCert), 0644); err != nil {
		t.Fatalf("write cert1: %v", err)
	}
	if err := os.WriteFile(key1Path, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key1: %v", err)
	}

	// Build first secret
	secret1, err := controlplane.MakeSecret("test-cert", cert1Path, key1Path, "")
	if err != nil {
		t.Fatalf("MakeSecret failed: %v", err)
	}

	hash1, err := controlplane.HashSecret(secret1)
	if err != nil {
		t.Fatalf("HashSecret failed: %v", err)
	}

	// Create second certificate (different content)
	cert2Path := filepath.Join(tmpDir, "cert2.pem")
	key2Path := filepath.Join(tmpDir, "key2.pem")

	if err := os.WriteFile(cert2Path, []byte(testCert+" "), 0644); err != nil {
		t.Fatalf("write cert2: %v", err)
	}
	if err := os.WriteFile(key2Path, []byte(testKey), 0600); err != nil {
		t.Fatalf("write key2: %v", err)
	}

	// Build second secret
	secret2, err := controlplane.MakeSecret("test-cert", cert2Path, key2Path, "")
	if err != nil {
		t.Fatalf("MakeSecret failed: %v", err)
	}

	hash2, err := controlplane.HashSecret(secret2)
	if err != nil {
		t.Fatalf("HashSecret failed: %v", err)
	}

	// Verify hashes are different
	if hash1 == hash2 {
		t.Error("different certificate content should produce different hashes")
	}

	t.Logf("✓ Secret hashes differ when content changes:")
	t.Logf("  hash1: %s", hash1[:16])
	t.Logf("  hash2: %s", hash2[:16])
}

// Helper functions

func getSecretNames(secrets map[string]types.Resource) []string {
	names := make([]string, 0, len(secrets))
	for name := range secrets {
		names = append(names, name)
	}
	return names
}
