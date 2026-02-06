package identity

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// PrincipalID is an opaque, stable identifier for users, services, or other principals.
// It MUST NOT contain domain, host, or any network-derived information.
// Format: "usr_" or "svc_" prefix + 16 hex characters
type PrincipalID string

// Principal represents an authenticated entity in the system.
// Identity is stable and independent of network configuration.
type Principal struct {
	ID        PrincipalID `json:"principal_id"`
	Username  string      `json:"username"`        // Display name, not identity
	Email     string      `json:"email,omitempty"` // Contact, not identity
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ClusterID is an opaque, stable identifier for a cluster.
type ClusterID string

// ClusterIdentity represents cluster network configuration.
// Domain values are routing labels only, never used for identity or authorization.
type ClusterIdentity struct {
	ID             ClusterID `json:"cluster_id"`
	ClusterDomain  string    `json:"cluster_domain"`  // Internal DNS only (e.g., "cluster.local")
	IngressDomains []string  `json:"ingress_domains"` // Public routing only
	CreatedAt      time.Time `json:"created_at"`
}

// ServiceID is an opaque identifier for a service instance.
type ServiceID string

// NewPrincipalID generates a new opaque principal identifier.
// Prefix should be "usr" for users, "svc" for services.
func NewPrincipalID(prefix string) (PrincipalID, error) {
	b := make([]byte, 8) // 16 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate principal ID: %w", err)
	}
	return PrincipalID(fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(b))), nil
}

// NewClusterID generates a new opaque cluster identifier.
func NewClusterID() (ClusterID, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate cluster ID: %w", err)
	}
	return ClusterID(fmt.Sprintf("cls_%s", hex.EncodeToString(b))), nil
}

// NewServiceID generates a new opaque service identifier.
func NewServiceID() (ServiceID, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate service ID: %w", err)
	}
	return ServiceID(fmt.Sprintf("svc_%s", hex.EncodeToString(b))), nil
}

// Validate checks if a PrincipalID has valid format.
func (id PrincipalID) Validate() error {
	if len(id) < 20 { // "usr_" + 16 hex chars
		return fmt.Errorf("invalid principal ID: too short")
	}
	// Basic format check - should have prefix and hex suffix
	return nil
}

// String returns the string representation.
func (id PrincipalID) String() string {
	return string(id)
}

// String returns the string representation.
func (id ClusterID) String() string {
	return string(id)
}

// String returns the string representation.
func (id ServiceID) String() string {
	return string(id)
}
