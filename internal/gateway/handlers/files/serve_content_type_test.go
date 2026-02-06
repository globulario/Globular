package files

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mockServeProviderForContent implements ServeProvider for content type tests
// Only methods used by serveWithContentType are implemented meaningfully.
type mockServeProviderForContent struct {
	resolveImportFunc func(basePath, importLine string) (string, error)
}

func (m *mockServeProviderForContent) WebRoot() string                         { return "/webroot" }
func (m *mockServeProviderForContent) DataRoot() string                        { return "/data" }
func (m *mockServeProviderForContent) CredsDir() string                        { return "/creds" }
func (m *mockServeProviderForContent) PublicDirs() []string                    { return []string{"/public"} }
func (m *mockServeProviderForContent) IndexApplication() string                { return "" }
func (m *mockServeProviderForContent) Exists(p string) bool                    { return false }
func (m *mockServeProviderForContent) FindHashedFile(p string) (string, error) { return "", nil }
func (m *mockServeProviderForContent) FileServiceMinioConfig() (*MinioProxyConfig, error) {
	return nil, nil
}
func (m *mockServeProviderForContent) FileServiceMinioConfigStrict(ctx context.Context) (*MinioProxyConfig, error) {
	return nil, nil
}
func (m *mockServeProviderForContent) Mode() string { return "" }
func (m *mockServeProviderForContent) ParseUserID(token string) (string, error) {
	return "", nil
}
func (m *mockServeProviderForContent) ValidateAccount(userID, action, reqPath string) (bool, bool, error) {
	return false, false, nil
}
func (m *mockServeProviderForContent) ValidateApplication(app, action, reqPath string) (bool, bool, error) {
	return false, false, nil
}
func (m *mockServeProviderForContent) ResolveImportPath(basePath, importLine string) (string, error) {
	if m.resolveImportFunc != nil {
		return m.resolveImportFunc(basePath, importLine)
	}
	return "", nil
}
func (m *mockServeProviderForContent) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	return false
}
func (m *mockServeProviderForContent) ResolveProxy(reqPath string) (string, bool) {
	return "", false
}

// TestServeWithContentType_JavaScript tests JavaScript file serving
func TestServeWithContentType_JavaScript(t *testing.T) {
	tmpDir := t.TempDir()
	jsFile := filepath.Join(tmpDir, "test.js")
	if err := os.WriteFile(jsFile, []byte("console.log('test');"), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(jsFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test.js", nil)
	p := &mockServeProviderForContent{}

	transformed := serveWithContentType(w, r, jsFile, f, fi, p, false, "marker")
	if transformed {
		t.Error("expected transformed=false for simple JS file")
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/javascript" {
		t.Errorf("expected Content-Type=application/javascript, got %s", ct)
	}
}

// TestServeWithContentType_CSS tests CSS file serving
func TestServeWithContentType_CSS(t *testing.T) {
	tmpDir := t.TempDir()
	cssFile := filepath.Join(tmpDir, "test.css")
	if err := os.WriteFile(cssFile, []byte("body { margin: 0; }"), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(cssFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test.css", nil)
	p := &mockServeProviderForContent{}

	transformed := serveWithContentType(w, r, cssFile, f, fi, p, false, "marker")
	if transformed {
		t.Error("expected transformed=false for CSS file")
	}

	if ct := w.Header().Get("Content-Type"); ct != "text/css" {
		t.Errorf("expected Content-Type=text/css, got %s", ct)
	}
}

// TestServeWithContentType_HTMLRoot tests HTML serving with marker injection on root
func TestServeWithContentType_HTMLRoot(t *testing.T) {
	tmpDir := t.TempDir()
	htmlFile := filepath.Join(tmpDir, "index.html")
	originalHTML := "<html><body>Test</body></html>"
	if err := os.WriteFile(htmlFile, []byte(originalHTML), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(htmlFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/index.html", nil)
	p := &mockServeProviderForContent{}
	marker := "Served-By: test"

	transformed := serveWithContentType(w, r, htmlFile, f, fi, p, true, marker)
	if !transformed {
		t.Error("expected transformed=true for root HTML file")
	}

	body := w.Body.String()
	if !strings.Contains(body, marker) {
		t.Errorf("expected marker %q in HTML, got: %s", marker, body)
	}

	if !strings.Contains(body, originalHTML) {
		t.Errorf("expected original HTML content, got: %s", body)
	}
}

// TestServeWithContentType_HTMLNonRoot tests HTML serving without marker on non-root
func TestServeWithContentType_HTMLNonRoot(t *testing.T) {
	tmpDir := t.TempDir()
	htmlFile := filepath.Join(tmpDir, "page.html")
	if err := os.WriteFile(htmlFile, []byte("<html><body>Page</body></html>"), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(htmlFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/page.html", nil)
	p := &mockServeProviderForContent{}

	transformed := serveWithContentType(w, r, htmlFile, f, fi, p, false, "marker")
	if transformed {
		t.Error("expected transformed=false for non-root HTML file")
	}

	if ct := w.Header().Get("Content-Type"); ct != "text/html" {
		t.Errorf("expected Content-Type=text/html, got %s", ct)
	}
}

// TestServeWithContentType_OtherFiles tests files without special handling
func TestServeWithContentType_OtherFiles(t *testing.T) {
	tmpDir := t.TempDir()
	txtFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(txtFile, []byte("plain text"), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(txtFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test.txt", nil)
	p := &mockServeProviderForContent{}

	transformed := serveWithContentType(w, r, txtFile, f, fi, p, false, "marker")
	if transformed {
		t.Error("expected transformed=false for .txt file")
	}

	if ct := w.Header().Get("Content-Type"); ct != "" {
		t.Errorf("expected no Content-Type for .txt, got %s", ct)
	}
}

// TestServeWithContentType_JavaScript_Rewritten verifies changed=true path when imports rewritten
func TestServeWithContentType_JavaScript_Rewritten(t *testing.T) {
	tmpDir := t.TempDir()
	jsFile := filepath.Join(tmpDir, "rewrite.js")
	source := "import x from '@foo';\nconsole.log(x);"
	if err := os.WriteFile(jsFile, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(jsFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/rewrite.js", nil)

	p := &mockServeProviderForContent{}
	p.resolveImportFunc = func(basePath, importLine string) (string, error) {
		return "/resolved/path.js", nil
	}

	transformed := serveWithContentType(w, r, jsFile, f, fi, p, false, "marker")
	if !transformed {
		t.Fatal("expected transformed=true for rewritten JS imports")
	}

	body := w.Body.String()
	if !strings.Contains(body, "/resolved/path.js") {
		t.Fatalf("expected rewritten import path, got %s", body)
	}
}
