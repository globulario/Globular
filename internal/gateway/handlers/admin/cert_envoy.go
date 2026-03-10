package admin

import (
	"fmt"
	"path/filepath"
	"strings"

	coreConfig "github.com/globulario/Globular/internal/config"
)

// ── Enriched Envoy TLS types ─────────────────────────────────────────────────

// EnvoyListenerTLS describes TLS usage for an Envoy listener (downstream).
type EnvoyListenerTLS struct {
	Name           string `json:"name"`
	Address        string `json:"address,omitempty"`
	TLSMode        string `json:"tlsMode"` // "public" | "internal"
	SecretName     string `json:"secretName,omitempty"`
	CertPath       string `json:"certPath,omitempty"`
	KeyPath        string `json:"keyPath,omitempty"`
	Exists         bool   `json:"exists"`
	Status         string `json:"status"` // "ok" | "missing" | "invalid"
	CertificateRef string `json:"certificateRef,omitempty"`
}

// EnvoyUpstreamTLS describes TLS usage for an Envoy upstream (to backend services).
type EnvoyUpstreamTLS struct {
	Name           string `json:"name"`
	TLSMode        string `json:"tlsMode"` // "internal" | "none"
	ServerCertPath string `json:"serverCertPath,omitempty"`
	KeyPath        string `json:"keyPath,omitempty"`
	CAPath         string `json:"caPath,omitempty"`
	Exists         bool   `json:"exists"`
	Status         string `json:"status"` // "ok" | "missing" | "no_tls"
	CertificateRef string `json:"certificateRef,omitempty"`
	CARef          string `json:"caRef,omitempty"`
}

// EnvoySecretRef describes an SDS secret with its consumers.
type EnvoySecretRef struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"` // "tls_certificate" | "validation_context"
	CertPath  string   `json:"certPath,omitempty"`
	KeyPath   string   `json:"keyPath,omitempty"`
	CAPath    string   `json:"caPath,omitempty"`
	Exists    bool     `json:"exists"`
	Status    string   `json:"status"` // "ok" | "missing" | "invalid"
	Consumers []string `json:"consumers"`
}

// XDSClientTLS describes the xDS control plane client TLS material.
type XDSClientTLS struct {
	Enabled  bool   `json:"enabled"`
	CertPath string `json:"certPath,omitempty"`
	KeyPath  string `json:"keyPath,omitempty"`
	CAPath   string `json:"caPath,omitempty"`
	Exists   bool   `json:"exists"`
	Status   string `json:"status"` // "ok" | "missing"
}

// EnvoyStateEnriched is the enriched Envoy TLS state for PR3.
type EnvoyStateEnriched struct {
	SDSEnabled bool               `json:"sdsEnabled"`
	Usage      []EnvoyTLSUsage    `json:"usage"`
	Listeners  []EnvoyListenerTLS `json:"listeners"`
	Upstreams  []EnvoyUpstreamTLS `json:"upstreams"`
	Secrets    []EnvoySecretRef   `json:"secrets"`
	XDSClient  XDSClientTLS       `json:"xdsClient"`
}

// ── Builder functions ────────────────────────────────────────────────────────

// buildEnrichedEnvoyState builds the full Envoy TLS introspection state.
func buildEnrichedEnvoyState(
	runtimeDir string,
	extDomains []ExternalDomainTLS,
	envoyUsage []EnvoyTLSUsage,
	prov CertProvider,
) EnvoyStateEnriched {
	xdsServerCert, xdsServerKey := coreConfig.GetXDSServerCertPaths(runtimeDir)
	xdsClientCert, xdsClientKey := coreConfig.GetEnvoyXDSClientCertPaths(runtimeDir)
	caBundle := coreConfig.GetClusterCABundlePath(runtimeDir)
	cp := prov.CertPaths()

	// ── Listeners ────────────────────────────────────────────────────────
	listeners := collectListeners(xdsServerCert, xdsServerKey, extDomains, cp, prov)

	// ── Upstreams (from service configs) ─────────────────────────────────
	upstreams := collectUpstreams(prov, caBundle)

	// ── Secrets ──────────────────────────────────────────────────────────
	secrets := collectSecrets(xdsServerCert, xdsServerKey, caBundle, extDomains, upstreams)

	// ── xDS Client TLS ──────────────────────────────────────────────────
	xdsClient := XDSClientTLS{
		Enabled:  true,
		CertPath: xdsClientCert,
		KeyPath:  xdsClientKey,
		CAPath:   caBundle,
	}
	xdsClient.Exists = fileExists(xdsClientCert) && fileExists(xdsClientKey) && fileExists(caBundle)
	if xdsClient.Exists {
		xdsClient.Status = "ok"
	} else {
		xdsClient.Status = "missing"
	}

	return EnvoyStateEnriched{
		SDSEnabled: true,
		Usage:      envoyUsage,
		Listeners:  listeners,
		Upstreams:  upstreams,
		Secrets:    secrets,
		XDSClient:  xdsClient,
	}
}

// collectListeners builds listener TLS entries.
func collectListeners(
	xdsServerCert, xdsServerKey string,
	extDomains []ExternalDomainTLS,
	cp *coreConfig.CertPaths,
	prov CertProvider,
) []EnvoyListenerTLS {
	var listeners []EnvoyListenerTLS

	// Internal HTTPS listener (downstream TLS for *.globular.internal)
	internalExists := fileExists(xdsServerCert) && fileExists(xdsServerKey)
	internalStatus := "ok"
	if !internalExists {
		internalStatus = "missing"
	}
	listeners = append(listeners, EnvoyListenerTLS{
		Name:           "internal_https",
		Address:        ":443",
		TLSMode:        "internal",
		SecretName:     "internal-server-cert",
		CertPath:       xdsServerCert,
		KeyPath:        xdsServerKey,
		Exists:         internalExists,
		Status:         internalStatus,
		CertificateRef: "internal_service",
	})

	// External domain listeners (one per domain, SNI-routed)
	for _, ext := range extDomains {
		certPath := ""
		keyPath := ""
		if ext.LeafCert != nil {
			certPath = ext.LeafCert.Path
		}
		keyPath = ext.KeyPath

		exists := fileExists(certPath) && fileExists(keyPath)
		status := "ok"
		if !exists {
			status = "missing"
		}

		listeners = append(listeners, EnvoyListenerTLS{
			Name:           fmt.Sprintf("ext_%s", ext.FQDN),
			Address:        ":443",
			TLSMode:        "public",
			SecretName:     fmt.Sprintf("ext-cert/%s", ext.FQDN),
			CertPath:       certPath,
			KeyPath:        keyPath,
			Exists:         exists,
			Status:         status,
			CertificateRef: "public_leaf",
		})
	}

	return listeners
}

// collectUpstreams builds upstream TLS entries from service configs.
func collectUpstreams(prov CertProvider, caBundle string) []EnvoyUpstreamTLS {
	var upstreams []EnvoyUpstreamTLS

	configs, err := prov.AllServiceConfigs()
	if err != nil {
		return upstreams
	}

	// Group by service name to avoid duplicates
	seen := make(map[string]bool)
	for _, cfg := range configs {
		name, _ := cfg["Name"].(string)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true

		// Determine if service uses TLS
		protocol, _ := cfg["Protocol"].(string)
		tlsEnabled, _ := cfg["TLS"].(bool)
		usesTLS := strings.EqualFold(protocol, "https") || tlsEnabled

		if !usesTLS {
			upstreams = append(upstreams, EnvoyUpstreamTLS{
				Name:    name,
				TLSMode: "none",
				Exists:  true,
				Status:  "no_tls",
			})
			continue
		}

		caExists := fileExists(caBundle)
		status := "ok"
		if !caExists {
			status = "missing"
		}

		upstreams = append(upstreams, EnvoyUpstreamTLS{
			Name:    name,
			TLSMode: "internal",
			CAPath:  caBundle,
			Exists:  caExists,
			Status:  status,
			CARef:   "internal_ca",
		})
	}

	return upstreams
}

// collectSecrets builds the SDS secret inventory with consumer lists.
func collectSecrets(
	xdsServerCert, xdsServerKey, caBundle string,
	extDomains []ExternalDomainTLS,
	upstreams []EnvoyUpstreamTLS,
) []EnvoySecretRef {
	var secrets []EnvoySecretRef

	// Internal server cert secret
	internalConsumers := []string{"internal_https listener"}
	for _, u := range upstreams {
		if u.TLSMode == "internal" {
			internalConsumers = append(internalConsumers, u.Name)
		}
	}

	secrets = append(secrets, EnvoySecretRef{
		Name:      "internal-server-cert",
		Type:      "tls_certificate",
		CertPath:  xdsServerCert,
		KeyPath:   xdsServerKey,
		Exists:    fileExists(xdsServerCert) && fileExists(xdsServerKey),
		Status:    secretStatus(xdsServerCert, xdsServerKey, ""),
		Consumers: internalConsumers,
	})

	// Internal CA bundle secret
	caConsumers := []string{}
	for _, u := range upstreams {
		if u.TLSMode == "internal" {
			caConsumers = append(caConsumers, u.Name)
		}
	}

	secrets = append(secrets, EnvoySecretRef{
		Name:      "internal-ca-bundle",
		Type:      "validation_context",
		CAPath:    caBundle,
		Exists:    fileExists(caBundle),
		Status:    secretStatus("", "", caBundle),
		Consumers: caConsumers,
	})

	// External domain cert secrets
	for _, ext := range extDomains {
		certPath := ""
		if ext.LeafCert != nil {
			certPath = ext.LeafCert.Path
		}
		secretName := fmt.Sprintf("ext-cert/%s", ext.FQDN)
		consumers := []string{fmt.Sprintf("ext_%s listener", ext.FQDN)}

		secrets = append(secrets, EnvoySecretRef{
			Name:      secretName,
			Type:      "tls_certificate",
			CertPath:  certPath,
			KeyPath:   ext.KeyPath,
			CAPath:    ext.ChainPath,
			Exists:    fileExists(certPath) && fileExists(ext.KeyPath),
			Status:    secretStatus(certPath, ext.KeyPath, ext.ChainPath),
			Consumers: consumers,
		})
	}

	return secrets
}

// secretStatus returns the status for an SDS secret based on file existence.
func secretStatus(certPath, keyPath, caPath string) string {
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
	if allExist {
		return "ok"
	}
	return "missing"
}

// ── Enriched Envoy warnings ─────────────────────────────────────────────────

// collectEnvoyEnrichedWarnings generates warnings from the enriched Envoy state.
func collectEnvoyEnrichedWarnings(state EnvoyStateEnriched) []Warning {
	var ws []Warning

	// Listener warnings
	for _, l := range state.Listeners {
		if l.Status == "missing" {
			ws = append(ws, Warning{
				Severity: "error",
				Message:  fmt.Sprintf("Listener %q: certificate files missing", l.Name),
			})
		}
	}

	// Upstream warnings
	for _, u := range state.Upstreams {
		if u.Status == "missing" {
			ws = append(ws, Warning{
				Severity: "warning",
				Message:  fmt.Sprintf("Upstream %q: CA bundle missing for TLS verification", u.Name),
			})
		}
	}

	// Secret warnings
	for _, s := range state.Secrets {
		if s.Status == "missing" {
			ws = append(ws, Warning{
				Severity: "error",
				Message: fmt.Sprintf("SDS secret %q: backing files missing (consumers: %s)",
					s.Name, strings.Join(s.Consumers, ", ")),
			})
		}
	}

	// xDS client warnings
	if !state.XDSClient.Exists {
		ws = append(ws, Warning{
			Severity: "error",
			Message:  "xDS client TLS material missing — Envoy cannot connect to xDS server",
		})
	}

	// Cross-check: public listener using internal cert
	for _, l := range state.Listeners {
		if l.TLSMode == "public" && l.CertificateRef == "internal_service" {
			ws = append(ws, Warning{
				Severity: "warning",
				Message:  fmt.Sprintf("Listener %q: public listener using internal certificate", l.Name),
			})
		}
	}

	// Cross-check: internal cert path inside domains/ dir (probably wrong)
	for _, l := range state.Listeners {
		if l.TLSMode == "internal" && l.CertPath != "" && strings.Contains(l.CertPath, filepath.Join("domains", "")) {
			ws = append(ws, Warning{
				Severity: "warning",
				Message:  fmt.Sprintf("Listener %q: internal listener using external domain cert path", l.Name),
			})
		}
	}

	return ws
}
