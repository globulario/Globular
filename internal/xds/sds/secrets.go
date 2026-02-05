package sds

import (
	"fmt"
	"log"
	"os"

	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/globulario/Globular/internal/controlplane"
)

// SecretName constants define the logical names for SDS secrets.
const (
	InternalServerCert = "internal-server-cert" // Server cert for *.globular.internal
	InternalCABundle   = "internal-ca-bundle"   // CA bundle for validating internal services
	InternalClientCert = "internal-client-cert" // (Future) mTLS client identity
	PublicServerCert   = "public-server-cert"   // ACME cert for public domain
)

// CertPaths holds file paths for a certificate bundle.
type CertPaths struct {
	CertFile string // Certificate chain (PEM)
	KeyFile  string // Private key (PEM)
	CAFile   string // CA bundle (PEM)
}

// BuildInternalServerSecret creates the internal-server-cert secret from canonical TLS paths.
// This secret is used for downstream TLS (clients → Envoy) on internal domains.
func BuildInternalServerSecret(paths CertPaths) (*tls_v3.Secret, error) {
	if paths.CertFile == "" || paths.KeyFile == "" {
		return nil, fmt.Errorf("cert and key files required for server secret")
	}

	// Verify files exist before building
	if !fileExists(paths.CertFile) {
		return nil, fmt.Errorf("cert file not found: %s", paths.CertFile)
	}
	if !fileExists(paths.KeyFile) {
		return nil, fmt.Errorf("key file not found: %s", paths.KeyFile)
	}

	secret, err := controlplane.MakeSecret(InternalServerCert, paths.CertFile, paths.KeyFile, "")
	if err != nil {
		return nil, fmt.Errorf("build internal server secret: %w", err)
	}

	log.Printf("sds: built internal-server-cert from %s", paths.CertFile)
	return secret, nil
}

// BuildInternalCASecret creates the internal-ca-bundle secret from canonical CA path.
// This secret is used for upstream TLS validation (Envoy → services).
func BuildInternalCASecret(caPath string) (*tls_v3.Secret, error) {
	if caPath == "" {
		return nil, fmt.Errorf("CA file required for CA bundle secret")
	}

	if !fileExists(caPath) {
		return nil, fmt.Errorf("CA file not found: %s", caPath)
	}

	secret, err := controlplane.MakeSecret(InternalCABundle, "", "", caPath)
	if err != nil {
		return nil, fmt.Errorf("build internal CA secret: %w", err)
	}

	log.Printf("sds: built internal-ca-bundle from %s", caPath)
	return secret, nil
}

// BuildPublicServerSecret creates the public-server-cert secret from ACME paths.
// This secret is used for public-facing HTTPS (external clients → Envoy).
func BuildPublicServerSecret(paths CertPaths) (*tls_v3.Secret, error) {
	if paths.CertFile == "" || paths.KeyFile == "" {
		return nil, fmt.Errorf("cert and key files required for public server secret")
	}

	if !fileExists(paths.CertFile) {
		return nil, fmt.Errorf("public cert file not found: %s", paths.CertFile)
	}
	if !fileExists(paths.KeyFile) {
		return nil, fmt.Errorf("public key file not found: %s", paths.KeyFile)
	}

	secret, err := controlplane.MakeSecret(PublicServerCert, paths.CertFile, paths.KeyFile, "")
	if err != nil {
		return nil, fmt.Errorf("build public server secret: %w", err)
	}

	log.Printf("sds: built public-server-cert from %s", paths.CertFile)
	return secret, nil
}

// BuildAllSecrets builds all available secrets from the provided paths.
// Returns a map of secret name to Secret resource.
// Skips secrets where paths are not provided (e.g., public cert if ACME disabled).
func BuildAllSecrets(internalPaths, publicPaths CertPaths, caPath string) (map[string]*tls_v3.Secret, error) {
	secrets := make(map[string]*tls_v3.Secret)

	// Internal server cert (required for internal TLS)
	if internalPaths.CertFile != "" && internalPaths.KeyFile != "" {
		secret, err := BuildInternalServerSecret(internalPaths)
		if err != nil {
			return nil, fmt.Errorf("build internal server secret: %w", err)
		}
		secrets[InternalServerCert] = secret
	}

	// Internal CA bundle (required for upstream validation)
	if caPath != "" {
		secret, err := BuildInternalCASecret(caPath)
		if err != nil {
			return nil, fmt.Errorf("build internal CA secret: %w", err)
		}
		secrets[InternalCABundle] = secret
	}

	// Public server cert (optional, only if ACME enabled)
	if publicPaths.CertFile != "" && publicPaths.KeyFile != "" {
		secret, err := BuildPublicServerSecret(publicPaths)
		if err != nil {
			// Log but don't fail - public cert is optional
			log.Printf("sds: skipping public cert (ACME disabled?): %v", err)
		} else {
			secrets[PublicServerCert] = secret
		}
	}

	if len(secrets) == 0 {
		return nil, fmt.Errorf("no secrets could be built from provided paths")
	}

	log.Printf("sds: built %d secrets", len(secrets))
	return secrets, nil
}

// fileExists checks if a file exists and is readable.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
