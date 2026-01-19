package files_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	files "github.com/globulario/Globular/internal/gateway/handlers/files"
	"github.com/globulario/services/golang/config"
)

type proxyServe struct{ target string }

func (p proxyServe) WebRoot() string                       { return "" }
func (p proxyServe) DataRoot() string                      { return "" }
func (p proxyServe) CredsDir() string                      { return config.GetConfigDir() + "/tls" }
func (p proxyServe) IndexApplication() string              { return "" }
func (p proxyServe) PublicDirs() []string                  { return nil }
func (p proxyServe) Exists(string) bool                    { return false }
func (p proxyServe) FindHashedFile(string) (string, error) { return "", nil }
func (p proxyServe) FileServiceMinioConfig() (*files.MinioProxyConfig, error) {
	return nil, nil
}
func (p proxyServe) FileServiceMinioConfigStrict(ctx context.Context) (*files.MinioProxyConfig, error) {
	return nil, nil
}
func (p proxyServe) Mode() string                       { return "direct" }
func (p proxyServe) ParseUserID(string) (string, error) { return "", nil }
func (p proxyServe) ValidateAccount(string, string, string) (bool, bool, error) {
	return false, false, nil
}
func (p proxyServe) ValidateApplication(string, string, string) (bool, bool, error) {
	return false, false, nil
}
func (p proxyServe) ResolveImportPath(string, string) (string, error)            { return "", nil }
func (p proxyServe) MaybeStream(string, http.ResponseWriter, *http.Request) bool { return false }
func (p proxyServe) ResolveProxy(reqPath string) (string, bool) {
	// Proxy when path starts with /proxy
	if len(reqPath) >= 6 && reqPath[:6] == "/proxy" {
		return p.target + "/proxied", true
	}
	return "", false
}

func TestServe_ReverseProxy_ShortCircuits(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("listen not permitted: %v", err)
	}
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxied" {
			t.Fatalf("expected /proxied, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("PROXIED"))
	}))
	backend.Listener = ln
	backend.Start()
	t.Cleanup(backend.Close)

	p := proxyServe{target: backend.URL}
	h := files.NewServeFile(p)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/proxy/anything", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Body.String(); got != "PROXIED" {
		t.Fatalf("expected backend body, got %q", got)
	}
}
