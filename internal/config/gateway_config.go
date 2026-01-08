package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// GatewayConfig defines the settings used by the gateway when bootstrapping from
// an explicit configuration file.
type GatewayConfig struct {
	Mode          string `json:"mode"`
	EnvoyHTTPAddr string `json:"envoy_http_addr"`
	MaxUpload     int64  `json:"max_upload"`
	RateRPS       int    `json:"rate_rps"`
	RateBurst     int    `json:"rate_burst"`
	HTTPPort      int    `json:"http_port"`
	HTTPSPort     int    `json:"https_port"`
	Domain        string `json:"domain,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

// DefaultGatewayConfig returns the defaults used by the gateway CLI.
func DefaultGatewayConfig() GatewayConfig {
	return GatewayConfig{
		Mode:          "direct",
		EnvoyHTTPAddr: "127.0.0.1:8080",
		MaxUpload:     2 << 30,
		RateRPS:       50,
		RateBurst:     200,
		HTTPPort:      8080,
		HTTPSPort:     8443,
	}
}

// LoadGatewayConfig reads the JSON config at path and merges it with the defaults.
func LoadGatewayConfig(path string) (GatewayConfig, error) {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return GatewayConfig{}, fmt.Errorf("config path is empty")
	}
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return GatewayConfig{}, err
	}
	cfg := DefaultGatewayConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return GatewayConfig{}, fmt.Errorf("parse %s: %w", cleanPath, err)
	}
	return cfg, nil
}

// Validate ensures the configuration contains the fields required for startup.
func (c GatewayConfig) Validate() error {
	mode := strings.ToLower(strings.TrimSpace(c.Mode))
	switch mode {
	case "direct", "mesh":
	default:
		return fmt.Errorf("mode must be \"direct\" or \"mesh\"")
	}

	if mode == "mesh" && strings.TrimSpace(c.EnvoyHTTPAddr) == "" {
		return fmt.Errorf("envoy_http_addr is required when mode=mesh")
	}

	if c.MaxUpload <= 0 {
		return fmt.Errorf("max_upload must be > 0")
	}
	if c.RateRPS < 0 {
		return fmt.Errorf("rate_rps must be >= 0")
	}
	if c.RateBurst < 0 {
		return fmt.Errorf("rate_burst must be >= 0")
	}
	if c.HTTPPort < 1 || c.HTTPPort > 65535 {
		return fmt.Errorf("http_port must be between 1 and 65535")
	}
	if c.HTTPSPort < 1 || c.HTTPSPort > 65535 {
		return fmt.Errorf("https_port must be between 1 and 65535")
	}
	return nil
}
