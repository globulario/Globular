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
	"path"
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
	mode    string
}

func (p serveProvider) WebRoot() string                           { return config_.GetWebRootDir() }
func (p serveProvider) DataRoot() string                          { return "" }
func (p serveProvider) CredsDir() string                          { return config_.GetConfigDir() + "/tls" }
func (p serveProvider) IndexApplication() string                  { return p.globule.IndexApplication }
func (p serveProvider) PublicDirs() []string                      { return config_.GetPublicDirs() }
func (p serveProvider) Exists(pth string) bool                    { return Utility.Exists(pth) }
func (p serveProvider) FindHashedFile(pth string) (string, error) { return findHashedFile(pth) }
func (p serveProvider) Mode() string                              { return p.mode }
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
		// Safe type assertion with check
		str, ok := v.(string)
		if !ok {
			continue
		}
		parts := strings.SplitN(strings.TrimSpace(str), "|", 2)
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
	// v1 Conformance: Return PrincipalID (security violation INV-1.1)
	// REMOVED: claims.ID + "@" + claims.UserDomain
	// Identity MUST NOT include domain - domain is routing label, not identity
	// Return opaque PrincipalID for stable, domain-independent identity
	principalID := claims.PrincipalID
	if principalID == "" {
		// Fallback for legacy tokens without PrincipalID
		principalID = claims.ID
	}
	return principalID, nil
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

func (cfgSaver) Validate(tok string) error {
	_, err := security.ValidateToken(tok)
	return err
}

func (cfgSaver) SaveServiceConfig(cfg map[string]any) error {
	return config_.SaveServiceConfiguration(cfg)
}

func (cfgProvider) AllServiceConfigs() ([]map[string]any, error) {
	cfgs, err := config_.GetServicesConfigurations()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, len(cfgs))
	for i, c := range cfgs {
		result[i] = c
	}
	return result, nil
}

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

	// cold-start strict probe controls
	strictOnce  bool
	strictUntil time.Time
	strictProbe func(context.Context) (*filesHandlers.MinioProxyConfig, error)
}

func (c *minioConfigCache) get() (*filesHandlers.MinioProxyConfig, error) {
	now := time.Now()
	c.mu.RLock()
	cfg := c.cfg
	err := c.err
	loaded := c.loadedAt
	ttl := c.ttl
	strictOnce := c.strictOnce
	strictUntil := c.strictUntil
	strictProbe := c.strictProbe
	c.mu.RUnlock()
	if ttl <= 0 {
		ttl = minioConfigCacheTTL
	}

	if strictProbe == nil {
		strictProbe = c.getStrict
	}

	if !strictOnce && (strictUntil.IsZero() || now.After(strictUntil)) {
		c.mu.Lock()
		if strictProbe == nil {
			strictProbe = c.getStrict
		}
		if !c.strictOnce && (c.strictUntil.IsZero() || now.After(c.strictUntil)) {
			strictCfg, strictErr := strictProbe(context.Background())
			if strictErr == nil && strictCfg != nil {
				c.cfg = strictCfg
				c.err = nil
				c.loadedAt = now
				c.strictOnce = true
				c.mu.Unlock()
				return strictCfg, nil
			}
			if strictCfg != nil {
				c.cfg = strictCfg
			}
			c.err = strictErr
			c.strictUntil = now.Add(strictProbeBackoff)
			c.loadedAt = now
		}
		cfg = c.cfg
		err = c.err
		c.mu.Unlock()
		return cfg, err
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
	return buildFilesMinioProxyConfigWithContext(probeCtx, cfg, true)
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
	strictProbeBackoff  = 10 * time.Second
)

func buildFilesMinioProxyConfig(cfg *config_.MinioProxyConfig) (*filesHandlers.MinioProxyConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), minioHealthTimeout)
	defer cancel()
	return buildFilesMinioProxyConfigWithContext(ctx, cfg, true)
}

func buildFilesMinioProxyConfigWithContext(ctx context.Context, cfg *config_.MinioProxyConfig, allowPartial bool) (*filesHandlers.MinioProxyConfig, error) {
	opts, err := buildMinioOptions(cfg)
	if err != nil {
		return nil, err
	}
	client, err := minio.New(cfg.Endpoint, &opts)
	if err != nil {
		return nil, err
	}
	layout := deriveMinioLayout(cfg)
	if err := validateBucketExists(ctx, client, layout.usersBucket); err != nil {
		if !allowPartial {
			return nil, err
		}
		return buildPartialMinioConfig(cfg, layout, client, err), err
	}
	if err := validateBucketExists(ctx, client, layout.webrootBucket); err != nil {
		if !allowPartial {
			return nil, err
		}
		return buildPartialMinioConfig(cfg, layout, client, err), err
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
		Endpoint:      cfg.Endpoint,
		Bucket:        cfg.Bucket,
		Prefix:        prefix,
		UsersPrefix:   layout.usersPrefix,
		WebrootPrefix: layout.webrootPrefix,
		UsersBucket:   layout.usersBucket,
		WebrootBucket: layout.webrootBucket,
		Domain:        layout.domain,
		UseSSL:        cfg.Secure,
		Client:        client,
		Stat:          statFn,
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

func buildPartialMinioConfig(cfg *config_.MinioProxyConfig, layout minioLayout, client *minio.Client, cause error) *filesHandlers.MinioProxyConfig {
	prefix := strings.Trim(cfg.Prefix, "/")
	return &filesHandlers.MinioProxyConfig{
		Endpoint:      cfg.Endpoint,
		Bucket:        cfg.Bucket,
		Prefix:        prefix,
		UsersPrefix:   layout.usersPrefix,
		WebrootPrefix: layout.webrootPrefix,
		UsersBucket:   layout.usersBucket,
		WebrootBucket: layout.webrootBucket,
		Domain:        layout.domain,
		UseSSL:        cfg.Secure,
		Client:        client,
	}
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

type minioLayout struct {
	usersPrefix   string
	webrootPrefix string
	usersBucket   string
	webrootBucket string
	domain        string
}

func deriveMinioLayout(cfg *config_.MinioProxyConfig) minioLayout {
	domain, _ := config_.GetDomain()
	domain = strings.TrimSpace(domain)
	if domain == "" {
		domain = "localhost"
	}
	layout := minioLayout{
		usersBucket:   strings.TrimSpace(cfg.Bucket),
		webrootBucket: strings.TrimSpace(cfg.Bucket),
		domain:        domain,
	}
	if ub := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_USERS_BUCKET")); ub != "" {
		layout.usersBucket = ub
	}
	if wb := strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_WEBROOT_BUCKET")); wb != "" {
		layout.webrootBucket = wb
	}

	basePrefix := strings.Trim(cfg.Prefix, "/")
	layout.usersPrefix = strings.Trim(strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_USERS_PREFIX")), "/")
	layout.webrootPrefix = strings.Trim(strings.TrimSpace(os.Getenv("GLOBULAR_MINIO_WEBROOT_PREFIX")), "/")
	if layout.usersPrefix == "" && basePrefix != "" {
		layout.usersPrefix = path.Join(basePrefix, "users")
	}
	if layout.webrootPrefix == "" && basePrefix != "" {
		layout.webrootPrefix = path.Join(basePrefix, "webroot")
	}
	if layout.usersPrefix == "" {
		// v1 Conformance: Use stable prefix (security violation INV-1.3)
		// REMOVED: path.Join(layout.domain, "users") - Domain MUST NOT determine storage paths
		// Domain is routing configuration, not identity - using it breaks on domain changes
		// For multi-tenancy, use explicit prefix config with clusterID or principalID
		layout.usersPrefix = "users" // Stable prefix, independent of domain config
	}
	if layout.webrootPrefix == "" {
		// v1 Conformance: Use stable prefix (security violation INV-1.3)
		// REMOVED: path.Join(layout.domain, "webroot")
		layout.webrootPrefix = "webroot" // Stable prefix, independent of domain config
	}
	return layout
}

func validateBucketExists(ctx context.Context, client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%w: objectstore not provisioned: %v (run node-agent plan ensure-objectstore-layout)", ErrObjectStoreUnavailable, err)
	}
	if exists {
		return nil
	}
	return fmt.Errorf("%w: objectstore not provisioned: missing bucket %s (run node-agent plan ensure-objectstore-layout)", ErrObjectStoreUnavailable, bucket)
}

func (h *GatewayHandlers) newServeProvider() serveProvider {
	return serveProvider{globule: h.globule, mode: h.cfg.Mode}
}

func (h *GatewayHandlers) newUploadProvider() uploadProvider {
	return uploadProvider{globule: h.globule}
}
