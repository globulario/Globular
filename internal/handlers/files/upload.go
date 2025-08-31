package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
}

type UploadOptions struct {
	MaxBytes    int64    // maximum request size; defaults to 100 MiB
	AllowedExts []string // optional allow-list (lowercase, with dots), e.g. []string{".png",".jpg",".pdf"}
}

// NewUploadFileWithOptions lets you control size & allowlist.
func NewUploadFileWithOptions(p UploadProvider, opt UploadOptions) http.Handler {
	if opt.MaxBytes <= 0 {
		opt.MaxBytes = 100 << 20 // 100 MiB default
	}
	allowed := make(map[string]struct{}, len(opt.AllowedExts))
	for _, e := range opt.AllowedExts {
		allowed[strings.ToLower(e)] = struct{}{}
	}

	type resp struct {
		Path string `json:"path"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Enforce max body size early
		r.Body = http.MaxBytesReader(w, r.Body, opt.MaxBytes)

		// token/application
		app := firstNonEmpty(r.Header.Get("application"), r.URL.Query().Get("application"), r.FormValue("application"))
		token := firstNonEmpty(r.Header.Get("token"), r.URL.Query().Get("token"), r.FormValue("token"))

		// target dir
		dir := firstNonEmpty(r.URL.Query().Get("dir"), r.FormValue("dir"))
		if dir == "" {
			http.Error(w, "missing 'dir' to upload into", http.StatusBadRequest)
			return
		}
		dir = pathCleanOS(dir)

		// Only allow protected writable roots
		if !(strings.HasPrefix(dir, "/users/") ||
			strings.HasPrefix(dir, "/applications/") ||
			strings.HasPrefix(dir, "/templates/") ||
			strings.HasPrefix(dir, "/projects/")) {
			http.Error(w, "dir not allowed", http.StatusBadRequest)
			return
		}

		// Map to actual filesystem under data/files
		root := filepath.Join(p.DataRoot(), "files")
		targetDir := filepath.Join(root, strings.TrimPrefix(dir, "/"))
		if !strings.HasPrefix(filepath.Clean(targetDir)+string(filepath.Separator), filepath.Clean(root)+string(filepath.Separator)) {
			http.Error(w, "invalid dir", http.StatusBadRequest)
			return
		}

		// Validate write access
		has, denied, err := false, false, error(nil)
		if token != "" {
			if uid, e := p.ParseUserID(token); e == nil && uid != "" {
				has, denied, err = p.ValidateAccount(uid, "write", dir)
			}
		}
		if !has && !denied && app != "" {
			has, denied, err = p.ValidateApplication(app, "write", dir)
		}
		if !has || denied || err != nil {
			http.Error(w, "write access denied", http.StatusUnauthorized)
			return
		}

		// Parse multipart form with the same max bound
		// Note: ParseMultipartForm stores small parts in memory & larger parts on disk.
		if err := r.ParseMultipartForm(opt.MaxBytes); err != nil {
			// Body too large → 413
			status := http.StatusBadRequest
			if strings.Contains(err.Error(), "request body too large") {
				status = http.StatusRequestEntityTooLarge
			}
			http.Error(w, "multipart error: "+err.Error(), status)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "form file error: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Optional allowlist by extension
		if len(allowed) > 0 {
			ext := strings.ToLower(filepath.Ext(header.Filename))
			if _, ok := allowed[ext]; !ok {
				http.Error(w, "file type not allowed", http.StatusBadRequest)
				return
			}
		}

		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			http.Error(w, "mkdir error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		dst := filepath.Join(targetDir, header.Filename)
		out, err := os.Create(dst)
		if err != nil {
			http.Error(w, "create error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer out.Close()

		n, err := io.Copy(out, file)
		if err != nil {
			// If we hit the bytes cap during io.Copy, it surfaces here
			if strings.Contains(err.Error(), "request body too large") {
				http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "write error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if n == 0 {
			http.Error(w, "empty file", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"path":%q}`, dst)))
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

// keep imports tidy
var _ = multipart.FileHeader{} // ensure the package is retained if tools strip unused
