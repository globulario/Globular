package controlplane

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"google.golang.org/protobuf/types/known/anypb"
)

// MakeSecret builds an xDS Secret resource from certificate files.
// This enables Envoy SDS (Secret Discovery Service) for dynamic TLS certificate delivery.
//
// Parameters:
//   - name: Logical secret name (e.g., "internal-server-cert", "internal-ca-bundle")
//   - certPath: Path to certificate file (PEM format), empty string if CA-only secret
//   - keyPath: Path to private key file (PEM format), empty string if CA-only secret
//   - caPath: Path to CA bundle (PEM format), empty string if cert-only secret
//
// Returns an xDS Secret resource ready to be included in snapshot cache.
func MakeSecret(name string, certPath, keyPath, caPath string) (*tls_v3.Secret, error) {
	if name == "" {
		return nil, fmt.Errorf("secret name is required")
	}

	secret := &tls_v3.Secret{
		Name: name,
	}

	// Determine secret type based on provided paths
	hasCert := certPath != "" && keyPath != ""
	hasCA := caPath != ""

	if !hasCert && !hasCA {
		return nil, fmt.Errorf("secret %q must have either cert+key or CA", name)
	}

	if hasCert {
		// TLS certificate secret (server or client cert with private key)
		tlsCert, err := MakeTLSCertificate(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("build TLS certificate for %q: %w", name, err)
		}

		secret.Type = &tls_v3.Secret_TlsCertificate{
			TlsCertificate: tlsCert,
		}
	} else {
		// Validation context secret (CA bundle only)
		validationCtx, err := MakeCertificateValidationContext(caPath)
		if err != nil {
			return nil, fmt.Errorf("build validation context for %q: %w", name, err)
		}

		secret.Type = &tls_v3.Secret_ValidationContext{
			ValidationContext: validationCtx,
		}
	}

	return secret, nil
}

// MakeTLSCertificate builds a TlsCertificate from PEM-encoded certificate and key files.
// The certificate chain and private key are read from disk and embedded inline.
func MakeTLSCertificate(certPath, keyPath string) (*tls_v3.TlsCertificate, error) {
	// Read certificate chain (may contain intermediate certs)
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("read cert %s: %w", certPath, err)
	}
	if len(certPEM) == 0 {
		return nil, fmt.Errorf("cert file %s is empty", certPath)
	}

	// Read private key
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read key %s: %w", keyPath, err)
	}
	if len(keyPEM) == 0 {
		return nil, fmt.Errorf("key file %s is empty", keyPath)
	}

	return &tls_v3.TlsCertificate{
		CertificateChain: &core_v3.DataSource{
			Specifier: &core_v3.DataSource_InlineBytes{
				InlineBytes: certPEM,
			},
		},
		PrivateKey: &core_v3.DataSource{
			Specifier: &core_v3.DataSource_InlineBytes{
				InlineBytes: keyPEM,
			},
		},
	}, nil
}

// MakeCertificateValidationContext builds a validation context from a CA bundle.
// This is used for validating peer certificates (upstream TLS or mTLS).
func MakeCertificateValidationContext(caPath string) (*tls_v3.CertificateValidationContext, error) {
	// Read CA bundle (may contain multiple CA certs)
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read CA %s: %w", caPath, err)
	}
	if len(caPEM) == 0 {
		return nil, fmt.Errorf("CA file %s is empty", caPath)
	}

	return &tls_v3.CertificateValidationContext{
		TrustedCa: &core_v3.DataSource{
			Specifier: &core_v3.DataSource_InlineBytes{
				InlineBytes: caPEM,
			},
		},
	}, nil
}

// HashSecret computes a content hash of a Secret resource for version tracking.
// This enables change detection: same cert = same hash, different cert = different hash.
//
// The hash is computed over the secret's serialized content, ensuring any change
// (cert rotation, key change, CA update) results in a new version.
func HashSecret(secret *tls_v3.Secret) (string, error) {
	if secret == nil {
		return "", fmt.Errorf("secret is nil")
	}

	// Serialize secret to binary protobuf
	data, err := anypb.New(secret)
	if err != nil {
		return "", fmt.Errorf("serialize secret: %w", err)
	}

	// Compute SHA256 hash
	h := sha256.New()
	h.Write(data.Value)
	hash := h.Sum(nil)

	return hex.EncodeToString(hash), nil
}

// SecretVersion generates a version string for xDS snapshot versioning.
// Format: "sds-{hash[:16]}" (first 16 chars of hash for brevity)
func SecretVersion(secret *tls_v3.Secret) (string, error) {
	hash, err := HashSecret(secret)
	if err != nil {
		return "", err
	}

	// Use first 16 chars for readability in logs
	if len(hash) > 16 {
		hash = hash[:16]
	}

	return fmt.Sprintf("sds-%s", hash), nil
}

// ValidateCertKeyPair performs basic validation that cert and key are valid PEM.
// This catches obvious errors before pushing to Envoy.
func ValidateCertKeyPair(certPEM, keyPEM []byte) error {
	if len(certPEM) == 0 {
		return fmt.Errorf("certificate is empty")
	}
	if len(keyPEM) == 0 {
		return fmt.Errorf("private key is empty")
	}

	// Basic PEM format checks
	if !isPEMFormat(certPEM) {
		return fmt.Errorf("certificate is not in PEM format")
	}
	if !isPEMFormat(keyPEM) {
		return fmt.Errorf("private key is not in PEM format")
	}

	return nil
}

// isPEMFormat checks if data looks like PEM (starts with -----BEGIN)
func isPEMFormat(data []byte) bool {
	return len(data) > 0 && (data[0] == '-' || (len(data) > 10 && string(data[:10]) == "-----BEGIN"))
}
