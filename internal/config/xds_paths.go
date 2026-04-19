package config

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	globconfig "github.com/globulario/services/golang/config"
)

// Canonical PKI locations for xDS mTLS.
const (
	xdsServerSubdir      = "xds/current"
	envoyXDSClientSubdir = "envoy-xds-client/current"
)

func runtimeRoot(runtimeConfigDir string) string {
	if strings.TrimSpace(runtimeConfigDir) != "" {
		return filepath.Clean(runtimeConfigDir)
	}
	return globconfig.GetStateRootDir()
}

// GetXDSServerCertPaths returns (cert, key) for the xDS server identity.
// Paths are stable across rotations and rooted at the canonical PKI directory.
func GetXDSServerCertPaths(runtimeConfigDir string) (certPath, keyPath string) {
	root := runtimeRoot(runtimeConfigDir)
	dir := filepath.Join(root, "pki", xdsServerSubdir)
	return filepath.Join(dir, "tls.crt"), filepath.Join(dir, "tls.key")
}

// GetEnvoyXDSClientCertPaths returns (cert, key) for the Envoy xDS client identity.
func GetEnvoyXDSClientCertPaths(runtimeConfigDir string) (certPath, keyPath string) {
	root := runtimeRoot(runtimeConfigDir)
	dir := filepath.Join(root, "pki", envoyXDSClientSubdir)
	return filepath.Join(dir, "tls.crt"), filepath.Join(dir, "tls.key")
}

// GetClusterCABundlePath returns the canonical CA bundle path shared across control plane components.
func GetClusterCABundlePath(runtimeConfigDir string) string {
	return filepath.Join(runtimeRoot(runtimeConfigDir), "pki", "ca.pem")
}

// EnsureXDSMTLSMaterials verifies the xDS server/client certificates exist and generates
// them from the local CA if they are missing. When insecureAllowed is true, the function
// is a no-op to avoid writing TLS assets during explicit plaintext development runs.
func EnsureXDSMTLSMaterials(runtimeConfigDir string, insecureAllowed bool) error {
	if insecureAllowed {
		return nil
	}

	caKeyPath, caCertPath, caBundlePath := CanonicalCAPaths(runtimeConfigDir)
	if !fileExists(caCertPath) {
		return fmt.Errorf("cluster CA cert missing: %s", caCertPath)
	}
	if !fileExists(caBundlePath) {
		if err := copyFile(caCertPath, caBundlePath, 0o444); err != nil {
			return fmt.Errorf("ensure CA bundle %s: %w", caBundlePath, err)
		}
	}

	hasCAKey := fileExists(caKeyPath)

	// Ensure xDS server and Envoy client identities exist.
	if hasCAKey {
		// Day-0: sign locally with CA key.
		if err := ensureXDSCertificate(true, caCertPath, caKeyPath, runtimeConfigDir); err != nil {
			return err
		}
		if err := ensureXDSCertificate(false, caCertPath, caKeyPath, runtimeConfigDir); err != nil {
			return err
		}
	} else {
		// Day-1: no CA key — request signing from gateway.
		if err := ensureXDSCertificateViaGateway(true, caCertPath, runtimeConfigDir); err != nil {
			return err
		}
		if err := ensureXDSCertificateViaGateway(false, caCertPath, runtimeConfigDir); err != nil {
			return err
		}
	}
	return nil
}

func ensureXDSCertificate(server bool, caCertPath, caKeyPath, runtimeConfigDir string) error {
	var certPath, keyPath, commonName string
	if server {
		certPath, keyPath = GetXDSServerCertPaths(runtimeConfigDir)
		commonName = "xds-server"
	} else {
		certPath, keyPath = GetEnvoyXDSClientCertPaths(runtimeConfigDir)
		commonName = "envoy-xds-client"
	}

	if fileExists(certPath) && fileExists(keyPath) {
		return nil
	}

	caCert, signer, err := loadCA(caCertPath, caKeyPath)
	if err != nil {
		return err
	}

	// v1 Conformance (INV-5.3): Use stable certificate identities
	// Internal xDS mTLS certificates should not depend on domain (routing label)
	// The commonName is the stable identity; cluster_domain SANs are optional for DNS resolution
	tpl, err := xdsCertTemplate(commonName, server)
	if err != nil {
		return err
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate xds key: %w", err)
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, caCert, &privKey.PublicKey, signer)
	if err != nil {
		return fmt.Errorf("sign %s certificate: %w", commonName, err)
	}

	if err := os.MkdirAll(filepath.Dir(certPath), 0o750); err != nil {
		return fmt.Errorf("create xds cert dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(keyPath), 0o700); err != nil {
		return fmt.Errorf("create xds key dir: %w", err)
	}

	if err := writePEMFile(keyPath, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)}, 0o400); err != nil {
		return err
	}
	if err := writePEMFile(certPath, &pem.Block{Type: "CERTIFICATE", Bytes: der}, 0o644); err != nil {
		return err
	}
	return nil
}

// ensureXDSCertificateViaGateway generates a keypair and CSR, sends it to the
// gateway's /sign_ca_certificate endpoint, and saves the signed cert.
// Used on Day-1 nodes where ca.key is not available.
func ensureXDSCertificateViaGateway(server bool, caCertPath, runtimeConfigDir string) error {
	var certPath, keyPath, commonName string
	if server {
		certPath, keyPath = GetXDSServerCertPaths(runtimeConfigDir)
		commonName = "xds-server"
	} else {
		certPath, keyPath = GetEnvoyXDSClientCertPaths(runtimeConfigDir)
		commonName = "envoy-xds-client"
	}

	if fileExists(certPath) && fileExists(keyPath) {
		return nil
	}

	// Discover gateway address from local config.
	gatewayAddr := discoverGatewayForCertSigning(runtimeConfigDir)
	if gatewayAddr == "" {
		return fmt.Errorf("cannot sign %s cert: no gateway address (Day-1 node without config.json?)", commonName)
	}

	// Generate keypair.
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate %s key: %w", commonName, err)
	}

	// Create CSR.
	csrTpl := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"globular.internal"},
		},
		DNSNames: []string{commonName, "localhost"},
	}
	if server {
		csrTpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, csrTpl, privKey)
	if err != nil {
		return fmt.Errorf("create %s CSR: %w", commonName, err)
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})

	// Send CSR to gateway.
	signedCert, err := signCSRViaGateway(gatewayAddr, csrPEM, caCertPath)
	if err != nil {
		return fmt.Errorf("sign %s cert via gateway: %w", commonName, err)
	}

	// Save keypair and signed cert.
	if err := os.MkdirAll(filepath.Dir(certPath), 0o750); err != nil {
		return fmt.Errorf("create xds cert dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(keyPath), 0o700); err != nil {
		return fmt.Errorf("create xds key dir: %w", err)
	}
	if err := writePEMFile(keyPath, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)}, 0o400); err != nil {
		return err
	}
	if err := os.WriteFile(certPath, signedCert, 0o644); err != nil {
		return fmt.Errorf("write %s cert: %w", commonName, err)
	}
	return nil
}

// signCSRViaGateway sends a PEM-encoded CSR to the gateway's /sign_ca_certificate
// endpoint and returns the signed certificate PEM.
func signCSRViaGateway(gatewayAddr string, csrPEM []byte, caCertPath string) ([]byte, error) {
	csrB64 := base64Encode(csrPEM)

	url := fmt.Sprintf("https://%s/sign_ca_certificate?csr=%s", gatewayAddr, csrB64)

	client := httpClientWithCA(caCertPath)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to gateway: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("gateway returned %d: %s", resp.StatusCode, string(body))
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("gateway returned empty certificate")
	}
	return body, nil
}

// discoverGatewayForCertSigning finds the Day-0 gateway address for CSR signing.
// Reads the etcd_endpoints file written by the join script — the first non-local
// endpoint's host is the bootstrap node which runs the gateway.
func discoverGatewayForCertSigning(runtimeConfigDir string) string {
	root := runtimeRoot(runtimeConfigDir)

	// Read etcd endpoints — first entry is the bootstrap node.
	data, err := os.ReadFile(filepath.Join(root, "config", "etcd_endpoints"))
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Extract host from "https://10.0.0.63:2379"
			line = strings.TrimPrefix(line, "https://")
			line = strings.TrimPrefix(line, "http://")
			host := strings.Split(line, ":")[0]
			if host != "" && host != "127.0.0.1" && host != "localhost" {
				return host + ":8443"
			}
		}
	}

	// Fallback: try config.json PortHTTPS (works on Day-0 node).
	cfg, err := globconfig.GetLocalConfig(false)
	if err == nil {
		if addr, ok := cfg["Address"].(string); ok && addr != "" {
			return addr + ":8443"
		}
	}
	return ""
}

// httpClientWithCA returns an HTTP client that trusts the given CA cert.
func httpClientWithCA(caCertPath string) *http.Client {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return http.DefaultClient
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: true, // gateway cert may not match IP
			},
		},
		Timeout: 30 * time.Second,
	}
}

// base64Encode returns standard base64 encoding of data.
func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// v1 Conformance (INV-5.3): Certificate template with stable identity
// commonName is the stable identity (e.g., "xds-server", "envoy-xds-client")
// Does NOT include domain-based SANs which would tie identity to routing configuration
func xdsCertTemplate(commonName string, isServer bool) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("serial: %w", err)
	}

	// Stable DNSNames: only the commonName itself
	// This ensures certificates remain valid across domain configuration changes
	// For xDS mTLS, the CN is sufficient for identity verification
	dnsNames := []string{commonName}

	// Localhost IP for local development and testing
	ips := []net.IP{net.ParseIP("127.0.0.1")}

	extUsages := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	if isServer {
		extUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	}

	return &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-5 * time.Minute),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           extUsages,
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
		IPAddresses:           ips,
	}, nil
}

func loadCA(certPath, keyPath string) (*x509.Certificate, crypto.Signer, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read CA cert: %w", err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read CA key: %w", err)
	}

	caCert, err := parseFirstCert(certPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA cert: %w", err)
	}
	signer, err := parsePrivateKey(keyPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA key: %w", err)
	}
	return caCert, signer, nil
}

func parseFirstCert(pemData []byte) (*x509.Certificate, error) {
	rest := pemData
	for {
		block, r := pem.Decode(rest)
		rest = r
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			return x509.ParseCertificate(block.Bytes)
		}
	}
	return nil, fmt.Errorf("no certificate block found")
}

func parsePrivateKey(pemData []byte) (crypto.Signer, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}
	if strings.Contains(strings.ToUpper(block.Type), "ENCRYPTED") || strings.Contains(strings.ToUpper(block.Headers["DEK-Info"]), "ENCRYPTED") {
		return nil, fmt.Errorf("encrypted CA keys are not supported")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key, nil
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key, nil
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		switch key := k.(type) {
		case *rsa.PrivateKey:
			return key, nil
		case *ecdsa.PrivateKey:
			return key, nil
		case ed25519.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("unsupported PKCS#8 key type %T", k)
		}
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}

// CanonicalCAPaths returns the CA key, cert, and bundle paths under the canonical PKI root.
func CanonicalCAPaths(runtimeConfigDir string) (keyPath, certPath, bundlePath string) {
	pkiDir := filepath.Join(runtimeRoot(runtimeConfigDir), "pki")
	keyPath = filepath.Join(pkiDir, "ca.key")
	certPath = filepath.Join(pkiDir, "ca.crt")
	bundlePath = filepath.Join(pkiDir, "ca.pem")
	return
}

func writePEMFile(path string, block *pem.Block, mode os.FileMode) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	if err := pem.Encode(f, block); err != nil {
		_ = f.Close()
		return fmt.Errorf("encode %s: %w", path, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(dst, data, mode); err != nil {
		return err
	}
	return nil
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
