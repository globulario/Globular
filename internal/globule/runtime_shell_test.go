package globule

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	servicesconfig "github.com/globulario/services/golang/config"
)

func newTestGlobule() *Globule {
	return New(slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))
}

func TestSetConfigIgnoresControllerManagedNetworkIdentity(t *testing.T) {
	stateDir := filepath.Join(t.TempDir(), "state")
	cfgDir := filepath.Join(t.TempDir(), "etc", "globular")
	t.Setenv("GLOBULAR_STATE_DIR", stateDir)
	t.Setenv("GLOBULAR_CONFIG_DIR", cfgDir)
	t.Setenv("GLOBULAR_RUNTIME_CONFIG_DIR", filepath.Join(stateDir, "config"))
	reset := servicesconfig.ResetLocalConfigForTesting()
	defer reset()
	g := newTestGlobule()
	g.Domain = "cluster.example"
	g.Protocol = "https"
	g.PortsRange = "10000-10100"
	g.PortHTTP = 80
	g.PortHTTPS = 443

	cfg := map[string]interface{}{
		"Domain":     "local.override",
		"Protocol":   "http",
		"PortsRange": "20000-20100",
		"PortHTTP":   8080,
		"PortHTTPS":  8443,
	}
	if err := g.SetConfig(cfg); err != nil {
		t.Fatalf("SetConfig: %v", err)
	}
	if g.Domain != "cluster.example" {
		t.Fatalf("Domain=%q want cluster.example", g.Domain)
	}
	if g.Protocol != "https" {
		t.Fatalf("Protocol=%q want https", g.Protocol)
	}
	if g.PortsRange != "10000-10100" {
		t.Fatalf("PortsRange=%q want 10000-10100", g.PortsRange)
	}
	if g.PortHTTP != 8080 {
		t.Fatalf("PortHTTP=%d want 8080", g.PortHTTP)
	}
	if g.PortHTTPS != 8443 {
		t.Fatalf("PortHTTPS=%d want 8443", g.PortHTTPS)
	}
}

func TestNodeAndControllerAddressUseStableDefaults(t *testing.T) {
	t.Setenv(nodeAgentEnvKey, "")
	t.Setenv(controllerAgentEnvKey, "")
	if got := NodeAgentAddress(); got != defaultNodeAgentAddr {
		t.Fatalf("NodeAgentAddress=%q want %q", got, defaultNodeAgentAddr)
	}
	if got := ControllerAddress(); got != defaultControllerAddr {
		t.Fatalf("ControllerAddress=%q want %q", got, defaultControllerAddr)
	}

	t.Setenv(nodeAgentEnvKey, "10.0.0.9:11000")
	t.Setenv(controllerAgentEnvKey, "10.0.0.8:12000")
	if got := NodeAgentAddress(); got != "10.0.0.9:11000" {
		t.Fatalf("NodeAgentAddress override=%q", got)
	}
	if got := ControllerAddress(); got != "10.0.0.8:12000" {
		t.Fatalf("ControllerAddress override=%q", got)
	}
}
