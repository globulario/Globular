package main

import (
	"io"
	"log/slog"
	"testing"

	gatewayconfig "github.com/globulario/Globular/internal/config"
	globpkg "github.com/globulario/Globular/internal/globule"
)

func testGatewayLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
}

func TestApplyGatewayConfigToGlobuleIgnoresDomainAndProtocol(t *testing.T) {
	g := globpkg.New(testGatewayLogger())
	g.Domain = "cluster.example"
	g.Protocol = "https"
	g.PortHTTP = 80
	g.PortHTTPS = 443

	applyGatewayConfigToGlobule(testGatewayLogger(), g, gatewayconfig.GatewayConfig{
		Domain:    "local.override",
		Protocol:  "http",
		HTTPPort:  8080,
		HTTPSPort: 8443,
	})

	if g.Domain != "cluster.example" {
		t.Fatalf("Domain=%q want cluster.example", g.Domain)
	}
	if g.Protocol != "https" {
		t.Fatalf("Protocol=%q want https", g.Protocol)
	}
}

func TestApplyGatewayConfigToGlobuleAppliesPortOverrides(t *testing.T) {
	g := globpkg.New(testGatewayLogger())
	g.PortHTTP = 80
	g.PortHTTPS = 443

	applyGatewayConfigToGlobule(testGatewayLogger(), g, gatewayconfig.GatewayConfig{
		HTTPPort:  18080,
		HTTPSPort: 18443,
	})

	if g.PortHTTP != 18080 {
		t.Fatalf("PortHTTP=%d want 18080", g.PortHTTP)
	}
	if g.PortHTTPS != 18443 {
		t.Fatalf("PortHTTPS=%d want 18443", g.PortHTTPS)
	}
}
