package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultGatewayConfigValidation(t *testing.T) {
	if err := DefaultGatewayConfig().Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestLoadGatewayConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "gateway.json")
	content := []byte(`{
		"mode": "mesh",
		"envoy_http_addr": "127.0.0.1:9901",
		"http_port": 9090,
		"https_port": 9443
	}`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadGatewayConfig(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Mode != "mesh" {
		t.Fatalf("expected mode=mesh; got %q", cfg.Mode)
	}
	if cfg.EnvoyHTTPAddr != "127.0.0.1:9901" {
		t.Fatalf("expected envoy address; got %q", cfg.EnvoyHTTPAddr)
	}
	if cfg.HTTPPort != 9090 || cfg.HTTPSPort != 9443 {
		t.Fatalf("expected ports 9090/9443; got %d/%d", cfg.HTTPPort, cfg.HTTPSPort)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}
}

func TestGatewayConfigValidateErrors(t *testing.T) {
	tests := []struct {
		name string
		cfg  GatewayConfig
	}{
		{
			name: "invalid mode",
			cfg:  GatewayConfig{Mode: "unknown"},
		},
		{
			name: "mesh requires envoy",
			cfg: GatewayConfig{
				Mode:      "mesh",
				HTTPPort:  8080,
				HTTPSPort: 8443,
			},
		},
		{
			name: "bad http port",
			cfg: GatewayConfig{
				Mode:          "direct",
				EnvoyHTTPAddr: "127.0.0.1:8080",
				HTTPPort:      0,
				HTTPSPort:     8443,
			},
		},
		{
			name: "bad https port",
			cfg: GatewayConfig{
				Mode:          "direct",
				EnvoyHTTPAddr: "127.0.0.1:8080",
				HTTPPort:      8080,
				HTTPSPort:     70000,
			},
		},
		{
			name: "invalid max upload",
			cfg: GatewayConfig{
				Mode:          "direct",
				EnvoyHTTPAddr: "127.0.0.1:8080",
				HTTPPort:      8080,
				HTTPSPort:     8443,
				MaxUpload:     0,
			},
		},
		{
			name: "negative rate",
			cfg: GatewayConfig{
				Mode:          "direct",
				EnvoyHTTPAddr: "127.0.0.1:8080",
				HTTPPort:      8080,
				HTTPSPort:     8443,
				MaxUpload:     1,
				RateRPS:       -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); err == nil {
				t.Fatalf("expected validation error for %q", tt.name)
			}
		})
	}
}
