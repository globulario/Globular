package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	globpkg "github.com/globulario/Globular/internal/globule"
)

func TestHealthObjectStoreStrictUnavailable(t *testing.T) {
	t.Setenv("GLOBULAR_MINIO_ENDPOINT", "127.0.0.1:1")
	t.Setenv("GLOBULAR_MINIO_BUCKET", "bucket")
	t.Setenv("GLOBULAR_MINIO_ACCESS_KEY", "ak")
	t.Setenv("GLOBULAR_MINIO_SECRET_KEY", "sk")

	glob := &globpkg.Globule{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
		Protocol:       "http",
		Domain:         "localhost",
	}
	h := &GatewayHandlers{globule: glob}
	mux := http.NewServeMux()
	h.wireObjectStoreHealth(mux, func(h http.Handler) http.Handler { return h })

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/objectstore?strict=1", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rr.Code)
	}
}
