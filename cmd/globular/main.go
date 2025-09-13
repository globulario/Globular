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
	"syscall"
	"time"

	globpkg "github.com/globulario/Globular/internal/globule"
	cfgHandlers "github.com/globulario/Globular/internal/handlers/config"
	filesHandlers "github.com/globulario/Globular/internal/handlers/files"
	mediaHandlers "github.com/globulario/Globular/internal/handlers/media"
	httplib "github.com/globulario/Globular/internal/http"
	middleware "github.com/globulario/Globular/internal/http/middleware"
	"github.com/globulario/Globular/internal/server"

	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
)

var (
	// 200 MiB cap by default (tweak with --max-upload)
	maxUpload = flag.Int64("max-upload", 200<<20, "max upload size in bytes")

	// Process-wide Globule instance (keeps adapters simple)
	globule *globpkg.Globule
)

// ============================================================
// main()
// ============================================================
func main() {

	// Keep flags around for future expansion, but ports will come from Globule config.
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

		// Bootstrap DNS+ACME and write the certs
		if err := globule.BootstrapTLSAndDNS(context.Background()); err != nil {
			logger.Error("tls/dns bootstrap failed", "err", err)
			os.Exit(1)
		}
		httpsAddr = fmt.Sprintf(":%d", globule.PortHTTPS)
		logger.Info("starting HTTPS (from globule config)", "addr", httpsAddr, "domain", globule.Domain)
	default:
		// keep basic FS for HTTP-only
		if err := globule.InitFS(); err != nil {
			logger.Error("bootstrap failed", "err", err)
			os.Exit(1)
		}
		httpAddr = fmt.Sprintf(":%d", globule.PortHTTP)
		logger.Info("starting HTTP (from globule config)", "addr", httpAddr, "domain", globule.Domain)
	}

	// router wiring unchanged...
	mux := httplib.NewRouter(logger, httplib.Config{ /* ... */ })
	wireConfig(mux)
	wireFiles(mux)
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

	// start services tied to that context
	servicesCtx, servicesCancel := context.WithCancel(ctx)
	go globule.StartServices(servicesCtx)

	// one-shot shutdown path
	go func() {
		<-ctx.Done()
		slog.Info("shutdown requested; stopping services...")

		// stop service supervisor first (cancels StartServices loop)
		servicesCancel()

		// ask Globule to stop all children (use your robust killers inside)
		globule.StopServices()

		// stop HTTP/S (best-effort; add a timeout if you want)
		_ = sup.Stop(context.Background())
	}()

	// block main until shutdown happened
	<-ctx.Done()
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
	return claims.Id + "@" + claims.UserDomain, nil
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
func (serveProvider) DataRoot() string                        { return config_.GetDataDir() }
func (serveProvider) CredsDir() string                        { return config_.GetConfigDir() + "/tls" }
func (serveProvider) IndexApplication() string                { return globule.IndexApplication }
func (serveProvider) PublicDirs() []string                    { return config_.GetPublicDirs() }
func (serveProvider) Exists(p string) bool                    { return Utility.Exists(p) }
func (serveProvider) FindHashedFile(p string) (string, error) { return findHashedFile(p) }
func (serveProvider) ParseUserID(tok string) (string, error)  { return tokenParser{}.ParseUserID(tok) }
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
func (serveProvider) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	streamHandler(name, w, r)
	return true
}

// Reverse proxy pass-through.
func (serveProvider) ResolveProxy(reqPath string) (string, bool) {
	return proxyResolver{}.ResolveProxy(reqPath)
}

type uploadProvider struct{}

func (uploadProvider) DataRoot() string                       { return config_.GetDataDir() }
func (uploadProvider) PublicDirs() []string                   { return config_.GetPublicDirs() }
func (uploadProvider) ParseUserID(tok string) (string, error) { return tokenParser{}.ParseUserID(tok) }
func (uploadProvider) ValidateAccount(u, action, p string) (bool, bool, error) {
	return accessControl{}.ValidateAccount(u, action, p)
}
func (uploadProvider) ValidateApplication(app, action, p string) (bool, bool, error) {
	return accessControl{}.ValidateApplication(app, action, p)
}

// ---------------------
// Config adapters
// ---------------------

type tokenValidator struct{}

func (tokenValidator) Validate(tok string) error {
	_, err := security.ValidateToken(tok)
	return err
}

type cfgSaver struct{}

// Save persists config through Globule (keeps compatible shape: map[string]any).
func (cfgSaver) Save(m map[string]any) error { return globule.SetConfig(m) }

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
	path, _ := cfg["ConfigPath"].(string)
	if path == "" {
		return nil, fmt.Errorf("missing ConfigPath for service %q", serviceID)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var disk map[string]any
	if err := json.Unmarshal(data, &disk); err != nil {
		return nil, err
	}
	raw := disk["Permissions"]
	switch v := raw.(type) {
	case nil:
		return []any{}, nil
	case []interface{}:
		return v, nil
	default:
		return []any{}, nil
	}
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
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
}

// ---------------------
// Config wiring
// ---------------------

type cfgProvider struct{}

func (cfgProvider) Address() (string, error) { return config_.GetAddress() }
func (cfgProvider) RemoteConfig(host string, port int) (map[string]any, error) {
	return config_.GetRemoteConfig(host, port)
}
func (cfgProvider) MyIP() string                { return Utility.MyIP() }
func (cfgProvider) LocalConfig() map[string]any { return globule.GetConfig() }
func (cfgProvider) RootDir() string             { return config_.GetRootDir() }
func (cfgProvider) DataDir() string             { return config_.GetDataDir() }
func (cfgProvider) ConfigDir() string           { return config_.GetConfigDir() }
func (cfgProvider) WebRootDir() string          { return config_.GetWebRootDir() }
func (cfgProvider) PublicDirs() []string        { return config_.GetPublicDirs() }

// ============================================================
// Media adapters (IMDb)
// ============================================================

// -- Season/Episode scraping (best-effort) --

type imdbSeasonEpisode struct{ Client *http.Client }

func (s imdbSeasonEpisode) ResolveSeasonEpisode(titleID string) (int, int, string, error) {
	if s.Client == nil {
		s.Client = http.DefaultClient
	}
	resp, err := s.Client.Get("https://www.imdb.com/title/" + titleID)
	if err != nil {
		return -1, -1, "", err
	}
	defer resp.Body.Close()

	page, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, "", err
	}

	// Example snippet: >S02<!-- -->.<!-- -->E05<
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

	// Series link near hero title block
	reSeries := regexp.MustCompile(`(?s)data-testid="hero-title-block__series-link".*?href="/title/(tt\d{7,8})/`)
	seriesID := ""
	if m := reSeries.FindSubmatch(page); len(m) == 2 {
		seriesID = string(m[1])
	}
	return season, episode, seriesID, nil
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

// -- Titles (search) using IMDb suggestion API --

type imdbTitles struct{}

func (imdbTitles) SearchIMDBTitles(q mediaHandlers.TitlesQuery) ([]map[string]any, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	query := strings.TrimSpace(q.Q)
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	first := 'a'
	for _, r := range strings.ToLower(query) {
		first = r
		break
	}
	slug := strings.ReplaceAll(query, " ", "_")
	u := fmt.Sprintf("https://v2.sg.media-imdb.com/suggestion/%c/%s.json", first, url.PathEscape(slug))

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload struct {
		D []struct {
			ID string `json:"id"` // "tt0133093"
			L  string `json:"l"`  // title
			Y  int    `json:"y"`  // year (optional)
			Q  string `json:"q"`  // type-ish string (e.g., "feature", "TV series")
		} `json:"d"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	var out []map[string]any
	for _, it := range payload.D {
		if it.ID == "" || it.L == "" {
			continue
		}
		if q.Year > 0 && it.Y != q.Year {
			continue
		}
		if qt := strings.TrimSpace(q.Type); qt != "" &&
			!strings.Contains(strings.ToLower(it.Q), strings.ToLower(qt)) {
			continue
		}

		m := map[string]any{"id": it.ID, "title": it.L}
		if it.Y != 0 {
			m["year"] = it.Y
		}
		if it.Q != "" {
			m["type"] = it.Q
		}
		if poster, err := GetIMDBPoster(it.ID); err == nil && poster != "" {
			m["poster"] = poster
		}
		out = append(out, m)
	}

	// offset/limit window
	start := q.Offset
	if start < 0 {
		start = 0
	}
	if start > len(out) {
		return []map[string]any{}, nil
	}
	end := start + q.Limit
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], nil
}

// ============================================================
// HTTP Wiring (handlers + middleware)
// ============================================================

func wireConfig(mux *http.ServeMux) {
	getConfig := cfgHandlers.NewGetConfig(cfgProvider{})
	saveConfig := cfgHandlers.NewSaveConfig(cfgSaver{}, tokenValidator{})
	getSvcPerms := cfgHandlers.NewGetServicePermissions(svcPermsProvider{})

	// Redirect (peer) + Preflight wrapper
	wrap := middleware.WithRedirectAndPreflight(redirector{}, setHeaders)

	cfgHandlers.Mount(mux, cfgHandlers.Deps{
		GetConfig:             wrap(getConfig),
		SaveConfig:            wrap(saveConfig),
		GetServicePermissions: wrap(getSvcPerms),
	})
}

func wireFiles(mux *http.ServeMux) {
	getImages := filesHandlers.NewGetImages(imgLister{})
	serve := filesHandlers.NewServeFile(serveProvider{})

	upload := filesHandlers.NewUploadFileWithOptions(
		uploadProvider{},
		filesHandlers.UploadOptions{
			MaxBytes:    *maxUpload,
			AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt"},
		},
	)

	wrap := middleware.WithRedirectAndPreflight(redirector{}, setHeaders)
	filesHandlers.Mount(mux, filesHandlers.Deps{
		GetImages: wrap(getImages),
		Serve:     wrap(serve),
		Upload:    wrap(upload),
	})
}

func wireMedia(mux *http.ServeMux) {
	titles := mediaHandlers.NewGetIMDBTitles(imdbTitles{})
	poster := mediaHandlers.NewGetIMDBPoster(imdbPoster{})

	wrap := middleware.WithRedirectAndPreflight(redirector{}, setHeaders)
	mediaHandlers.Mount(mux, mediaHandlers.Deps{
		GetIMDBTitles: wrap(titles),
		GetIMDBPoster: wrap(poster),
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

// streamHandler serves local files with Range support (placeholder for real streaming).
func streamHandler(name string, w http.ResponseWriter, r *http.Request) {
	clean := filepath.Clean(name)
	if strings.Contains(clean, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	f, err := os.Open(clean)
	if err != nil {
		http.Error(w, fmt.Sprintf("open %s: %v", clean, err), http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("stat %s: %v", clean, err), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
}

// ----- IMDb helpers (shared by media adapters) -----

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

func defaultUA() string {
	// Modern desktop UA helps IMDb return richer markup
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36"
}
