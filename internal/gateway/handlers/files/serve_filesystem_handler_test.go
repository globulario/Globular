package files

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type fsStubProvider struct {
	existing      map[string]bool
	maybeStreamFn func(name string, w http.ResponseWriter, r *http.Request) bool
}

func (s *fsStubProvider) WebRoot() string                       { return "" }
func (s *fsStubProvider) DataRoot() string                      { return "" }
func (s *fsStubProvider) CredsDir() string                      { return "" }
func (s *fsStubProvider) IndexApplication() string              { return "" }
func (s *fsStubProvider) PublicDirs() []string                  { return nil }
func (s *fsStubProvider) Exists(p string) bool                  { return s.existing[p] }
func (s *fsStubProvider) FindHashedFile(string) (string, error) { return "", nil }
func (s *fsStubProvider) FileServiceMinioConfig() (*MinioProxyConfig, error) {
	return nil, nil
}
func (s *fsStubProvider) FileServiceMinioConfigStrict(ctx context.Context) (*MinioProxyConfig, error) {
	return nil, nil
}
func (s *fsStubProvider) Mode() string { return "" }
func (s *fsStubProvider) ParseUserID(token string) (string, error) {
	return "", nil
}
func (s *fsStubProvider) ValidateAccount(userID, action, reqPath string) (bool, bool, error) {
	return false, false, nil
}
func (s *fsStubProvider) ValidateApplication(app, action, reqPath string) (bool, bool, error) {
	return false, false, nil
}
func (s *fsStubProvider) ResolveImportPath(basePath, importLine string) (string, error) {
	return "", nil
}
func (s *fsStubProvider) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	if s.maybeStreamFn != nil {
		return s.maybeStreamFn(name, w, r)
	}
	return false
}
func (s *fsStubProvider) ResolveProxy(reqPath string) (string, bool) {
	return "", false
}

func TestFilesystemHandler_ServeSuccess(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "servefs-*.txt")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, _ = tmpFile.WriteString("hello")
	_ = tmpFile.Close()

	p := &pathInfo{name: tmpFile.Name(), marker: "marker"}
	provider := &fsStubProvider{existing: map[string]bool{tmpFile.Name(): true}}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/file.txt", nil)

	handler := &FilesystemHandler{Provider: provider}
	if !handler.Serve(rr, req, p) {
		t.Fatalf("Serve returned false")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if body := rr.Body.String(); body != "hello" {
		t.Fatalf("expected body 'hello', got %q", body)
	}
}

func TestFilesystemHandler_NotFound(t *testing.T) {
	p := &pathInfo{name: "/does/not/exist"}
	provider := &fsStubProvider{existing: map[string]bool{}}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)

	handler := &FilesystemHandler{Provider: provider}
	if !handler.Serve(rr, req, p) {
		t.Fatalf("Serve returned false")
	}
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestFilesystemHandler_StreamShortCircuit(t *testing.T) {
	provider := &fsStubProvider{
		existing: map[string]bool{},
		maybeStreamFn: func(name string, w http.ResponseWriter, r *http.Request) bool {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("stream"))
			return true
		},
	}
	p := &pathInfo{name: "/videos/movie.mkv"}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/videos/movie.mkv", nil)

	handler := &FilesystemHandler{Provider: provider}
	if !handler.Serve(rr, req, p) {
		t.Fatalf("Serve returned false")
	}
	if rr.Body.String() != "stream" {
		t.Fatalf("expected stream body, got %q", rr.Body.String())
	}
}
