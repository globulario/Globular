package config

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	config_ "github.com/globulario/services/golang/config"
)

// ============================================================
// HTTP Wiring (handlers + middleware)
// ============================================================

// caProvider gives access to the CA cert, SAN config, and CSR signing.
// It satisfies the CACertReader, CASigner, and SANConfReader interfaces.
//
// INV-PKI-1: CA material and service certs both live in canonical PKI directory.
type caProvider struct {
	pkiDir         string // Canonical PKI location: /var/lib/globular/pki/ (ca.key, ca.crt)
	serviceCertDir string // Service certificates: /var/lib/globular/pki/issued/services/ (service.crt, service.key)
}

func NewCAProvider() caProvider {
	// INV-PKI-1: Single source of truth for CA material and service certificates
	// Canonical PKI location: /var/lib/globular/pki/
	// The Gateway CA provider loads CA key/cert from canonical PKI paths.
	//
	// CA material: /var/lib/globular/pki/ca.key, ca.crt
	// Service certs: /var/lib/globular/pki/issued/services/service.crt, service.key
	baseDir := config_.GetConfigDir()
	return caProvider{
		pkiDir:         filepath.Join(baseDir, "pki"),
		serviceCertDir: filepath.Join(baseDir, "pki", "issued", "services"),
	}
}

// ReadCACertificate loads ca.crt from the canonical PKI directory.
// INV-PKI-1: CA certificate MUST be read from /var/lib/globular/pki/ca.crt.
func (p caProvider) ReadCACertificate() ([]byte, error) {
	path := filepath.Join(p.pkiDir, "ca.crt")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read CA certificate from canonical PKI path %q: %w", path, err)
	}
	return data, nil
}

// SignCSR signs a CSR PEM string with the local CA in p.credsDir.
func (p caProvider) SignCSR(csrPEM []byte) (string, error) {
	// Parse CSR
	csrBlock, _ := pem.Decode(csrPEM)
	if csrBlock == nil {
		return "", errors.New("invalid CSR: not PEM")
	}
	if csrBlock.Type != "CERTIFICATE REQUEST" && csrBlock.Type != "NEW CERTIFICATE REQUEST" {
		return "", fmt.Errorf("invalid CSR: unexpected PEM type %q", csrBlock.Type)
	}
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse CSR: %w", err)
	}
	if err := csr.CheckSignature(); err != nil {
		return "", fmt.Errorf("CSR signature invalid: %w", err)
	}

	// Load CA cert/key from canonical PKI directory
	// INV-PKI-1: CA material MUST be read from /var/lib/globular/pki/
	caCrtPEM, err := os.ReadFile(filepath.Join(p.pkiDir, "ca.crt"))
	if err != nil {
		return "", fmt.Errorf("read ca.crt from canonical PKI path %q: %w", filepath.Join(p.pkiDir, "ca.crt"), err)
	}
	caCert, err := func() (*x509.Certificate, error) {
		rest := caCrtPEM
		for {
			b, r := pem.Decode(rest)
			if b == nil {
				break
			}
			if b.Type == "CERTIFICATE" {
				return x509.ParseCertificate(b.Bytes)
			}
			rest = r
		}
		return nil, errors.New("ca.crt: no CERTIFICATE block found")
	}()
	if err != nil {
		return "", err
	}

	// INV-PKI-1: CA key MUST be read from /var/lib/globular/pki/ca.key
	caKeyPEM, err := os.ReadFile(filepath.Join(p.pkiDir, "ca.key"))
	if err != nil {
		return "", fmt.Errorf("read ca.key from canonical PKI path %q: %w", filepath.Join(p.pkiDir, "ca.key"), err)
	}
	keyBlock, _ := pem.Decode(caKeyPEM)
	if keyBlock == nil {
		return "", errors.New("ca.key: invalid PEM")
	}
	if strings.Contains(keyBlock.Headers["DEK-Info"], "DES-CBC") || strings.Contains(keyBlock.Type, "ENCRYPTED") {
		return "", errors.New("encrypted PEM blocks are not supported")
	}
	var signer crypto.Signer
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		k, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse RSA key: %w", err)
		}
		signer = k
	case "EC PRIVATE KEY":
		k, err := x509.ParseECPrivateKey(keyBlock.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse EC key: %w", err)
		}
		signer = k
	case "PRIVATE KEY": // PKCS#8
		kAny, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse PKCS#8 key: %w", err)
		}
		switch k := kAny.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			signer = k.(crypto.Signer)
		default:
			return "", fmt.Errorf("unsupported PKCS#8 key type %T", k)
		}
	default:
		return "", fmt.Errorf("unsupported key PEM type %q", keyBlock.Type)
	}

	// Template: copy SANs, emails, subject, add EKUs
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return "", fmt.Errorf("serial: %w", err)
	}
	nb := time.Now().Add(-5 * time.Minute)
	na := nb.Add(360 * 24 * time.Hour)

	tpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               csr.Subject,
		NotBefore:             nb,
		NotAfter:              na,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              csr.DNSNames,
		IPAddresses:           csr.IPAddresses,
		URIs:                  csr.URIs,
		EmailAddresses:        csr.EmailAddresses,
		AuthorityKeyId:        caCert.SubjectKeyId,
		ExtraExtensions:       append([]pkix.Extension(nil), csr.Extensions...),
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, caCert, csr.PublicKey, signer)
	if err != nil {
		return "", fmt.Errorf("create certificate: %w", err)
	}

	var out bytes.Buffer
	if err := pem.Encode(&out, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return "", fmt.Errorf("PEM encode: %w", err)
	}
	return out.String(), nil
}

// ReadSANConfiguration returns a synthesized san.conf using the SANs embedded
// in <configDir>/tls/<domain>/server.crt. It renders an OpenSSL-style layout:
//
// [ req ]
// distinguished_name = req_distinguished_name
// req_extensions = v3_req
//
// [ v3_req ]
// subjectAltName = @alt_names
//
// [ alt_names ]
// DNS.1 = example.com
// IP.1  = 203.0.113.10
// URI.1 = spiffe://example/service
// email.1 = admin@example.com
func (p caProvider) ReadSANConfiguration() ([]byte, error) {
	// INV-PKI-1: Service certificates live in canonical PKI issued/services directory
	crtPath := filepath.Join(p.serviceCertDir, "service.crt")
	pemData, err := os.ReadFile(crtPath)
	if err != nil {
		return nil, fmt.Errorf("read server certificate %q: %w", crtPath, err)
	}

	// decode the first CERTIFICATE block
	var cert *x509.Certificate
	rest := pemData
	for {
		var b *pem.Block
		b, rest = pem.Decode(rest)
		if b == nil {
			break
		}
		if b.Type == "CERTIFICATE" {
			cert, err = x509.ParseCertificate(b.Bytes)
			if err != nil {
				return nil, fmt.Errorf("parse certificate: %w", err)
			}
			break
		}
	}
	if cert == nil {
		return nil, errors.New("server.crt: no CERTIFICATE block found")
	}

	// build san.conf text
	var b strings.Builder
	b.WriteString("[ req ]\n")
	b.WriteString("distinguished_name = req_distinguished_name\n")
	b.WriteString("req_extensions = v3_req\n\n")

	b.WriteString("[ v3_req ]\n")
	b.WriteString("subjectAltName = @alt_names\n\n")

	b.WriteString("[ alt_names ]\n")
	idx := 0
	for i, dns := range cert.DNSNames {
		fmt.Fprintf(&b, "DNS.%d = %s\n", i+1, dns)
		idx++
	}
	for i, ip := range cert.IPAddresses {
		fmt.Fprintf(&b, "IP.%d = %s\n", i+1, ip.String())
		idx++
	}
	for i, uri := range cert.URIs {
		fmt.Fprintf(&b, "URI.%d = %s\n", i+1, uri.String())
		idx++
	}
	for i, email := range cert.EmailAddresses {
		fmt.Fprintf(&b, "email.%d = %s\n", i+1, email)
		idx++
	}

	// if the cert had no SANs, still return a valid skeleton
	if idx == 0 {
		b.WriteString("# (no subjectAltName entries present in certificate)\n")
	}

	return []byte(b.String()), nil
}

// -------- Dependency contracts (small + testable) --------

// CACertReader returns the current CA certificate (PEM bytes).
type CACertReader interface {
	ReadCACertificate() ([]byte, error)
}

// SANConfReader returns the server SAN configuration (san.conf bytes).
type SANConfReader interface {
	ReadSANConfiguration() ([]byte, error)
}

// CASigner signs a CSR (PEM) and returns a signed certificate (PEM string).
type CASigner interface {
	SignCSR(csrPEM []byte) (string, error)
}

// -------- Handlers constructors --------

// NewGetCACertificate — delay WriteHeader until success
func NewGetCACertificate(ca CACertReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pem, err := ca.ReadCACertificate()
		if err != nil {
			http.Error(w, "Client ca cert not found: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated) // match legacy behavior
		_, _ = w.Write(pem)
	}
}

// NewSignCACertificate — delay WriteHeader until success
func NewSignCACertificate(signer CASigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrB64 := r.URL.Query().Get("csr")
		if csrB64 == "" {
			http.Error(w, "missing csr query parameter", http.StatusBadRequest)
			return
		}
		csrPEM, err := base64.StdEncoding.DecodeString(csrB64)
		if err != nil {
			http.Error(w, "fail to decode csr base64 string", http.StatusBadRequest)
			return
		}
		crt, err := signer.SignCSR(csrPEM)
		if err != nil {
			http.Error(w, "fail to sign certificate! "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated) // match legacy behavior
		_, _ = io.WriteString(w, crt)
	}
}

func NewGetSANConf(san SANConfReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := san.ReadSANConfiguration()
		if err != nil {
			http.Error(w, "Client Subject Alternate Name configuration not found: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write(data) // 200 OK by default
	}
}
