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

	httplib "github.com/globulario/Globular/internal/http"
)

// MinioProxyConfig captures the subset of FileService MinIO settings required
// to proxy /users/ requests directly to the MinIO gateway.
type MinioProxyConfig struct {
	Endpoint string
	Bucket   string
	Prefix   string
	UseSSL   bool
	Fetch    MinioFetchFunc
	Put      MinioPutFunc
	Delete   MinioDeleteFunc
}

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
	FileServiceMinioConfig() (*MinioProxyConfig, bool)

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

		// --- reverse proxy by path prefix ---
		if target, ok := p.ResolveProxy(reqPath); ok {
			u, _ := url.Parse(target)
			hostURL, _ := url.Parse(u.Scheme + "://" + u.Host)
			rp := httputil.NewSingleHostReverseProxy(hostURL)

			// forward to full target
			r.URL, _ = url.Parse(target)
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

		minioCfg, hasMinio := p.FileServiceMinioConfig()
		useMinioUsers := hasMinio && minioCfg != nil && strings.HasPrefix(rqstPath, "/users/")

		// Base root from webRoot, with vhost subdir if present
		dir := p.WebRoot()
		if p.Exists(filepath.Join(dir, r.Host)) {
			dir = filepath.Join(dir, r.Host)
		}

		// token & application from header or query
		app := headerOrQuery(r, "application")
		token := headerOrQuery(r, "token")
		if token == "null" || token == "undefined" {
			token = ""
		}

		// Index resolution (keep legacy behavior)
		if rqstPath == "/" {
			if idx := p.IndexApplication(); idx != "" {
				rqstPath = "/" + idx
				app = idx
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
			if !strings.Contains(rqstPath, "/.hidden/") {
				hasAccess = false
			}
		}

		// Windows drive quirk: "/C:..." -> "C:..."
		if len(rqstPath) > 3 && runtime.GOOS == "windows" && rqstPath[0] == '/' && rqstPath[2] == ':' {
			rqstPath = rqstPath[1:]
		}

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

		// Streaming allow list
		if strings.Contains(rqstPath, "/.hidden/") ||
			strings.HasSuffix(rqstPath, ".ts") ||
			strings.HasSuffix(rqstPath, "240p.m3u8") ||
			strings.HasSuffix(rqstPath, "360p.m3u8") ||
			strings.HasSuffix(rqstPath, "480p.m3u8") ||
			strings.HasSuffix(rqstPath, "720p.m3u8") ||
			strings.HasSuffix(rqstPath, "1080p.m3u8") ||
			strings.HasSuffix(rqstPath, "2160p.m3u8") {
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

		if useMinioUsers {
			if serveUsersFromMinio(w, r, minioCfg, rqstPath) {
				return
			}
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
			httplib.WriteJSONError(w, http.StatusNoContent, ferr.Error())
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
		}

		// Default: Range + Last-Modified supported by stdlib
		http.ServeFile(w, r, name)
	})
}

func serveUsersFromMinio(w http.ResponseWriter, r *http.Request, cfg *MinioProxyConfig, rqstPath string) bool {
	if cfg == nil || cfg.Fetch == nil {
		return false
	}
	key := serveMinioObjectKey(cfg, rqstPath)
	reader, info, err := cfg.Fetch(r.Context(), cfg.Bucket, key)
	if err != nil {
		if isMinioNoSuchKey(err) && isDirCandidate(rqstPath) {
			playlistPath := path.Join(rqstPath, "playlist.m3u8")
			if minioObjectExists(r.Context(), cfg, playlistPath) {
				redirectToPlaylist(w, r, playlistPath)
				return true
			}
		}
		httplib.WriteJSONError(w, http.StatusBadGateway, "failed to reach storage: "+err.Error())
		return true
	}
	defer reader.Close()
	mod := info.ModTime
	if mod.IsZero() {
		mod = time.Now().UTC()
	}
	name := path.Base(key)
	http.ServeContent(w, r, name, mod, reader)
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

func minioObjectExists(ctx context.Context, cfg *MinioProxyConfig, rqstPath string) bool {
	if cfg == nil || cfg.Fetch == nil {
		return false
	}
	key := serveMinioObjectKey(cfg, rqstPath)
	reader, _, err := cfg.Fetch(ctx, cfg.Bucket, key)
	if err != nil {
		return false
	}
	_ = reader.Close()
	return true
}

func serveMinioObjectKey(cfg *MinioProxyConfig, rqstPath string) string {
	key := strings.TrimPrefix(rqstPath, "/")
	prefix := strings.Trim(cfg.Prefix, "/")
	if prefix != "" {
		return path.Join(prefix, key)
	}
	return key
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
