package files_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	files "github.com/globulario/Globular/internal/handlers/files"
)

type cacheServe struct {
	webRoot string
}

func (c cacheServe) WebRoot() string                       { return c.webRoot }
func (c cacheServe) DataRoot() string                      { return c.webRoot }
func (c cacheServe) CredsDir() string                      { return c.webRoot }
func (c cacheServe) IndexApplication() string              { return "" }
func (c cacheServe) PublicDirs() []string                  { return nil }
func (c cacheServe) Exists(p string) bool                  { _, err := os.Stat(p); return err == nil }
func (c cacheServe) FindHashedFile(string) (string, error) { return "", nil }
func (c cacheServe) ParseUserID(string) (string, error)    { return "", nil }
func (c cacheServe) ValidateAccount(string, string, string) (bool, bool, error) {
	return true, false, nil
}
func (c cacheServe) ValidateApplication(string, string, string) (bool, bool, error) {
	return false, false, nil
}
func (c cacheServe) ResolveImportPath(string, string) (string, error)            { return "", nil }
func (c cacheServe) MaybeStream(string, http.ResponseWriter, *http.Request) bool { return false }
func (c cacheServe) ResolveProxy(string) (string, bool)                          { return "", false }

func TestServe_ETag_304(t *testing.T) {
	tmp := t.TempDir()
	fname := filepath.Join(tmp, "foo.txt")
	if err := os.WriteFile(fname, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	// ensure stable mtime
	_ = os.Chtimes(fname, time.Now().Add(-2*time.Hour), time.Now().Add(-2*time.Hour))

	h := files.NewServeFile(cacheServe{webRoot: tmp})

	// First request to get the ETag
	rr1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/foo.txt", nil)
	req1.Host = "localhost"
	h.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr1.Code)
	}
	etag := rr1.Header().Get("ETag")
	if etag == "" {
		t.Fatalf("missing ETag")
	}

	// Second request with If-None-Match
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/foo.txt", nil)
	req2.Header.Set("If-None-Match", etag)
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusNotModified {
		t.Fatalf("expected 304, got %d", rr2.Code)
	}
}
