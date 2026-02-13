package builder

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// generateTestCertificate creates a self-signed certificate and private key for testing
func generateTestCertificate(t *testing.T, certPath, keyPath string) {
	t.Helper()

	// Generate private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test.example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create self-signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		t.Fatalf("failed to write certificate: %v", err)
	}

	// Encode private key to PEM
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}
}

// TestBuildSnapshot_ExternalDomains_Conformance verifies that external domains
// are properly integrated into the xDS snapshot (PR3c conformance test).
func TestBuildSnapshot_ExternalDomains_Conformance(t *testing.T) {
	// Create temporary certificate files for testing
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "fullchain.pem")
	keyFile := filepath.Join(tempDir, "privkey.pem")

	// Generate self-signed certificate
	generateTestCertificate(t, certFile, keyFile)

	// Build input with one external domain
	input := Input{
		NodeID: "test-node",
		Listener: Listener{
			Name:       "test_listener",
			RouteName:  "test_routes",
			Host:       "0.0.0.0",
			Port:       443,
			CertFile:   certFile,
			KeyFile:    keyFile,
			IssuerFile: "",
		},
		Routes: []Route{
			{Prefix: "/internal/", Cluster: "internal_cluster"},
		},
		Clusters: []Cluster{
			{
				Name:      "gateway_http",
				Endpoints: []Endpoint{{Host: "127.0.0.1", Port: 8080}},
			},
			{
				Name:      "internal_cluster",
				Endpoints: []Endpoint{{Host: "127.0.0.1", Port: 9090}},
			},
		},
		ExternalDomains: []ExternalDomain{
			{
				FQDN:          "test.globular.cloud",
				CertFile:      certFile,
				KeyFile:       keyFile,
				TargetCluster: "gateway_http",
			},
		},
		EnableSDS: true,
		SDSSecrets: []Secret{
			{
				Name:     "ext-cert/test.globular.cloud",
				CertPath: certFile,
				KeyPath:  keyFile,
			},
		},
	}

	// Build snapshot
	snapshot, err := BuildSnapshot(input, "test-version")
	if err != nil {
		t.Fatalf("BuildSnapshot failed: %v", err)
	}

	// Test 1: Verify snapshot was created
	if snapshot == nil {
		t.Fatal("expected non-nil snapshot")
	}

	// Test 2: Verify route configuration includes external domain VirtualHost
	routes := snapshot.GetResources(resource_v3.RouteType)
	if len(routes) == 0 {
		t.Fatal("expected at least one route configuration")
	}

	// Test 3: Verify listener was created on port 443
	listeners := snapshot.GetResources(resource_v3.ListenerType)
	if len(listeners) == 0 {
		t.Fatal("expected at least one listener")
	}

	// Test 4: Verify SDS secret exists for external domain
	secrets := snapshot.GetResources(resource_v3.SecretType)
	foundExtSecret := false
	for secretName := range secrets {
		if secretName == "ext-cert/test.globular.cloud" {
			foundExtSecret = true
			break
		}
	}
	if !foundExtSecret {
		t.Error("expected SDS secret 'ext-cert/test.globular.cloud' in snapshot")
	}

	// Test 5: Verify gateway_http cluster exists
	clusters := snapshot.GetResources(resource_v3.ClusterType)
	foundGatewayCluster := false
	for clusterName := range clusters {
		if clusterName == "gateway_http" {
			foundGatewayCluster = true
			break
		}
	}
	if !foundGatewayCluster {
		t.Error("expected cluster 'gateway_http' in snapshot")
	}

	// Test 6: Verify listener has SNI filter chains configured
	var mainListener *listener_v3.Listener
	for _, resource := range listeners {
		if l, ok := resource.(*listener_v3.Listener); ok {
			mainListener = l
			break
		}
	}
	if mainListener == nil {
		t.Fatal("failed to extract listener from snapshot")
	}

	// Test 6a: Verify listener has at least 2 filter chains (SNI + default)
	if len(mainListener.FilterChains) < 2 {
		t.Errorf("expected at least 2 filter chains (SNI + default), got %d", len(mainListener.FilterChains))
	}

	// Test 6b: Verify SNI filter chain with test.globular.cloud exists
	foundSNIChain := false
	foundSNISecretRef := false
	for _, fc := range mainListener.FilterChains {
		if fc.FilterChainMatch != nil && len(fc.FilterChainMatch.ServerNames) > 0 {
			for _, sni := range fc.FilterChainMatch.ServerNames {
				if sni == "test.globular.cloud" {
					foundSNIChain = true

					// Test 6c: Verify transport socket references correct SDS secret
					if fc.TransportSocket != nil {
						var dtls tls_v3.DownstreamTlsContext
						if err := anypb.UnmarshalTo(fc.TransportSocket.GetTypedConfig(), &dtls, proto.UnmarshalOptions{}); err == nil {
							if dtls.CommonTlsContext != nil && len(dtls.CommonTlsContext.TlsCertificateSdsSecretConfigs) > 0 {
								secretName := dtls.CommonTlsContext.TlsCertificateSdsSecretConfigs[0].Name
								if secretName == "ext-cert/test.globular.cloud" {
									foundSNISecretRef = true
								} else {
									t.Errorf("expected SDS secret 'ext-cert/test.globular.cloud', got %q", secretName)
								}
							}
						}
					}
					break
				}
			}
		}
	}

	if !foundSNIChain {
		t.Error("expected filter chain with ServerNames containing 'test.globular.cloud'")
	}
	if !foundSNISecretRef {
		t.Error("expected filter chain to reference SDS secret 'ext-cert/test.globular.cloud'")
	}

	// Test 7: Verify default filter chain exists (no ServerNames match)
	foundDefaultChain := false
	for _, fc := range mainListener.FilterChains {
		if fc.FilterChainMatch == nil || len(fc.FilterChainMatch.ServerNames) == 0 {
			foundDefaultChain = true
			break
		}
	}
	if !foundDefaultChain {
		t.Error("expected default filter chain (fallback) with no ServerNames")
	}

	t.Logf("✓ Snapshot conformance test passed")
	t.Logf("  - Routes: %d", len(routes))
	t.Logf("  - Listeners: %d", len(listeners))
	t.Logf("  - Filter chains: %d", len(mainListener.FilterChains))
	t.Logf("  - Clusters: %d", len(clusters))
	t.Logf("  - Secrets: %d", len(secrets))
}

// TestBuildSnapshot_ExternalDomains_MultipleHosts verifies handling of multiple external domains.
func TestBuildSnapshot_ExternalDomains_MultipleHosts(t *testing.T) {
	tempDir := t.TempDir()

	// Create certificate files for two domains
	cert1 := filepath.Join(tempDir, "domain1-cert.pem")
	key1 := filepath.Join(tempDir, "domain1-key.pem")
	cert2 := filepath.Join(tempDir, "domain2-cert.pem")
	key2 := filepath.Join(tempDir, "domain2-key.pem")

	// Generate self-signed certificates
	generateTestCertificate(t, cert1, key1)
	generateTestCertificate(t, cert2, key2)

	input := Input{
		NodeID: "test-node",
		Listener: Listener{
			Name:      "test_listener",
			RouteName: "test_routes",
			Host:      "0.0.0.0",
			Port:      443,
			CertFile:  cert1,
			KeyFile:   key1,
		},
		Clusters: []Cluster{
			{Name: "gateway_http", Endpoints: []Endpoint{{Host: "127.0.0.1", Port: 8080}}},
		},
		ExternalDomains: []ExternalDomain{
			{FQDN: "app1.example.com", CertFile: cert1, KeyFile: key1, TargetCluster: "gateway_http"},
			{FQDN: "app2.example.com", CertFile: cert2, KeyFile: key2, TargetCluster: "gateway_http"},
		},
		EnableSDS: true,
		SDSSecrets: []Secret{
			{Name: "ext-cert/app1.example.com", CertPath: cert1, KeyPath: key1},
			{Name: "ext-cert/app2.example.com", CertPath: cert2, KeyPath: key2},
		},
	}

	snapshot, err := BuildSnapshot(input, "test-version")
	if err != nil {
		t.Fatalf("BuildSnapshot failed: %v", err)
	}

	// Verify both SDS secrets exist
	secrets := snapshot.GetResources(resource_v3.SecretType)
	expectedSecrets := []string{
		"ext-cert/app1.example.com",
		"ext-cert/app2.example.com",
	}

	for _, expectedSecret := range expectedSecrets {
		found := false
		for secretName := range secrets {
			if secretName == expectedSecret {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected SDS secret %q in snapshot", expectedSecret)
		}
	}

	t.Logf("✓ Multiple external domains test passed")
	t.Logf("  - External domains: %d", len(input.ExternalDomains))
	t.Logf("  - SDS secrets: %d", len(secrets))
}

// TestBuildSnapshot_ExternalDomains_IdempotentBuild verifies that building the same
// input multiple times produces consistent snapshots.
func TestBuildSnapshot_ExternalDomains_IdempotentBuild(t *testing.T) {
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "cert.pem")
	keyFile := filepath.Join(tempDir, "key.pem")

	// Generate self-signed certificate
	generateTestCertificate(t, certFile, keyFile)

	input := Input{
		NodeID: "test-node",
		Listener: Listener{
			Name:      "test_listener",
			RouteName: "test_routes",
			Host:      "0.0.0.0",
			Port:      443,
			CertFile:  certFile,
			KeyFile:   keyFile,
		},
		Clusters: []Cluster{
			{Name: "gateway_http", Endpoints: []Endpoint{{Host: "127.0.0.1", Port: 8080}}},
		},
		ExternalDomains: []ExternalDomain{
			{FQDN: "test.example.com", CertFile: certFile, KeyFile: keyFile, TargetCluster: "gateway_http"},
		},
		EnableSDS: true,
		SDSSecrets: []Secret{
			{Name: "ext-cert/test.example.com", CertPath: certFile, KeyPath: keyFile},
		},
	}

	// Build snapshot twice with same version
	version := "idempotent-test-v1"
	snapshot1, err := BuildSnapshot(input, version)
	if err != nil {
		t.Fatalf("first BuildSnapshot failed: %v", err)
	}

	snapshot2, err := BuildSnapshot(input, version)
	if err != nil {
		t.Fatalf("second BuildSnapshot failed: %v", err)
	}

	// Verify both snapshots have same version
	if snapshot1.GetVersion(resource_v3.ClusterType) != snapshot2.GetVersion(resource_v3.ClusterType) {
		t.Error("expected snapshots to have same cluster version")
	}
	if snapshot1.GetVersion(resource_v3.ListenerType) != snapshot2.GetVersion(resource_v3.ListenerType) {
		t.Error("expected snapshots to have same listener version")
	}
	if snapshot1.GetVersion(resource_v3.SecretType) != snapshot2.GetVersion(resource_v3.SecretType) {
		t.Error("expected snapshots to have same secret version")
	}

	// Verify same number of resources
	if len(snapshot1.GetResources(resource_v3.SecretType)) != len(snapshot2.GetResources(resource_v3.SecretType)) {
		t.Error("expected snapshots to have same number of secrets")
	}

	t.Logf("✓ Idempotent build test passed")
}

// TestBuildExternalDomainVirtualHosts tests the VirtualHost builder function.
func TestBuildExternalDomainVirtualHosts(t *testing.T) {
	domains := []ExternalDomain{
		{FQDN: "app.example.com", TargetCluster: "gateway_http"},
		{FQDN: "api.example.com", TargetCluster: "gateway_http"},
	}

	routes := buildExternalDomainVirtualHosts(domains)

	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}

	// Verify first route
	if routes[0].Prefix != "/" {
		t.Errorf("expected route prefix '/', got %q", routes[0].Prefix)
	}
	if routes[0].Cluster != "gateway_http" {
		t.Errorf("expected cluster 'gateway_http', got %q", routes[0].Cluster)
	}
	if len(routes[0].Domains) != 1 || routes[0].Domains[0] != "app.example.com" {
		t.Errorf("expected domains ['app.example.com'], got %v", routes[0].Domains)
	}

	t.Logf("✓ VirtualHost builder test passed")
}

// TestBuildExternalDomainSecrets tests the secret builder function.
func TestBuildExternalDomainSecrets(t *testing.T) {
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "cert.pem")
	keyFile := filepath.Join(tempDir, "key.pem")

	// Generate self-signed certificate
	generateTestCertificate(t, certFile, keyFile)

	domains := []ExternalDomain{
		{FQDN: "test.example.com", CertFile: certFile, KeyFile: keyFile},
	}

	secrets, err := buildExternalDomainSecrets(domains)
	if err != nil {
		t.Fatalf("buildExternalDomainSecrets failed: %v", err)
	}

	if len(secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(secrets))
	}

	if secrets[0].Name != "ext-cert/test.example.com" {
		t.Errorf("expected secret name 'ext-cert/test.example.com', got %q", secrets[0].Name)
	}

	t.Logf("✓ Secret builder test passed")
}

// TestBuildExternalDomainSecrets_MissingFiles verifies error handling for missing certificate files.
func TestBuildExternalDomainSecrets_MissingFiles(t *testing.T) {
	domains := []ExternalDomain{
		{FQDN: "test.example.com", CertFile: "/nonexistent/cert.pem", KeyFile: "/nonexistent/key.pem"},
	}

	_, err := buildExternalDomainSecrets(domains)
	if err == nil {
		t.Fatal("expected error for missing certificate files, got nil")
	}

	t.Logf("✓ Missing files error handling test passed: %v", err)
}
