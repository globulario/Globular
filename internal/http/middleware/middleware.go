// Package middleware provides HTTP middleware functions for common tasks such as
// security headers, CORS, logging, recovery, and rate limiting.
package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Chain func(http.Handler) http.Handler

func (c Chain) Then(h http.Handler) http.Handler { return c(h) }

// Compose creates a middleware chain from a list of middleware functions.
// The middleware are applied in the order they are provided, such that the first
// middleware wraps the subsequent ones. It returns a Chain, which is a function
// that takes an http.Handler and returns a new http.Handler with all middleware applied.
//
// Example usage:
//
//	handler := Compose(middleware1, middleware2, middleware3)(finalHandler)
func Compose(mw ...func(http.Handler) http.Handler) Chain {
	return func(h http.Handler) http.Handler {
		for i := len(mw) - 1; i >= 0; i-- {
			h = mw[i](h)
		}
		return h
	}
}

// Recoverer returns a middleware that recovers from panics in HTTP handlers.
// If a panic occurs, it logs the error using the provided slog.Logger and responds
// with a 500 Internal Server Error. The logger receives the request path and the
// recovered error value for diagnostic purposes.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic", "path", r.URL.Path, "err", rec)
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders is a middleware that sets common security-related HTTP headers
// on the response to enhance protection against common web vulnerabilities.
// It sets the following headers:
//   - X-Content-Type-Options: nosniff (prevents MIME type sniffing)
//   - Referrer-Policy: no-referrer (prevents referrer information from being sent)
//   - X-Frame-Options: DENY (prevents the page from being displayed in a frame)
//   - Content-Security-Policy: restricts sources for content, images, and objects
//
// Use this middleware to help secure your HTTP handlers.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data: blob:; object-src 'none'")
		next.ServeHTTP(w, r)
	})
}

// CORS returns a middleware that sets CORS (Cross-Origin Resource Sharing) headers
// for HTTP requests. It allows you to specify allowed origins, methods, and headers.
//
// Parameters:
//
//	origins - a slice of allowed origin strings. If the first element is "*", all origins are allowed.
//	methods - a slice of allowed HTTP methods (e.g., "GET", "POST").
//	headers - a slice of allowed HTTP headers.
//
// The middleware checks the "Origin" header of incoming requests and, if allowed,
// sets the appropriate CORS headers. For OPTIONS requests, it responds with 204 No Content.
//
// Example usage:
//
//	handler := CORS([]string{"https://example.com"}, []string{"GET", "POST"}, []string{"Content-Type"})(yourHandler)
func CORS(origins, methods, headers []string) func(http.Handler) http.Handler {
	join := func(xs []string) string { return strings.Join(xs, ",") }
	allowOrigin := func(origin string) bool {
		if len(origins) == 0 {
			return false
		}
		if origins[0] == "*" {
			return true
		}
		for _, o := range origins {
			if o == origin {
				return true
			}
		}
		return false
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && allowOrigin(origin) {
				h := w.Header()
				h.Set("Access-Control-Allow-Origin", origin)
				h.Set("Access-Control-Allow-Credentials", "true")
				if len(methods) > 0 {
					h.Set("Access-Control-Allow-Methods", join(methods))
				}
				if len(headers) > 0 {
					h.Set("Access-Control-Allow-Headers", join(headers))
				}
				h.Set("Access-Control-Max-Age", "3600")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Logger returns a middleware that logs HTTP requests using the provided slog.Logger.
// It logs the request method, path, remote client IP, and the duration taken to process the request.
// Usage: Wrap your HTTP handler with Logger to enable request logging.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("req", "method", r.Method, "path", r.URL.Path, "remote", clientIP(r), "dur", time.Since(start).String())
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// LimiterStore manages rate limiters for multiple clients identified by a key (e.g., IP address).
// It stores the rate limit configuration (limit and burst) and a map of rate limiters per client.
// This allows for efficient tracking and enforcement of rate limits across different clients.
type LimiterStore struct {
	limit  rate.Limit
	burst  int
	bucket map[string]*rate.Limiter
}

// NewLimiterStore creates and returns a new LimiterStore instance configured with the specified
// requests per second (rps) rate and burst size. The LimiterStore maintains a map of rate limiters
// for different keys, allowing for per-key rate limiting.
// Parameters:
//
//	rps   - The allowed requests per second.
//	burst - The maximum burst size allowed.
//
// Returns:
//
//	A pointer to the initialized LimiterStore.
func NewLimiterStore(rps float64, burst int) *LimiterStore {
	return &LimiterStore{limit: rate.Limit(rps), burst: burst, bucket: make(map[string]*rate.Limiter)}
}

// get retrieves the rate limiter for the given IP address from the LimiterStore.
// If no limiter exists for the IP, a new one is created and added to the store.
func (s *LimiterStore) get(ip string) *rate.Limiter {
	if l, ok := s.bucket[ip]; ok {
		return l
	}
	l := rate.NewLimiter(s.limit, s.burst)
	s.bucket[ip] = l
	return l
}

// RateLimiter returns a middleware that limits the rate of incoming HTTP requests per client IP.
// It uses the provided LimiterStore to track and enforce rate limits.
// If a client exceeds the allowed rate, the middleware responds with HTTP 429 (Too Many Requests).
// Otherwise, it passes the request to the next handler.
//
// store: The LimiterStore instance used to manage rate limiting per client.
// Returns: A middleware function to be used with HTTP handlers.
func RateLimiter(store *LimiterStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !store.get(clientIP(r)).Allow() {
				http.Error(w, "rate limit", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
