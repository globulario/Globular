package admin

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	coreConfig "github.com/globulario/Globular/internal/config"
)

// ── Provider interface ──────────────────────────────────────────────────────

// CertProvider is the minimal surface the certificates handler needs.
type CertProvider interface {
	AdminProvider
	Protocol() string
	Domain() string
	AlternateDomains() []string
	CertPaths() *coreConfig.CertPaths
	RuntimeConfigDir() string
}

// ── JSON response types ─────────────────────────────────────────────────────

// CertOverview is the top-level response for GET /admin/certificates.
type CertOverview struct {
	InternalPKI InternalPKIState `json:"internalPKI"`
	PublicTLS   PublicTLSState   `json:"publicTLS"`
	Envoy       EnvoyState       `json:"envoy"`
	Warnings    []Warning        `json:"warnings"`
	DebugGraph  []DebugNode      `json:"debugGraph"`
}

// InternalPKIState describes the internal CA and service certificate.
type InternalPKIState struct {
	CA          CertRecord `json:"ca"`
	ServiceCert CertRecord `json:"serviceCert"`
	Bundle      CertRecord `json:"bundle"`
	SANConfig   string     `json:"sanConfig"`
	Consumers   []string   `json:"consumers"`
}

// PublicTLSState describes the public-facing TLS certificate.
type PublicTLSState struct {
	LeafCert         *CertRecord         `json:"leafCert"`
	IssuerBundle     *CertRecord         `json:"issuerBundle"`
	Protocol         string              `json:"protocol"`
	Domain           string              `json:"domain"`
	AlternateDomains []string            `json:"alternateDomains"`
	ExternalDomains  []ExternalDomainTLS `json:"externalDomains"`
}

// ExternalDomainTLS describes a certificate obtained for an external domain
// (via the domain reconciler / ACME).
type ExternalDomainTLS struct {
	FQDN      string      `json:"fqdn"`
	LeafCert  *CertRecord `json:"leafCert"`
	KeyPath   string      `json:"keyPath"`
	ChainPath string      `json:"chainPath"`
}

// EnvoyTLSUsage describes a single Envoy TLS secret.
type EnvoyTLSUsage struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
	CAPath   string `json:"caPath"`
	Exists   bool   `json:"exists"`
	Status   string `json:"status"`
}

// EnvoyState describes xDS/SDS TLS consumption.
type EnvoyState struct {
	SDSEnabled bool            `json:"sdsEnabled"`
	Usage      []EnvoyTLSUsage `json:"usage"`
}

// CertRecord is the metadata for a single certificate file.
type CertRecord struct {
	Name              string   `json:"name"`
	Scope             string   `json:"scope"`
	Kind              string   `json:"kind"`
	Subject           string   `json:"subject"`
	Issuer            string   `json:"issuer"`
	SANs              []string `json:"sans"`
	NotBefore         string   `json:"notBefore,omitempty"`
	NotAfter          string   `json:"notAfter,omitempty"`
	DaysUntilExpiry   *int     `json:"daysUntilExpiry,omitempty"`
	FingerprintSHA256 string   `json:"fingerprintSha256,omitempty"`
	Path              string   `json:"path,omitempty"`
	Exists            bool     `json:"exists"`
	Status            string   `json:"status"`
	Source            string   `json:"source"`
}

// Warning is a human-readable issue detected during the certificate audit.
type Warning struct {
	Severity string `json:"severity"` // "error" | "warning" | "info"
	Message  string `json:"message"`
}

// DebugNode is a chain-of-trust graph node for the debug panel.
type DebugNode struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Path    string `json:"path,omitempty"`
	Exists  *bool  `json:"exists,omitempty"`
	Status  string `json:"status,omitempty"`
	Details string `json:"details,omitempty"`
	Chain   string `json:"chain"`
}

// ── Expiry threshold ────────────────────────────────────────────────────────

const expiryWarningDays = 30

// ── Handler ─────────────────────────────────────────────────────────────────

// NewCertificatesHandler returns a GET-only handler for /admin/certificates.
func NewCertificatesHandler(prov CertProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cp := prov.CertPaths()
		domain := prov.Domain()
		runtimeDir := prov.RuntimeConfigDir()

		// Internal PKI
		caPath, _, caBundlePath := canonicalCAPaths(runtimeDir)
		caCert := parsePEMCertificate(caPath, "Internal CA", "internal", "ca", "internal-ca")
		svcCert := parsePEMCertificate(cp.InternalServerCert(), "Internal Service", "internal", "service", "internal-ca")
		bundleCert := parsePEMCertificate(caBundlePath, "CA Bundle", "internal", "issuer_bundle", "internal-ca")

		// SAN config from existing endpoint logic
		sanConfig := ""
		if sanData, err := readSANConf(cp.InternalServerCert()); err == nil {
			sanConfig = string(sanData)
		}

		internalConsumers := []string{"gRPC services", "Envoy upstream TLS", "Node-to-node trust"}

		// Public TLS (ACME) — check primary domain via CertPaths
		acmePath := cp.ACMECert(domain)
		var leafCertPtr *CertRecord
		if acmePath != "" && fileExists(acmePath) {
			rec := parsePEMCertificate(acmePath, "Public Leaf (ACME)", "public", "public_leaf", "acme")
			leafCertPtr = &rec
		}

		// Scan /var/lib/globular/domains/ for all external domain certs
		// (the domain reconciler stores ACME certs here)
		extDomains := scanExternalDomainCerts(runtimeDir)

		// hasPrimaryACME tracks whether the primary domain has its own ACME cert
		// (as opposed to only having external domain certs from the reconciler).
		hasPrimaryACME := leafCertPtr != nil

		// Issuer bundle — check external domain chain.pem files
		var issuerBundlePtr *CertRecord
		if leafCertPtr == nil && len(extDomains) > 0 {
			// Promote the first external domain cert as primary leaf
			leafCertPtr = extDomains[0].LeafCert
		}
		// Look for a chain.pem (issuer bundle) from external domains
		for _, ext := range extDomains {
			if fileExists(ext.ChainPath) {
				rec := parsePEMCertificate(ext.ChainPath, fmt.Sprintf("Issuer Chain (%s)", ext.FQDN), "public", "issuer_bundle", "acme")
				issuerBundlePtr = &rec
				break
			}
		}

		// Envoy xDS
		xdsServerCertPath, xdsServerKeyPath := coreConfig.GetXDSServerCertPaths(runtimeDir)
		xdsClientCertPath, xdsClientKeyPath := coreConfig.GetEnvoyXDSClientCertPaths(runtimeDir)
		caBundle := coreConfig.GetClusterCABundlePath(runtimeDir)

		envoyUsage := []EnvoyTLSUsage{
			makeEnvoyUsage("internal-server-cert", "listener", xdsServerCertPath, xdsServerKeyPath, ""),
			makeEnvoyUsage("internal-ca-bundle", "upstream", "", "", caBundle),
			makeEnvoyUsage("xds-server", "xds_server", xdsServerCertPath, xdsServerKeyPath, caBundle),
			makeEnvoyUsage("envoy-xds-client", "xds_client", xdsClientCertPath, xdsClientKeyPath, caBundle),
		}

		// Only add public-ingress-cert for primary domain ACME (not promoted external certs)
		if hasPrimaryACME {
			acmeKeyPath := cp.ACMEKey(domain)
			envoyUsage = append(envoyUsage, makeEnvoyUsage("public-ingress-cert", "listener", acmePath, acmeKeyPath, ""))
		}

		// Add Envoy secrets for each external domain cert
		for _, ext := range extDomains {
			if ext.LeafCert != nil && ext.LeafCert.Exists {
				secretName := fmt.Sprintf("ext-cert/%s", ext.FQDN)
				envoyUsage = append(envoyUsage, makeEnvoyUsage(secretName, "listener", ext.LeafCert.Path, ext.KeyPath, ext.ChainPath))
			}
		}

		overview := CertOverview{
			InternalPKI: InternalPKIState{
				CA:          caCert,
				ServiceCert: svcCert,
				Bundle:      bundleCert,
				SANConfig:   sanConfig,
				Consumers:   internalConsumers,
			},
			PublicTLS: PublicTLSState{
				LeafCert:         leafCertPtr,
				IssuerBundle:     issuerBundlePtr,
				Protocol:         prov.Protocol(),
				Domain:           domain,
				AlternateDomains: prov.AlternateDomains(),
				ExternalDomains:  extDomains,
			},
			Envoy: EnvoyState{
				SDSEnabled: true, // SDS is always enabled in current architecture
				Usage:      envoyUsage,
			},
		}

		overview.Warnings = collectWarnings(&overview, prov)
		overview.DebugGraph = buildDebugGraph(prov)

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(overview)
	}
}

// ── Certificate parsing ─────────────────────────────────────────────────────

// parsePEMCertificate reads a PEM file and extracts x509 metadata.
// It NEVER returns private key data.
func parsePEMCertificate(path, name, scope, kind, source string) CertRecord {
	rec := CertRecord{
		Name:   name,
		Scope:  scope,
		Kind:   kind,
		Path:   path,
		Source: source,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		rec.Status = "missing"
		rec.Exists = false
		return rec
	}
	rec.Exists = true

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		rec.Status = "parse_error"
		return rec
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		rec.Status = "parse_error"
		return rec
	}

	rec.Subject = cert.Subject.String()
	rec.Issuer = cert.Issuer.String()
	rec.SANs = collectSANs(cert)
	rec.NotBefore = cert.NotBefore.UTC().Format(time.RFC3339)
	rec.NotAfter = cert.NotAfter.UTC().Format(time.RFC3339)

	days := int(math.Floor(time.Until(cert.NotAfter).Hours() / 24))
	rec.DaysUntilExpiry = &days

	hash := sha256.Sum256(cert.Raw)
	rec.FingerprintSHA256 = formatFingerprint(hash[:])

	rec.Status = checkCertStatus(cert)
	return rec
}

// collectSANs gathers DNS names and IP addresses from the certificate.
func collectSANs(cert *x509.Certificate) []string {
	var sans []string
	sans = append(sans, cert.DNSNames...)
	for _, ip := range cert.IPAddresses {
		sans = append(sans, ip.String())
	}
	return sans
}

// formatFingerprint formats a SHA-256 hash as colon-separated hex.
func formatFingerprint(b []byte) string {
	parts := make([]string, len(b))
	for i, v := range b {
		parts[i] = fmt.Sprintf("%02X", v)
	}
	return strings.Join(parts, ":")
}

// checkCertStatus returns "valid", "expiring", or "expired".
func checkCertStatus(cert *x509.Certificate) string {
	now := time.Now()
	if now.After(cert.NotAfter) {
		return "expired"
	}
	if time.Until(cert.NotAfter).Hours()/24 < expiryWarningDays {
		return "expiring"
	}
	return "valid"
}

// ── Envoy usage helpers ──────────────────────────────────────────────────────

// makeEnvoyUsage builds an EnvoyTLSUsage entry, checking file existence.
func makeEnvoyUsage(name, typ, certPath, keyPath, caPath string) EnvoyTLSUsage {
	u := EnvoyTLSUsage{
		Name:     name,
		Type:     typ,
		CertPath: certPath,
		KeyPath:  keyPath,
		CAPath:   caPath,
	}
	// Check existence of referenced files
	allExist := true
	if certPath != "" && !fileExists(certPath) {
		allExist = false
	}
	if keyPath != "" && !fileExists(keyPath) {
		allExist = false
	}
	if caPath != "" && !fileExists(caPath) {
		allExist = false
	}
	u.Exists = allExist
	if allExist {
		u.Status = "ok"
	} else {
		u.Status = "missing"
	}
	return u
}

// scanExternalDomainCerts scans /var/lib/globular/domains/ for subdirectories
// containing fullchain.pem (certs placed by the domain reconciler).
func scanExternalDomainCerts(runtimeDir string) []ExternalDomainTLS {
	domainsDir := filepath.Join(runtimeDir, "domains")
	entries, err := os.ReadDir(domainsDir)
	if err != nil {
		return nil
	}
	var result []ExternalDomainTLS
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		fqdn := entry.Name()
		certPath := filepath.Join(domainsDir, fqdn, "fullchain.pem")
		keyPath := filepath.Join(domainsDir, fqdn, "privkey.pem")
		chainPath := filepath.Join(domainsDir, fqdn, "chain.pem")

		if !fileExists(certPath) {
			continue
		}

		rec := parsePEMCertificate(certPath, fmt.Sprintf("Public Leaf (%s)", fqdn), "public", "public_leaf", "acme")
		ext := ExternalDomainTLS{
			FQDN:      fqdn,
			LeafCert:  &rec,
			KeyPath:   keyPath,
			ChainPath: chainPath,
		}
		result = append(result, ext)
	}
	return result
}

// readSANConf extracts SANs from a certificate file and returns san.conf text.
func readSANConf(certPath string) ([]byte, error) {
	pemData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("no certificate block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString("[ alt_names ]\n")
	for i, dns := range cert.DNSNames {
		fmt.Fprintf(&b, "DNS.%d = %s\n", i+1, dns)
	}
	for i, ip := range cert.IPAddresses {
		fmt.Fprintf(&b, "IP.%d = %s\n", i+1, ip.String())
	}
	return []byte(b.String()), nil
}

// ── Warning collection ──────────────────────────────────────────────────────

func collectWarnings(overview *CertOverview, prov CertProvider) []Warning {
	var ws []Warning
	domain := prov.Domain()
	altDomains := prov.AlternateDomains()
	protocol := prov.Protocol()

	// Internal CA
	warnCert(&ws, "Internal CA", overview.InternalPKI.CA)

	// Internal service cert
	warnCert(&ws, "Internal service", overview.InternalPKI.ServiceCert)

	// Service cert SAN coverage
	svcStatus := overview.InternalPKI.ServiceCert.Status
	if svcStatus == "valid" || svcStatus == "expiring" {
		if domain != "" && !sanCovers(overview.InternalPKI.ServiceCert.SANs, domain) {
			ws = append(ws, Warning{Severity: "warning", Message: fmt.Sprintf("Internal service certificate SANs do not cover domain %q", domain)})
		}
	}

	// Public cert checks — only for primary domain ACME cert (not promoted external certs).
	// External domain certs are validated separately below.
	if protocol == "https" {
		hasAnyPublicCert := overview.PublicTLS.LeafCert != nil || len(overview.PublicTLS.ExternalDomains) > 0
		if !hasAnyPublicCert {
			ws = append(ws, Warning{Severity: "error", Message: "No public TLS certificates found (protocol=https)"})
		}

		// Check primary domain ACME cert (only if it's a real primary cert, not promoted)
		if overview.PublicTLS.LeafCert != nil && overview.PublicTLS.LeafCert.Source == "acme" {
			cp := prov.CertPaths()
			acmePath := cp.ACMECert(domain)
			// Only validate SAN/key if the primary domain actually has its own ACME cert
			if acmePath != "" && fileExists(acmePath) {
				warnCert(&ws, "Public TLS (primary)", *overview.PublicTLS.LeafCert)
				leafStatus := overview.PublicTLS.LeafCert.Status
				if leafStatus == "valid" || leafStatus == "expiring" {
					allDomains := append([]string{domain}, altDomains...)
					for _, d := range allDomains {
						if d != "" && !sanCovers(overview.PublicTLS.LeafCert.SANs, d) {
							ws = append(ws, Warning{Severity: "warning", Message: fmt.Sprintf("Primary ACME certificate SANs do not cover domain %q", d)})
						}
					}
				}
				keyPath := cp.ACMEKey(domain)
				if _, err := os.Stat(keyPath); err != nil {
					ws = append(ws, Warning{Severity: "error", Message: fmt.Sprintf("ACME private key missing at %s", keyPath)})
				}
			}
		}
	}

	// External domain cert warnings
	for _, ext := range overview.PublicTLS.ExternalDomains {
		if ext.LeafCert == nil {
			continue
		}
		warnCert(&ws, fmt.Sprintf("External cert (%s)", ext.FQDN), *ext.LeafCert)
		if ext.LeafCert.Exists && !fileExists(ext.KeyPath) {
			ws = append(ws, Warning{Severity: "error", Message: fmt.Sprintf("Private key missing for %s at %s", ext.FQDN, ext.KeyPath)})
		}
	}

	// Envoy usage warnings
	for _, u := range overview.Envoy.Usage {
		if u.Status == "missing" {
			ws = append(ws, Warning{Severity: "warning", Message: fmt.Sprintf("Envoy secret %q has missing files", u.Name)})
		}
	}

	return ws
}

// warnCert adds warnings for a single CertRecord.
func warnCert(ws *[]Warning, label string, cr CertRecord) {
	switch cr.Status {
	case "missing":
		*ws = append(*ws, Warning{Severity: "error", Message: fmt.Sprintf("%s certificate is missing", label)})
	case "expired":
		*ws = append(*ws, Warning{Severity: "error", Message: fmt.Sprintf("%s certificate has expired", label)})
	case "expiring":
		*ws = append(*ws, Warning{Severity: "warning", Message: fmt.Sprintf("%s certificate expires in %d days", label, safeExpiry(cr.DaysUntilExpiry))})
	case "parse_error":
		*ws = append(*ws, Warning{Severity: "error", Message: fmt.Sprintf("%s certificate could not be parsed", label)})
	}
}

func safeExpiry(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// sanCovers returns true if the SAN list covers the given domain (exact or wildcard).
func sanCovers(sans []string, domain string) bool {
	for _, s := range sans {
		if strings.EqualFold(s, domain) {
			return true
		}
		// Wildcard match: *.example.com covers foo.example.com
		if strings.HasPrefix(s, "*.") {
			suffix := s[1:] // ".example.com"
			if strings.HasSuffix(strings.ToLower(domain), strings.ToLower(suffix)) && !strings.Contains(domain[:len(domain)-len(suffix)], ".") {
				return true
			}
		}
	}
	return false
}

// ── Debug graph ─────────────────────────────────────────────────────────────

func buildDebugGraph(prov CertProvider) []DebugNode {
	cp := prov.CertPaths()
	domain := prov.Domain()
	runtimeDir := prov.RuntimeConfigDir()

	caKeyPath, caCertPath, caBundlePath := canonicalCAPaths(runtimeDir)
	xdsServerCert, xdsServerKey := coreConfig.GetXDSServerCertPaths(runtimeDir)
	xdsClientCert, xdsClientKey := coreConfig.GetEnvoyXDSClientCertPaths(runtimeDir)

	var nodes []DebugNode

	// Config node
	nodes = append(nodes, DebugNode{
		Type:    "config",
		Name:    "Protocol",
		Details: prov.Protocol(),
		Chain:   "public-leaf",
	})
	nodes = append(nodes, DebugNode{
		Type:    "config",
		Name:    "Domain",
		Details: domain,
		Chain:   "public-leaf",
	})

	// Internal CA chain
	nodes = append(nodes, fileNode("CA Certificate", caCertPath, "internal-ca"))
	nodes = append(nodes, fileNode("CA Private Key", caKeyPath, "internal-ca"))
	nodes = append(nodes, fileNode("CA Bundle", caBundlePath, "internal-ca"))

	// Internal service chain
	nodes = append(nodes, fileNode("Service Certificate", cp.InternalServerCert(), "internal-service"))
	nodes = append(nodes, fileNode("Service Key", cp.InternalServerKey(), "internal-service"))

	// Public leaf chain
	nodes = append(nodes, fileNode("ACME Certificate", cp.ACMECert(domain), "public-leaf"))
	nodes = append(nodes, fileNode("ACME Key", cp.ACMEKey(domain), "public-leaf"))

	// xDS server chain
	nodes = append(nodes, fileNode("xDS Server Cert", xdsServerCert, "xds-server"))
	nodes = append(nodes, fileNode("xDS Server Key", xdsServerKey, "xds-server"))

	// xDS client chain
	nodes = append(nodes, fileNode("xDS Client Cert", xdsClientCert, "xds-client"))
	nodes = append(nodes, fileNode("xDS Client Key", xdsClientKey, "xds-client"))

	// Consumer nodes
	nodes = append(nodes, DebugNode{
		Type:  "consumer",
		Name:  "gRPC Services",
		Chain: "internal-service",
	})
	nodes = append(nodes, DebugNode{
		Type:  "consumer",
		Name:  "Envoy Ingress",
		Chain: "public-leaf",
	})
	nodes = append(nodes, DebugNode{
		Type:  "consumer",
		Name:  "xDS Server",
		Chain: "xds-server",
	})
	nodes = append(nodes, DebugNode{
		Type:  "consumer",
		Name:  "Envoy xDS Client",
		Chain: "xds-client",
	})

	return nodes
}

// fileNode creates a DebugNode for a file path with existence check.
func fileNode(name, path, chain string) DebugNode {
	exists := fileExists(path)
	status := "ok"
	if !exists {
		status = "missing"
	}
	return DebugNode{
		Type:   "file",
		Name:   name,
		Path:   path,
		Exists: &exists,
		Status: status,
		Chain:  chain,
	}
}

// canonicalCAPaths returns (certPath, keyPath, bundlePath) for the internal CA.
// Delegates to the config package to avoid path duplication.
func canonicalCAPaths(runtimeConfigDir string) (certPath, keyPath, bundlePath string) {
	// CanonicalCAPaths returns (key, cert, bundle) — reorder to (cert, key, bundle).
	keyPath_, certPath_, bundlePath_ := coreConfig.CanonicalCAPaths(runtimeConfigDir)
	return certPath_, keyPath_, bundlePath_
}

// fileExists checks whether a regular file exists at path.
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
