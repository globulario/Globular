package identity

import (
	"strings"
	"testing"
)

func TestGeneratedIDsAreOpaqueAndPrefixed(t *testing.T) {
	principal, err := NewPrincipalID("usr")
	if err != nil {
		t.Fatalf("NewPrincipalID: %v", err)
	}
	if !strings.HasPrefix(principal.String(), "usr_") {
		t.Fatalf("principal prefix=%q want usr_", principal)
	}
	if strings.Contains(principal.String(), "@") {
		t.Fatalf("principal contains domain separator: %q", principal)
	}

	cluster, err := NewClusterID()
	if err != nil {
		t.Fatalf("NewClusterID: %v", err)
	}
	if !strings.HasPrefix(cluster.String(), "cls_") {
		t.Fatalf("cluster prefix=%q want cls_", cluster)
	}
	if strings.Contains(cluster.String(), "@") {
		t.Fatalf("cluster contains domain separator: %q", cluster)
	}

	service, err := NewServiceID()
	if err != nil {
		t.Fatalf("NewServiceID: %v", err)
	}
	if !strings.HasPrefix(service.String(), "svc_") {
		t.Fatalf("service prefix=%q want svc_", service)
	}
	if strings.Contains(service.String(), "@") {
		t.Fatalf("service contains domain separator: %q", service)
	}
}

func TestPrincipalIDValidateRejectsTooShort(t *testing.T) {
	for _, id := range []PrincipalID{
		"usr_short",
		"bad_0123456789abcdef",
		"usr_0123456789abcdeg",
	} {
		if err := id.Validate(); err == nil {
			t.Fatalf("Validate should reject malformed ID %q", id)
		}
	}
	if err := PrincipalID("usr_0123456789abcdef").Validate(); err != nil {
		t.Fatalf("Validate rejected well-formed user ID: %v", err)
	}
	if err := PrincipalID("svc_0123456789abcdef").Validate(); err != nil {
		t.Fatalf("Validate rejected well-formed service ID: %v", err)
	}
}
