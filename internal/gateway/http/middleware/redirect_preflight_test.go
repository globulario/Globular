package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/globulario/Globular/internal/gateway/http/middleware"
)

func TestWithRedirectAndPreflight_OptionsShortCircuit(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // should never run
	})
	h := middleware.WithRedirectAndPreflight(nil, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CORS", "1")
	})(next)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/getConfig", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Header().Get("X-CORS") != "1" {
		t.Fatalf("expected CORS header set by setHeaders")
	}
}
