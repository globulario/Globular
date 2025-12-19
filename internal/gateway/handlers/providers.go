package handlers

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	globpkg "github.com/globulario/Globular/internal/globule"
	filesHandlers "github.com/globulario/Globular/internal/handlers/files"
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
func (p serveProvider) FileServiceMinioConfig() (*filesHandlers.MinioProxyConfig, bool) {
	return fileServiceMinioConfigCache.get()
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
func (u uploadProvider) AddResourceOwner(token, path, owner, resourceType string) error {
	return u.globule.AddResourceOwner(token, path, owner, resourceType, rbacpb.SubjectType_ACCOUNT)
}
func (uploadProvider) FileServiceMinioConfig() (*filesHandlers.MinioProxyConfig, bool) {
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
}

func (cfgProvider) Address() (string, error)      { return config_.GetAddress() }
func (cfgProvider) MyIP() string                  { return Utility.MyIP() }
func (p cfgProvider) LocalConfig() map[string]any { return p.globule.GetConfig() }
func (cfgProvider) ServiceConfig(idOrName string) (map[string]any, error) {
	return config_.GetServiceConfigurationById(idOrName)
}
func (cfgProvider) RootDir() string      { return config_.GetRootDir() }
func (cfgProvider) DataDir() string      { return config_.GetDataDir() }
func (cfgProvider) ConfigDir() string    { return config_.GetConfigDir() }
func (cfgProvider) WebRootDir() string   { return config_.GetWebRootDir() }
func (cfgProvider) PublicDirs() []string { return config_.GetPublicDirs() }

type describeProvider struct {
	globule *globpkg.Globule
}

func (p describeProvider) DescribeService(name string, timeout time.Duration) (config_.ServiceDesc, string, error) {
	return p.globule.DescribeService(name, timeout)
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

var fileServiceMinioConfigCache = minioConfigCache{ttl: 30 * time.Second}

type minioConfigCache struct {
	mu       sync.RWMutex
	cfg      *filesHandlers.MinioProxyConfig
	loadedAt time.Time
	ttl      time.Duration
}

func (c *minioConfigCache) get() (*filesHandlers.MinioProxyConfig, bool) {
	now := time.Now()
	c.mu.RLock()
	cfg := c.cfg
	loaded := c.loadedAt
	ttl := c.ttl
	c.mu.RUnlock()
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	if cfg != nil && now.Sub(loaded) < ttl {
		return cfg, true
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg != nil && now.Sub(c.loadedAt) < ttl {
		return c.cfg, true
	}
	cfg, err := loadFileServiceMinioConfig()
	c.loadedAt = time.Now()
	if err != nil || cfg == nil {
		c.cfg = nil
		return nil, false
	}
	c.cfg = cfg
	return cfg, true
}

func loadFileServiceMinioConfig() (*filesHandlers.MinioProxyConfig, error) {
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
		return nil, fmt.Errorf("file service missing MinIO endpoint or bucket")
	}
	accessKey := strings.TrimSpace(Utility.ToString(cfg["MinioAccessKey"]))
	secretKey := strings.TrimSpace(Utility.ToString(cfg["MinioSecretKey"]))
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("file service missing MinIO credentials")
	}
	prefix := strings.Trim(Utility.ToString(cfg["MinioPrefix"]), "/")
	useSSL := Utility.ToBool(cfg["MinioUseSSL"])
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &filesHandlers.MinioProxyConfig{
		Endpoint: endpoint,
		Bucket:   bucket,
		Prefix:   prefix,
		UseSSL:   useSSL,
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

func (h *GatewayHandlers) newServeProvider() serveProvider {
	return serveProvider{globule: h.globule}
}

func (h *GatewayHandlers) newUploadProvider() uploadProvider {
	return uploadProvider{globule: h.globule}
}
