package files

import (
	"context"
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

	httplib "github.com/globulario/Globular/internal/gateway/http"
)

// UploadProvider abstracts write roots and access checks for uploads.
type UploadProvider interface {
	DataRoot() string
	PublicDirs() []string

	ParseUserID(token string) (string, error)
	ValidateAccount(userID, action, reqPath string) (has, denied bool, err error)
	ValidateApplication(app, action, reqPath string) (has, denied bool, err error)
	AddResourceOwner(path, resourceType, owner string) error
	FileServiceMinioConfig() (*MinioProxyConfig, error)
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
			httplib.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
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
			httplib.WriteJSONError(w, http.StatusBadRequest, "missing 'path' to upload into")
			return
		}
		dir = pathCleanOS(dir)

		publicDirs := p.PublicDirs()

		var (
			targetDir       string
			cleanTarget     string
			addOwner        bool
			isWebrootTarget bool
			minioTargetPath bool
		)

		filesRoot := filepath.Join(p.DataRoot(), "files")
		if isRBACProtected(dir) {
			targetDir = filepath.Join(filesRoot, filepath.FromSlash(strings.TrimPrefix(dir, "/")))
			cleanTarget = filepath.Clean(targetDir) + string(filepath.Separator)
			addOwner = true
			minioTargetPath = strings.HasPrefix(dir, "/users/")
		} else if resolved, ok := resolvePublicDir(dir, publicDirs); ok {
			targetDir = resolved
			cleanTarget = filepath.Clean(targetDir) + string(filepath.Separator)
		} else {
			targetDir = filepath.Join(p.DataRoot(), "webroot", filepath.FromSlash(strings.TrimPrefix(dir, "/")))
			cleanTarget = filepath.Clean(targetDir) + string(filepath.Separator)
			isWebrootTarget = true
			minioTargetPath = true
		}

		minioCfg, minioErr := p.FileServiceMinioConfig()
		if minioErr != nil && strings.HasPrefix(dir, "/users/") {
			httplib.WriteJSONError(w, http.StatusServiceUnavailable, objectStoreErrMsg(minioErr))
			return
		}
		useMinio := minioErr == nil && minioCfg != nil && minioCfg.Put != nil && minioTargetPath

		// Validate write access (account first via token; fallback to application).
		has, denied, err := false, false, error(nil)
		var ownerID string
		action := "write"
		if isWebrootTarget {
			action = "deploy"
		}
		if token != "" {
			if uid, e := p.ParseUserID(token); e == nil && uid != "" {
				ownerID = uid
				has, denied, err = p.ValidateAccount(uid, action, dir)
			}
		}
		if !has && !denied && app != "" {
			has, denied, err = p.ValidateApplication(app, action, dir)
		}
		if !has || denied || err != nil {
			fmt.Println("UploadFile: access denied for", ownerID, "app", app, "to", dir, "has:", has, "denied:", denied, "err:", err)
			httplib.WriteJSONError(w, http.StatusUnauthorized, "write access denied")
			return
		}

		// Parse multipart (use same cap as MaxBytes; large parts spill to temp files).
		if err := r.ParseMultipartForm(opt.MaxBytes); err != nil {
			status := http.StatusBadRequest
			if strings.Contains(strings.ToLower(err.Error()), "request body too large") {
				status = http.StatusRequestEntityTooLarge
			}
			httplib.WriteJSONError(w, status, "multipart error: "+err.Error())
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
			httplib.WriteJSONError(w, http.StatusBadRequest, "no files provided (multiplefiles)")
			return
		}

		// Ensure destination directory exists.
		if !useMinio {
			if err := os.MkdirAll(targetDir, 0o755); err != nil {
				httplib.WriteJSONError(w, http.StatusInternalServerError, "mkdir error: "+err.Error())
				return
			}
		}

		writtenPaths := make([]string, 0, len(fhs))

		for _, fh := range fhs {
			// Optional allowlist by extension
			if len(allowed) > 0 {
				ext := strings.ToLower(filepath.Ext(fh.Filename))
				if _, ok := allowed[ext]; !ok {
					httplib.WriteJSONError(w, http.StatusBadRequest, "file type not allowed: "+fh.Filename)
					return
				}
			}

			logicalPath := path.Join(dir, fh.Filename)

			if useMinio {
				var (
					objKey string
					bucket string
					err    error
				)
				if strings.HasPrefix(dir, "/users/") {
					bucket = minioCfg.usersBucket()
					objKey, err = usersObjectKey(minioCfg, logicalPath)
				} else {
					bucket = minioCfg.webrootBucket()
					host := cleanRequestHost(r.Host, minioCfg)
					objKey, err = webrootObjectKey(minioCfg, host, logicalPath)
				}
				if err != nil {
					httplib.WriteJSONError(w, http.StatusBadRequest, "invalid path")
					return
				}
				if err := uploadFileToMinio(r.Context(), minioCfg, bucket, objKey, fh); err != nil {
					httplib.WriteJSONError(w, http.StatusInternalServerError, err.Error())
					return
				}
				if ownerID != "" && addOwner {
					if err := p.AddResourceOwner(logicalPath, "file", ownerID); err != nil {
						if minioCfg.Delete != nil {
							_ = minioCfg.Delete(r.Context(), bucket, objKey)
						}
						httplib.WriteJSONError(w, http.StatusInternalServerError, "ownership error: "+err.Error())
						return
					}
				}
				writtenPaths = append(writtenPaths, logicalPath)
				continue
			}

			src, err := fh.Open()
			if err != nil {
				httplib.WriteJSONError(w, http.StatusBadRequest, "form file open error: "+err.Error())
				return
			}

			// Create destination file
			dst := filepath.Join(targetDir, fh.Filename)

			// Prevent escaping target dir by re-checking after join
			if !strings.HasPrefix(filepath.Clean(dst)+string(filepath.Separator), cleanTarget) &&
				!strings.HasPrefix(filepath.Clean(dst), strings.TrimSuffix(cleanTarget, string(filepath.Separator))) {
				_ = src.Close()
				httplib.WriteJSONError(w, http.StatusBadRequest, "invalid file name")
				return
			}

			out, err := os.Create(dst)
			if err != nil {
				_ = src.Close()
				httplib.WriteJSONError(w, http.StatusInternalServerError, "create error: "+err.Error())
				return
			}

			n, err := io.Copy(out, src)
			_ = src.Close()
			closeErr := out.Close()
			if err != nil {
				// io.Copy can surface body-too-large here as well
				if strings.Contains(strings.ToLower(err.Error()), "request body too large") {
					httplib.WriteJSONError(w, http.StatusRequestEntityTooLarge, "payload too large")
					return
				}
				httplib.WriteJSONError(w, http.StatusInternalServerError, "write error: "+err.Error())
				return
			}
			if closeErr != nil {
				httplib.WriteJSONError(w, http.StatusInternalServerError, "close error: "+closeErr.Error())
				return
			}
			if n == 0 {
				httplib.WriteJSONError(w, http.StatusBadRequest, "empty file: "+fh.Filename)
				return
			}

			if ownerID != "" && addOwner {
				if err := p.AddResourceOwner(logicalPath, "file", ownerID); err != nil {
					_ = os.Remove(dst)
					httplib.WriteJSONError(w, http.StatusInternalServerError, "ownership error: "+err.Error())
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

func uploadFileToMinio(ctx context.Context, cfg *MinioProxyConfig, bucket, key string, fh *multipart.FileHeader) error {
	if cfg == nil || cfg.Put == nil {
		return fmt.Errorf("storage backend not configured")
	}
	if fh.Size == 0 {
		return fmt.Errorf("empty file: %s", fh.Filename)
	}
	src, err := fh.Open()
	if err != nil {
		return fmt.Errorf("form file open error: %w", err)
	}
	defer src.Close()
	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return cfg.Put(ctx, bucket, key, src, fh.Size, contentType)
}
