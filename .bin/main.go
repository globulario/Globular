package main

// ------------------------------------------------------------
// Globular HTTP entrypoint
//
// - Boots a refactored Globule (internal/globule)
// - Wires HTTP handlers (config, files, media)
// - Runs the HTTP server with graceful shutdown
//
// Notes:
// - TLS/ACME integration hooks are present but optional (commented).
// - Adapters keep legacy surface intact (ReverseProxies, permissions, etc.).
// - Helper funcs for hashed assets, import resolution and basic streaming
//   live at the bottom for now (can be moved to internal packages later).
// ------------------------------------------------------------

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/StalkR/imdb"
	httplib "github.com/globulario/Globular/internal/gateway/http"
	middleware "github.com/globulario/Globular/internal/gateway/http/middleware"
	globpkg "github.com/globulario/Globular/internal/globule"
	cfgHandlers "github.com/globulario/Globular/internal/handlers/config"
	filesHandlers "github.com/globulario/Globular/internal/handlers/files"
	mediaHandlers "github.com/globulario/Globular/internal/handlers/media"
	"github.com/globulario/Globular/internal/server"

	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/security"
	"github.com/globulario/services/golang/title/titlepb"
	Utility "github.com/globulario/utility"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	// 2 GiB cap by default (tweak with --max-upload)
	maxUpload = flag.Int64("max-upload", 2<<30, "max upload size in bytes")
	// Rate limit defaults (can be relaxed/disabled per deployment).
	rateRPS   = flag.Int("rate-rps", 50, "max requests/sec per client; <=0 disables throttling")
	rateBurst = flag.Int("rate-burst", 200, "max burst per client; <=0 disables throttling")

	// Process-wide Globule instance (keeps adapters simple)
	globule *globpkg.Globule
)

// ============================================================
// main()
// ============================================================
func main() {
	_ = flag.String("http", "", "ignored: HTTP port is taken from Globule config")
	_ = flag.String("https", "", "ignored: HTTPS port is taken from Globule config")
	flag.Parse()

	// Structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// 1) Build the refactored Globule
	globule = globpkg.New(logger)

	// If we’re going HTTPS, prepare DNS+ACME now.
	var httpAddr, httpsAddr string
	switch strings.ToLower(globule.Protocol) {
	case "https":
		if err := globule.BootstrapTLSAndDNS(context.Background()); err != nil {
			logger.Error("tls/dns bootstrap failed", "err", err)
			os.Exit(1)
		}
		httpsAddr = fmt.Sprintf(":%d", globule.PortHTTPS)
		logger.Info("starting HTTPS (from globule config)", "addr", httpsAddr, "domain", globule.Domain)
	default:
		if err := globule.InitFS(); err != nil {
			logger.Error("bootstrap failed", "err", err)
			os.Exit(1)
		}
		httpAddr = fmt.Sprintf(":%d", globule.PortHTTP)
		logger.Info("starting HTTP (from globule config)", "addr", httpAddr, "domain", globule.Domain)
	}

	// Build the static file handler once (DRY)
	serve := filesHandlers.NewServeFile(serveProvider{})

	// Same wrapper you already use (redirect host + CORS setHeaders + preflight)
	wrap := middleware.WithRedirectAndPreflight(redirector{}, setHeaders)

	limiterDisabled := *rateRPS <= 0 || *rateBurst <= 0
	if limiterDisabled {
		logger.Warn("http rate limiter disabled", "rateRPS", *rateRPS, "rateBurst", *rateBurst)
	}

	// 2) Router: inject the static handler as the root handler so "/" resolves to index.html
	mux := httplib.NewRouter(logger, httplib.Config{
		AllowedOrigins: []string{"*"}, // or your domains
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		RateRPS:        *rateRPS,
		RateBurst:      *rateBurst,
	}, wrap(serve))

	// Optional: keep a second entry point at /serve/*
	mux.Handle("/serve/", wrap(http.StripPrefix("/serve", serve)))

	// Mount the rest
	wireConfig(mux)
	wireFiles(mux) // images + upload endpoints (no extra "/" registrations)
	wireMedia(mux)

	// supervisor with TLS files when needed
	sup := server.Supervisor{
		Logger:            logger,
		HTTPAddr:          httpAddr,
		HTTPSAddr:         httpsAddr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if httpsAddr != "" {
		credDir := config_.GetConfigDir() + "/tls/" + globule.LocalDomain()
		sup.TLS = &server.TLSFiles{
			CertFile: filepath.Join(credDir, "fullchain.pem"),
			KeyFile:  filepath.Join(credDir, "server.key"),
		}
	}

	if err := sup.Start(mux); err != nil {
		logger.Error("start failed", "err", err)
		os.Exit(1)
	}

	// Wait until a listener is bound (https preferred in your config)
	<-sup.Ready

	// single parent shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// start services tied to that context (and wait for them on shutdown)
	servicesCtx, servicesCancel := context.WithCancel(ctx)
	var svcWG sync.WaitGroup
	svcWG.Add(1)
	go func() {
		defer svcWG.Done()
		globule.StartServices(servicesCtx)
	}()

	// --- block main until we get a signal ---
	<-ctx.Done()

	// Synchronous, ordered shutdown (so we don't race with process exit)
	logger.Info("shutdown requested; stopping services...")

	// 1) stop service supervisor first (cancels StartServices loop)
	servicesCancel()
	// 2) ask Globule to stop all children (use your robust killers inside)
	globule.StopServices()
	// 3) wait for StartServices() to unwind
	svcWG.Wait()

	// 4) stop HTTP/S with a timeout
	shCtx, shCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shCancel()
	if err := sup.Stop(shCtx); err != nil {
		logger.Error("HTTP/S shutdown error", "err", err)
	}

	// 5) stop etcd server (if any)
	if err := process.StopEtcdServer(); err != nil {
		logger.Error("etcd shutdown error", "err", err)
	}
}

// ============================================================
// Adapters / Providers
//   Keep existing HTTP layer decoupled from globule internals.
// ============================================================

// ---------------------
// Access & token parsing
// ---------------------

type tokenParser struct{}

// ParseUserID extracts "<id>@<domain>" from a JWT.
func (tokenParser) ParseUserID(tok string) (string, error) {
	claims, err := security.ValidateToken(tok)
	if err != nil {
		return "", err
	}
	return claims.ID + "@" + claims.UserDomain, nil
}

type accessControl struct{}

// ValidateAccount defers to Globule RBAC for account subjects.
func (accessControl) ValidateAccount(userID, action, reqPath string) (bool, bool, error) {
	return globule.ValidateAccess(userID, rbacpb.SubjectType_ACCOUNT, action, reqPath)
}

// ValidateApplication defers to Globule RBAC for application subjects.
func (accessControl) ValidateApplication(app, action, reqPath string) (bool, bool, error) {
	return globule.ValidateAccess(app, rbacpb.SubjectType_APPLICATION, action, reqPath)
}

// ---------------------
// Reverse proxies
// ---------------------

type proxyResolver struct{}

// ResolveProxy matches the first "<url>|<pathPrefix>" whose prefix matches reqPath.
func (proxyResolver) ResolveProxy(reqPath string) (string, bool) {
	for _, v := range globule.ReverseProxies {
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

// ---------------------
// Static file serving / uploads
// ---------------------

type serveProvider struct{}

func (serveProvider) WebRoot() string                         { return config_.GetWebRootDir() }
func (serveProvider) DataRoot() string                        { return /*config_.GetDataDir()*/ "" }
func (serveProvider) CredsDir() string                        { return config_.GetConfigDir() + "/tls" }
func (serveProvider) IndexApplication() string                { return globule.IndexApplication }
func (serveProvider) PublicDirs() []string                    { return config_.GetPublicDirs() }
func (serveProvider) Exists(p string) bool                    { return Utility.Exists(p) }
func (serveProvider) FindHashedFile(p string) (string, error) { return findHashedFile(p) }
func (serveProvider) FileServiceMinioConfig() (*filesHandlers.MinioProxyConfig, bool) {
	return fileServiceMinioConfigCache.get()
}
func (serveProvider) ParseUserID(tok string) (string, error) { return tokenParser{}.ParseUserID(tok) }
func (serveProvider) ValidateAccount(u, action, p string) (bool, bool, error) {
	return accessControl{}.ValidateAccount(u, action, p)
}
func (serveProvider) ValidateApplication(app, action, p string) (bool, bool, error) {
	return accessControl{}.ValidateApplication(app, action, p)
}

// ResolveImportPath optionally rewrites JS imports to hashed filenames.
func (serveProvider) ResolveImportPath(base, line string) (string, error) {
	return resolveImportPath(base, line)
}

// MaybeStream serves media with Range support (placeholder streaming).
// Return true only if we actually streamed, so normal ServeFile can still run.
func (serveProvider) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	return streamHandlerMaybe(name, w, r)
}

// Reverse proxy pass-through.
func (serveProvider) ResolveProxy(reqPath string) (string, bool) {
	return proxyResolver{}.ResolveProxy(reqPath)
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

type uploadProvider struct{}

func (uploadProvider) DataRoot() string                       { return /*config_.GetDataDir()*/ "" }
func (uploadProvider) PublicDirs() []string                   { return config_.GetPublicDirs() }
func (uploadProvider) ParseUserID(tok string) (string, error) { return tokenParser{}.ParseUserID(tok) }
func (uploadProvider) ValidateAccount(u, action, p string) (bool, bool, error) {
	return accessControl{}.ValidateAccount(u, action, p)
}
func (uploadProvider) ValidateApplication(app, action, p string) (bool, bool, error) {
	return accessControl{}.ValidateApplication(app, action, p)
}
func (uploadProvider) AddResourceOwner(token, path, owner, resourceType string) error {
	return globule.AddResourceOwner(token, path, owner, resourceType, rbacpb.SubjectType_ACCOUNT)
}
func (uploadProvider) FileServiceMinioConfig() (*filesHandlers.MinioProxyConfig, bool) {
	return fileServiceMinioConfigCache.get()
}

// ---------------------
// Config wiring
// ---------------------

type cfgProvider struct{}

func (cfgProvider) Address() (string, error)    { return config_.GetAddress() }
func (cfgProvider) MyIP() string                { return Utility.MyIP() }
func (cfgProvider) LocalConfig() map[string]any { return globule.GetConfig() }
func (cfgProvider) ServiceConfig(idOrName string) (map[string]any, error) {
	return config_.GetServiceConfigurationById(idOrName)
}
func (cfgProvider) RootDir() string      { return config_.GetRootDir() }
func (cfgProvider) DataDir() string      { return config_.GetDataDir() }
func (cfgProvider) ConfigDir() string    { return config_.GetConfigDir() }
func (cfgProvider) WebRootDir() string   { return config_.GetWebRootDir() }
func (cfgProvider) PublicDirs() []string { return config_.GetPublicDirs() }

type describeProvider struct{}

func (describeProvider) DescribeService(name string, timeout time.Duration) (config_.ServiceDesc, string, error) {
	return globule.DescribeService(name, timeout)
}

// ---------------------
// Service permissions
// ---------------------

type svcPermsProvider struct{}

// LoadPermissions returns the "Permissions" array from a service config file.
func (svcPermsProvider) LoadPermissions(serviceID string) ([]any, error) {
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

// ---------------------
// get-images (public-only)
// ---------------------

type imgLister struct{}

// ListImages recursively collects common image files under allowed public roots.
func (imgLister) ListImages(dir string) ([]string, error) {
	roots := config_.GetPublicDirs()
	ok := false
	cleanDir := filepath.Clean(dir)
	for _, root := range roots {
		root = filepath.Clean(root)
		if cleanDir == root ||
			strings.HasPrefix(cleanDir+string(os.PathSeparator), root+string(os.PathSeparator)) ||
			strings.HasPrefix(cleanDir, root+string(os.PathSeparator)) {
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
			return nil // skip unreadable entries
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

// ---------------------
// Redirector (host-based)
// ---------------------

type redirector struct{}

// RedirectTo maps an incoming Host header to a peer target (via Globule).
func (redirector) RedirectTo(host string) (bool, *middleware.Target) {
	ok, p := globule.RedirectTo(host)
	if !ok || p == nil {
		return false, nil
	}
	return true, &middleware.Target{
		Hostname:   p.Hostname,
		Domain:     p.Domain,
		Protocol:   p.Protocol,
		PortHTTP:   int(p.PortHttp),
		PortHTTPS:  int(p.PortHttps),
		LocalIP:    p.LocalIpAddress,
		ExternalIP: p.ExternalIpAddress,
		Raw:        p,
	}
}

// HandleRedirect proxies the current request to the chosen target.
func (redirector) HandleRedirect(to *middleware.Target, w http.ResponseWriter, r *http.Request) {
	addr := to.Domain
	scheme := "http"
	if to.Protocol == "https" {
		addr += ":" + Utility.ToString(to.PortHTTPS)
	} else {
		addr += ":" + Utility.ToString(to.PortHTTP)
	}
	addr = strings.ReplaceAll(addr, ".localhost", "")

	u, _ := url.Parse(scheme + "://" + addr)
	proxy := httputil.NewSingleHostReverseProxy(u)

	// Forwarded headers (aligns with legacy)
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	proxy.ServeHTTP(w, r)
}

// setHeaders is available if you want to force specific CORS headers here.
// Router middleware may already handle these.
func setHeaders(w http.ResponseWriter, r *http.Request) {
	// Determine allowed origin (echo the request origin if it’s allowed or "*" is configured)
	origin := r.Header.Get("Origin")
	allowedOrigin := globule.Protocol + "://" + globule.Domain
	if origin != "" {
		for _, allowed := range globule.AllowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = origin
				break
			}
		}
	}

	// Join allowed methods/headers from globule config
	allowedMethods := strings.Join(globule.AllowedMethods, ",")
	allowedHeaders := strings.Join(globule.AllowedHeaders, ",")

	h := w.Header()
	if allowedOrigin != "" {
		h.Set("Access-Control-Allow-Origin", allowedOrigin)
		// Only send Allow-Credentials when not using "*"
		if allowedOrigin != "*" {
			h.Set("Access-Control-Allow-Credentials", "true")
		}
		h.Add("Vary", "Origin")
	}
	h.Set("Access-Control-Allow-Methods", allowedMethods)
	h.Set("Access-Control-Allow-Headers", allowedHeaders)
	h.Set("Access-Control-Allow-Private-Network", "true")

	// Short-circuit preflight
	if r.Method == http.MethodOptions {
		h.Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// ---------------------
// Config wiring
// ---------------------

func wireConfig(mux *http.ServeMux) {
	getConfig := cfgHandlers.NewGetConfig(cfgProvider{})
	getServiceConfig := cfgHandlers.NewGetServiceConfig(cfgProvider{})
	saveConfig := cfgHandlers.NewSaveConfig(cfgSaver{}, tokenValidator{})
	getSvcPerms := cfgHandlers.NewGetServicePermissions(svcPermsProvider{})
	describeService := cfgHandlers.NewDescribeService(describeProvider{})

	ca := cfgHandlers.NewCAProvider()
	wrap := middleware.WithRedirectAndPreflight(redirector{}, setHeaders)

	cfgHandlers.Mount(mux, cfgHandlers.Deps{
		GetConfig:             wrap(getConfig),
		GetServiceConfig:      wrap(getServiceConfig),
		SaveConfig:            wrap(saveConfig),
		GetServicePermissions: wrap(getSvcPerms),
		DescribeService:       wrap(describeService),

		GetCACertificate:  wrap(cfgHandlers.NewGetCACertificate(ca)),
		SignCACertificate: wrap(cfgHandlers.NewSignCACertificate(ca)),
		GetSANConf:        wrap(cfgHandlers.NewGetSANConf(ca)),
	})
}

// ---------------------
// Files wiring (no "/" registrations here)
// ---------------------

type tokenValidator struct{}

func (tokenValidator) Validate(tok string) error {
	_, err := security.ValidateToken(tok)
	return err
}

type cfgSaver struct{}

func (cfgSaver) Save(m map[string]any) error { return globule.SetConfig(m) }

func wireFiles(mux *http.ServeMux) {
	getImages := filesHandlers.NewGetImages(imgLister{})

	upload := filesHandlers.NewUploadFileWithOptions(
		uploadProvider{},
		filesHandlers.UploadOptions{
			MaxBytes:    *maxUpload,
			AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".mp4", ".webm", ".mov", ".avi", ".mkv", ".mp3", ".wav", ".zip", ".rar", ".7z", ".tar", ".gz", ".csv", ".json", ".xml", ".md", ".html", ".css", ".js", ".svg", ".ttf", ".otf", ".woff", ".woff2", ".eot", ".tgz"},
		},
	)

	wrap := middleware.WithRedirectAndPreflight(redirector{}, setHeaders)

	// Mount imgs + uploads (the public /serve/* mapping is set in main)
	filesHandlers.Mount(mux, filesHandlers.Deps{
		GetImages: wrap(getImages),
		Upload:    wrap(upload),
	})
}

// ---------------------
// Media wiring (IMDb)
// ---------------------

func wireMedia(mux *http.ServeMux) {
	titles := mediaHandlers.NewGetIMDBTitles(imdbTitles{})
	poster := mediaHandlers.NewGetIMDBPoster(imdbPoster{})
	seasonEpisode := mediaHandlers.NewGetIMDBSeasonEpisode(imdbSeasonEpisode{})
	trailer := mediaHandlers.NewGetIMDBTrailer(imdbTrailer{})

	wrap := middleware.WithRedirectAndPreflight(redirector{}, setHeaders)
	mediaHandlers.Mount(mux, mediaHandlers.Deps{
		GetIMDBTitles:        wrap(titles),
		GetIMDBPoster:        wrap(poster),
		GetIMDBSeasonEpisode: wrap(seasonEpisode),
		GetIMDBTrailer:       wrap(trailer),
	})
}

// ============================================================
// Helpers (assets, imports, basic streaming, IMDb internals)
// ============================================================

// findHashedFile resolves a file "name.ext" to "name.<hash>.ext" within the same dir.
func findHashedFile(p string) (string, error) {
	dir := filepath.Dir(p)
	base := filepath.Base(p)

	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fname := e.Name()
		// minimal heuristic: "name.<something>.ext"
		if strings.HasPrefix(fname, name+".") && strings.HasSuffix(fname, ext) {
			return filepath.Join(dir, fname), nil
		}
	}
	return "", errors.New("hashed file not found for " + base)
}

// resolveImportPath optionally rewrites a JS import line to a hashed target
// (only for relative imports inside the same base directory).
func resolveImportPath(base, line string) (string, error) {
	re := regexp.MustCompile(`from\s+['"]([^'"]+)['"]`)
	m := re.FindStringSubmatch(line)
	if len(m) != 2 {
		return line, nil // no import path found
	}
	importPath := m[1]

	// only rewrite relative imports
	if strings.HasPrefix(importPath, ".") {
		target := filepath.Join(base, importPath)
		if hashed, err := findHashedFile(target); err == nil {
			hashedRel := strings.TrimPrefix(hashed, base+string(filepath.Separator))
			return strings.Replace(line, importPath, "./"+hashedRel, 1), nil
		}
	}
	return line, nil
}

// streamHandlerMaybe serves local files with Range support.
// Returns true if it actually streamed something.
func streamHandlerMaybe(name string, w http.ResponseWriter, r *http.Request) bool {
	clean := filepath.Clean(name)
	if strings.Contains(clean, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return true
	}

	f, err := os.Open(clean)
	if err != nil {
		// Let caller fall back to normal static handling if we couldn't open.
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("stat %s: %v", clean, err), http.StatusInternalServerError)
		return true
	}
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
	return true
}

// ----- IMDb helpers (shared by media adapters) -----

// -- Season/Episode scraping (best-effort) --

type imdbSeasonEpisode struct{ Client *http.Client }

func (s imdbSeasonEpisode) ResolveSeasonEpisode(titleID string) (int, int, string, error) {
	if s.Client == nil {
		s.Client = &http.Client{Timeout: 10 * time.Second}
	}
	req, _ := http.NewRequest(http.MethodGet, "https://www.imdb.com/title/"+titleID+"/", nil)
	req.Header.Set("User-Agent", defaultUA())
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := s.Client.Do(req)
	if err != nil {
		return -1, -1, "", err
	}
	defer resp.Body.Close()

	page, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, "", err
	}

	reSE := regexp.MustCompile(`>S(\d{1,2})<!-- -->\.<!-- -->E(\d{1,2})<`)
	season, episode := 0, 0
	if m := reSE.FindSubmatch(page); len(m) == 3 {
		if v, e := strconv.Atoi(string(m[1])); e == nil {
			season = v
		}
		if v, e := strconv.Atoi(string(m[2])); e == nil {
			episode = v
		}
	}

	reSeries := regexp.MustCompile(`(?s)data-testid="hero-title-block__series-link".*?href="/title/(tt\d{7,8})/`)
	seriesID := ""
	if m := reSeries.FindSubmatch(page); len(m) == 2 {
		seriesID = string(m[1])
	}
	return season, episode, seriesID, nil
}

type imdbTrailer struct{}

func (imdbTrailer) FetchIMDBTrailer(id string) (string, string, string, error) {
	return fetchIMDBTrailer(id)
}

// -- Poster fetching (URL only, no bytes to proxy) --

type imdbPoster struct{}

func (imdbPoster) FetchIMDBPoster(id, size string) ([]byte, string, string, error) {
	u, err := fetchIMDBPosterURL(id)
	if err != nil {
		return nil, "", "", err
	}
	if size != "" {
		u = rewriteIMDBImageSize(u, size) // "small"|"medium"|"large"|"orig"
	}
	return nil, "", u, nil
}

// Convenience for search adapter
func GetIMDBPoster(imdbID string) (string, error) { return fetchIMDBPosterURL(imdbID) }
func GetIMDBPosterSized(imdbID, size string) (string, error) {
	u, err := fetchIMDBPosterURL(imdbID)
	if err != nil {
		return "", err
	}
	return rewriteIMDBImageSize(u, size), nil
}

func fetchIMDBTrailer(imdbID string) (string, string, string, error) {
	page, err := fetchIMDBHTML("https://www.imdb.com/title/" + imdbID + "/")
	if err != nil {
		return "", "", "", err
	}

	if u, img := extractTrailerFromTitle(page); u != "" {
		if strings.Contains(u, "/video/") {
			videoSrc, err := fetchVideoSource(u)
			if err != nil {
				return u, img, "", nil
			}
			return u, img, videoSrc, nil
		}
		if videoURL, videoImg, err := findTrailerInGallery(u); err == nil && videoURL != "" {
			if img == "" {
				img = videoImg
			}
			videoSrc, err := fetchVideoSource(videoURL)
			if err != nil {
				return videoURL, img, "", nil
			}
			return videoURL, img, videoSrc, nil
		}
		return u, img, "", nil
	}

	return "", "", "", nil
}

func fetchIMDBHTML(url string) (string, error) {
	client := &http.Client{Timeout: 12 * time.Second}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", defaultUA())
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func extractTrailerFromTitle(page string) (string, string) {
	reOG := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["']og:video(?::secure_url)?["']\s+content=["']([^"']+)["']`)
	if m := reOG.FindStringSubmatch(page); len(m) == 2 && m[1] != "" {
		img := ""
		if m2 := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["']og:image(?::secure_url)?["']\s+content=["']([^"']+)["']`).FindStringSubmatch(page); len(m2) == 2 {
			img = m2[1]
		}
		return m[1], img
	}

	reMedia := regexp.MustCompile(`(?s)data-testid="[^"]*hero-media__slate[^"]*".*?<a[^>]+href="([^"]+)"[^>]*>.*?<img[^>]+src="([^"]+)"`)
	if m := reMedia.FindStringSubmatch(page); len(m) == 3 {
		return absoluteIMDBURL(m[1]), m[2]
	}

	cands := collectVideoCandidates(page)
	return chooseCandidate(cands)
}

func findTrailerInGallery(galleryURL string) (string, string, error) {
	page, err := fetchIMDBHTML(absoluteIMDBURL(galleryURL))
	if err != nil {
		return "", "", err
	}
	if url, img := chooseCandidate(collectVideoCandidates(page)); url != "" {
		return url, img, nil
	}
	return "", "", nil
}

type videoCandidate struct {
	url   string
	image string
	label string
}

func collectVideoCandidates(page string) []videoCandidate {
	var out []videoCandidate
	reBlock := regexp.MustCompile(`(?s)<div[^>]+class="[^"]*ipc-media__img[^"]*".*?<img[^>]+src="([^"]+)"[^>]*>.*?</div>\s*<a[^>]+href="(/video/[^"]+)"[^>]*?(?:aria-label|title)="([^"]+)"`)
	for _, m := range reBlock.FindAllStringSubmatch(page, -1) {
		out = append(out, videoCandidate{url: absoluteIMDBURL(m[2]), image: m[1], label: m[3]})
	}
	reGeneric := regexp.MustCompile(`(?s)<a[^>]+href="(/video/[^"]+)"[^>]*?(?:aria-label|title)="([^"]+)"`)
	for _, m := range reGeneric.FindAllStringSubmatch(page, -1) {
		out = append(out, videoCandidate{url: absoluteIMDBURL(m[1]), label: m[2]})
	}
	return out
}

func chooseCandidate(cands []videoCandidate) (string, string) {
	for _, c := range cands {
		if strings.Contains(strings.ToLower(c.label), "trailer") {
			return c.url, c.image
		}
	}
	if len(cands) > 0 {
		return cands[0].url, cands[0].image
	}
	return "", ""
}

func fetchVideoSource(videoPage string) (string, error) {
	page, err := fetchIMDBHTML(absoluteIMDBURL(videoPage))
	if err != nil {
		return "", err
	}
	reNext := regexp.MustCompile(`(?s)<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	if m := reNext.FindStringSubmatch(page); len(m) == 2 {
		var data map[string]any
		if json.Unmarshal([]byte(m[1]), &data) == nil {
			if url := pickVideoURL(lookupPath(data, "props", "pageProps", "videoPlaybackData", "videoLegacyEncodings")); url != "" {
				return url, nil
			}
			if url := pickVideoURL(lookupPath(data, "props", "pageProps", "videoPlaybackData", "playbackURLs")); url != "" {
				return url, nil
			}
		}
	}
	reURL := regexp.MustCompile(`"videoUrl":"([^"]+)"`)
	if m := reURL.FindStringSubmatch(page); len(m) == 2 {
		return strings.ReplaceAll(m[1], `\/`, `/`), nil
	}
	return "", fmt.Errorf("video source not found")
}

func pickVideoURL(node any) string {
	arr, ok := node.([]any)
	if !ok {
		return ""
	}
	var fallback string
	for _, it := range arr {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		url := Utility.ToString(m["url"])
		mime := strings.ToLower(Utility.ToString(m["mimeType"]))
		if strings.Contains(mime, "mp4") && url != "" {
			return strings.ReplaceAll(url, `\/`, `/`)
		}
		if fallback == "" && url != "" {
			fallback = strings.ReplaceAll(url, `\/`, `/`)
		}
	}
	return fallback
}

// -- Titles (search) using IMDb suggestion API --

type imdbTitles struct{}

type imdbHeaderTransport struct {
	base http.RoundTripper
}

func (t imdbHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	if r.Header.Get("User-Agent") == "" {
		r.Header.Set("User-Agent", defaultUA())
	}
	if r.Header.Get("Accept-Language") == "" {
		r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}
	if r.Header.Get("Accept") == "" {
		r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	}
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(r)
}

func newIMDBClient(timeout time.Duration) *http.Client {
	base := http.DefaultTransport
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		base = t.Clone()
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: imdbHeaderTransport{base: base},
	}
}

var imdbIDRE = regexp.MustCompile(`^tt\d+$`)

func (imdbTitles) SearchIMDBTitles(q mediaHandlers.TitlesQuery) ([]map[string]any, error) {
	query := strings.TrimSpace(q.Q)
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	client := newIMDBClient(10 * time.Second)
	resolver := imdbSeasonEpisode{Client: client}

	if imdbIDRE.MatchString(query) {
		it, err := imdb.NewTitle(client, query)
		if err != nil {
			return nil, err
		}
		return []map[string]any{titleProtoToMap(buildTitleProto(*it, resolver))}, nil
	}

	results, err := imdb.SearchTitle(client, query)
	if err != nil {
		return nil, err
	}

	filtered := make([]imdb.Title, 0, len(results))
	for _, it := range results {
		if it.ID == "" || it.Name == "" {
			continue
		}
		if q.Year > 0 && it.Year != q.Year {
			continue
		}
		if qt := strings.TrimSpace(q.Type); qt != "" &&
			!strings.Contains(strings.ToLower(it.Type), strings.ToLower(qt)) {
			continue
		}
		filtered = append(filtered, it)
	}
	if len(filtered) == 0 {
		filtered = results
	}

	start := q.Offset
	if start < 0 {
		start = 0
	}
	if start > len(filtered) {
		return []map[string]any{}, nil
	}
	end := start + q.Limit
	if end > len(filtered) {
		end = len(filtered)
	}
	selected := filtered[start:end]

	out := make([]map[string]any, 0, len(selected))
	for _, it := range selected {
		out = append(out, titleProtoToMap(buildTitleProto(it, resolver)))
	}
	return out, nil
}

func titleProtoToMap(t *titlepb.Title) map[string]any {
	if t == nil {
		return nil
	}
	out := map[string]any{
		"ID":          t.GetID(),
		"URL":         t.GetURL(),
		"Name":        t.GetName(),
		"Type":        t.GetType(),
		"Year":        t.GetYear(),
		"Rating":      t.GetRating(),
		"RatingCount": t.GetRatingCount(),
		"Description": t.GetDescription(),
		"Duration":    t.GetDuration(),
	}
	if v := t.GetSerie(); v != "" {
		out["Serie"] = v
	}
	if v := t.GetSeason(); v > 0 {
		out["Season"] = v
	}
	if v := t.GetEpisode(); v > 0 {
		out["Episode"] = v
	}
	if v := t.GetUUID(); v != "" {
		out["UUID"] = v
	}
	if langs := t.GetLanguage(); len(langs) > 0 {
		out["Languages"] = append([]string(nil), langs...)
	}
	if genres := t.GetGenres(); len(genres) > 0 {
		out["Genres"] = append([]string(nil), genres...)
	}
	if nat := t.GetNationalities(); len(nat) > 0 {
		out["Nationalities"] = append([]string(nil), nat...)
	}
	if aka := t.GetAKA(); len(aka) > 0 {
		out["AKA"] = append([]string(nil), aka...)
	}
	if dirs := t.GetDirectors(); len(dirs) > 0 {
		out["Directors"] = personsToMaps(dirs)
	}
	if writers := t.GetWriters(); len(writers) > 0 {
		out["Writers"] = personsToMaps(writers)
	}
	if actors := t.GetActors(); len(actors) > 0 {
		out["Actors"] = personsToMaps(actors)
	}
	if poster := t.GetPoster(); poster != nil {
		pm := map[string]any{}
		if v := poster.GetID(); v != "" {
			pm["ID"] = v
		}
		if v := poster.GetTitleId(); v != "" {
			pm["TitleID"] = v
		}
		if v := poster.GetURL(); v != "" {
			pm["URL"] = v
		}
		if v := poster.GetContentUrl(); v != "" {
			pm["ContentURL"] = v
		}
		if len(pm) > 0 {
			out["Poster"] = pm
		}
	}
	return out
}

func buildTitleProto(it imdb.Title, resolver imdbSeasonEpisode) *titlepb.Title {
	var rating float32
	if it.Rating != "" {
		if v, err := strconv.ParseFloat(it.Rating, 32); err == nil {
			rating = float32(v)
		}
	}

	title := &titlepb.Title{
		ID:            it.ID,
		URL:           it.URL,
		Name:          it.Name,
		Type:          it.Type,
		Year:          int32(it.Year),
		Rating:        rating,
		RatingCount:   int32(it.RatingCount),
		Description:   it.Description,
		Genres:        append([]string(nil), it.Genres...),
		Language:      append([]string(nil), it.Languages...),
		Nationalities: append([]string(nil), it.Nationalities...),
		Duration:      it.Duration,
		Directors:     namesToPersons(it.Directors),
		Writers:       namesToPersons(it.Writers),
		Actors:        namesToPersons(it.Actors),
	}

	if strings.EqualFold(it.Type, "TVEpisode") {
		if season, episode, serie, err := resolver.ResolveSeasonEpisode(it.ID); err == nil {
			title.Season = int32(season)
			title.Episode = int32(episode)
			title.Serie = serie
		}
	}

	if poster, err := GetIMDBPoster(it.ID); err == nil && poster != "" {
		title.Poster = &titlepb.Poster{URL: poster}
	}

	return title
}

func namesToPersons(names []imdb.Name) []*titlepb.Person {
	out := make([]*titlepb.Person, 0, len(names))
	for _, n := range names {
		out = append(out, &titlepb.Person{
			ID:       n.ID,
			URL:      n.URL,
			FullName: n.FullName,
		})
	}
	return out
}

func personsToMaps(list []*titlepb.Person) []map[string]any {
	out := make([]map[string]any, 0, len(list))
	for _, p := range list {
		if p == nil {
			continue
		}
		m := map[string]any{}
		if v := p.GetID(); v != "" {
			m["ID"] = v
		}
		if v := p.GetURL(); v != "" {
			m["URL"] = v
		}
		if v := p.GetFullName(); v != "" {
			m["FullName"] = v
		}
		if aliases := p.GetAliases(); len(aliases) > 0 {
			m["Aliases"] = append([]string(nil), aliases...)
		}
		if len(m) > 0 {
			out = append(out, m)
		}
	}
	return out
}

// ----- IMDb helpers -----

func fetchIMDBPosterURL(imdbID string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest(http.MethodGet, "https://www.imdb.com/title/"+imdbID+"/", nil)
	req.Header.Set("User-Agent", defaultUA())
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	page := string(b)

	// 1) og:image
	reOG := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["']og:image["']\s+content=["']([^"']+)["']`)
	if m := reOG.FindStringSubmatch(page); len(m) == 2 && m[1] != "" {
		return m[1], nil
	}

	// 2) JSON-LD image inside any <script type="application/ld+json">
	reLD := regexp.MustCompile(`(?s)<script[^>]+type=["']application/ld\+json["'][^>]*>(.*?)</script>`)
	for _, m := range reLD.FindAllStringSubmatch(page, -1) {
		var v any
		if json.Unmarshal([]byte(m[1]), &v) == nil {
			if u := extractImageFromLDJSON(v); u != "" {
				return u, nil
			}
		}
	}

	// 3) Fallback through mediaviewer page
	rePosterLink := regexp.MustCompile(`(?s)data-testid="hero-media__poster".*?href="([^"]+mediaviewer[^"]+)"`)
	viewerPath := ""
	if m := rePosterLink.FindStringSubmatch(page); len(m) == 2 {
		viewerPath = m[1]
	}
	if viewerPath == "" {
		reViewer := regexp.MustCompile(`href="(/title/tt\d+/mediaviewer/[^"]+)"`)
		if m := reViewer.FindStringSubmatch(page); len(m) == 2 {
			viewerPath = m[1]
		}
	}
	if viewerPath == "" {
		return "", fmt.Errorf("poster not found")
	}

	vReq, _ := http.NewRequest(http.MethodGet, "https://www.imdb.com"+viewerPath, nil)
	vReq.Header.Set("User-Agent", defaultUA())
	vReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
	vResp, err := client.Do(vReq)
	if err != nil {
		return "", err
	}
	defer vResp.Body.Close()

	vb, err := io.ReadAll(vResp.Body)
	if err != nil {
		return "", err
	}
	view := string(vb)

	// Prefer largest srcset entry
	reSrcset := regexp.MustCompile(`(?s)<img[^>]+srcset="([^"]+)"[^>]*>`)
	if m := reSrcset.FindStringSubmatch(view); len(m) == 2 {
		if u := pickMaxWidthFromSrcset(m[1]); u != "" {
			return u, nil
		}
	}
	// Fallback: plain src
	reSrc := regexp.MustCompile(`(?s)<img[^>]+src="([^"]+)"[^>]*>`)
	if m := reSrc.FindStringSubmatch(view); len(m) == 2 && m[1] != "" {
		return m[1], nil
	}

	return "", fmt.Errorf("poster not found")
}

func pickMaxWidthFromSrcset(srcset string) string {
	var maxURL string
	maxW := -1
	for _, part := range strings.Split(srcset, ",") {
		part = strings.TrimSpace(part)
		items := strings.Fields(part)
		if len(items) != 2 {
			continue
		}
		u := items[0]
		w := strings.TrimSuffix(items[1], "w")
		if n, err := strconv.Atoi(w); err == nil && n > maxW {
			maxW, maxURL = n, u
		}
	}
	return maxURL
}

// rewriteIMDBImageSize attempts to upgrade an IMDb image URL size via _V1_ tag.
func rewriteIMDBImageSize(u, size string) string {
	size = strings.ToLower(size)
	if size == "" || size == "orig" {
		return u
	}
	target := map[string]string{
		"small":  "UX400",
		"medium": "UX800",
		"large":  "UX1200",
	}[size]
	if target == "" {
		return u
	}
	re := regexp.MustCompile(`\._V1_[^_.]*(_[A-Z]\w+)*(\.[a-zA-Z0-9]+)$`)
	return re.ReplaceAllString(u, "._V1_"+target+"$2")
}

// extractImageFromLDJSON looks for "image" inside JSON-LD (string|object|array).
func extractImageFromLDJSON(v any) string {
	switch x := v.(type) {
	case map[string]any:
		if img := x["image"]; img != nil {
			if u := imageURLFromAny(img); u != "" {
				return u
			}
		}
		// also try nested graph
		if g, ok := x["@graph"]; ok {
			if u := extractImageFromLDJSON(g); u != "" {
				return u
			}
		}
	case []any:
		for _, it := range x {
			if u := extractImageFromLDJSON(it); u != "" {
				return u
			}
		}
	}
	return ""
}

func imageURLFromAny(a any) string {
	switch y := a.(type) {
	case string:
		return y
	case map[string]any:
		if u, _ := y["url"].(string); u != "" {
			return u
		}
	case []any:
		for _, it := range y {
			if u := imageURLFromAny(it); u != "" {
				return u
			}
		}
	}
	return ""
}

func absoluteIMDBURL(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if strings.HasPrefix(path, "//") {
		return "https:" + path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "https://www.imdb.com" + path
}

func lookupPath(v any, keys ...string) any {
	cur := v
	for _, k := range keys {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[k]
	}
	return cur
}

func defaultUA() string {
	// Modern desktop UA helps IMDb return richer markup
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36"
}
