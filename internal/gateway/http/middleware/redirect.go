package middleware

import (
	"net/http"
)

// Target is the minimal info your redirect code needs.
type Target struct {
	Hostname   string
	Domain     string
	Protocol   string // "http" | "https"
	PortHTTP   int
	PortHTTPS  int
	LocalIP    string
	ExternalIP string

	// Optional: keep the original object if you want to call legacy code
	Raw any
}

// Redirector lets you plug in your existing redirectTo + handleRequestAndRedirect.
type Redirector interface {
	RedirectTo(host string) (bool, *Target)
	HandleRedirect(to *Target, w http.ResponseWriter, r *http.Request)
}

// WithRedirectAndPreflight handles:
//   - CORS preflight (OPTIONS)
//   - host redirection (local/external)
//
// Then calls next.
func WithRedirectAndPreflight(rdr Redirector, setHeaders func(http.ResponseWriter, *http.Request)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if setHeaders != nil {
				setHeaders(w, r) // your legacy setupResponse(&w, r) equivalent
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			if rdr != nil {
				if ok, to := rdr.RedirectTo(r.Host); ok && to != nil {
					rdr.HandleRedirect(to, w, r)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
