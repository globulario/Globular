package admin

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	coreConfig "github.com/globulario/Globular/internal/config"
)

// ── Action provider interface ────────────────────────────────────────────────

// CertActionsProvider is the surface the cert action handlers need.
type CertActionsProvider interface {
	Protocol() string
	Domain() string
	AlternateDomains() []string
	CertPaths() *coreConfig.CertPaths
	RuntimeConfigDir() string
}

// ── Action response ──────────────────────────────────────────────────────────

// CertActionResponse is the JSON body returned by certificate action endpoints.
type CertActionResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func writeActionResponse(w http.ResponseWriter, status int, ok bool, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(CertActionResponse{OK: ok, Message: msg})
}

// ── Renew public certificate ─────────────────────────────────────────────────

// NewRenewPublicHandler returns a POST handler that triggers public cert renewal
// by invalidating existing cert files so the domain reconciler re-obtains them.
func NewRenewPublicHandler(prov CertActionsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		protocol := prov.Protocol()
		if protocol != "https" {
			writeActionResponse(w, http.StatusBadRequest, false,
				"Cannot renew public certificate: protocol is not HTTPS")
			return
		}

		runtimeDir := prov.RuntimeConfigDir()
		domainsDir := filepath.Join(runtimeDir, "domains")

		// Scan for external domain cert directories
		entries, err := os.ReadDir(domainsDir)
		if err != nil {
			writeActionResponse(w, http.StatusBadRequest, false,
				fmt.Sprintf("Cannot renew: no external domains directory found (%s)", domainsDir))
			return
		}

		renewed := 0
		var errs []string
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			fqdn := entry.Name()
			certFile := filepath.Join(domainsDir, fqdn, "fullchain.pem")

			if !fileExists(certFile) {
				continue
			}

			// Rename the cert file to force the reconciler to re-obtain.
			// The reconciler checks isCertificateValid() which returns false
			// if the file is missing, triggering a new ACME obtain.
			backupFile := certFile + ".renew-backup"
			if err := os.Rename(certFile, backupFile); err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", fqdn, err))
				continue
			}
			renewed++
		}

		if renewed == 0 && len(errs) == 0 {
			writeActionResponse(w, http.StatusBadRequest, false,
				"No external domain certificates found to renew")
			return
		}

		if len(errs) > 0 {
			writeActionResponse(w, http.StatusInternalServerError, false,
				fmt.Sprintf("Failed to invalidate certificates for renewal: %v", errs))
			return
		}

		writeActionResponse(w, http.StatusOK, true,
			fmt.Sprintf("Public certificate renewal queued for %d domain(s). "+
				"The domain reconciler will re-obtain certificates within ~60 seconds.", renewed))
	}
}

// ── Regenerate internal certificates ─────────────────────────────────────────

// NewRegenerateInternalHandler returns a POST handler that regenerates
// internal PKI certificates (service cert + xDS mTLS certs).
func NewRegenerateInternalHandler(prov CertActionsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		runtimeDir := prov.RuntimeConfigDir()
		cp := prov.CertPaths()

		// Validate CA exists
		caKeyPath, caCertPath, _ := coreConfig.CanonicalCAPaths(runtimeDir)
		if !fileExists(caCertPath) || !fileExists(caKeyPath) {
			writeActionResponse(w, http.StatusBadRequest, false,
				"Cannot regenerate: internal CA is missing")
			return
		}

		var regenerated []string
		var errs []string

		// 1. Regenerate service certificate
		if err := regenerateServiceCert(cp, caCertPath, caKeyPath, prov); err != nil {
			errs = append(errs, fmt.Sprintf("service cert: %v", err))
		} else {
			regenerated = append(regenerated, "service certificate")
		}

		// 2. Regenerate xDS certificates (delete + re-ensure)
		if err := regenerateXDSCerts(runtimeDir); err != nil {
			errs = append(errs, fmt.Sprintf("xDS certs: %v", err))
		} else {
			regenerated = append(regenerated, "xDS server cert", "xDS client cert")
		}

		if len(errs) > 0 && len(regenerated) == 0 {
			writeActionResponse(w, http.StatusInternalServerError, false,
				fmt.Sprintf("Failed to regenerate internal certificates: %v", errs))
			return
		}

		msg := fmt.Sprintf("Regenerated: %v.", regenerated)
		if len(errs) > 0 {
			msg += fmt.Sprintf(" Errors: %v.", errs)
		}
		msg += " Services may need restart to pick up new certificates."

		writeActionResponse(w, http.StatusOK, true, msg)
	}
}

// regenerateServiceCert generates a new service key+cert signed by the internal CA.
func regenerateServiceCert(cp *coreConfig.CertPaths, caCertPath, caKeyPath string, prov CertActionsProvider) error {
	certPath := cp.InternalServerCert()
	keyPath := cp.InternalServerKey()

	// Load CA
	caCert, caSigner, err := loadCAForSigning(caCertPath, caKeyPath)
	if err != nil {
		return fmt.Errorf("load CA: %w", err)
	}

	// Collect SANs from current config
	domain := prov.Domain()
	var dnsNames []string
	var ips []net.IP

	if domain != "" {
		dnsNames = append(dnsNames, domain)
	}
	dnsNames = append(dnsNames, "localhost")
	for _, alt := range prov.AlternateDomains() {
		if alt != "" {
			dnsNames = append(dnsNames, alt)
		}
	}
	ips = append(ips, net.ParseIP("127.0.0.1"), net.ParseIP("::1"))

	// Also preserve SANs from the existing cert if present
	if existing := readExistingSANs(certPath); existing != nil {
		for _, dns := range existing.dnsNames {
			if !containsStr(dnsNames, dns) {
				dnsNames = append(dnsNames, dns)
			}
		}
		for _, ip := range existing.ips {
			if !containsIP(ips, ip) {
				ips = append(ips, ip)
			}
		}
	}

	// Generate new key
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}

	// Create certificate template
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("serial: %w", err)
	}

	now := time.Now()
	tpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: domain,
		},
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
		IPAddresses:           ips,
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, caCert, &privKey.PublicKey, caSigner)
	if err != nil {
		return fmt.Errorf("sign certificate: %w", err)
	}

	// Ensure directories exist
	if err := os.MkdirAll(filepath.Dir(certPath), 0o750); err != nil {
		return fmt.Errorf("create cert dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(keyPath), 0o700); err != nil {
		return fmt.Errorf("create key dir: %w", err)
	}

	// Write key
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	if err := os.WriteFile(keyPath, keyPEM, 0o400); err != nil {
		return fmt.Errorf("write key: %w", err)
	}

	// Write cert
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	if err := os.WriteFile(certPath, certPEM, 0o644); err != nil {
		return fmt.Errorf("write cert: %w", err)
	}

	return nil
}

// regenerateXDSCerts deletes existing xDS certs and recreates them.
func regenerateXDSCerts(runtimeDir string) error {
	// Delete existing xDS server certs
	serverCert, serverKey := coreConfig.GetXDSServerCertPaths(runtimeDir)
	_ = os.Remove(serverCert)
	_ = os.Remove(serverKey)

	// Delete existing xDS client certs
	clientCert, clientKey := coreConfig.GetEnvoyXDSClientCertPaths(runtimeDir)
	_ = os.Remove(clientCert)
	_ = os.Remove(clientKey)

	// Regenerate via existing logic
	return coreConfig.EnsureXDSMTLSMaterials(runtimeDir, false)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// loadCAForSigning loads the CA cert and private key for certificate signing.
func loadCAForSigning(certPath, keyPath string) (*x509.Certificate, crypto.Signer, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read CA cert: %w", err)
	}
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("invalid CA cert PEM")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA cert: %w", err)
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read CA key: %w", err)
	}
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("invalid CA key PEM")
	}

	var signer crypto.Signer
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		signer, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		signer, err = x509.ParseECPrivateKey(keyBlock.Bytes)
	case "PRIVATE KEY":
		var k any
		k, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if err == nil {
			switch typed := k.(type) {
			case *rsa.PrivateKey:
				signer = typed
			case *ecdsa.PrivateKey:
				signer = typed
			case ed25519.PrivateKey:
				signer = typed
			default:
				err = fmt.Errorf("unsupported PKCS#8 key type %T", k)
			}
		}
	default:
		err = fmt.Errorf("unsupported key PEM type %q", keyBlock.Type)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA key: %w", err)
	}
	return caCert, signer, nil
}

type existingSANs struct {
	dnsNames []string
	ips      []net.IP
}

func readExistingSANs(certPath string) *existingSANs {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return nil
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil
	}
	return &existingSANs{dnsNames: cert.DNSNames, ips: cert.IPAddresses}
}

func containsStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func containsIP(ips []net.IP, ip net.IP) bool {
	for _, v := range ips {
		if v.Equal(ip) {
			return true
		}
	}
	return false
}
