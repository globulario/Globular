package http

// Config defines the configuration settings for HTTP server middleware,
// including CORS policies and rate limiting parameters.
//
// AllowedOrigins specifies the list of origins permitted to access the server.
// AllowedMethods specifies the HTTP methods allowed for cross-origin requests.
// AllowedHeaders specifies the HTTP headers allowed in cross-origin requests.
//
// RateRPS sets the maximum number of requests per second allowed.
// RateBurst sets the maximum burst size for rate limiting.
type Config struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string

	RateRPS   float64
	RateBurst int
}
