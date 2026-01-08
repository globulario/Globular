package handlers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	coreConfig "github.com/globulario/Globular/internal/config"
	cfgHandlers "github.com/globulario/Globular/internal/gateway/handlers/config"
	filesHandlers "github.com/globulario/Globular/internal/gateway/handlers/files"
	globpkg "github.com/globulario/Globular/internal/globule"
	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// serveProvider satisfies files.ServeProvider for the gateway surface.
type serveProvider struct {
	globule *globpkg.Globule
}

func (p serveProvider) WebRoot() string                           { return config_.GetWebRootDir() }
func (p serveProvider) DataRoot() string                          { return "" }
func (p serveProvider) CredsDir() string                          { return config_.GetConfigDir() + "/tls" }
func (p serveProvider) IndexApplication() string                  { return p.globule.IndexApplication }
func (p serveProvider) PublicDirs() []string                      { return config_.GetPublicDirs() }
func (p serveProvider) Exists(pth string) bool                    { return Utility.Exists(pth) }
func (p serveProvider) FindHashedFile(pth string) (string, error) { return findHashedFile(pth) }
func (p serveProvider) FileServiceMinioConfig() (*filesHandlers.MinioProxyConfig, error) {
	return fileServiceMinioConfigCache.get()
}
func (p serveProvider) FileServiceMinioConfigStrict(ctx context.Context) (*filesHandlers.MinioProxyConfig, error) {
	return fileServiceMinioConfigCache.getStrict(ctx)
}
func (p serveProvider) ParseUserID(tok string) (string, error) { return tokenParser{}.ParseUserID(tok) }
func (p serveProvider) ValidateAccount(userID, action, reqPath string) (bool, bool, error) {
	return accessControl{globule: p.globule}.ValidateAccount(userID, action, reqPath)
}
func (p serveProvider) ValidateApplication(app, action, reqPath string) (bool, bool, error) {
	return accessControl{globule: p.globule}.ValidateApplication(app, action, reqPath)
}
func (p serveProvider) ResolveImportPath(base, line string) (string, error) {
	return resolveImportPath(base, line)
}
func (p serveProvider) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	return streamHandlerMaybe(name, w, r)
}
func (p serveProvider) ResolveProxy(reqPath string) (string, bool) {
	return proxyResolver{globule: p.globule}.ResolveProxy(reqPath)
}

type proxyResolver struct {
	globule *globpkg.Globule
}

func (r proxyResolver) ResolveProxy(reqPath string) (string, bool) {
	for _, v := range r.globule.ReverseProxies {
		parts := strings.SplitN(strings.TrimSpace(v.(string)), "|", 2)
		if len(parts) != 2 {
			continue
		}
		proxyURLStr := strings.TrimSpace(parts[0])
		proxyPath := strings.TrimSpace(parts[1])
		if strings.HasPrefix(reqPath, proxyPath) {
			return proxyURLStr, true
		}
	}
	return "", false
}

type uploadProvider struct {
	globule *globpkg.Globule
}

func (uploadProvider) DataRoot() string                       { return "" }
func (uploadProvider) PublicDirs() []string                   { return config_.GetPublicDirs() }
func (uploadProvider) ParseUserID(tok string) (string, error) { return tokenParser{}.ParseUserID(tok) }
func (u uploadProvider) ValidateAccount(uID, action, path string) (bool, bool, error) {
	return accessControl{globule: u.globule}.ValidateAccount(uID, action, path)
}
func (u uploadProvider) ValidateApplication(app, action, path string) (bool, bool, error) {
	return accessControl{globule: u.globule}.ValidateApplication(app, action, path)
}
func (u uploadProvider) AddResourceOwner(path, resourceType, owner string) error {
	return u.globule.AddResourceOwner(path, resourceType, owner, rbacpb.SubjectType_ACCOUNT)
}
func (uploadProvider) FileServiceMinioConfig() (*filesHandlers.MinioProxyConfig, error) {
	return fileServiceMinioConfigCache.get()
}

type tokenParser struct{}

func (tokenParser) ParseUserID(tok string) (string, error) {
	claims, err := security.ValidateToken(tok)
	if err != nil {
		return "", err
	}
	return claims.ID + "@" + claims.UserDomain, nil
}

type accessControl struct {
	globule *globpkg.Globule
}

func (a accessControl) ValidateAccount(userID, action, reqPath string) (bool, bool, error) {
	return a.globule.ValidateAccess(userID, rbacpb.SubjectType_ACCOUNT, action, reqPath)
}

func (a accessControl) ValidateApplication(app, action, reqPath string) (bool, bool, error) {
	return a.globule.ValidateAccess(app, rbacpb.SubjectType_APPLICATION, action, reqPath)
}

type cfgProvider struct {
	globule *globpkg.Globule
	cache   *cfgHandlers.ServiceConfigCache
}

func (cfgProvider) Address() (string, error)      { return config_.GetAddress() }
func (cfgProvider) MyIP() string                  { return Utility.MyIP() }
func (p cfgProvider) LocalConfig() map[string]any { return p.globule.GetConfig() }
func (p cfgProvider) ServiceConfig(idOrName string) (map[string]any, error) {
	if p.cache != nil {
		if svc, ok := p.cache.Get(idOrName); ok {
			return svc, nil
		}
	}
	return config_.GetServiceConfigurationById(idOrName)
}
func (cfgProvider) RootDir() string      { return config_.GetRootDir() }
func (cfgProvider) DataDir() string      { return config_.GetDataDir() }
func (cfgProvider) ConfigDir() string    { return config_.GetConfigDir() }
func (cfgProvider) WebRootDir() string   { return config_.GetWebRootDir() }
func (cfgProvider) PublicDirs() []string { return config_.GetPublicDirs() }

type describeProvider struct{}

func (describeProvider) DescribeService(name string, timeout time.Duration) (config_.ServiceDesc, string, error) {
	return config_.ServiceDesc{}, "", fmt.Errorf("describe service not supported in gateway (name=%s)", name)
}

type tokenValidator struct{}

type cfgSaver struct {
	globule *globpkg.Globule
}

func (tokenValidator) Validate(tok string) error {
	_, err := security.ValidateToken(tok)
	return err
}

func (s cfgSaver) Save(m map[string]any) error { return s.globule.SetConfig(m) }

type svcPermsProvider struct {
	globule *globpkg.Globule
}

func (p svcPermsProvider) LoadPermissions(serviceID string) ([]any, error) {
	cfg, err := config_.GetServiceConfigurationById(serviceID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("service not found")
	}
	perms, ok := cfg["Permissions"].([]any)
	if !ok {
		return nil, fmt.Errorf("invalid Permissions format")
	}
	return perms, nil
}

type imgLister struct{}

func (imgLister) ListImages(dir string) ([]string, error) {
	roots := config_.GetPublicDirs()
	ok := false
	cleanDir := filepath.Clean(dir)
	for _, root := range roots {
		root = filepath.Clean(root)
		if cleanDir == root || strings.HasPrefix(cleanDir+string(os.PathSeparator), root+string(os.PathSeparator)) || strings.HasPrefix(cleanDir, root+string(os.PathSeparator)) {
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("dir not allowed")
	}
	var out []string
	err := filepath.WalkDir(cleanDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		switch strings.ToLower(filepath.Ext(p)) {
		case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".tiff", ".tif":
			out = append(out, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

var fileServiceMinioConfigCache = minioConfigCache{ttl: minioConfigCacheTTL}

type minioConfigCache struct {
	mu       sync.RWMutex
	cfg      *filesHandlers.MinioProxyConfig
	err      error
	loadedAt time.Time
	lastWarn time.Time
	ttl      time.Duration
}

func (c *minioConfigCache) get() (*filesHandlers.MinioProxyConfig, error) {
	now := time.Now()
	c.mu.RLock()
	cfg := c.cfg
	err := c.err
	loaded := c.loadedAt
	ttl := c.ttl
	c.mu.RUnlock()
	if ttl <= 0 {
		ttl = minioConfigCacheTTL
	}
	if now.Sub(loaded) < ttl {
		return cfg, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if now.Sub(c.loadedAt) < ttl {
		return c.cfg, c.err
	}
	cfg, err = c.refresh()
	c.loadedAt = time.Now()
	c.cfg = cfg
	c.err = err
	if err != nil && errors.Is(err, ErrObjectStoreUnavailable) {
		c.logUnavailable(cfg, err)
	}
	return cfg, err
}

func (c *minioConfigCache) getStrict(ctx context.Context) (*filesHandlers.MinioProxyConfig, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg, err := coreConfig.LoadMinioProxyConfig()
	if err != nil || cfg == nil {
		return nil, err
	}
	probeCtx, cancel := context.WithTimeout(ctx, strictProbeTimeout)
	defer cancel()
	return buildFilesMinioProxyConfigWithContext(probeCtx, cfg)
}

func (c *minioConfigCache) refresh() (*filesHandlers.MinioProxyConfig, error) {
	cfg, err := coreConfig.LoadMinioProxyConfig()
	if err != nil || cfg == nil {
		return nil, err
	}
	return buildFilesMinioProxyConfig(cfg)
}

func (c *minioConfigCache) logUnavailable(cfg *filesHandlers.MinioProxyConfig, err error) {
	ttl := c.ttl
	if ttl <= 0 {
		ttl = minioConfigCacheTTL
	}
	now := time.Now()
	if now.Sub(c.lastWarn) < ttl {
		return
	}
	c.lastWarn = now
	if cfg != nil {
		slog.Warn("object store unavailable", "endpoint", cfg.Endpoint, "bucket", cfg.Bucket, "err", err)
		return
	}
	slog.Warn("object store unavailable", "err", err)
}

var ErrObjectStoreUnavailable = errors.New("object store unavailable")

const (
	minioHealthTimeout  = 5 * time.Second
	minioConfigCacheTTL = 30 * time.Second
	strictProbeTimeout  = 3 * time.Second
)

func buildFilesMinioProxyConfig(cfg *config_.MinioProxyConfig) (*filesHandlers.MinioProxyConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), minioHealthTimeout)
	defer cancel()
	return buildFilesMinioProxyConfigWithContext(ctx, cfg)
}

func buildFilesMinioProxyConfigWithContext(ctx context.Context, cfg *config_.MinioProxyConfig) (*filesHandlers.MinioProxyConfig, error) {
	opts, err := buildMinioOptions(cfg)
	if err != nil {
		return nil, err
	}
	client, err := minio.New(cfg.Endpoint, &opts)
	if err != nil {
		return nil, err
	}
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrObjectStoreUnavailable, err)
	}
	if !exists {
		return nil, fmt.Errorf("%w: bucket %s not found", ErrObjectStoreUnavailable, cfg.Bucket)
	}
	prefix := strings.Trim(cfg.Prefix, "/")
	statFn := func(ctx context.Context, bucket, key string) (filesHandlers.MinioObjectInfo, error) {
		st, err := client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
		if err != nil {
			return filesHandlers.MinioObjectInfo{}, err
		}
		return filesHandlers.MinioObjectInfo{
			Size:    st.Size,
			ModTime: st.LastModified,
		}, nil
	}
	return &filesHandlers.MinioProxyConfig{
		Endpoint: cfg.Endpoint,
		Bucket:   cfg.Bucket,
		Prefix:   prefix,
		UseSSL:   cfg.Secure,
		Client:   client,
		Stat:     statFn,
		Fetch: func(ctx context.Context, bucket, key string) (io.ReadSeekCloser, filesHandlers.MinioObjectInfo, error) {
			obj, err := client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
			if err != nil {
				return nil, filesHandlers.MinioObjectInfo{}, err
			}
			st, err := obj.Stat()
			if err != nil {
				_ = obj.Close()
				return nil, filesHandlers.MinioObjectInfo{}, err
			}
			return obj, filesHandlers.MinioObjectInfo{
				Size:    st.Size,
				ModTime: st.LastModified,
			}, nil
		},
		Put: func(ctx context.Context, bucket, key string, src io.Reader, size int64, contentType string) error {
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			_, err := client.PutObject(ctx, bucket, key, src, size, minio.PutObjectOptions{ContentType: contentType})
			return err
		},
		Delete: func(ctx context.Context, bucket, key string) error {
			return client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
		},
	}, nil
}

func buildMinioOptions(cfg *config_.MinioProxyConfig) (minio.Options, error) {
	opts := minio.Options{
		Secure: cfg.Secure,
	}
	creds, err := buildMinioCredentials(cfg.Auth)
	if err != nil {
		return opts, err
	}
	opts.Creds = creds
	if cfg.CABundlePath != "" {
		pool, err := loadCABundle(cfg.CABundlePath)
		if err != nil {
			return opts, err
		}
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{RootCAs: pool}
		opts.Transport = transport
	}
	return opts, nil
}

func buildMinioCredentials(auth *config_.MinioProxyAuth) (*credentials.Credentials, error) {
	if auth == nil {
		return credentials.NewStatic("", "", "", credentials.SignatureAnonymous), nil
	}
	switch auth.Mode {
	case config_.MinioProxyAuthModeFile:
		ak, sk, err := readMinioCredentialsFile(auth.CredFile)
		if err != nil {
			return nil, err
		}
		return credentials.NewStaticV4(ak, sk, ""), nil
	case config_.MinioProxyAuthModeAccessKey, "":
		ak := strings.TrimSpace(auth.AccessKey)
		sk := strings.TrimSpace(auth.SecretKey)
		if ak == "" || sk == "" {
			return nil, fmt.Errorf("missing MinIO access key/secret")
		}
		return credentials.NewStaticV4(ak, sk, ""), nil
	case config_.MinioProxyAuthModeNone:
		return credentials.NewStatic("", "", "", credentials.SignatureAnonymous), nil
	default:
		return nil, fmt.Errorf("unsupported auth mode %q", auth.Mode)
	}
}

func readMinioCredentialsFile(path string) (string, string, error) {
	if path == "" {
		return "", "", fmt.Errorf("%w: credential file path is empty", ErrObjectStoreUnavailable)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", "", fmt.Errorf("%w: stat %s: %v", ErrObjectStoreUnavailable, path, err)
	}
	if !info.Mode().IsRegular() {
		return "", "", fmt.Errorf("%w: credential file %s is not a regular file", ErrObjectStoreUnavailable, path)
	}
	if info.Mode().Perm()&0o077 != 0 {
		return "", "", fmt.Errorf("%w: credential file %s must be private (600)", ErrObjectStoreUnavailable, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("%w: read %s: %v", ErrObjectStoreUnavailable, path, err)
	}
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return "", "", fmt.Errorf("%w: credential file %s is empty", ErrObjectStoreUnavailable, path)
	}
	var payload struct {
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
	}
	if err := json.Unmarshal(data, &payload); err == nil {
		if payload.AccessKey != "" && payload.SecretKey != "" {
			return strings.TrimSpace(payload.AccessKey), strings.TrimSpace(payload.SecretKey), nil
		}
	}
	parts := strings.Split(raw, ":")
	if len(parts) == 2 {
		ak := strings.TrimSpace(parts[0])
		sk := strings.TrimSpace(parts[1])
		if ak != "" && sk != "" {
			return ak, sk, nil
		}
	}
	return "", "", fmt.Errorf("%w: credential file %s missing keys", ErrObjectStoreUnavailable, path)
}

func loadCABundle(path string) (*x509.CertPool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("failed to parse CA bundle %s", path)
	}
	return pool, nil
}

func (h *GatewayHandlers) newServeProvider() serveProvider {
	return serveProvider{globule: h.globule}
}

func (h *GatewayHandlers) newUploadProvider() uploadProvider {
	return uploadProvider{globule: h.globule}
}
