package files_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"

	handlers "github.com/globulario/Globular/internal/gateway/handlers"
	files "github.com/globulario/Globular/internal/gateway/handlers/files"
)

type fakeServe struct {
	webRoot    string
	dataRoot   string
	credsDir   string
	indexApp   string
	publicDirs []string
	allowRead  bool
	streamed   bool
	minioCfg   *files.MinioProxyConfig
	minioErr   error
	mode       string
}

type nopSeekCloser struct {
	*strings.Reader
}

func (n nopSeekCloser) Close() error { return nil }

func (f fakeServe) WebRoot() string                       { return f.webRoot }
func (f fakeServe) DataRoot() string                      { return f.dataRoot }
func (f fakeServe) CredsDir() string                      { return f.credsDir }
func (f fakeServe) IndexApplication() string              { return f.indexApp }
func (f fakeServe) PublicDirs() []string                  { return f.publicDirs }
func (f fakeServe) Exists(p string) bool                  { _, err := os.Stat(p); return err == nil }
func (f fakeServe) FindHashedFile(string) (string, error) { return "", fmt.Errorf("no-hash") }
func (f fakeServe) FileServiceMinioConfig() (*files.MinioProxyConfig, error) {
	return f.minioCfg, f.minioErr
}
func (f fakeServe) FileServiceMinioConfigStrict(ctx context.Context) (*files.MinioProxyConfig, error) {
	return f.minioCfg, f.minioErr
}
func (f fakeServe) Mode() string {
	if f.mode == "" {
		return "direct"
	}
	return f.mode
}
func (f fakeServe) ParseUserID(tok string) (string, error) {
	if tok == "ok" {
		return "u@d", nil
	}
	return "", fmt.Errorf("bad token")
}
func (f fakeServe) ValidateAccount(string, string, string) (bool, bool, error) {
	if f.allowRead {
		return true, false, nil
	}
	return false, false, nil
}
func (f fakeServe) ValidateApplication(string, string, string) (bool, bool, error) {
	return false, false, nil
}
func (f fakeServe) ResolveImportPath(string, string) (string, error) { return "", fmt.Errorf("no") }
func (f *fakeServe) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	// pretend we streamed successfully
	f.streamed = true
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("STREAM"))
	return true
}
func (f fakeServe) ResolveProxy(string) (string, bool) { return "", false }

func TestServe_DenyWithoutAccess_401(t *testing.T) {
	tmp := t.TempDir()
	p := &fakeServe{webRoot: tmp, dataRoot: tmp, allowRead: false}

	h := files.NewServeFile(p)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/file.txt", nil)
	req.Header.Set("token", "ok") // triggers account validation

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (body: %s)", rr.Code, rr.Body.String())
	}
}

func TestServe_StreamMKV_AllowsAndShortCircuits_200(t *testing.T) {
	tmp := t.TempDir()
	// The path is under /users => requires read permission; grant it.
	p := &fakeServe{webRoot: tmp, dataRoot: tmp, allowRead: true}

	h := files.NewServeFile(p)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/video.mkv", nil)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	if !p.streamed {
		t.Fatalf("expected streaming path to be taken")
	}
}

func TestServe_MinioUsersProxy(t *testing.T) {
	tmp := t.TempDir()
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg: &files.MinioProxyConfig{
			Endpoint: "ignored",
			Bucket:   "bucket",
			Prefix:   "files",
			Fetch: func(_ context.Context, bucket, key string) (io.ReadSeekCloser, files.MinioObjectInfo, error) {
				if bucket != "bucket" {
					t.Fatalf("expected bucket bucket, got %s", bucket)
				}
				if key != "files/users/alice/avatar.png" {
					t.Fatalf("unexpected key %s", key)
				}
				return nopSeekCloser{Reader: strings.NewReader("MINIO")}, files.MinioObjectInfo{ModTime: time.Unix(10, 0)}, nil
			},
			Stat: func(_ context.Context, bucket, key string) (files.MinioObjectInfo, error) {
				if bucket != "bucket" {
					t.Fatalf("expected bucket bucket, got %s stat", bucket)
				}
				switch key {
				case "files/users/alice/dir":
					return files.MinioObjectInfo{}, minio.ErrorResponse{Code: "NoSuchKey", Message: "missing"}
				case "files/users/alice/dir/playlist.m3u8":
					return files.MinioObjectInfo{ModTime: time.Unix(10, 0)}, nil
				default:
					t.Fatalf("unexpected stat key %s", key)
				}
				return files.MinioObjectInfo{}, fmt.Errorf("unreachable")
			},
		},
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/avatar.png", nil)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if body := rr.Body.String(); body != "MINIO" {
		t.Fatalf("expected backend payload, got %q", body)
	}
}

func TestServe_MinioDirectoryRedirectsToPlaylist(t *testing.T) {
	tmp := t.TempDir()
	var keys []string
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg: &files.MinioProxyConfig{
			Endpoint: "ignored",
			Bucket:   "bucket",
			Prefix:   "files",
			Fetch: func(_ context.Context, bucket, key string) (io.ReadSeekCloser, files.MinioObjectInfo, error) {
				if bucket != "bucket" {
					t.Fatalf("expected bucket bucket, got %s", bucket)
				}
				keys = append(keys, key)
				switch key {
				case "files/users/alice/dir":
					return nil, files.MinioObjectInfo{}, minio.ErrorResponse{Code: "NoSuchKey", Message: "missing"}
				case "files/users/alice/dir/playlist.m3u8":
					return nopSeekCloser{Reader: strings.NewReader("PLAYLIST")}, files.MinioObjectInfo{ModTime: time.Unix(10, 0)}, nil
				default:
					t.Fatalf("unexpected key %s", key)
				}
				return nil, files.MinioObjectInfo{}, fmt.Errorf("unreachable")
			},
		},
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/dir", nil)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/users/alice/dir/playlist.m3u8" {
		t.Fatalf("unexpected location %s", loc)
	}
	if got, want := keys, []string{"files/users/alice/dir", "files/users/alice/dir/playlist.m3u8"}; !equalStringSlices(got, want) {
		t.Fatalf("unexpected keys: %v", got)
	}
}

func TestServe_MinioUnavailable503(t *testing.T) {
	tmp := t.TempDir()
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg:  &files.MinioProxyConfig{Bucket: "bucket"},
		minioErr:  handlers.ErrObjectStoreUnavailable,
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/avatar.png", nil)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d (body: %s)", rr.Code, rr.Body.String())
	}
}

func TestServe_MinioNotProvisionedMessage(t *testing.T) {
	tmp := t.TempDir()
	errMsg := "objectstore not provisioned: missing bucket globular (run node-agent plan ensure-objectstore-layout)"
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg:  &files.MinioProxyConfig{Bucket: "bucket"},
		minioErr:  fmt.Errorf("%w: %s", handlers.ErrObjectStoreUnavailable, errMsg),
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/avatar.png", nil)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ensure-objectstore-layout") {
		t.Fatalf("expected actionable error, got %s", rr.Body.String())
	}
}

func TestServe_WebrootUnavailableWhenMinioConfigured(t *testing.T) {
	tmp := t.TempDir()
	minioCfg := &files.MinioProxyConfig{Bucket: "bucket"}
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg:  minioCfg,
		minioErr:  handlers.ErrObjectStoreUnavailable,
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	req.Host = "globular.io"

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rr.Code)
	}
}

func TestServe_DirectModeStillRequiresMinioWhenConfigured(t *testing.T) {
	tmp := t.TempDir()
	// place a file locally to detect unintended disk fallback
	if err := os.WriteFile(filepath.Join(tmp, "index.html"), []byte("LOCAL"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg:  &files.MinioProxyConfig{Bucket: "bucket"},
		minioErr:  handlers.ErrObjectStoreUnavailable,
		mode:      "direct",
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	req.Host = "globular.io"

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	if body := rr.Body.String(); strings.Contains(body, "LOCAL") {
		t.Fatalf("should not fall back to disk when minio configured; body=%s", body)
	}
}

func TestServe_MinioErrorWithoutConfigStill503(t *testing.T) {
	tmp := t.TempDir()
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg:  nil,
		minioErr:  handlers.ErrObjectStoreUnavailable,
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	req.Host = "globular.io"

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rr.Code)
	}
}

func TestServe_MinioWebroot(t *testing.T) {
	tmp := t.TempDir()
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg: &files.MinioProxyConfig{
			Domain:        "example.com",
			WebrootPrefix: "",
			WebrootBucket: "bucket",
			Fetch: func(_ context.Context, bucket, key string) (io.ReadSeekCloser, files.MinioObjectInfo, error) {
				if bucket != "bucket" {
					t.Fatalf("expected bucket bucket, got %s", bucket)
				}
				if key != "globular.io/webroot/index.html" {
					t.Fatalf("unexpected key %s", key)
				}
				return nopSeekCloser{Reader: strings.NewReader("WEBROOT")}, files.MinioObjectInfo{ModTime: time.Unix(20, 0)}, nil
			},
		},
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	req.Host = "globular.io"

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "WEBROOT" {
		t.Fatalf("unexpected body %q", rr.Body.String())
	}
}

func TestWebrootObjectKey_DefaultPrefixUsesHost(t *testing.T) {
	cfg := &files.MinioProxyConfig{Domain: "localhost"}
	key, err := files.WebrootObjectKeyForTest(cfg, "localhost", "/")
	if err != nil {
		t.Fatalf("webroot key err: %v", err)
	}
	if key != "localhost/webroot/index.html" {
		t.Fatalf("unexpected key %s", key)
	}
}

func TestWebrootObjectKey_ExplicitPrefix(t *testing.T) {
	cfg := &files.MinioProxyConfig{WebrootPrefix: "customprefix"}
	key, err := files.WebrootObjectKeyForTest(cfg, "localhost", "/")
	if err != nil {
		t.Fatalf("webroot key err: %v", err)
	}
	if key != "customprefix/index.html" {
		t.Fatalf("unexpected key %s", key)
	}
}

func TestServe_MinioWebrootRootPath(t *testing.T) {
	tmp := t.TempDir()
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
		minioCfg: &files.MinioProxyConfig{
			Domain:        "example.com",
			WebrootPrefix: "",
			WebrootBucket: "bucket",
			Fetch: func(_ context.Context, bucket, key string) (io.ReadSeekCloser, files.MinioObjectInfo, error) {
				if bucket != "bucket" {
					t.Fatalf("expected bucket bucket, got %s", bucket)
				}
				if key != "globular.io/webroot/index.html" {
					t.Fatalf("unexpected key %s", key)
				}
				return nopSeekCloser{Reader: strings.NewReader("ROOT")}, files.MinioObjectInfo{ModTime: time.Unix(20, 0)}, nil
			},
		},
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "globular.io"

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() == "" {
		t.Fatalf("expected body served from minio")
	}
}

func TestServe_NotFound(t *testing.T) {
	tmp := t.TempDir()
	p := &fakeServe{
		webRoot:   tmp,
		dataRoot:  tmp,
		allowRead: true,
	}

	h := files.NewServeFile(p)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing.txt", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
