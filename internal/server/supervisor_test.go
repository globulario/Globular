package server

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
)

func TestReloadPopulatesLeafCertificate(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := writeTestKeyPair(t, dir, "server.local")

	r, err := newCertReloader(certPath, keyPath)
	if err != nil {
		t.Fatalf("newCertReloader: %v", err)
	}
	if r.cert == nil || r.cert.Leaf == nil {
		t.Fatal("expected loaded certificate leaf metadata")
	}
	if got := r.cert.Leaf.Subject.CommonName; got != "server.local" {
		t.Fatalf("leaf CN=%q want server.local", got)
	}
}

func TestMaybeReloadKeepsPreviousCertWhenFilesDisappear(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := writeTestKeyPair(t, dir, "server.local")

	r, err := newCertReloader(certPath, keyPath)
	if err != nil {
		t.Fatalf("newCertReloader: %v", err)
	}
	loaded := r.cert
	if err := os.Remove(certPath); err != nil {
		t.Fatalf("remove cert: %v", err)
	}
	if err := os.Remove(keyPath); err != nil {
		t.Fatalf("remove key: %v", err)
	}

	if err := r.maybeReload(); err != nil {
		t.Fatalf("maybeReload: %v", err)
	}
	if r.cert != loaded {
		t.Fatal("expected previous certificate to remain active when files disappear")
	}
}

func writeTestKeyPair(t *testing.T, dir, commonName string) (string, string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{commonName},
	}
	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("CreateCertificate: %v", err)
	}
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")
	if err := os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o644); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}), 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	return certPath, keyPath
}
