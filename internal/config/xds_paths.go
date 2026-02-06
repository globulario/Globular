package config

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
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

	caKeyPath, caCertPath, caBundlePath := canonicalCAPaths(runtimeConfigDir)
	if !fileExists(caKeyPath) || !fileExists(caCertPath) {
		return fmt.Errorf("cluster CA missing: %s or %s", caKeyPath, caCertPath)
	}
	if !fileExists(caBundlePath) {
		// Prefer to expose the CA bundle even if rotation tooling has not created it yet.
		if err := copyFile(caCertPath, caBundlePath, 0o444); err != nil {
			return fmt.Errorf("ensure CA bundle %s: %w", caBundlePath, err)
		}
	}

	// Ensure xDS server and Envoy client identities exist.
	if err := ensureXDSCertificate(true, caCertPath, caKeyPath, runtimeConfigDir); err != nil {
		return err
	}
	if err := ensureXDSCertificate(false, caCertPath, caKeyPath, runtimeConfigDir); err != nil {
		return err
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

// canonicalCAPaths returns the CA key, cert, and bundle paths under the canonical PKI root.
func canonicalCAPaths(runtimeConfigDir string) (keyPath, certPath, bundlePath string) {
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
