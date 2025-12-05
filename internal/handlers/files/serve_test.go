package files_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	files "github.com/globulario/Globular/internal/handlers/files"
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
func (f fakeServe) FileServiceMinioConfig() (*files.MinioProxyConfig, bool) {
	if f.minioCfg == nil {
		return nil, false
	}
	return f.minioCfg, true
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
