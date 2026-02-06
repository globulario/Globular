package secrets

// Secret name constants for SDS (Secret Discovery Service).
// These names are referenced in:
// - xDS snapshots (secret resources)
// - Envoy TLS contexts (SDS secret references)
// - Certificate rotation logic
const (
	// InternalServerCert is the server certificate for *.globular.internal domains.
	// Used for downstream TLS (clients → Envoy ingress).
	// Source: Cluster local CA
	InternalServerCert = "internal-server-cert"

	// InternalCABundle is the CA bundle for validating internal services.
	// Used for upstream TLS (Envoy → internal services).
	// Source: Cluster local CA
	InternalCABundle = "internal-ca-bundle"

	// InternalClientCert is the client certificate for mTLS to internal services.
	// Used for upstream client authentication (Envoy → services requiring mTLS).
	// Source: Cluster local CA
	// Status: Optional, reserved for future use
	InternalClientCert = "internal-client-cert"

	// PublicServerCert is the ACME certificate for public domains.
	// Used for downstream TLS on public-facing ingress.
	// Source: Let's Encrypt / ACME provider
	// Status: Optional, for clusters with public ingress
	PublicServerCert = "public-server-cert"

	// PublicIngressCert is the primary ACME certificate for public HTTPS ingress.
	// Used for downstream TLS on the main public listener.
	// Source: ACME (Let's Encrypt) fullchain.pem + privkey.pem
	// Rotation: Triggered by ACME renewal (file hash change)
	PublicIngressCert = "public-ingress-cert"
)
