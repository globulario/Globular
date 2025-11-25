package files

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// UploadProvider abstracts write roots and access checks for uploads.
type UploadProvider interface {
	DataRoot() string
	PublicDirs() []string

	ParseUserID(token string) (string, error)
	ValidateAccount(userID, action, reqPath string) (has, denied bool, err error)
	ValidateApplication(app, action, reqPath string) (has, denied bool, err error)
	AddResourceOwner(token, path, owner, resourceType string) error
}

type UploadOptions struct {
	MaxBytes    int64    // maximum request size; defaults to 100 MiB
	AllowedExts []string // optional allow-list (lowercase, with dots), e.g. []string{".png",".jpg",".pdf"}
}

// NewUploadFileWithOptions lets you control size & allowlist.
// NewUploadFileWithOptions lets you control size & allowlist.
// Compatible with client FormData:
//
//	multiplefiles: <File> (one or many)
//	path:          "/users/xxx/dir" (or "/applications/...")  <-- used here
//
// Auth (same as before):
//
//	token header or ?token=..., application header or ?application=...
func NewUploadFileWithOptions(p UploadProvider, opt UploadOptions) http.Handler {
	if opt.MaxBytes <= 0 {
		opt.MaxBytes = 100 << 20 // 100 MiB default
	}
	allowed := make(map[string]struct{}, len(opt.AllowedExts))
	for _, e := range opt.AllowedExts {
		allowed[strings.ToLower(e)] = struct{}{}
	}

	type resp struct {
		Paths []string `json:"paths"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Cap request size early.
		r.Body = http.MaxBytesReader(w, r.Body, opt.MaxBytes)

		// Read auth context the same way as before.
		app := firstNonEmpty(r.Header.Get("application"), r.URL.Query().Get("application"), r.FormValue("application"))
		token := firstNonEmpty(r.Header.Get("token"), r.URL.Query().Get("token"), r.FormValue("token"))

		// Client sends "path" (not "dir"). Keep backward-compat on query "dir" just in case.
		dir := firstNonEmpty(r.FormValue("path"), r.URL.Query().Get("path"), r.URL.Query().Get("dir"), r.FormValue("dir"))
		if dir == "" {
			http.Error(w, "missing 'path' to upload into", http.StatusBadRequest)
			return
		}
		dir = pathCleanOS(dir)

		publicDirs := p.PublicDirs()

		var (
			targetDir   string
			cleanTarget string
			addOwner    bool
		)

		if isRBACProtected(dir) {
			root := filepath.Join(p.DataRoot(), "files")
			targetDir = filepath.Join(root, strings.TrimPrefix(dir, "/"))

			cleanRoot := filepath.Clean(root) + string(filepath.Separator)
			cleanTarget = filepath.Clean(targetDir) + string(filepath.Separator)
			if !strings.HasPrefix(cleanTarget, cleanRoot) {
				http.Error(w, "invalid path", http.StatusBadRequest)
				return
			}
			addOwner = true
		} else if resolved, ok := resolvePublicDir(dir, publicDirs); ok {
			targetDir = resolved
			cleanTarget = filepath.Clean(targetDir) + string(filepath.Separator)
		} else {
			http.Error(w, "path not allowed", http.StatusBadRequest)
			return
		}

		// Validate write access (account first via token; fallback to application).
		has, denied, err := false, false, error(nil)
		var ownerID string
		if token != "" {
			if uid, e := p.ParseUserID(token); e == nil && uid != "" {
				ownerID = uid
				has, denied, err = p.ValidateAccount(uid, "write", dir)
			}
		}
		if !has && !denied && app != "" {
			has, denied, err = p.ValidateApplication(app, "write", dir)
		}
		if !has || denied || err != nil {
			fmt.Println("UploadFile: access denied for", ownerID, "app", app, "to", dir, "has:", has, "denied:", denied, "err:", err)
			http.Error(w, "write access denied", http.StatusUnauthorized)
			return
		}

		// Parse multipart (use same cap as MaxBytes; large parts spill to temp files).
		if err := r.ParseMultipartForm(opt.MaxBytes); err != nil {
			status := http.StatusBadRequest
			if strings.Contains(strings.ToLower(err.Error()), "request body too large") {
				status = http.StatusRequestEntityTooLarge
			}
			http.Error(w, "multipart error: "+err.Error(), status)
			return
		}

		// Read files from "multiplefiles" (this matches your client)
		mf := r.MultipartForm
		fhs := mf.File["multiplefiles"]
		if len(fhs) == 0 {
			// Some browsers use a single part but same key; try FormFile as a fallback
			if f, h, e := r.FormFile("multiplefiles"); e == nil && h != nil {
				_ = f.Close()
				fhs = []*multipart.FileHeader{h}
			}
		}
		if len(fhs) == 0 {
			http.Error(w, "no files provided (multiplefiles)", http.StatusBadRequest)
			return
		}

		// Ensure destination directory exists.
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			http.Error(w, "mkdir error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		writtenPaths := make([]string, 0, len(fhs))

		for _, fh := range fhs {
			// Optional allowlist by extension
			if len(allowed) > 0 {
				ext := strings.ToLower(filepath.Ext(fh.Filename))
				if _, ok := allowed[ext]; !ok {
					http.Error(w, "file type not allowed: "+fh.Filename, http.StatusBadRequest)
					return
				}
			}

			src, err := fh.Open()
			if err != nil {
				http.Error(w, "form file open error: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Create destination file
			dst := filepath.Join(targetDir, fh.Filename)

			// Prevent escaping target dir by re-checking after join
			if !strings.HasPrefix(filepath.Clean(dst)+string(filepath.Separator), cleanTarget) &&
				!strings.HasPrefix(filepath.Clean(dst), strings.TrimSuffix(cleanTarget, string(filepath.Separator))) {
				_ = src.Close()
				http.Error(w, "invalid file name", http.StatusBadRequest)
				return
			}

			out, err := os.Create(dst)
			if err != nil {
				_ = src.Close()
				http.Error(w, "create error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			n, err := io.Copy(out, src)
			_ = src.Close()
			closeErr := out.Close()
			if err != nil {
				// io.Copy can surface body-too-large here as well
				if strings.Contains(strings.ToLower(err.Error()), "request body too large") {
					http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
					return
				}
				http.Error(w, "write error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if closeErr != nil {
				http.Error(w, "close error: "+closeErr.Error(), http.StatusInternalServerError)
				return
			}
			if n == 0 {
				http.Error(w, "empty file: "+fh.Filename, http.StatusBadRequest)
				return
			}

			if ownerID != "" && addOwner {
				logicalPath := path.Join(dir, fh.Filename)
				if err := p.AddResourceOwner(token, logicalPath, ownerID, "file"); err != nil {
					_ = os.Remove(dst)
					http.Error(w, "ownership error: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			writtenPaths = append(writtenPaths, dst)
		}

		// Minimal JSON response (client only checks status; returning paths is harmless)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp{Paths: writtenPaths})
	})
}

// Default (100 MiB, allow all)
func NewUploadFile(p UploadProvider) http.Handler {
	return NewUploadFileWithOptions(p, UploadOptions{})
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" && v != "null" && v != "undefined" {
			return v
		}
	}
	return ""
}

func pathCleanOS(p string) string {
	c := filepath.ToSlash(filepath.Clean(p))
	// Windows "/C:..." quirk -> "C:..."
	if len(c) > 3 && runtime.GOOS == "windows" && c[0] == '/' && c[2] == ':' {
		c = c[1:]
	}
	return c
}

func isRBACProtected(dir string) bool {
	return strings.HasPrefix(dir, "/users/") ||
		strings.HasPrefix(dir, "/applications/") ||
		strings.HasPrefix(dir, "/templates/") ||
		strings.HasPrefix(dir, "/projects/")
}

func resolvePublicDir(dir string, publicRoots []string) (string, bool) {
	cleanDir := filepath.Clean(dir)
	sep := string(filepath.Separator)
	for _, root := range publicRoots {
		cleanRoot := filepath.Clean(root)
		if cleanRoot == "" {
			continue
		}
		if cleanDir == cleanRoot ||
			strings.HasPrefix(cleanDir+sep, cleanRoot+sep) {
			return cleanDir, true
		}
	}
	return "", false
}

// keep imports tidy
var _ = multipart.FileHeader{} // ensure the package is retained if tools strip unused
