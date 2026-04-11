package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	servicesConfig "github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
)

const (
	MinioContractPathVar = "/var/lib/globular/objectstore/minio.json"
)

const minioContractLogTTL = 30 * time.Second

var (
	minioContractPaths    = []string{MinioContractPathVar}
	minioContractSavePath = MinioContractPathVar
	minioContractLogState = struct {
		mu   sync.Mutex
		last time.Time
	}{}
	getServiceConfigurationByID = servicesConfig.GetServiceConfigurationById

	// loadMinioEtcdFallback returns the cluster MinIO config from etcd, which
	// is the authoritative source of truth when the on-disk contract is
	// missing or corrupt. Exposed as a package variable so unit tests can
	// replace it with a stub (the default reads from a live etcd cluster).
	loadMinioEtcdFallback = func() (*servicesConfig.MinioProxyConfig, error) {
		return servicesConfig.BuildMinioProxyConfig()
	}
)

func init() {
	if custom := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_CONTRACT_PATH")); custom != "" {
		minioContractPaths = []string{custom}
		minioContractSavePath = custom
	}
}

// LoadMinioProxyConfig locates the MinIO contract, falls back to etcd / env
// / legacy config, and validates input.
//
// The on-disk contract at /var/lib/globular/objectstore/minio.json is a
// convenience written by the installer and the MinIO package pre-start hook,
// but the authoritative source of MinIO connection info is etcd (the cluster
// config). On nodes where the contract file got corrupted — e.g. overwritten
// with plain-text credentials instead of the JSON shape — we previously
// returned the parse error directly, which bubbled up as a 503 on the admin
// gateway. Treat a parse failure the same way we treat a missing file: log
// it, then fall through to etcd and env fallbacks. This mirrors the design
// principle that etcd is the single source of truth — the file is only a
// hint.
func LoadMinioProxyConfig() (*servicesConfig.MinioProxyConfig, error) {
	if cfg, err := loadMinioContract(); err == nil {
		return cfg, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		slog.Warn("objectstore contract invalid, falling back to etcd/env",
			"paths", strings.Join(minioContractPaths, ","), "err", err)
		// fall through to etcd / env instead of returning the error
	}

	// Fallback 1: etcd cluster config (authoritative).
	if loadMinioEtcdFallback != nil {
		if cfg, err := loadMinioEtcdFallback(); err == nil && cfg != nil {
			slog.Info("objectstore contract loaded from etcd fallback",
				"endpoint", cfg.Endpoint, "bucket", cfg.Bucket, "secure", cfg.Secure)
			return cfg, nil
		} else if err != nil {
			slog.Debug("etcd fallback unavailable", "err", err)
		}
	}

	// Fallback 2: environment variables (legacy).
	slog.Warn("objectstore contract not found; trying env/legacy",
		"paths", strings.Join(minioContractPaths, ","))
	if cfg, err := loadMinioEnvConfig(); err != nil {
		return nil, err
	} else if cfg != nil {
		return cfg, nil
	}
	return nil, os.ErrNotExist
}

func loadMinioContract() (*servicesConfig.MinioProxyConfig, error) {
	for _, path := range minioContractPaths {
		file, err := os.Open(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read object store contract %s: %w", path, err)
		}
		cfg, err := servicesConfig.LoadMinioProxyConfigFrom(file)
		_ = file.Close()
		if err != nil {
			if errors.Is(err, servicesConfig.ErrInvalidObjectStoreContract) {
				logContractParseError(path, err)
			}
			return nil, err
		}
		slog.Info("objectstore contract loaded", "path", path, "endpoint", cfg.Endpoint, "bucket", cfg.Bucket, "secure", cfg.Secure)
		return cfg, nil
	}
	logContractNotFound(minioContractPaths)
	return nil, os.ErrNotExist
}

func logContractParseError(path string, err error) {
	minioContractLogState.mu.Lock()
	defer minioContractLogState.mu.Unlock()
	if time.Since(minioContractLogState.last) < minioContractLogTTL {
		return
	}
	minioContractLogState.last = time.Now()
	slog.Warn("failed to parse object store contract", "path", path, "err", err)
}

func logContractNotFound(paths []string) {
	minioContractLogState.mu.Lock()
	defer minioContractLogState.mu.Unlock()
	if time.Since(minioContractLogState.last) < minioContractLogTTL {
		return
	}
	minioContractLogState.last = time.Now()
	slog.Warn("object store contract not found; serving disk webroot", "paths", strings.Join(paths, ","))
}

func loadMinioEnvConfig() (*servicesConfig.MinioProxyConfig, error) {
	endpoint := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_ENDPOINT"))
	bucket := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_BUCKET"))
	if endpoint == "" || bucket == "" {
		return nil, nil
	}
	secure := true
	if raw := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_SECURE")); raw != "" {
		secure = Utility.ToBool(raw)
	}
	cfg := &servicesConfig.MinioProxyConfig{
		Endpoint:     endpoint,
		Bucket:       bucket,
		Prefix:       strings.Trim(strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_PREFIX")), "/"),
		Secure:       secure,
		CABundlePath: strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_CA_BUNDLE")),
		Auth:         buildEnvMinioAuth(),
	}
	cfg = servicesConfig.NormalizeMinioProxyConfig(cfg)
	if err := servicesConfig.ValidateMinioProxyConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func buildEnvMinioAuth() *servicesConfig.MinioProxyAuth {
	accessKey := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_ACCESS_KEY"))
	secretKey := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_SECRET_KEY"))
	if accessKey == "" || secretKey == "" {
		return nil
	}
	return &servicesConfig.MinioProxyAuth{
		Mode:      servicesConfig.MinioProxyAuthModeAccessKey,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func loadLegacyMinioConfig() (*servicesConfig.MinioProxyConfig, error) {
	cfg, err := getServiceConfigurationByID("file.FileService")
	if err != nil || cfg == nil {
		return nil, err
	}
	if !Utility.ToBool(cfg["UseMinio"]) {
		return nil, nil
	}
	endpoint := strings.TrimSpace(Utility.ToString(cfg["MinioEndpoint"]))
	bucket := strings.TrimSpace(Utility.ToString(cfg["MinioBucket"]))
	if endpoint == "" || bucket == "" {
		return nil, fmt.Errorf("legacy file service missing MinIO endpoint or bucket")
	}
	legacy := &servicesConfig.MinioProxyConfig{
		Endpoint: endpoint,
		Bucket:   bucket,
		Prefix:   strings.Trim(Utility.ToString(cfg["MinioPrefix"]), "/"),
		Secure:   Utility.ToBool(cfg["MinioUseSSL"]),
		Auth:     buildLegacyMinioAuth(Utility.ToString(cfg["MinioAccessKey"]), Utility.ToString(cfg["MinioSecretKey"])),
	}
	legacy = servicesConfig.NormalizeMinioProxyConfig(legacy)
	if err := servicesConfig.ValidateMinioProxyConfig(legacy); err != nil {
		return nil, err
	}
	return legacy, nil
}

func buildLegacyMinioAuth(accessKey, secretKey string) *servicesConfig.MinioProxyAuth {
	ak := strings.TrimSpace(accessKey)
	sk := strings.TrimSpace(secretKey)
	if ak == "" || sk == "" {
		return nil
	}
	return &servicesConfig.MinioProxyAuth{
		Mode:      servicesConfig.MinioProxyAuthModeAccessKey,
		AccessKey: ak,
		SecretKey: sk,
	}
}

// SaveMinioProxyConfig persists the MinIO contract to the canonical location.
func SaveMinioProxyConfig(cfg *servicesConfig.MinioProxyConfig) error {
	if cfg == nil {
		return fmt.Errorf("minio config is nil")
	}
	norm := servicesConfig.NormalizeMinioProxyConfig(cfg)
	if err := servicesConfig.ValidateMinioProxyConfig(norm); err != nil {
		return fmt.Errorf("validate minio config: %w", err)
	}

	dir := filepath.Dir(minioContractSavePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create contract directory %s: %w", dir, err)
	}

	tmpFile, err := os.CreateTemp(dir, "minio.json.*")
	if err != nil {
		return fmt.Errorf("create tmp contract file: %w", err)
	}
	tmpPath := tmpFile.Name()
	done := false
	defer func() {
		if tmpFile != nil {
			_ = tmpFile.Close()
		}
		if !done {
			_ = os.Remove(tmpPath)
		}
	}()

	if err := tmpFile.Chmod(0o600); err != nil {
		return fmt.Errorf("set tmp contract permissions: %w", err)
	}
	if err := servicesConfig.SaveMinioProxyConfigTo(tmpFile, norm); err != nil {
		return fmt.Errorf("write tmp contract file: %w", err)
	}
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("sync tmp contract file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close tmp contract file: %w", err)
	}
	tmpFile = nil

	if err := os.Rename(tmpPath, minioContractSavePath); err != nil {
		return fmt.Errorf("rename contract file: %w", err)
	}
	done = true

	dirFile, err := os.Open(dir)
	if err == nil {
		_ = dirFile.Sync()
		_ = dirFile.Close()
	}
	return nil
}
