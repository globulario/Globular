package files

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMinIOHandler_CanServeRequiresConfig(t *testing.T) {
	h := &MinIOHandler{useWeb: true, cfg: nil}
	if h.CanServe(&pathInfo{}) {
		t.Fatalf("expected CanServe=false when cfg nil")
	}
	h.cfg = &MinioProxyConfig{}
	if !h.CanServe(&pathInfo{}) {
		t.Fatalf("expected CanServe=true when cfg present and useWeb set")
	}
}

func TestMinIOHandler_ServeWebSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)

	h := &MinIOHandler{
		cfg:    &MinioProxyConfig{},
		useWeb: true,
		serveWeb: func(w http.ResponseWriter, r *http.Request, cfg *MinioProxyConfig, rq, host string, fallback bool) bool {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return true
		},
	}

	served := h.Serve(rr, req, &pathInfo{reqPath: "/index.html"})
	if !served {
		t.Fatalf("expected Serve to return true")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if body := rr.Body.String(); body != "ok" {
		t.Fatalf("expected body 'ok', got %q", body)
	}
}

func TestMinIOHandler_ServeUsersFallsThrough(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/file.txt", nil)

	h := &MinIOHandler{
		cfg:      &MinioProxyConfig{},
		useUsers: true,
		serveUsers: func(w http.ResponseWriter, r *http.Request, cfg *MinioProxyConfig, rq string) bool {
			return false // simulate missing object -> fall through
		},
	}

	served := h.Serve(rr, req, &pathInfo{reqPath: "/users/alice/file.txt"})
	if served {
		t.Fatalf("expected Serve to return false when MinIO not handling")
	}
	if rr.Body.Len() != 0 {
		t.Fatalf("expected no body written, got %q", rr.Body.String())
	}
}
