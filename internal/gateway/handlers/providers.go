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
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	coreConfig "github.com/globulario/Globular/internal/config"
	adminHandlers "github.com/globulario/Globular/internal/gateway/handlers/admin"
	cfgHandlers "github.com/globulario/Globular/internal/gateway/handlers/config"
	filesHandlers "github.com/globulario/Globular/internal/gateway/handlers/files"
	globpkg "github.com/globulario/Globular/internal/globule"
	"github.com/globulario/Globular/internal/journal"
	clustercontroller_client "github.com/globulario/services/golang/cluster_controller/cluster_controller_client"
	cluster_controllerpb "github.com/globulario/services/golang/cluster_controller/cluster_controllerpb"
	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/domain"
	"github.com/globulario/services/golang/globular_client"
	nodeagent_client "github.com/globulario/services/golang/node_agent/node_agent_client"
	planpb "github.com/globulario/services/golang/plan/planpb"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/repository/repository_client"
	repopb "github.com/globulario/services/golang/repository/repositorypb"
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

func (s cfgSaver) SetCorsPolicy(p *cfgHandlers.CorsPolicy) error {
	return s.globule.SetCorsPolicy(p)
}

func (cfgSaver) SaveServiceConfig(cfg map[string]any) error {
	return config_.SaveServiceConfiguration(cfg)
}

func (cfgSaver) GetServiceConfig(idOrName string) (map[string]any, error) {
	// Fast path: direct etcd key lookup (no full scan fallback).
	// The save handler always has the exact Id from the frontend,
	// so the expensive Name-based scan is unnecessary.
	if cfg, err := config_.GetServiceConfigurationByExactId(idOrName); err == nil {
		return cfg, nil
	}
	// Fallback for Name-based lookups (rare).
	return config_.GetServiceConfigurationById(idOrName)
}

func (p cfgProvider) GetCorsPolicy() *cfgHandlers.CorsPolicy {
	return p.globule.GetCorsPolicy()
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

// adminProvider satisfies admin.AdminProvider for the /admin/metrics/* surface.
type adminProvider struct {
	globule *globpkg.Globule
}

func (adminProvider) AllServiceConfigs() ([]map[string]any, error) {
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

func (adminProvider) PublicDirs() []string { return config_.GetPublicDirs() }
func (adminProvider) DataDir() string      { return config_.GetDataDir() }
func (adminProvider) StateDir() string     { return config_.GetStateRootDir() }
func (p adminProvider) Hostname() string   { return p.globule.Name }
func (adminProvider) IP() string           { return Utility.MyIP() }

// certProvider satisfies admin.CertProvider for the /admin/certificates surface.
type certProvider struct {
	globule *globpkg.Globule
}

func (certProvider) AllServiceConfigs() ([]map[string]any, error) {
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

func (certProvider) PublicDirs() []string { return config_.GetPublicDirs() }
func (certProvider) DataDir() string      { return config_.GetDataDir() }
func (certProvider) StateDir() string     { return config_.GetStateRootDir() }
func (p certProvider) Hostname() string   { return p.globule.Name }
func (certProvider) IP() string           { return Utility.MyIP() }

func (p certProvider) Protocol() string { return p.globule.Protocol }
func (p certProvider) Domain() string   { return p.globule.Domain }
func (p certProvider) AlternateDomains() []string {
	var out []string
	for _, v := range p.globule.AlternateDomains {
		if s, ok := v.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}
func (certProvider) CertPaths() *coreConfig.CertPaths {
	return coreConfig.NewCertPaths(config_.GetStateRootDir())
}
func (certProvider) RuntimeConfigDir() string {
	return config_.GetRuntimeConfigDir()
}

// journalAdapter bridges internal/journal → admin.JournalReader.
type journalAdapter struct{}

func (journalAdapter) ReadUnit(ctx context.Context, unit string, lines int, sinceSec int) adminHandlers.JournalResult {
	r := journal.ReadUnit(ctx, unit, lines, sinceSec)
	return adminHandlers.JournalResult{
		Unit:      r.Unit,
		Lines:     r.Lines,
		Truncated: r.Truncated,
		Error:     r.Error,
	}
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
	// Loopback endpoints always skip TLS verification — traffic is local-only
	// and after a backup restore the CA may not match MinIO's current cert.
	host, _, _ := net.SplitHostPort(cfg.Endpoint)
	if host == "127.0.0.1" || host == "::1" || host == "localhost" {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // loopback only
		opts.Transport = transport
	} else if cfg.CABundlePath != "" {
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

// upgradesProvider satisfies admin.UpgradesProvider by querying the repository service.
type upgradesProvider struct {
	mu       sync.Mutex
	bundles  []adminHandlers.BundleInfo
	loadedAt time.Time
}

const upgradesCacheTTL = 60 * time.Second

func (p *upgradesProvider) AvailableBundles(ctx context.Context) ([]adminHandlers.BundleInfo, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Return cached result if fresh.
	if time.Since(p.loadedAt) < upgradesCacheTTL && p.bundles != nil {
		return p.bundles, nil
	}

	// Discover the repository service address.
	cfgs, err := config_.GetServicesConfigurationsByName("repository.PackageRepository")
	if err != nil || len(cfgs) == 0 {
		return nil, fmt.Errorf("repository service not found")
	}

	cfg := cfgs[0]
	address := ""
	port := 0
	if p, ok := cfg["Port"].(float64); ok {
		port = int(p)
	}
	if port > 0 {
		// Use localhost for local service — the cluster domain (e.g.
		// "globular.internal") may not resolve via system DNS.
		address = fmt.Sprintf("localhost:%d", port)
	}
	if address == "" {
		return nil, fmt.Errorf("repository service address not available")
	}

	// Create repository client via globular_client.GetClient.
	Utility.RegisterFunction("NewRepositoryService_Client", repository_client.NewRepositoryService_Client)
	clientIface, err := globular_client.GetClient(address, "repository.PackageRepository", "NewRepositoryService_Client")
	if err != nil {
		return nil, fmt.Errorf("connect to repository: %w", err)
	}
	repoClient, ok := clientIface.(*repository_client.Repository_Service_Client)
	if !ok {
		return nil, fmt.Errorf("unexpected client type for repository")
	}

	// Query legacy bundles.
	// TODO(migration): Remove once all packages are published via UploadArtifact.
	// The artifact query below already picks up dual-written packages.
	summaries, err := repoClient.ListBundles()
	if err != nil {
		return nil, fmt.Errorf("list bundles: %w", err)
	}

	var bundles []adminHandlers.BundleInfo
	for _, s := range summaries {
		bundles = append(bundles, adminHandlers.BundleInfo{
			Name:        s.GetName(),
			Version:     s.GetVersion(),
			Platform:    s.GetPlatform(),
			PublisherID: s.GetPublisherId(),
			SizeBytes:   s.GetSizeBytes(),
			SHA256:      s.GetSha256(),
			Kind:        "SERVICE",
		})
	}

	// Supplement with modern artifacts (picks up APPLICATION and INFRASTRUCTURE packages).
	artifacts, artErr := repoClient.ListArtifacts()
	if artErr == nil {
		for _, m := range artifacts {
			ref := m.GetRef()
			if ref == nil {
				continue
			}
			kind := ref.GetKind().String()
			bundles = append(bundles, adminHandlers.BundleInfo{
				Name:        ref.GetName(),
				Version:     ref.GetVersion(),
				BuildNumber: m.GetBuildNumber(),
				Platform:    ref.GetPlatform(),
				PublisherID: ref.GetPublisherId(),
				SizeBytes:   m.GetSizeBytes(),
				SHA256:      m.GetChecksum(),
				Kind:        kind,
			})
		}
	}

	p.bundles = bundles
	p.loadedAt = time.Now()
	return bundles, nil
}

// nodeAgentProvider satisfies admin.NodeAgentProvider by calling the local node-agent gRPC service.
type nodeAgentProvider struct {
	mu     sync.Mutex
	client *nodeagent_client.NodeAgentClient
}

func (p *nodeAgentProvider) getClient() (*nodeagent_client.NodeAgentClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		return p.client, nil
	}

	// Discover node-agent address from config.
	cfgs, err := config_.GetServicesConfigurationsByName("node_agent.NodeAgentService")
	if err != nil || len(cfgs) == 0 {
		return nil, fmt.Errorf("node-agent service not found in config")
	}
	cfg := cfgs[0]
	domain := ""
	port := 0
	if d, ok := cfg["Domain"].(string); ok {
		domain = d
	}
	if pp, ok := cfg["Port"].(float64); ok {
		port = int(pp)
	}
	if domain == "" || port == 0 {
		return nil, fmt.Errorf("node-agent address not available (domain=%q port=%d)", domain, port)
	}
	address := fmt.Sprintf("%s:%d", domain, port)

	// Connect using globular_client for proper TLS handling.
	Utility.RegisterFunction("NewNodeAgentClient", nodeagent_client.NewNodeAgentClient)
	clientIface, err := globular_client.GetClient(address, "node_agent.NodeAgentService", "NewNodeAgentClient")
	if err != nil {
		return nil, fmt.Errorf("connect to node-agent: %w", err)
	}
	naClient, ok := clientIface.(*nodeagent_client.NodeAgentClient)
	if !ok {
		return nil, fmt.Errorf("unexpected client type for node-agent")
	}
	p.client = naClient
	return p.client, nil
}

func (p *nodeAgentProvider) ApplyPlanV1(ctx context.Context, plan *planpb.NodePlan) (string, error) {
	c, err := p.getClient()
	if err != nil {
		return "", err
	}
	resp, err := c.ApplyPlanV1(ctx, plan)
	if err != nil {
		return "", err
	}
	return resp.GetOperationId(), nil
}

func (p *nodeAgentProvider) GetPlanStatus(ctx context.Context, operationID string) (*planpb.NodePlanStatus, error) {
	c, err := p.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := c.GetPlanStatusV1(ctx, operationID)
	if err != nil {
		return nil, err
	}
	return resp.GetStatus(), nil
}

// controllerProvider satisfies admin.ControllerProvider.
// It lazily connects to the cluster_controller's PlanServiceUpgrades/ApplyServiceUpgrades RPCs.
type controllerProvider struct {
	mu     sync.Mutex
	client *clustercontroller_client.ClusterControllerClient
}

func (p *controllerProvider) getClient() (*clustercontroller_client.ClusterControllerClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client != nil {
		return p.client, nil
	}

	cfgs, err := config_.GetServicesConfigurationsByName("cluster_controller.ClusterControllerService")
	if err != nil || len(cfgs) == 0 {
		return nil, fmt.Errorf("cluster_controller service not found in config")
	}
	cfg := cfgs[0]
	domain := ""
	port := 0
	if d, ok := cfg["Domain"].(string); ok {
		domain = d
	}
	if pp, ok := cfg["Port"].(float64); ok {
		port = int(pp)
	}
	if domain == "" || port == 0 {
		return nil, fmt.Errorf("cluster_controller address not available (domain=%q port=%d)", domain, port)
	}
	address := fmt.Sprintf("%s:%d", domain, port)

	Utility.RegisterFunction("NewClusterControllerClient", clustercontroller_client.NewClusterControllerClient)
	clientIface, err := globular_client.GetClient(address, "cluster_controller.ClusterControllerService", "NewClusterControllerClient")
	if err != nil {
		return nil, fmt.Errorf("connect to cluster_controller: %w", err)
	}
	ccClient, ok := clientIface.(*clustercontroller_client.ClusterControllerClient)
	if !ok {
		return nil, fmt.Errorf("unexpected client type for cluster_controller")
	}
	p.client = ccClient
	return p.client, nil
}

func (p *controllerProvider) PlanServiceUpgrades(ctx context.Context, services []string) (*adminHandlers.ControllerPlanResult, error) {
	c, err := p.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := c.PlanServiceUpgrades(ctx, &cluster_controllerpb.PlanServiceUpgradesRequest{
		Services: services,
	})
	if err != nil {
		return nil, err
	}

	result := &adminHandlers.ControllerPlanResult{
		RepositoryStatus: resp.GetRepositoryStatus(),
	}
	for _, item := range resp.GetItems() {
		result.Items = append(result.Items, adminHandlers.ControllerUpgradePlanItem{
			Service:         item.GetService(),
			FromVersion:     item.GetFromVersion(),
			FromBuildNumber: item.GetFromBuildNumber(),
			ToVersion:       item.GetToVersion(),
			ToBuildNumber:   item.GetToBuildNumber(),
			PackageName:     item.GetPackageName(),
			SHA256:          item.GetSha256(),
			RestartRequired: item.GetRestartRequired(),
			Impacts:         item.GetImpacts(),
		})
	}
	return result, nil
}

func (p *controllerProvider) ApplyServiceUpgrades(ctx context.Context, services []string) (*adminHandlers.ControllerApplyResult, error) {
	c, err := p.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := c.ApplyServiceUpgrades(ctx, &cluster_controllerpb.ApplyServiceUpgradesRequest{
		Services: services,
	})
	if err != nil {
		return nil, err
	}
	return &adminHandlers.ControllerApplyResult{
		OK:          resp.GetOk(),
		OperationID: resp.GetOperationId(),
		Message:     resp.GetMessage(),
	}, nil
}

// domainStoreProvider satisfies domains.StoreProvider.
// It lazily creates the domain store from the etcd client.
type domainStoreProvider struct{}

func (domainStoreProvider) DomainStore() domain.DomainStore {
	cli, err := config_.GetEtcdClient()
	if err != nil || cli == nil {
		return nil
	}
	return domain.NewEtcdDomainStore(cli)
}

// repositoryProvider satisfies admin.RepositoryProvider by connecting to the local repository service.
type repositoryProvider struct{}

func (repositoryProvider) repoClient() (*repository_client.Repository_Service_Client, error) {
	cfgs, err := config_.GetServicesConfigurationsByName("repository.PackageRepository")
	if err != nil || len(cfgs) == 0 {
		return nil, fmt.Errorf("repository service not found")
	}
	cfg := cfgs[0]
	var port int
	if p, ok := cfg["Port"].(float64); ok {
		port = int(p)
	}
	if port == 0 {
		return nil, fmt.Errorf("repository service address not available")
	}
	address := fmt.Sprintf("localhost:%d", port)

	Utility.RegisterFunction("NewRepositoryService_Client", repository_client.NewRepositoryService_Client)
	clientIface, err := globular_client.GetClient(address, "repository.PackageRepository", "NewRepositoryService_Client")
	if err != nil {
		return nil, fmt.Errorf("connect to repository: %w", err)
	}
	client, ok := clientIface.(*repository_client.Repository_Service_Client)
	if !ok {
		return nil, fmt.Errorf("unexpected client type for repository")
	}
	return client, nil
}

func (p repositoryProvider) SearchArtifacts(query, kind, publisher, platform, pageToken string, pageSize int32) ([]*repopb.ArtifactManifest, string, int32, error) {
	client, err := p.repoClient()
	if err != nil {
		return nil, "", 0, err
	}
	defer client.Close()

	req := &repopb.SearchArtifactsRequest{
		Query:       query,
		PublisherId: publisher,
		Platform:    platform,
		PageToken:   pageToken,
		PageSize:    pageSize,
	}
	if kind != "" {
		if k, ok := repopb.ArtifactKind_value[kind]; ok {
			req.Kind = repopb.ArtifactKind(k)
		}
	}

	resp, err := client.SearchArtifacts(req)
	if err != nil {
		return nil, "", 0, err
	}
	return resp.GetArtifacts(), resp.GetNextPageToken(), resp.GetTotalCount(), nil
}

func (p repositoryProvider) GetArtifactManifest(ref *repopb.ArtifactRef, buildNumber int64) (*repopb.ArtifactManifest, error) {
	client, err := p.repoClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.GetArtifactManifest(ref, buildNumber)
}

func (p repositoryProvider) GetArtifactVersions(publisherID, name, platform string) ([]*repopb.ArtifactManifest, error) {
	client, err := p.repoClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.GetArtifactVersions(publisherID, name, platform)
}

func (p repositoryProvider) DeleteArtifact(ref *repopb.ArtifactRef) error {
	client, err := p.repoClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return client.DeleteArtifact(ref)
}
