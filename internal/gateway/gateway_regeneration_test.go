package gateway

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/globulario/Globular/internal/globule"
	"github.com/globulario/services/golang/config"
)

func writeGatewayConfig(t *testing.T, cfgDir, domain, protocol string) {
	t.Helper()
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	payload := map[string]any{
		"Domain":   domain,
		"Protocol": protocol,
	}
	data, _ := json.Marshal(payload)
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), data, 0o644); err != nil {
		t.Fatalf("write gateway config: %v", err)
	}
}

func ensureGatewayTLS(t *testing.T) {
	t.Helper()
	_, certPath, keyPath, caPath := config.CanonicalTLSPaths(config.GetRuntimeConfigDir())
	for _, p := range []string{certPath, keyPath, caPath} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("mkdir tls dir: %v", err)
		}
		if err := os.WriteFile(p, []byte("dummy"), 0o644); err != nil {
			t.Fatalf("write tls file: %v", err)
		}
	}
}

func TestGatewayRespectsNetworkConfigOnRestart(t *testing.T) {
	stateDir := filepath.Join(t.TempDir(), "state")
	cfgDir := filepath.Join(t.TempDir(), "etc", "globular")
	t.Setenv("GLOBULAR_STATE_DIR", stateDir)
	t.Setenv("GLOBULAR_CONFIG_DIR", cfgDir)
	t.Setenv("GLOBULAR_RUNTIME_CONFIG_DIR", filepath.Join(stateDir, "config"))
	t.Setenv("GLOBULAR_NODE_ID", "test-node")
	t.Setenv("GLOBULAR_SKIP_PEER_KEYS", "1")

	writeGatewayConfig(t, cfgDir, "globular.io", "https")
	writeGatewayConfig(t, stateDir, "globular.io", "https")
	ensureGatewayTLS(t)
	if _, err := config.GetLocalConfig(false); err != nil {
		t.Fatalf("load config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	g1 := globule.New(logger)
	if err := g1.InitFS(); err != nil {
		t.Fatalf("init fs: %v", err)
	}
	if g1.Domain != "globular.io" || strings.ToLower(g1.Protocol) != "https" {
		t.Fatalf("unexpected domain/protocol: %s %s", g1.Domain, g1.Protocol)
	}

	// Change network config and re-init (simulates restart)
	writeGatewayConfig(t, cfgDir, "localhost", "http")
	writeGatewayConfig(t, stateDir, "localhost", "http")
	if _, err := config.GetLocalConfig(false); err != nil {
		t.Fatalf("reload config: %v", err)
	}
	cfg, err := config.GetLocalConfig(false)
	if err != nil {
		t.Fatalf("final reload: %v", err)
	}
	if cfg["Domain"] != "localhost" || strings.ToLower(cfg["Protocol"].(string)) != "http" {
		t.Fatalf("expected updated domain/protocol, got %v/%v", cfg["Domain"], cfg["Protocol"])
	}
}
