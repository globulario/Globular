package files

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"log/slog"

	httplib "github.com/globulario/Globular/internal/gateway/http"
)

// MinioProxyConfig captures the subset of FileService MinIO settings required
// to proxy /users/ requests directly to the MinIO gateway.
type MinioProxyConfig struct {
	Endpoint      string
	Bucket        string
	Prefix        string
	UsersPrefix   string
	WebrootPrefix string
	Domain        string
	UsersBucket   string
	WebrootBucket string
	UseSSL        bool
	Client        *minio.Client
	Stat          MinioStatFunc
	Fetch         MinioFetchFunc
	Put           MinioPutFunc
	Delete        MinioDeleteFunc
}

// MinioStatFunc checks whether a MinIO object exists.
type MinioStatFunc func(ctx context.Context, bucket, key string) (MinioObjectInfo, error)

// MinioObjectInfo describes an object served from MinIO.
type MinioObjectInfo struct {
	Size    int64
	ModTime time.Time
}

// MinioFetchFunc fetches an object reader + metadata for a given bucket/key.
type MinioFetchFunc func(ctx context.Context, bucket, key string) (io.ReadSeekCloser, MinioObjectInfo, error)

// MinioPutFunc uploads a new object into MinIO.
type MinioPutFunc func(ctx context.Context, bucket, key string, src io.Reader, size int64, contentType string) error

// MinioDeleteFunc removes an object from MinIO.
type MinioDeleteFunc func(ctx context.Context, bucket, key string) error

// ServeProvider abstracts all platform-specific bits needed to serve files.
type ServeProvider interface {
	// Roots & config
	WebRoot() string
	DataRoot() string
	CredsDir() string
	IndexApplication() string
	PublicDirs() []string
	Exists(p string) bool
	FindHashedFile(p string) (string, error)
	FileServiceMinioConfig() (*MinioProxyConfig, error)
	FileServiceMinioConfigStrict(ctx context.Context) (*MinioProxyConfig, error)
	Mode() string

	// Security
	ParseUserID(token string) (string, error) // returns "id@domain" or ""
	ValidateAccount(userID, action, reqPath string) (has, denied bool, err error)
	ValidateApplication(app, action, reqPath string) (has, denied bool, err error)

	// Optional hooks
	ResolveImportPath(basePath string, importLine string) (string, error)
	MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool

	// Reverse proxy
	ResolveProxy(reqPath string) (targetURL string, ok bool)
}

// hlsExtensions lists HLS (HTTP Live Streaming) file extensions that should be
// allowed for streaming access even in protected directories.
var hlsExtensions = []string{
	".ts",        // MPEG-TS segment files
	"240p.m3u8",  // Low quality playlist
	"360p.m3u8",  // SD quality playlist
	"480p.m3u8",  // SD quality playlist
	"720p.m3u8",  // HD quality playlist
	"1080p.m3u8", // Full HD quality playlist
	"2160p.m3u8", // 4K quality playlist
}

// isHLSFile returns true if the path is an HLS streaming file that should be
// accessible for playback. HLS files include .ts segments and resolution-specific
// m3u8 playlists.
func isHLSFile(path string) bool {
	for _, ext := range hlsExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

// NewServeFile implements GET /serve/* with:
// - Reverse-proxy passthrough for configured prefixes
// - Host-based subroot under WebRoot()
// - Protected bases: /users/, /applications/, /templates/, /projects/ under DataRoot()/files
// - RBAC checks (account/application) for "read"
// - Windows drive-path quirk
// - HLS/mkv streaming via MaybeStream
// - JS import rewrite hook
// - Caching: ETag + Last-Modified + 304
// - Range support via http.ServeFile/ServeContent
func NewServeFile(p ServeProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Normalize to exactly one leading slash (works with/without StripPrefix)
		reqPath := r.URL.Path
		if reqPath == "" || reqPath[0] != '/' {
			reqPath = "/" + reqPath
		}
		isRoot := reqPath == "/" || reqPath == ""
		xfp := r.Header.Get("X-Forwarded-Proto")
		xff := r.Header.Get("X-Forwarded-For")
		if isRoot && r.Method == http.MethodGet {
			slog.Info("gateway root request", "host", r.Host, "xfp", xfp, "xff", xff)
		}

		// --- reverse proxy by path prefix ---
		if target, ok := p.ResolveProxy(reqPath); ok {
			u, err := url.Parse(target)
			if err != nil {
				httplib.WriteJSONError(w, http.StatusInternalServerError, "invalid proxy target URL")
				return
			}
			hostURL, err := url.Parse(u.Scheme + "://" + u.Host)
			if err != nil {
				httplib.WriteJSONError(w, http.StatusInternalServerError, "invalid proxy host URL")
				return
			}
			rp := httputil.NewSingleHostReverseProxy(hostURL)

			// forward to full target
			targetURL, err := url.Parse(target)
			if err != nil {
				httplib.WriteJSONError(w, http.StatusInternalServerError, "invalid target URL")
				return
			}
			r.URL = targetURL
			r.Host = u.Host
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

			rp.ServeHTTP(w, r)
			return
		}

		// Local path handling (cleaned)
		rqstPath := path.Clean(reqPath)
		if rqstPath == "/null" {
			httplib.WriteJSONError(w, http.StatusBadRequest, "No file path was given in the file url path!")
			return
		}

		rqstPath, cleanErr := sanitizeRequestPath(rqstPath)
		if cleanErr != nil {
			httplib.WriteJSONError(w, http.StatusBadRequest, "invalid path")
			return
		}
		protoHint := strings.TrimSpace(xfp)
		if protoHint == "" {
			protoHint = strings.ToLower(strings.TrimSpace(r.URL.Scheme))
		}
		if protoHint == "" {
			protoHint = strings.ToLower(strings.TrimSpace(strings.Split(r.Proto, "/")[0]))
		}
		if protoHint == "" {
			protoHint = "http"
		}
		marker := fmt.Sprintf("Served-By: gateway at %s host=%s proto=%s", time.Now().UTC().Format(time.RFC3339), r.Host, protoHint)
		w.Header().Set("X-Served-By", marker)

		// v1 Conformance: Use stable webroot path
		// REMOVED: Host header-based directory selection (security violation INV-1.2)
		// Host header is untrusted client input and MUST NOT determine file access paths
		// For multi-tenancy, use authenticated principalID instead of Host header
		dir := p.WebRoot()

		// v2 Conformance: Token from Authorization header ONLY (security violation INV-3.1)
		// REMOVED: Query parameter token extraction - tokens in URLs leak via logs
		// Tokens MUST be in Authorization header to prevent exposure in:
		// - Access logs, proxy logs, browser history
		// - Referer headers, URL sharing
		// - Server-side request logging
		app := r.Header.Get("application") // Application can still use header fallback
		token := r.Header.Get("token")     // Token: Header ONLY (no query fallback)
		if token == "null" || token == "undefined" {
			token = ""
		}

		// Index resolution (keep legacy behavior)
		if rqstPath == "/" {
			if idx := p.IndexApplication(); idx != "" {
				rqstPath = "/" + idx
				app = idx
			} else {
				rqstPath = "/index.html"
			}
		} else if strings.Count(rqstPath, "/") == 1 {
			if hasExt(rqstPath, ".js", ".json", ".css", ".htm", ".html") && p.Exists(filepath.Join(dir, rqstPath)) {
				if idx := p.IndexApplication(); idx != "" {
					rqstPath = "/" + idx + rqstPath
				}
			}
		}

		// Protected areas under data/files
		hasAccess := true
		if strings.HasPrefix(rqstPath, "/users/") {
			hasAccess = false

		}

		// Allow access to hidden directories and HLS streaming files
		if strings.Contains(rqstPath, "/.hidden/") || isHLSFile(rqstPath) {
			hasAccess = true
		}

		// Windows drive quirk: "/C:..." -> "C:..."
		if len(rqstPath) > 3 && runtime.GOOS == "windows" && rqstPath[0] == '/' && rqstPath[2] == ':' {
			rqstPath = rqstPath[1:]
		}

		minioCfg, minioErr := p.FileServiceMinioConfig()
		isUsersPath := strings.HasPrefix(rqstPath, "/users/")
		minioConfigured := minioCfg != nil || minioErr != nil
		if minioConfigured && minioErr != nil {
			httplib.WriteJSONError(w, http.StatusServiceUnavailable, objectStoreErrMsg(minioErr))
			return
		}
		hasMinio := minioConfigured && minioErr == nil
		useMinioUsers := hasMinio && isUsersPath
		useMinioWeb := hasMinio && !isUsersPath
		// Only fallback to disk if MinIO is not configured at all
		fallbackToDisk := !minioConfigured

		// Compute filename; "public" paths are absolute and force validation
		name := filepath.Join(dir, rqstPath)
		// Directory request should redirect to a playlist manifest when available
		if info, err := os.Stat(name); err == nil && info.IsDir() && (r.Method == http.MethodGet || r.Method == http.MethodHead) {
			playlist := filepath.Join(name, "playlist.m3u8")
			if p.Exists(playlist) {
				redirectPath := path.Join(rqstPath, "playlist.m3u8")
				if q := r.URL.RawQuery; q != "" {
					redirectPath += "?" + q
				}
				http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
				return
			}
		}
		if isPublicLike(rqstPath, p.PublicDirs()) {
			name = rqstPath
			hasAccess = false
		}

		// Streaming allow list (hidden directories and HLS files)
		if strings.Contains(rqstPath, "/.hidden/") || isHLSFile(rqstPath) {
			hasAccess = true
		}

		// CA certificate special-case
		if rqstPath == "/ca.crt" {
			name = filepath.Join(p.CredsDir(), rqstPath)
		}

		// Access checks
		var (
			hasDenied bool
			err       error
		)
		if token != "" && !hasAccess {
			if uid, e := p.ParseUserID(token); e == nil && uid != "" {
				hasAccess, hasDenied, err = p.ValidateAccount(uid, "read", rqstPath)
			} else if e != nil {
				httplib.WriteJSONError(w, http.StatusUnauthorized, "invalid access token")
				return
			}
		}

		if isPublicLike(rqstPath, p.PublicDirs()) && !hasDenied && !hasAccess {
			hasAccess = true
		} else if !hasAccess && !hasDenied && app != "" {
			hasAccess, hasDenied, err = p.ValidateApplication(app, "read", rqstPath)
		}
		if !hasAccess || hasDenied || err != nil {
			httplib.WriteJSONError(w, http.StatusUnauthorized, "unable to read the file "+rqstPath+"; check your access privilege")
			return
		}

		if useMinioUsers && serveUsersFromMinio(w, r, minioCfg, rqstPath) {
			return
		}

		hostPart := cleanRequestHost(r.Host, minioCfg)
		if useMinioWeb && serveWebrootFromMinio(w, r, minioCfg, rqstPath, hostPart, fallbackToDisk) {
			return
		}

		// Streaming hook
		lname := strings.ToLower(name)
		if r.Method == http.MethodGet && strings.HasSuffix(lname, ".mkv") {
			if p.MaybeStream(name, w, r) {
				return
			}
		}

		// Open or hashed fallback
		f, ferr := openFile(p, name)
		if ferr != nil {
			httplib.WriteJSONError(w, http.StatusNotFound, ferr.Error())
			return
		}
		defer f.Close()

		// Stat once for caching/etag/range
		fi, _ := f.Stat()
		mod := fi.ModTime().UTC()
		etag := weakETag(fi)

		// Conditional GET (304)
		if inm := r.Header.Get("If-None-Match"); inm != "" && etagMatch(etag, inm) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if ims := r.Header.Get("If-Modified-Since"); ims != "" {
			if t, err := time.Parse(http.TimeFormat, ims); err == nil && !fi.ModTime().After(t) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		w.Header().Set("ETag", etag)
		w.Header().Set("Last-Modified", mod.Format(http.TimeFormat))

		// Content type & optional JS import rewriting
		switch {
		case strings.HasSuffix(lname, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
			code, changed := maybeRewriteJSImports(p, rqstPath, f)
			if changed {
				http.ServeContent(w, r, name, fi.ModTime(), strings.NewReader(code))
				return
			}
		case strings.HasSuffix(lname, ".css"):
			w.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(lname, ".html") || strings.HasSuffix(lname, ".htm"):
			w.Header().Set("Content-Type", "text/html")
			if isRoot {
				if data, err := os.ReadFile(name); err == nil {
					htmlWithMarker := fmt.Sprintf("<!-- %s -->\n<div id=\"served-by\">%s</div>\n%s", marker, marker, string(data))
					http.ServeContent(w, r, name, fi.ModTime(), strings.NewReader(htmlWithMarker))
					return
				}
			}
		}

		// Default: Range + Last-Modified supported by stdlib
		http.ServeFile(w, r, name)
	})
}

func serveUsersFromMinio(w http.ResponseWriter, r *http.Request, cfg *MinioProxyConfig, rqstPath string) bool {
	if cfg == nil || cfg.Fetch == nil {
		return false
	}
	key, err := usersObjectKey(cfg, rqstPath)
	if err != nil {
		httplib.WriteJSONError(w, http.StatusBadRequest, "invalid path")
		return true
	}
	reader, info, err := cfg.Fetch(r.Context(), cfg.usersBucket(), key)
	if err != nil {
		if isMinioNoSuchKey(err) {
			if isDirCandidate(rqstPath) {
				playlistPath := path.Join(rqstPath, "playlist.m3u8")
				playlistKey, keyErr := usersObjectKey(cfg, playlistPath)
				if keyErr != nil {
					httplib.WriteJSONError(w, http.StatusBadRequest, "invalid path")
					return true
				}
				playlistReader, _, perr := cfg.Fetch(r.Context(), cfg.usersBucket(), playlistKey)
				if perr == nil {
					_ = playlistReader.Close()
					redirectToPlaylist(w, r, playlistPath)
					return true
				}
				if !isMinioNoSuchKey(perr) {
					httplib.WriteJSONError(w, http.StatusServiceUnavailable, "object store unavailable")
					return true
				}
			}
			return false
		}
		httplib.WriteJSONError(w, http.StatusServiceUnavailable, "object store unavailable")
		return true
	}
	defer reader.Close()

	name := path.Base(key)
	serveMinioContent(w, r, name, info, reader)
	return true
}

func redirectToPlaylist(w http.ResponseWriter, r *http.Request, playlistPath string) {
	redirectPath := playlistPath
	if q := r.URL.RawQuery; q != "" {
		redirectPath += "?" + q
	}
	http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
}

func isDirCandidate(rqstPath string) bool {
	if rqstPath == "" || rqstPath == "/" {
		return false
	}
	base := path.Base(rqstPath)
	return !strings.Contains(base, ".")
}

func serveWebrootFromMinio(w http.ResponseWriter, r *http.Request, cfg *MinioProxyConfig, rqstPath, host string, fallbackOnError bool) bool {
	if cfg == nil || cfg.Fetch == nil {
		return false
	}
	key, err := webrootObjectKey(cfg, host, rqstPath)
	if err != nil {
		httplib.WriteJSONError(w, http.StatusBadRequest, "invalid path")
		return true
	}
	reader, info, err := cfg.Fetch(r.Context(), cfg.webrootBucket(), key)
	if err != nil {
		if isMinioNoSuchKey(err) {
			return false
		}
		if !fallbackOnError {
			httplib.WriteJSONError(w, http.StatusServiceUnavailable, "object store unavailable")
			return true
		}
		return false
	}
	defer reader.Close()
	name := path.Base(key)
	serveMinioContent(w, r, name, info, reader)
	return true
}

func serveMinioContent(w http.ResponseWriter, r *http.Request, name string, info MinioObjectInfo, reader io.ReadSeeker) {
	mod := info.ModTime.UTC()
	if mod.IsZero() {
		mod = time.Now().UTC()
	}
	etag := weakMinioETag(info)

	// Conditional GET (304)
	if inm := r.Header.Get("If-None-Match"); inm != "" && etagMatch(etag, inm) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	if ims := r.Header.Get("If-Modified-Since"); ims != "" {
		if t, err := time.Parse(http.TimeFormat, ims); err == nil && !mod.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("ETag", etag)
	w.Header().Set("Last-Modified", mod.Format(http.TimeFormat))
	setContentTypeFromName(w, name)
	http.ServeContent(w, r, name, mod, reader)
}

func isMinioNoSuchKey(err error) bool {
	if err == nil {
		return false
	}
	if resp := minio.ToErrorResponse(err); resp.Code == "NoSuchKey" {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "nosuchkey")
}

func (cfg *MinioProxyConfig) usersBucket() string {
	if cfg == nil {
		return ""
	}
	if bucket := strings.TrimSpace(cfg.UsersBucket); bucket != "" {
		return bucket
	}
	return strings.TrimSpace(cfg.Bucket)
}

func (cfg *MinioProxyConfig) webrootBucket() string {
	if cfg == nil {
		return ""
	}
	if bucket := strings.TrimSpace(cfg.WebrootBucket); bucket != "" {
		return bucket
	}
	return strings.TrimSpace(cfg.Bucket)
}

func (cfg *MinioProxyConfig) usersPrefixValue() string {
	if cfg == nil {
		return ""
	}
	if p := strings.Trim(cfg.UsersPrefix, "/"); p != "" {
		return p
	}
	if p := strings.Trim(cfg.Prefix, "/"); p != "" {
		return path.Join(p, "users")
	}
	return path.Join(defaultDomain(cfg), "users")
}

func (cfg *MinioProxyConfig) webrootPrefixValue() string {
	if cfg == nil {
		return ""
	}
	if p := strings.Trim(cfg.WebrootPrefix, "/"); p != "" {
		return p
	}
	if p := strings.Trim(cfg.Prefix, "/"); p != "" {
		return path.Join(p, "webroot")
	}
	return path.Join(defaultDomain(cfg), "webroot")
}

func defaultDomain(cfg *MinioProxyConfig) string {
	if cfg != nil && strings.TrimSpace(cfg.Domain) != "" {
		return strings.TrimSpace(cfg.Domain)
	}
	return "localhost"
}

func usersObjectKey(cfg *MinioProxyConfig, rqstPath string) (string, error) {
	cleanPath, err := sanitizeRequestPath(rqstPath)
	if err != nil {
		return "", err
	}
	logical := strings.TrimPrefix(cleanPath, "/users/")
	return joinKey(cfg.usersPrefixValue(), logical)
}

func webrootObjectKey(cfg *MinioProxyConfig, host, rqstPath string) (string, error) {
	cleanPath, err := sanitizeRequestPath(rqstPath)
	if err != nil {
		return "", err
	}
	host = strings.TrimSpace(host)
	if host == "" {
		host = defaultDomain(cfg)
	}
	host = strings.Split(host, ":")[0]
	host = strings.Trim(host, " /.")
	if host == "" {
		host = defaultDomain(cfg)
	}
	logical := strings.TrimPrefix(cleanPath, "/")
	if logical == "" {
		logical = "index.html"
	}
	switch {
	case strings.TrimSpace(cfg.WebrootPrefix) != "":
		return joinKey(cfg.webrootPrefixValue(), logical)
	case strings.TrimSpace(cfg.Prefix) != "":
		return joinKey(cfg.webrootPrefixValue(), logical)
	default:
		// v1 Conformance: Use stable prefix (security violation INV-1.4)
		// REMOVED: path.Join(host, "webroot") - Host header MUST NOT determine storage paths
		// Host is untrusted client input - using it creates tenant isolation bypass
		// For multi-tenancy, use authenticated principalID in explicit prefix config
		base := "webroot" // Stable prefix, independent of Host header
		return joinKey(base, logical)
	}
}

// Exposed for tests
func WebrootObjectKeyForTest(cfg *MinioProxyConfig, host, rqstPath string) (string, error) {
	return webrootObjectKey(cfg, host, rqstPath)
}

func joinKey(parts ...string) (string, error) {
	cleaned := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		cleanPart, err := safeKeyPart(p)
		if err != nil {
			return "", err
		}
		if cleanPart != "" {
			cleaned = append(cleaned, cleanPart)
		}
	}
	return strings.Join(cleaned, "/"), nil
}

func safeKeyPart(part string) (string, error) {
	normalized := path.Clean(strings.ReplaceAll(part, "\\", "/"))
	normalized = strings.Trim(normalized, "/")
	if normalized == "" || normalized == "." {
		return "", nil
	}
	if normalized == ".." || strings.HasPrefix(normalized, "../") || strings.Contains(normalized, "/../") {
		return "", fmt.Errorf("invalid path segment")
	}
	return normalized, nil
}

func sanitizeRequestPath(p string) (string, error) {
	if p == "" {
		return "/", nil
	}
	raw := strings.ReplaceAll(p, "\\", "/")
	for _, seg := range strings.Split(raw, "/") {
		if seg == ".." {
			return "", fmt.Errorf("invalid path traversal")
		}
	}
	cleaned := path.Clean("/" + strings.TrimLeft(strings.TrimSpace(p), "/"))
	if strings.HasPrefix(cleaned, "/..") || strings.Contains(cleaned, "/../") {
		return "", fmt.Errorf("invalid path traversal")
	}
	return cleaned, nil
}

func cleanRequestHost(host string, cfg *MinioProxyConfig) string {
	h := strings.TrimSpace(host)
	if h == "" && cfg != nil {
		h = cfg.Domain
	}
	h = strings.Split(h, ":")[0]
	h = strings.Trim(h, " /.")
	if h == "" && cfg != nil {
		h = cfg.Domain
	}
	if h == "" {
		h = "localhost"
	}
	return h
}

func setContentTypeFromName(w http.ResponseWriter, name string) {
	lname := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lname, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	case strings.HasSuffix(lname, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(lname, ".html") || strings.HasSuffix(lname, ".htm"):
		w.Header().Set("Content-Type", "text/html")
	}
}

func weakMinioETag(info MinioObjectInfo) string {
	mod := info.ModTime
	if mod.IsZero() {
		mod = time.Now().UTC()
	}
	return fmt.Sprintf(`W/"%d-%d"`, info.Size, mod.Unix())
}

func objectStoreErrMsg(err error) string {
	if err == nil {
		return "object store unavailable"
	}
	return err.Error()
}

// --- helpers (internal) ---
func headerOrQuery(r *http.Request, key string) string {
	if v := r.Header.Get(key); v != "" {
		return v
	}
	return r.URL.Query().Get(key)
}

func hasExt(p string, exts ...string) bool {
	lp := strings.ToLower(p)
	for _, e := range exts {
		if strings.HasSuffix(lp, e) {
			return true
		}
	}
	return false
}

func isPublicLike(reqPath string, publicRoots []string) bool {
	req := filepath.Clean(reqPath)
	for _, root := range publicRoots {
		root = filepath.Clean(root)
		if req == root || strings.HasPrefix(req+string(filepath.Separator), root+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func maybeRewriteJSImports(p ServeProvider, base string, f io.ReadSeeker) (code string, changed bool) {
	// reset to start
	_, _ = f.Seek(0, io.SeekStart)
	sc := bufio.NewScanner(f)
	var b strings.Builder

	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "import") && strings.Contains(line, `'@`) {
			if newPath, err := p.ResolveImportPath(base, line); err == nil && newPath != "" {
				line = line[:strings.Index(line, `'@`)] + `'` + newPath + `'`
				changed = true
			}
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String(), changed
}

func openFile(p ServeProvider, name string) (*os.File, error) {
	if p.Exists(name) {
		return os.Open(name) // #nosec G304
	}
	if p.Exists("/" + name) {
		return os.Open("/" + name)
	}
	if alt, err := p.FindHashedFile(name); err == nil && alt != "" && p.Exists(alt) {
		return os.Open(alt)
	}
	return nil, fmt.Errorf("file %s not found", name)
}

func weakETag(fi os.FileInfo) string {
	return fmt.Sprintf(`W/"%d-%d"`, fi.Size(), fi.ModTime().Unix())
}

func etagMatch(etag, list string) bool {
	// Very small matcher: "*", or exact match among comma-separated list
	for _, v := range strings.Split(list, ",") {
		if strings.TrimSpace(v) == "*" || strings.TrimSpace(v) == etag {
			return true
		}
	}
	return false
}
