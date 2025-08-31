package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/globulario/Globular/internal/http/middleware"
)

type fakeRedirector struct {
	redirect bool
}

func (f *fakeRedirector) RedirectTo(host string) (bool, *middleware.Target) {
	if f.redirect {
		return true, &middleware.Target{
			Domain:     "peer.example",
			Protocol:   "http",
			PortHTTP:   8081,
			LocalIP:    "10.0.0.2",
			ExternalIP: "203.0.113.2",
		}
	}
	return false, nil
}

func (f *fakeRedirector) HandleRedirect(to *middleware.Target, w http.ResponseWriter, r *http.Request) {
	// In real code you'd proxy; for the test just signal we redirected.
	w.WriteHeader(http.StatusNoContent) // 204
}

func TestWithRedirectAndPreflight_Redirects(t *testing.T) {
	// Next handler would return 418 if it ever ran (it shouldn't).
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	rdr := &fakeRedirector{redirect: true}

	h := middleware.WithRedirectAndPreflight(rdr, nil)(next)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "api/getConfig", nil)
	req.Host = "peer.example:8081"

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 from redirect, got %d", rr.Code)
	}
}
