package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/resource/resourcepb"
	Utility "github.com/globulario/utility"
)

const (
	// minioConfigCacheTTL controls caching and log rate-limits for object store metadata.
	minioConfigCacheTTL = 30 * time.Second
)

var (
	minioContractPaths = []string{
		"/var/lib/globular/objectstore/minio.json",
		"/etc/globular/objectstore.d/minio.json",
	}
	contractLogState = struct {
		mu   sync.Mutex
		last time.Time
	}{}
)

var ErrInvalidObjectStoreContract = errors.New("invalid object store contract")

// ObjectStoreContract mirrors the minimal JSON shape written by the installer.
type ObjectStoreContract struct {
	Type         string `json:"type"`
	Endpoint     string `json:"endpoint"`
	Bucket       string `json:"bucket"`
	Prefix       string `json:"prefix,omitempty"`
	Secure       bool   `json:"secure"`
	CABundlePath string `json:"caBundlePath,omitempty"`
	Auth         struct {
		Mode      string `json:"mode"`
		AccessKey string `json:"accessKey,omitempty"`
		SecretKey string `json:"secretKey,omitempty"`
		CredFile  string `json:"credFile,omitempty"`
	} `json:"auth"`
}

// LoadMinioProxyConfig locates the MinIO contract, falls back to env/legacy config, and validates input.
func LoadMinioProxyConfig() (*resourcepb.MinioProxyConfig, error) {
	if cfg, err := loadMinioContract(); err == nil {
		return cfg, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if cfg := loadMinioEnvConfig(); cfg != nil {
		return cfg, nil
	}
	return loadLegacyMinioConfig()
}

func loadMinioContract() (*resourcepb.MinioProxyConfig, error) {
	for _, path := range minioContractPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read object store contract %s: %w", path, err)
		}
		cfg, err := parseObjectStoreContract(data)
		if err != nil {
			logContractParseError(path, err)
			return nil, fmt.Errorf("%w: %s", ErrInvalidObjectStoreContract, err)
		}
		return cfg, nil
	}
	return nil, os.ErrNotExist
}

func parseObjectStoreContract(data []byte) (*resourcepb.MinioProxyConfig, error) {
	contract := ObjectStoreContract{Secure: true}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&contract); err != nil {
		return nil, err
	}
	if contract.Type != "" && !strings.EqualFold(strings.TrimSpace(contract.Type), "minio") {
		return nil, fmt.Errorf("unsupported contract type %q", contract.Type)
	}
	if strings.TrimSpace(contract.Endpoint) == "" || strings.TrimSpace(contract.Bucket) == "" {
		return nil, fmt.Errorf("endpoint and bucket are required")
	}
	return contract.toProxyConfig(), nil
}

func logContractParseError(path string, err error) {
	contractLogState.mu.Lock()
	defer contractLogState.mu.Unlock()
	if time.Since(contractLogState.last) < minioConfigCacheTTL {
		return
	}
	contractLogState.last = time.Now()
	slog.Warn("failed to parse object store contract", "path", path, "err", err)
}

func (c ObjectStoreContract) toProxyConfig() *resourcepb.MinioProxyConfig {
	prefix := strings.Trim(c.Prefix, "/")
	auth := toProxyAuth(strings.TrimSpace(strings.ToLower(c.Auth.Mode)), strings.TrimSpace(c.Auth.AccessKey), strings.TrimSpace(c.Auth.SecretKey), strings.TrimSpace(c.Auth.CredFile))
	cfg := &resourcepb.MinioProxyConfig{
		Endpoint:     strings.TrimSpace(c.Endpoint),
		Bucket:       strings.TrimSpace(c.Bucket),
		Prefix:       prefix,
		Secure:       c.Secure,
		CABundlePath: strings.TrimSpace(c.CABundlePath),
		Auth:         auth,
	}
	return cfg
}

func loadMinioEnvConfig() *resourcepb.MinioProxyConfig {
	endpoint := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_ENDPOINT"))
	bucket := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_BUCKET"))
	if endpoint == "" || bucket == "" {
		return nil
	}
	secure := true
	if raw := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_SECURE")); raw != "" {
		secure = Utility.ToBool(raw)
	}
	prefix := strings.Trim(strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_PREFIX")), "/")
	auth := toProxyAuth("", strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_ACCESS_KEY")), strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_SECRET_KEY")), "")
	return &resourcepb.MinioProxyConfig{
		Endpoint:     endpoint,
		Bucket:       bucket,
		Prefix:       prefix,
		Secure:       secure,
		CABundlePath: strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_CA_BUNDLE")),
		Auth:         auth,
	}
}

func loadLegacyMinioConfig() (*resourcepb.MinioProxyConfig, error) {
	cfg, err := config_.GetServiceConfigurationById("file.FileService")
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
	prefix := strings.Trim(Utility.ToString(cfg["MinioPrefix"]), "/")
	secure := Utility.ToBool(cfg["MinioUseSSL"])
	auth := toProxyAuth("", strings.TrimSpace(Utility.ToString(cfg["MinioAccessKey"])), strings.TrimSpace(Utility.ToString(cfg["MinioSecretKey"])), "")
	return &resourcepb.MinioProxyConfig{
		Endpoint: endpoint,
		Bucket:   bucket,
		Prefix:   prefix,
		Secure:   secure,
		Auth:     auth,
	}, nil
}

func toProxyAuth(mode, accessKey, secretKey, credFile string) *resourcepb.MinioProxyAuth {
	accessKey = strings.TrimSpace(accessKey)
	secretKey = strings.TrimSpace(secretKey)
	credFile = strings.TrimSpace(credFile)
	switch mode {
	case "":
		if accessKey != "" && secretKey != "" {
			mode = resourcepb.MinioProxyAuthModeAccessKey
		} else if credFile != "" {
			mode = resourcepb.MinioProxyAuthModeFile
		} else {
			mode = resourcepb.MinioProxyAuthModeNone
		}
	case resourcepb.MinioProxyAuthModeAccessKey, resourcepb.MinioProxyAuthModeFile, resourcepb.MinioProxyAuthModeNone:
	default:
		mode = resourcepb.MinioProxyAuthModeNone
	}
	auth := &resourcepb.MinioProxyAuth{
		Mode:      mode,
		AccessKey: accessKey,
		SecretKey: secretKey,
		CredFile:  credFile,
	}
	if mode == resourcepb.MinioProxyAuthModeNone {
		auth.AccessKey = ""
		auth.SecretKey = ""
		auth.CredFile = ""
	}
	return auth
}
