package logbridge

import (
	"log/slog"
	"testing"
)

func TestSplitAttrsSeparatesComponentAndFields(t *testing.T) {
	var r slog.Record
	r.AddAttrs(
		slog.String("component", "gateway"),
		slog.Bool("enabled", true),
	)
	component, fields := splitAttrs([]slog.Attr{slog.Int("workers", 3)}, &r)
	if component != "gateway" {
		t.Fatalf("component=%q want gateway", component)
	}
	if fields["workers"] != "3" {
		t.Fatalf("workers=%q want 3", fields["workers"])
	}
	if fields["enabled"] != "true" {
		t.Fatalf("enabled=%q want true", fields["enabled"])
	}
}

func TestDeriveIDStableForSamePayload(t *testing.T) {
	a := deriveID("globular", "gateway", "Handle", "server.go:10", "hello")
	b := deriveID("globular", "gateway", "Handle", "server.go:10", "hello")
	c := deriveID("globular", "gateway", "Handle", "server.go:10", "different")
	if a != b {
		t.Fatalf("deriveID unstable: %q != %q", a, b)
	}
	if a == c {
		t.Fatalf("deriveID should differ for different payloads: %q == %q", a, c)
	}
}
