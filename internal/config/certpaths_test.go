package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCertPathsBaseNormalization(t *testing.T) {
	tests := []struct {
		base     string
		expected string
	}{
		{"/var/lib/globular", "/var/lib/globular"},
		{"/var/lib/globular/", "/var/lib/globular"},
		{"/custom/path", "/custom/path"},
		{"/custom/path/", "/custom/path"},
	}
	for _, tt := range tests {
		cp := NewCertPaths(tt.base)
		if got := cp.BaseDir(); got != tt.expected {
			t.Fatalf("BaseDir(%q)=%q want %q", tt.base, got, tt.expected)
		}
	}
}

func TestCertPaths_InternalPaths(t *testing.T) {
	cp := NewCertPaths("/var/lib/globular")
	if got := cp.InternalServerCert(); got != filepath.Join("/var/lib/globular", "config", "tls", "fullchain.pem") {
		t.Fatalf("InternalServerCert=%q", got)
	}
	if got := cp.InternalServerKey(); got != filepath.Join("/var/lib/globular", "config", "tls", "privkey.pem") {
		t.Fatalf("InternalServerKey=%q", got)
	}
	if got := cp.InternalCABundle(); got != filepath.Join("/var/lib/globular", "config", "tls", "ca.pem") {
		t.Fatalf("InternalCABundle=%q", got)
	}
}

func TestCertPaths_PKIPaths(t *testing.T) {
	cp := NewCertPaths("/var/lib/globular")
	if got := cp.PKICABundle(); got != filepath.Join("/var/lib/globular", "pki", "ca.pem") {
		t.Fatalf("PKICABundle=%q", got)
	}
	if got := cp.PKICert("gateway"); got != filepath.Join("/var/lib/globular", "pki", "gateway", "cert.pem") {
		t.Fatalf("PKICert=%q", got)
	}
	if got := cp.PKIKey("gateway"); got != filepath.Join("/var/lib/globular", "pki", "gateway", "key.pem") {
		t.Fatalf("PKIKey=%q", got)
	}
}

func TestCertPaths_ACMEPaths(t *testing.T) {
	cp := NewCertPaths("/var/lib/globular")
	dn := "example.com"
	if got := cp.ACMECert(dn); got != filepath.Join("/var/lib/globular", "config", "tls", "acme", dn, "fullchain.pem") {
		t.Fatalf("ACMECert=%q", got)
	}
	if got := cp.ACMEKey(dn); got != filepath.Join("/var/lib/globular", "config", "tls", "acme", dn, "privkey.pem") {
		t.Fatalf("ACMEKey=%q", got)
	}
	if got := cp.ACMEDir(dn); got != filepath.Join("/var/lib/globular", "config", "tls", "acme", dn) {
		t.Fatalf("ACMEDir=%q", got)
	}
}

func TestCertPaths_Dirs(t *testing.T) {
	cp := NewCertPaths("/var/lib/globular")
	if got := cp.TLSConfigDir(); got != filepath.Join("/var/lib/globular", "config", "tls") {
		t.Fatalf("TLSConfigDir=%q", got)
	}
	if got := cp.PKIDir(); got != filepath.Join("/var/lib/globular", "pki") {
		t.Fatalf("PKIDir=%q", got)
	}
	if got := cp.CredsDir(); got != filepath.Join("/var/lib/globular", "config") {
		t.Fatalf("CredsDir=%q", got)
	}
}

func TestCertPaths_NoDoubleSlash(t *testing.T) {
	cp := NewCertPaths("/var/lib/globular")
	paths := []string{
		cp.InternalServerCert(),
		cp.InternalServerKey(),
		cp.InternalCABundle(),
		cp.ACMECert("example.com"),
		cp.ACMEKey("example.com"),
	}
	for _, p := range paths {
		if strings.Contains(p, "//") {
			t.Fatalf("path contains double slash: %q", p)
		}
	}
}
