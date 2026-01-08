package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	servicesConfig "github.com/globulario/services/golang/config"
)

func TestLoadMinioProxyConfigContractOrder(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.json")
	second := filepath.Join(dir, "second.json")
	writeContractFile(t, first, &servicesConfig.MinioProxyConfig{
		Endpoint: "first.example",
		Bucket:   "bucket",
	})
	writeContractFile(t, second, &servicesConfig.MinioProxyConfig{
		Endpoint: "second.example",
		Bucket:   "bucket",
	})

	oldPaths := minioContractPaths
	t.Cleanup(func() { minioContractPaths = oldPaths })
	minioContractPaths = []string{first, second}

	cfg, err := LoadMinioProxyConfig()
	if err != nil {
		t.Fatalf("load contract: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config")
	}
	if cfg.Endpoint != "first.example" {
		t.Fatalf("unexpected endpoint %s", cfg.Endpoint)
	}
}

func TestLoadMinioProxyConfigEnvFallback(t *testing.T) {
	tempDir := t.TempDir()
	missing := filepath.Join(tempDir, "missing.json")
	oldPaths := minioContractPaths
	t.Cleanup(func() { minioContractPaths = oldPaths })
	minioContractPaths = []string{missing}

	t.Setenv("GLOBULAR_MINIO_ENDPOINT", " https://env.example ")
	t.Setenv("GLOBULAR_MINIO_BUCKET", " bucket ")
	t.Setenv("GLOBULAR_MINIO_PREFIX", " /env/prefix/ ")
	t.Setenv("GLOBULAR_MINIO_SECURE", "false")
	t.Setenv("GLOBULAR_MINIO_ACCESS_KEY", "ak")
	t.Setenv("GLOBULAR_MINIO_SECRET_KEY", "sk")
	t.Setenv("GLOBULAR_MINIO_CA_BUNDLE", " /tmp/env.pem ")

	cfg, err := LoadMinioProxyConfig()
	if err != nil {
		t.Fatalf("env load: %v", err)
	}
	if cfg.Endpoint != "https://env.example" {
		t.Fatalf("unexpected endpoint %s", cfg.Endpoint)
	}
	if cfg.Prefix != "env/prefix" {
		t.Fatalf("unexpected prefix %s", cfg.Prefix)
	}
	if cfg.Secure {
		t.Fatalf("secure should be false")
	}
	if cfg.Auth == nil || cfg.Auth.Mode != servicesConfig.MinioProxyAuthModeAccessKey {
		t.Fatalf("unexpected auth %+v", cfg.Auth)
	}
	if cfg.CABundlePath != "/tmp/env.pem" {
		t.Fatalf("unexpected CA bundle %s", cfg.CABundlePath)
	}
}

func TestLoadMinioProxyConfigLegacyFallback(t *testing.T) {
	tempDir := t.TempDir()
	missing := filepath.Join(tempDir, "missing.json")
	oldPaths := minioContractPaths
	t.Cleanup(func() { minioContractPaths = oldPaths })
	minioContractPaths = []string{missing}

	origGet := getServiceConfigurationByID
	t.Cleanup(func() { getServiceConfigurationByID = origGet })
	getServiceConfigurationByID = func(id string) (map[string]any, error) {
		if id != "file.FileService" {
			t.Fatalf("unexpected service id %s", id)
		}
		return map[string]any{
			"UseMinio":       true,
			"MinioEndpoint":  "legacy.example",
			"MinioBucket":    "bucket",
			"MinioPrefix":    "/legacy/",
			"MinioUseSSL":    true,
			"MinioAccessKey": "ak",
			"MinioSecretKey": "sk",
		}, nil
	}

	t.Setenv("GLOBULAR_MINIO_ENDPOINT", "")
	t.Setenv("GLOBULAR_MINIO_BUCKET", "")

	cfg, err := LoadMinioProxyConfig()
	if err != nil {
		t.Fatalf("legacy load: %v", err)
	}
	if cfg.Endpoint != "legacy.example" {
		t.Fatalf("unexpected endpoint %s", cfg.Endpoint)
	}
	if cfg.Prefix != "legacy" {
		t.Fatalf("unexpected prefix %s", cfg.Prefix)
	}
	if cfg.Auth == nil || cfg.Auth.Mode != servicesConfig.MinioProxyAuthModeAccessKey {
		t.Fatalf("unexpected auth %+v", cfg.Auth)
	}
}

func TestLoadMinioProxyConfigInvalidContract(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "minio.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("write invalid contract: %v", err)
	}

	oldPaths := minioContractPaths
	t.Cleanup(func() { minioContractPaths = oldPaths })
	minioContractPaths = []string{path}

	_, err := LoadMinioProxyConfig()
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, servicesConfig.ErrInvalidObjectStoreContract) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveMinioProxyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	savePath := filepath.Join(tmpDir, "config", "minio.json")
	oldSave := minioContractSavePath
	oldPaths := minioContractPaths
	t.Cleanup(func() {
		minioContractSavePath = oldSave
		minioContractPaths = oldPaths
	})
	minioContractSavePath = savePath
	minioContractPaths = []string{savePath}

	cfg := &servicesConfig.MinioProxyConfig{
		Endpoint:     "https://proxy",
		Bucket:       "bucket",
		Prefix:       "/saved/",
		Secure:       true,
		CABundlePath: "/tmp/ca",
		Auth: &servicesConfig.MinioProxyAuth{
			Mode:      servicesConfig.MinioProxyAuthModeAccessKey,
			AccessKey: "ak",
			SecretKey: "sk",
		},
	}

	if err := SaveMinioProxyConfig(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	info, err := os.Stat(savePath)
	if err != nil {
		t.Fatalf("stat saved file: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("unexpected perms %o", info.Mode().Perm())
	}

	loaded, err := LoadMinioProxyConfig()
	if err != nil {
		t.Fatalf("load saved config: %v", err)
	}
	if loaded.Endpoint != "https://proxy" {
		t.Fatalf("endpoint mismatch %s", loaded.Endpoint)
	}
	if loaded.Prefix != "saved" {
		t.Fatalf("prefix mismatch %s", loaded.Prefix)
	}
}

func writeContractFile(t *testing.T, path string, cfg *servicesConfig.MinioProxyConfig) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir contract dir: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create contract file: %v", err)
	}
	defer f.Close()
	if err := servicesConfig.SaveMinioProxyConfigTo(f, cfg); err != nil {
		t.Fatalf("save contract to %s: %v", path, err)
	}
}
