package config

import (
	"path/filepath"
	"strings"
)

// CertPaths provides standardized paths for TLS certificates and keys.
// All certificate locations should be accessed through this provider to
// keep path construction consistent and centralized.
type CertPaths struct {
	baseDir string
}

// NewCertPaths creates a CertPaths provider with the given base directory.
// The base directory is typically /var/lib/globular for production systems.
func NewCertPaths(baseDir string) *CertPaths {
	baseDir = strings.TrimSuffix(baseDir, "/")
	return &CertPaths{baseDir: baseDir}
}

func sanitizeSegment(seg string) string {
	s := strings.ReplaceAll(seg, "\\", "/")
	s = strings.Trim(s, "/")
	s = strings.ReplaceAll(s, "..", "_")
	return s
}

// BaseDir returns the normalized base directory used for path construction.
func (c *CertPaths) BaseDir() string {
	return c.baseDir
}

// InternalServerCert returns the internal server certificate path.
// /var/lib/globular/config/tls/fullchain.pem
func (c *CertPaths) InternalServerCert() string {
	return filepath.Join(c.baseDir, "config", "tls", "fullchain.pem")
}

// InternalServerKey returns the internal server private key path.
// /var/lib/globular/config/tls/privkey.pem
func (c *CertPaths) InternalServerKey() string {
	return filepath.Join(c.baseDir, "config", "tls", "privkey.pem")
}

// InternalCABundle returns the internal CA bundle path.
// /var/lib/globular/config/tls/ca.pem
func (c *CertPaths) InternalCABundle() string {
	return filepath.Join(c.baseDir, "config", "tls", "ca.pem")
}

// PKICABundle returns the PKI root CA bundle path.
// /var/lib/globular/pki/ca.pem
func (c *CertPaths) PKICABundle() string {
	return filepath.Join(c.baseDir, "pki", "ca.pem")
}

// PKICert returns a PKI-generated certificate path for the given service.
// /var/lib/globular/pki/{serviceName}/cert.pem
func (c *CertPaths) PKICert(serviceName string) string {
	name := sanitizeSegment(serviceName)
	return filepath.Join(c.baseDir, "pki", name, "cert.pem")
}

// PKIKey returns a PKI-generated private key path for the given service.
// /var/lib/globular/pki/{serviceName}/key.pem
func (c *CertPaths) PKIKey(serviceName string) string {
	name := sanitizeSegment(serviceName)
	return filepath.Join(c.baseDir, "pki", name, "key.pem")
}

// ACMECert returns the ACME fullchain certificate path for the given domain.
// /var/lib/globular/config/tls/acme/{domain}/fullchain.pem
func (c *CertPaths) ACMECert(domain string) string {
	d := sanitizeSegment(domain)
	return filepath.Join(c.baseDir, "config", "tls", "acme", d, "fullchain.pem")
}

// ACMEKey returns the ACME private key path for the given domain.
// /var/lib/globular/config/tls/acme/{domain}/privkey.pem
func (c *CertPaths) ACMEKey(domain string) string {
	d := sanitizeSegment(domain)
	return filepath.Join(c.baseDir, "config", "tls", "acme", d, "privkey.pem")
}

// ACMEDir returns the ACME directory for the given domain.
func (c *CertPaths) ACMEDir(domain string) string {
	d := sanitizeSegment(domain)
	return filepath.Join(c.baseDir, "config", "tls", "acme", d)
}

// TLSConfigDir returns the TLS config directory (/config/tls).
func (c *CertPaths) TLSConfigDir() string {
	return filepath.Join(c.baseDir, "config", "tls")
}

// PKIDir returns the PKI directory (/pki).
func (c *CertPaths) PKIDir() string {
	return filepath.Join(c.baseDir, "pki")
}

// CredsDir returns the credentials/config directory (/config).
func (c *CertPaths) CredsDir() string {
	return filepath.Join(c.baseDir, "config")
}
