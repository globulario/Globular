package main

import (
	"bufio"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/httpcache"
	"github.com/StalkR/imdb"
	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
	colly "github.com/gocolly/colly/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	protoparser "github.com/yoheimuta/go-protoparser/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"

const cacheTTL = 24 * time.Hour

// client is used by tests to perform cached requests.
// If cache directory exists it is used as a persistent cache.
// Otherwise a volatile memory cache is used.
var client *http.Client

func init() {
	if _, err := os.Stat("cache"); err == nil {
		client, err = httpcache.NewPersistentClient("cache", cacheTTL)
		if err != nil {
			panic(err)
		}
	} else {
		client = httpcache.NewVolatileClient(cacheTTL, 1024)
	}
	client.Transport = &customTransport{client.Transport}
}

// customTransport implements http.RoundTripper interface to add some headers.
type customTransport struct {
	http.RoundTripper
}

// googleOauthConfig is the OAuth2 configuration for Google.
var googleOauthConfig *oauth2.Config

func getGoogleOauthConfig() *oauth2.Config {
	if googleOauthConfig == nil {
		googleOauthConfig = &oauth2.Config{
			ClientID:     globule.OAuth2ClientID,
			ClientSecret: globule.OAuth2ClientSecret,
			RedirectURL:  "postmessage",
			Scopes:       []string{"openid", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	}
	return googleOauthConfig
}

func exchangeAuthCodeForToken(code string) (*oauth2.Token, error) {
	conf := getGoogleOauthConfig()
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Token Exchange Error:", err) // Print the exact error
		return nil, err
	}

	return token, nil
}

// handleTokenRefresh handles the token refresh request.
func handleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	// Enable CORS if needed
	setupResponse(&w, r)

	// Parse the request body
	refreshToken := r.URL.Query().Get("refresh_token")

	// Get OAuth2 Config
	conf := getGoogleOauthConfig()

	// Create a new token object with the refresh token
	token := &oauth2.Token{RefreshToken: refreshToken}

	// Create a token source
	tokenSource := conf.TokenSource(context.Background(), token)

	// Get a new token
	newToken, err := tokenSource.Token()
	if err != nil {
		http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Send the new tokens to the frontend
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  newToken.AccessToken,
		"refresh_token": newToken.RefreshToken, // Only returned if Google issues a new one
		"expires_in":    newToken.Expiry.Unix(),
	})

	if err != nil {
		http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}
}

// Handles the OAuth2 callback from Google
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {

	// Handle the preflight options...
	setupResponse(&w, r)

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	token, err := exchangeAuthCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	// 🔹 Validate the ID Token
	idToken := token.Extra("id_token").(string)
	payload, err := verifyGoogleIDToken(idToken, getGoogleOauthConfig())
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	// Send the token's and user's info as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"id_token":      idToken,
		"expires_in":    token.Expiry,
		"user":          payload, // User profile info
	})

	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// JWK represents a JSON Web Key as per RFC 7517.
type JWK struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// GoogleKeys represents a set of JSON Web Keys.
type GoogleKeys struct {
	Keys []JWK `json:"keys"`
}

// TokenClaims represents the expected claims in the Google ID token.
type TokenClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Exp           int64  `json:"exp"`
	Sub           string `json:"sub"` // User's Google ID
	Aud           string `json:"aud"` // Audience (your client ID)
	Iss           string `json:"iss"` // Issuer
}

// Google's public keys URL
const googleJWTKeySetURL = "https://www.googleapis.com/oauth2/v3/certs"
const maxJWTLen = 16 * 1024 // 16KB is plenty for an ID token

// verifyGoogleIDToken verifies the Google ID token and returns the claims.
func verifyGoogleIDToken(idToken string, config *oauth2.Config) (map[string]interface{}, error) {

	if len(idToken) == 0 || len(idToken) > maxJWTLen {
		return nil, fmt.Errorf("invalid token size")
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}), // Google uses RS256
		// jwt.WithLeeway(30*time.Second), // optional clock skew
	)

	// Fetch Google's public keys
	resp, err := http.Get(googleJWTKeySetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Google's public keys: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("warning: failed to close response body: %v\n", err)
		}
	}()

	// Read and decode JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Google's keys response: %w", err)
	}

	var keySet GoogleKeys
	if err := json.Unmarshal(body, &keySet); err != nil {
		return nil, fmt.Errorf("failed to decode Google's keys: %w", err)
	}

	token, err := parser.Parse(idToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		for _, key := range keySet.Keys {
			if key.Kid == t.Header["kid"] {
				return convertToPublicKey(key)
			}
		}
		return nil, errors.New("key not found")
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Convert claims into TokenClaims struct
	var tokenClaims TokenClaims
	claimsJSON, _ := json.Marshal(claims) // Convert map to JSON
	err = json.Unmarshal(claimsJSON, &tokenClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token claims: %w", err)
	}

	// Validate token claims
	if tokenClaims.Aud != config.ClientID {
		return nil, errors.New("invalid audience")
	}
	if tokenClaims.Exp < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	if tokenClaims.Iss != "accounts.google.com" && tokenClaims.Iss != "https://accounts.google.com" {
		return nil, errors.New("invalid issuer")
	}

	// Return verified claims
	return claims, nil
}

// convertToPublicKey converts Google's modulus and exponent to an RSA public key.
func convertToPublicKey(key JWK) (*rsa.PublicKey, error) {
	nBytes, err := decodeBase64URL(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}
	eBytes, err := decodeBase64URL(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert exponent to int
	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	// Construct the RSA public key
	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}
	return pubKey, nil
}

// decodeBase64URL decodes base64 URL-encoded strings.
func decodeBase64URL(s string) ([]byte, error) {
	return jwt.NewParser().DecodeSegment(s)
}

func (e *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Accept-Language", "en") // avoid IP-based language detection
	r.Header.Set("User-Agent", userAgent)
	return e.RoundTripper.RoundTrip(r)
}

// Find the peer with a given name and redirect the
// the request to it.
func redirectTo(host string) (bool, *resourcepb.Peer) {

	if globule == nil {
		globule = NewGlobule()
	}

	if globule.Peers == nil {
		return false, nil
	}

	// read the actual configuration.
	localAddress, err := config_.GetAddress()
	if err == nil {
		// no redirection if the address is the same...
		if strings.HasPrefix(localAddress, host) {
			return false, nil
		}
	}

	var peer *resourcepb.Peer

	globule.peers.Range(func(_, value interface{}) bool {
		p := value.(*resourcepb.Peer)
		address := p.Hostname + "." + p.Domain
		if p.Protocol == "https" {
			address += ":" + Utility.ToString(p.PortHttps)
		} else {
			address += ":" + Utility.ToString(p.PortHttp)
		}

		if strings.HasPrefix(address, host) {
			peer = p
			return false // stop the iteration.
		}
		return true
	})

	return peer != nil, peer
}

// Redirect the query to a peer one the network
func handleRequestAndRedirect(to *resourcepb.Peer, res http.ResponseWriter, req *http.Request) {

	address := to.Domain
	scheme := "http"
	if to.Protocol == "https" {
		address += ":" + Utility.ToString(to.PortHttps)
	} else {
		address += ":" + Utility.ToString(to.PortHttp)
	}

	// Here I will remove the .localhost part of the address (if it exist)
	address = strings.ReplaceAll(address, ".localhost", "")
	ur, _ := url.Parse(scheme + "://" + address)

	proxy := httputil.NewSingleHostReverseProxy(ur)

	// Update the headers to allow for SSL redirection
	req.URL.Host = ur.Host
	req.URL.Scheme = ur.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	proxy.ErrorHandler = ErrHandle

	proxy.ServeHTTP(res, req)
}

// ErrHandle logs the provided error to the standard output.
// It is intended to be used as an HTTP error handler.
// Parameters:
//   - res: the HTTP response writer.
//   - _: the HTTP request (ignored).
//   - err: the error to be handled and logged.
func ErrHandle(res http.ResponseWriter, _ *http.Request, err error) {
	fmt.Fprintf(os.Stderr, "proxy error: %v\n", err)
	http.Error(res, "proxy error: "+err.Error(), http.StatusBadGateway)
}

/**
 * Create a checksum from a given path.
 */
func getChecksumHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)

	execPath := Utility.GetExecName(os.Args[0])
	if Utility.Exists("/usr/local/share/globular/Globular") {
		execPath = "/usr/local/share/globular/Globular"
	}

	_, err := fmt.Fprint(w, Utility.CreateFileChecksum(execPath))
	if err != nil {
		http.Error(w, "fail to write checksum with error "+err.Error(), http.StatusBadRequest)
		return
	}

}

/**
 * Save the configuration.
 */
func saveConfigHanldler(w http.ResponseWriter, r *http.Request) {

	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// I will try to get the token from the header.
	token := r.Header.Get("token")

	// I will validate the token.
	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	// If not token was given i will return an error (403).
	if len(token) == 0 {
		http.Error(w, "no token was given!", http.StatusUnauthorized)
		return
	}

	// I will validate the token.
	_, err := security.ValidateToken(token)
	if err != nil {
		http.Error(w, "fail to validate token with error "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Now I will get the configuration from the request.
	decoder := json.NewDecoder(r.Body)
	var config map[string]interface{}
	err = decoder.Decode(&config)
	if err != nil {
		http.Error(w, "fail to decode configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// I will set the globular configuration.
	err = globule.setConfig(config)
	if err != nil {
		http.Error(w, "fail to set configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

}

// allow only hex + underscore after replacing ':' with '_' (MAC-safe, no slashes/dots)
var macFileRe = regexp.MustCompile(`^[A-Fa-f0-9_]{1,32}$`)

func getPublicKeyHanldler(w http.ResponseWriter, r *http.Request) { // keep name if it's referenced elsewhere
	// Redirect if needed
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// CORS / preflight
	setupResponse(&w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Build an anchored, validated path: <config>/keys/<MAC>_public
	base := filepath.Join(config_.GetConfigDir(), "keys")
	mac := strings.ReplaceAll(globule.Mac, ":", "_")
	if !macFileRe.MatchString(mac) {
		http.Error(w, "invalid mac", http.StatusBadRequest)
		return
	}
	filename := mac + "_public"

	// Resolve base and full paths (symlink-safe)
	baseResolved, err := filepath.EvalSymlinks(base)
	if err != nil {
		http.Error(w, "server configuration error", http.StatusInternalServerError)
		return
	}
	full := filepath.Join(baseResolved, filename)
	fullResolved, err := filepath.EvalSymlinks(full)
	if err != nil {
		http.Error(w, "public key not found", http.StatusNotFound)
		return
	}

	// Ensure the resolved file is still under base (no traversal/symlink escape)
	rel, err := filepath.Rel(baseResolved, fullResolved)
	if err != nil || strings.HasPrefix(rel, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Ensure it's a regular file
	fi, err := os.Stat(fullResolved)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "public key not found", http.StatusNotFound)
		} else {
			http.Error(w, "cannot stat file", http.StatusInternalServerError)
		}
		return
	}
	if !fi.Mode().IsRegular() {
		http.Error(w, "invalid file type", http.StatusBadRequest)
		return
	}

	// Open and stream it
	f, err := os.Open(fullResolved) // #nosec G304 -- Path is anchored, allow-listed, symlinks resolved, and confined to base.
	if err != nil {
		http.Error(w, "cannot open file", http.StatusNotFound)
		return
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "fail to read public key", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func getCertificateHanldler(w http.ResponseWriter, _ *http.Request) {

	// ... [existing code] ...
	address, err := config_.GetAddress()
	if err != nil {
		http.Error(w, "fail to get address with error "+err.Error(), http.StatusBadRequest)
		return
	}

	domain := strings.Split(address, ":")[0]
	certFilename := config_.GetLocalCertificate()
	path := config_.GetConfigDir() + "/tls/" + domain + "/" + config_.GetLocalCertificate()

	if !Utility.Exists(path) {
		http.Error(w, "no issuer certificate found at path "+path, http.StatusBadRequest)
		return
	}

	// #nosec G304 -- File path is constructed from validated input and constant strings.
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "fail to read public key with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set the headers to suggest a download file name and indicate the file type.
	w.Header().Set("Content-Disposition", "attachment; filename=\""+certFilename+"\"")
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "fail to write public key with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

func getIssuerCertificateHandler(w http.ResponseWriter, _ *http.Request) {

	// ... [existing code] ...
	address, err := config_.GetAddress()
	if err != nil {
		http.Error(w, "fail to get address with error "+err.Error(), http.StatusBadRequest)
		return
	}

	domain := strings.Split(address, ":")[0]
	certFilename := config_.GetLocalCertificateAuthorityBundle()
	path := config_.GetConfigDir() + "/tls/" + domain + "/" + config_.GetLocalCertificateAuthorityBundle()

	if !Utility.Exists(path) {
		http.Error(w, "no issuer certificate found at path "+path, http.StatusBadRequest)
		return
	}

	// #nosec G304 -- File path is constructed from validated input and constant strings.
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "fail to read public key with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set the headers to suggest a download file name and indicate the file type.
	w.Header().Set("Content-Disposition", "attachment; filename=\""+certFilename+"\"")
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "fail to write public key with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

/**
 * Return services permissions configuration to be able to manage resources access from rpc request.
 */
func getServicePermissionsHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// so here I will read the service configuration from the service id given in the query
	serviceID := r.URL.Query().Get("id") // the csr in base64

	//add prefix and clean
	serviceConfig, err := config_.GetServiceConfigurationById(serviceID)
	if err != nil {
		http.Error(w, "fail to get service configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// from the configuration i will read the configuration file...
	data, err := os.ReadFile(serviceConfig["ConfigPath"].(string))
	if err != nil {
		http.Error(w, "fail to read service configuration file with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// reload the configuration with the permissions...
	err = json.Unmarshal(data, &serviceConfig)
	if err != nil {
		http.Error(w, "fail to get service configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// set empty array if not defined...
	if serviceConfig["Permissions"] == nil {
		serviceConfig["Permissions"] = []interface{}{}
	}

	gotJSON, err := json.MarshalIndent(serviceConfig["Permissions"], "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}

	_, err = w.Write(gotJSON)
	if err != nil {
		http.Error(w, "fail to write service permissions with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

/**
 * This function is use to return a json object containing the description of the service.
 */
func getServiceDescriptorHanldler(w http.ResponseWriter, r *http.Request) {

	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// so here I will read the service configuration from the service id given in the query
	serviceID := r.URL.Query().Get("id") // the csr in base64

	//add prefix and clean
	serviceConfig, err := config_.GetServiceConfigurationById(serviceID)
	if err != nil {
		http.Error(w, "fail to get service configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// from the service configuration I will read it proto file...
	protoFile := serviceConfig["Proto"].(string)

	// I will read the proto file and return it as a json object.
	// #nosec G304 -- File path is constructed from validated input and constant strings.
	reader, err := os.Open(protoFile)

	if err != nil {
		http.Error(w, "fail to open proto file with error "+err.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		err = reader.Close()
		if err != nil {
			http.Error(w, "fail to close proto file with error "+err.Error(), http.StatusBadRequest)
		}
	}()

	// parse the proto file.

	got, err := protoparser.Parse(
		reader,
		protoparser.WithDebug(false),
		protoparser.WithPermissive(false),
		protoparser.WithFilename(filepath.Base(protoFile)),
	)

	var v interface{} = got

	gotJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}

	_, err = w.Write(gotJSON)
	if err != nil {
		http.Error(w, "fail to write service descriptor with error "+err.Error(), http.StatusBadRequest)
		return
	}

}

/**
 * Return the service configuration
 */
func getConfigHanldler(w http.ResponseWriter, r *http.Request) {

	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// i will redirect to the given host if the host is not the same...
	address, _ := config_.GetAddress()

	// I will redirect the request if host is defined in the query...
	if !strings.HasPrefix(address, r.URL.Query().Get("host")) && len(r.URL.Query().Get("host")) > 0 {

		redirect, to := redirectTo(r.URL.Query().Get("host"))

		if redirect && to != nil {

			// I will get the remote configuration and return it...
			var remoteConfig map[string]interface{}
			var err error
			address := to.LocalIpAddress
			if to.ExternalIpAddress != Utility.MyIP() {
				address = to.ExternalIpAddress
			}

			remoteConfig, err = config_.GetRemoteConfig(address, 0)

			if err != nil {
				http.Error(w, "Fail to get remote configuration with error "+err.Error(), http.StatusBadRequest)
				return
			}

			//add prefix and clean
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			err = json.NewEncoder(w).Encode(remoteConfig)
			if err != nil {
				http.Error(w, "fail to write remote configuration with error "+err.Error(), http.StatusBadRequest)
				return
			}

			return
		}

		// I will get the remote configuration and return it.
		remoteConfig, err := config_.GetRemoteConfig(r.URL.Query().Get("host"), Utility.ToInt(r.URL.Query().Get("port")))
		if err != nil {
			// Try again with port 80...
			remoteConfig, err = config_.GetRemoteConfig(r.URL.Query().Get("host"), 80)
			if err != nil {
				http.Error(w, "Fail to get remote configuration with error "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		jsonStr, err := json.MarshalIndent(remoteConfig, "", "  ")
		if err != nil {
			http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)

		}

		_, err = w.Write(jsonStr)
		if err != nil {
			http.Error(w, "fail to write remote configuration with error "+err.Error(), http.StatusBadRequest)
			return
		}

		return

	}

	setupResponse(&w, r)

	// if the host is not the same...
	serviceID := r.URL.Query().Get("id") // the csr in base64

	//add prefix and clean
	config := globule.getConfig()

	// give list of path...
	config["Root"] = config_.GetRootDir()
	config["DataPath"] = config_.GetDataDir()
	config["ConfigPath"] = config_.GetConfigDir()
	config["WebRoot"] = config_.GetWebRootDir()
	config["Public"] = config_.GetPublicDirs()
	config["OAuth2ClientSecret"] = "********" // hide the secret...

	// ask for a service configuration...
	if len(serviceID) > 0 {
		services := config["Services"].(map[string]interface{})
		exist := false
		for _, service := range services {
			if service.(map[string]interface{})["Id"].(string) == serviceID || service.(map[string]interface{})["Name"].(string) == serviceID {
				config = service.(map[string]interface{})
				exist = true
				break
			}
		}
		if !exist {
			http.Error(w, "no service found with name or id "+serviceID, http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	jsonStr, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, "fail to write json with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

func dealwithErr(err error) {
	if err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func getHardwareData(w http.ResponseWriter, r *http.Request) {

	// Here I will test fi the request conain a host in the query.
	// If so I will redirect the request to the host.
	// If not I will return the hardware data of the current machine.
	// Handle the preflight options...
	hostname := r.URL.Query().Get("host")
	if len(hostname) == 0 {
		hostname = r.Host
	}

	// Receive http request...
	redirect, to := redirectTo(hostname)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	runtimeOS := runtime.GOOS

	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)

	stats := make(map[string]interface{})

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"
	diskStat, err := disk.Usage("/")
	dealwithErr(err)

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	dealwithErr(err)
	percentage, err := cpu.Percent(0, true)
	dealwithErr(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)

	// get interfaces MAC/hardware address
	interfStat, err := net.Interfaces()
	dealwithErr(err)

	stats["os"] = runtimeOS
	stats["memory"] = make(map[string]interface{}, 0)
	stats["memory"].(map[string]interface{})["total"] = strconv.FormatUint(vmStat.Total, 10)
	stats["memory"].(map[string]interface{})["free"] = strconv.FormatUint(vmStat.Free, 10)
	stats["memory"].(map[string]interface{})["used"] = strconv.FormatUint(vmStat.Used, 10)
	stats["memory"].(map[string]interface{})["used_percent"] = strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64)

	// get disk serial number.... strange... not available from disk package at compile time
	// undefined: disk.GetDiskSerialNumber
	//serial := disk.GetDiskSerialNumber("/dev/sda")
	stats["disk"] = make(map[string]interface{}, 0)
	stats["disk"].(map[string]interface{})["total"] = strconv.FormatUint(diskStat.Total, 10)
	stats["disk"].(map[string]interface{})["free"] = strconv.FormatUint(diskStat.Used, 10)
	stats["disk"].(map[string]interface{})["used_bytes"] = strconv.FormatUint(diskStat.Used, 10)

	// since my machine has one CPU, I'll use the 0 index
	// if your machine has more than 1 CPU, use the correct index
	// to get the proper data

	// cpu infos.
	stats["cpu"] = make(map[string]interface{}, 0)
	if len(cpuStat) > 0 {
		stats["cpu"].(map[string]interface{})["index_number"] = strconv.FormatInt(int64(cpuStat[0].CPU), 10)
		stats["cpu"].(map[string]interface{})["vendor_id"] = cpuStat[0].VendorID
		stats["cpu"].(map[string]interface{})["family"] = cpuStat[0].Family
		stats["cpu"].(map[string]interface{})["number_of_cores"] = strconv.FormatInt(int64(cpuStat[0].Cores), 10)
		stats["cpu"].(map[string]interface{})["model_name"] = cpuStat[0].ModelName
		stats["cpu"].(map[string]interface{})["speed"] = strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64)
		stats["cpu"].(map[string]interface{})["utilizations"] = make([]map[string]interface{}, 0)
		for idx, cpupercent := range percentage {
			stats["cpu"].(map[string]interface{})["utilizations"] = append(stats["cpu"].(map[string]interface{})["utilizations"].([]map[string]interface{}), map[string]interface{}{"idx": strconv.Itoa(idx), "utilization": strconv.FormatFloat(cpupercent, 'f', 2, 64)})
		}
	}

	stats["hostname"] = hostStat.Hostname
	stats["uptime"] = strconv.FormatUint(hostStat.Uptime, 10)
	stats["number_of_running_processes"] = strconv.FormatUint(hostStat.Procs, 10)

	// another way to get the operating system name
	// both darwin for Mac OSX, For Linux, can be ubuntu as platform
	// and linux for OS
	stats["os"] = hostStat.OS
	stats["platform"] = hostStat.Platform

	stats["  networkInterfacess"] = make([]map[string]interface{}, 0)

	// the unique hardware id for this machine
	for _, interf := range interfStat {
		networkInterfaces := make(map[string]interface{}, 0)
		networkInterfaces["mac"] = interf.HardwareAddr

		networkInterfaces["flags"] = interf.Flags
		networkInterfaces["addresses"] = make([]string, 0)
		for _, addr := range interf.Addrs {
			networkInterfaces["addresses"] = append(networkInterfaces["addresses"].([]string), addr.String())
		}

		stats["network_interfaces"] = append(stats["network_interfaces"].([]map[string]interface{}), networkInterfaces)

	}

	// generate a json output.
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)

	jsonStr, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, "fail to write json with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

/**
 * Return the ca certificate public key.
 */
func getCaCertificateHanldler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)

	crt, err := os.ReadFile(globule.creds + "/ca.crt")
	if err != nil {
		http.Error(w, "Client ca cert not found", http.StatusBadRequest)
		return
	}

	_, err = fmt.Fprint(w, string(crt))
	if err != nil {
		http.Error(w, "fail to write ca certificate with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

/**
 * Return the server SAN configuration file.
 */
func getSanConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)

	crt, err := os.ReadFile(globule.creds + "/san.conf")
	if err != nil {
		http.Error(w, "Client Subject Alernate Name configuration found!", http.StatusBadRequest)
		return
	}

	_, err = fmt.Fprint(w, string(crt))
	if err != nil {
		http.Error(w, "fail to write san configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

/**
 * Setup allow Cors policies.
 */
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	// Dynamically check if the origin is allowed
	origin := req.Header.Get("Origin")
	allowedOrigin := globule.Protocol + "://" + globule.Domain // Default to the globule domain
	if len(origin) > 0 {
		for _, allowed := range globule.AllowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = origin
				break
			}
		}
	}

	// Construct allowed methods
	allowedMethods := ""
	for i, method := range globule.AllowedMethods {
		allowedMethods += method
		if i < len(globule.AllowedMethods)-1 {
			allowedMethods += ","
		}
	}

	// Construct allowed headers
	allowedHeaders := ""
	for i, header := range globule.AllowedHeaders {
		allowedHeaders += header
		if i < len(globule.AllowedHeaders)-1 {
			allowedHeaders += ","
		}
	}

	header := (*w).Header()

	// Set the CORS headers
	if allowedOrigin != "" {
		header.Set("Access-Control-Allow-Origin", allowedOrigin)
		header.Set("Access-Control-Allow-Credentials", "true") // Required for credentials
	}
	header.Set("Access-Control-Allow-Methods", allowedMethods)
	header.Set("Access-Control-Allow-Headers", allowedHeaders)

	// Handle preflight requests
	if req.Method == http.MethodOptions {
		header.Set("Access-Control-Max-Age", "3600")
		header.Set("Access-Control-Allow-Private-Network", "true")
		(*w).WriteHeader(http.StatusNoContent)
		return
	}

	header.Set("Access-Control-Allow-Private-Network", "true")
}

/**
 * Sign ca certificate request and return a certificate.
 */
func signCaCertificateHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")

	w.WriteHeader(http.StatusCreated)

	// sign the certificate.
	csrStr := r.URL.Query().Get("csr") // the csr in base64
	csr, err := base64.StdEncoding.DecodeString(csrStr)
	if err != nil {
		http.Error(w, "Fail to decode csr base64 string", http.StatusBadRequest)
		return
	}

	// Now I will sign the certificate.
	crt, err := globule.signCertificate(string(csr))
	if err != nil {
		http.Error(w, "fail to sign certificate!", http.StatusBadRequest)
		return
	}

	// Return the result as text string.
	_, err = fmt.Fprint(w, crt)
	if err != nil {
		http.Error(w, "fail to write certificate with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

// Return true if the file is found in the public path...
func isPublic(path string) bool {
	public := config_.GetPublicDirs()
	path = strings.ReplaceAll(path, "\\", "/")

	for i := 0; i < len(public); i++ {
		if strings.HasPrefix(strings.ToLower(path), strings.ReplaceAll(strings.ToLower(public[i]), "\\", "/")) {
			return true
		}
	}

	return false
}

// ImageList is the structure for our response
type ImageList struct {
	Images []string `json:"images"`
}

// GetImagesHandler handles HTTP requests to retrieve a list of image files from a specified directory.
// It supports CORS preflight requests and redirects if necessary based on the request host.
// The handler expects a "path" query parameter indicating the directory to search for images.
// If the path is not provided or does not exist, it returns an error.
// On success, it responds with a JSON-encoded list of image file paths relative to the web root.
func GetImagesHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	dir := globule.webRoot

	// If a directory with the same name as the host in the request exist
	// it will be taken as root. Permission will be manage by the resource
	// manager and not simply the name of the directory. If you want to protect
	// a given you need to set permission on it.
	if Utility.Exists(dir + "/" + r.Host) {
		dir += "/" + r.Host
	}

	// so I will get the path from the query...
	path := r.URL.Query().Get("path")

	// If the path is not defined I will return an error.
	if len(path) == 0 {
		http.Error(w, "Failed to get images no path was given", http.StatusInternalServerError)
		return
	}

	// Be sure that the path start with a /.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if !Utility.Exists(dir + path) {

		http.Error(w, "Failed to get images path not found "+dir+path, http.StatusInternalServerError)
		return

	}

	// Get a list of images
	imageFiles, err := getListOfImages(dir + path)
	if err != nil {
		http.Error(w, "Failed to get images", http.StatusInternalServerError)
		return
	}

	// Create a response structure
	response := ImageList{Images: imageFiles}

	// I will replace all images path by the relative path.
	for i := 0; i < len(response.Images); i++ {
		response.Images[i] = strings.ReplaceAll(response.Images[i], dir, "")
	}

	// Marshal the response structure to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func getListOfImages(dirPath string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, walkErr error) error {
		// If Walk itself reports an error, propagate it
		if walkErr != nil {
			return walkErr
		}

		// Sanity check f (can be nil if walkErr != nil, but we already handled that)
		if f == nil {
			return fmt.Errorf("nil FileInfo for path: %s", path)
		}

		// Skip directories, only add files
		if !f.IsDir() {
			fileList = append(fileList, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %q: %w", dirPath, err)
	}

	return fileList, nil
}

//	GetFileSizeAtURL handles HTTP requests to retrieve the size of a file at a given URL.
//
// It expects a "url" query parameter specifying the file location.
// The handler sends a HEAD request to the provided URL and reads the "Content-Length" header
// to determine the file size. The result is returned as a JSON object with the "size" field.
// If the request fails or the file size cannot be determined, an appropriate error response is sent.
//
// Example request:
//
//	GET /get-file-size-at-url?url=https://example.com/file.mp4
//
// Response:
//
//	{ "size": 12345678 }
func GetFileSizeAtURL(w http.ResponseWriter, r *http.Request) {

	// here in case of file uploaded from other website like pornhub...
	url := r.URL.Query().Get("url")

	// we are interested in getting the file or object name
	// so take the last item from the slice

	// #nosec G107 -- URL is validated before use.
	resp, err := http.Head(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get file size at %s: %v\n", url, err)
		http.Error(w, "failed to get file size at "+url+": "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Is our request ok?

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return
	}

	// the Header "Content-Length" will let us know
	// the total file size to download
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	downloadSize := int64(size)

	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(&map[string]int64{"size": downloadSize})
	if err == nil {
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Fail to get file size at "+url+" with error "+err.Error(), http.StatusExpectationFailed)
	}
}

// FileUploadHandler handles HTTP file upload requests.
// It supports multipart form uploads, validates access permissions for users and applications,
// checks available storage space, and saves uploaded files to the appropriate location based on the path.
// The handler also sets resource ownership and triggers post-upload events such as video preview generation
// or file indexing based on the file type. Access control is enforced via token and application validation.
// Responds with appropriate HTTP status codes for errors such as invalid permissions, insufficient space,
// or file operation failures.
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // grab the multipart form
	if err != nil {
		fmt.Println("transfert error: ", err)
		http.Error(w, "failed to parse multipart message "+err.Error(), http.StatusBadRequest)
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames

	// Get the path where to upload the file.
	path := r.FormValue("path")
	path = strings.ReplaceAll(path, "\\", "/")

	// If application is defined.
	token := r.Header.Get("token")
	application := r.Header.Get("application")

	// If the header dosent contain the required values i I will try to get it from the
	// http query instead...
	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	if len(application) == 0 {
		// the token can be given by the url directly...
		application = r.URL.Query().Get("application")
	}

	user := ""
	hasAccess := false

	// TODO fix it and uncomment it...
	hasAccessDenied := false
	infos := []*rbacpb.ResourceInfos{}

	// Here I will validate applications...
	if len(application) != 0 {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service method
		if strings.HasPrefix(path, "/applications") {
			hasAccess, hasAccessDenied, err = globule.validateAction("/file.FileService/FileUploadHandler", application, rbacpb.SubjectType_APPLICATION, infos)
			if err != nil {
				http.Error(w, "fail to validate access with error "+err.Error(), http.StatusUnauthorized)
				return
			}
		}
	}

	// get the user id from the token...
	domain := r.URL.Query().Get("domain")
	if len(token) != 0 {
		var claims *security.Claims
		claims, err := security.ValidateToken(token)
		if err == nil {
			user = claims.Id + "@" + claims.UserDomain
			domain = claims.Domain
		} else {
			fmt.Println("fail to validate token with error ", err.Error())
			http.Error(w, "fail to validate token with error "+err.Error(), http.StatusUnauthorized)
			return
		}
	}

	if len(user) != 0 {
		if !hasAccess || hasAccessDenied {
			hasAccess, hasAccessDenied, err = globule.validateAction("/file.FileService/FileUploadHandler", user, rbacpb.SubjectType_ACCOUNT, infos)
			if err != nil {
				http.Error(w, "fail to validate access with error "+err.Error(), http.StatusUnauthorized)
				return
			}

			if hasAccess && !hasAccessDenied {
				hasAccess, hasAccessDenied, err = globule.validateAccess(user, rbacpb.SubjectType_ACCOUNT, "write", path)
				if err != nil {
					http.Error(w, "fail to validate access with error "+err.Error(), http.StatusUnauthorized)
					return
				}
			}
		}
	}

	// validate ressource access...
	if !hasAccess || hasAccessDenied {
		http.Error(w, "unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
		return
	}

	for _, f := range files { // loop through the files one by one
		file, err := f.Open()

		if err != nil {
			http.Error(w, "fail to open file with error "+err.Error(), http.StatusUnauthorized)
			return
		}

		defer func() {
			err = file.Close()
			if err != nil {
				http.Error(w, "fail to close file with error "+err.Error(), http.StatusInternalServerError)
			}
		}()

		// Create the file depending if the path is users, applications or something else...
		filePath := path + "/" + f.Filename
		size, _ := file.Seek(0, 2)
		if len(user) > 0 {
			// #nosec G115 -- Ok
			hasSpace, err := ValidateSubjectSpace(user, rbacpb.SubjectType_ACCOUNT, uint64(size))
			if !hasSpace || err != nil {
				http.Error(w, user+" has no space available to copy file "+filePath+" allocated space and try again.", http.StatusUnauthorized)
				return
			}
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, "fail to seek file with error "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Now if the os is windows I will remove the leading /
		if len(filePath) > 3 {
			if runtime.GOOS == "windows" && filePath[0] == '/' && filePath[2] == ':' {
				filePath = filePath[1:]
			}
		}

		if strings.HasPrefix(path, "/users") || strings.HasPrefix(path, "/applications") {
			filePath = strings.ReplaceAll(globule.data+"/files"+filePath, "\\", "/")
		} else if !isPublic(filePath) && !strings.HasPrefix(filePath, globule.webRoot) {
			filePath = strings.ReplaceAll(globule.webRoot+filePath, "\\", "/")
		}

		// #nosec G304 -- File path is constructed from validated input and constant strings.
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege error "+err.Error(), http.StatusUnauthorized)
			return
		}

		defer func() {
			err = out.Close()
			if err != nil {
				http.Error(w, "fail to close file with error "+err.Error(), http.StatusInternalServerError)
			}
		}()

		_, err = io.Copy(out, file) // file not files[i] !
		if err != nil {
			http.Error(w, "Unable to copy the file to the server. Check your write access privilege", http.StatusUnauthorized)
			return
		}

		// Here I will set the ressource owner.
		if len(user) > 0 {
			err = globule.addResourceOwner(path+"/"+f.Filename, "file", user+"@"+domain, rbacpb.SubjectType_ACCOUNT)
		} else if len(application) > 0 {
			err = globule.addResourceOwner(path+"/"+f.Filename, "file", application+"@"+domain, rbacpb.SubjectType_APPLICATION)
		}

		if err != nil {
			http.Error(w, "fail to set resource owner with error "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Now from the file extension i will read it mime type.
		if strings.LastIndex(filePath, ".") != -1 {
			fileExtension := filePath[strings.LastIndex(filePath, "."):]
			fileType := mime.TypeByExtension(fileExtension)
			filePath = strings.ReplaceAll(filePath, "\\", "/")
			if len(fileType) > 0 {
				if strings.HasPrefix(fileType, "video/") {
					// Here I will call convert video...
					err = globule.publish("generate_video_preview_event", []byte(filePath))
				} else if fileType == "application/pdf" || strings.HasPrefix(fileType, "text") {
					// Here I will call convert video...
					err = globule.publish("index_file_event", []byte(filePath))
				}

				if err != nil {
					http.Error(w, "fail to publish event with error "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}

}

// resolveImportPath resolves an import like "import 'foo/bar'" to a relative path from `path`.
// It searches under globule.webRoot/<top-level-segment-of-path>.
func resolveImportPath(path string, importPath string) (string, error) {
	// Pull the content between the first and last single quotes.
	i := strings.IndexByte(importPath, '\'')
	j := strings.LastIndexByte(importPath, '\'')
	if i == -1 || j <= i {
		return "", fmt.Errorf("malformed importPath %q", importPath)
	}
	want := importPath[i+1 : j]

	// Determine the search root: webRoot + first segment of `path` (before the first '/').
	firstSlash := strings.IndexByte(path, '/')
	if firstSlash == -1 {
		return "", fmt.Errorf("path %q has no '/' to determine search root", path)
	}
	searchRoot := filepath.Join(globule.webRoot, path[:firstSlash])

	var found string
	// Walk and stop as soon as we find a match.
	walkErr := filepath.WalkDir(searchRoot, func(p string, _ fs.DirEntry, err error) error {
		if err != nil {
			// Propagate real traversal errors.
			return err
		}

		// Normalize slashes for suffix check.
		pp := filepath.ToSlash(p)
		if strings.HasSuffix(pp, want) {
			found = p
			// Stop the entire walk without treating it as an error.
			return fs.SkipAll
		}
		return nil
	})

	// Intercept/handle errors from WalkDir.
	if walkErr != nil {
		// No special sentinel here because we used fs.SkipAll (which yields nil).
		return "", fmt.Errorf("walk %s: %w", searchRoot, walkErr)
	}
	if found == "" {
		return "", fmt.Errorf("could not resolve %q under %q", want, searchRoot)
	}

	// Make `found` relative to the directory containing `path`.
	// If `path` may be relative, normalize it to a base dir first.
	baseDir := filepath.Dir(filepath.Join(globule.webRoot, filepath.FromSlash(path)))
	rel, err := filepath.Rel(baseDir, found)
	if err != nil {
		return "", fmt.Errorf("relativizing %q to %q: %w", found, baseDir, err)
	}

	// Return with forward slashes for web usage.
	return filepath.ToSlash(rel), nil
}

// findHashedFile looks for a file with a hashed name based on the original file path.
func findHashedFile(originalPath string) (string, error) {
	// Get the directory of the original file
	dir := filepath.Dir(originalPath)

	// Get the base name and extension of the original file
	baseName := strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(originalPath))
	ext := filepath.Ext(originalPath)

	// Read the files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	// Regular expression to match hashed file names
	hashRegex := regexp.MustCompile(fmt.Sprintf(`^%s\.[a-f0-9]{8,}\%s$`, regexp.QuoteMeta(baseName), regexp.QuoteMeta(ext)))

	// Search for a matching hashed file
	for _, file := range files {
		if hashRegex.MatchString(file.Name()) {
			return filepath.Join(dir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("hashed file not found for %s", originalPath)
}

func streamHandler(path string, w http.ResponseWriter, r *http.Request) {
	// Set the appropriate response headers for streaming
	w.Header().Set("Content-Type", "video/webm")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Prepare FFmpeg command to decode and stream the MKV file
	cmd := exec.Command("ffmpeg", "-i", path, "-c:v", "libvpx", "-c:a", "libvorbis", "-f", "webm", "pipe:1")

	// Get the output of FFmpeg (streaming to stdout)
	ffmpegOut, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with FFmpeg: %v", err), http.StatusInternalServerError)
		return
	}

	// Start the FFmpeg process
	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Error starting FFmpeg: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a channel to detect if the connection is closed
	done := make(chan struct{})
	go func() {
		// Wait for the client to close the connection
		<-r.Context().Done()
		// Kill the FFmpeg process if the connection is closed
		err = cmd.Process.Kill()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error killing FFmpeg: %v", err), http.StatusInternalServerError)
		}
		close(done)
	}()

	// Stream the FFmpeg output to the HTTP response
	_, err = io.Copy(w, ffmpegOut)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error streaming video: %v", err), http.StatusInternalServerError)
	}

	// Wait for the FFmpeg process to finish or the connection to be closed
	<-done
	err = cmd.Wait()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error waiting for FFmpeg: %v", err), http.StatusInternalServerError)
	}
}

// ServeFileHandler handles HTTP requests for serving files from the server.
// It supports reverse proxying, access control based on tokens and application headers,
// and dynamic file path resolution. The handler checks if the requested file path should
// be redirected to a reverse proxy, validates access permissions for protected resources,
// and serves files with appropriate content types. Special handling is provided for
// streaming video files, CA certificates, and JavaScript imports. If access is denied
// or the file is not found, the handler responds with the appropriate HTTP error code.
//
// Request headers:
//   - "application": Specifies the application context for access control.
//   - "token": JWT token for user authentication and authorization.
//
// Query parameters:
//   - "application": Alternative way to specify the application context.
//   - "token": Alternative way to provide the JWT token.
//
// Access control:
//   - Validates user or application access to protected resources using RBAC.
//   - Public resources are served without access validation.
//   - Special handling for streaming files and hidden directories.
//
// Content types:
//   - Sets "Content-Type" header based on file extension (.js, .css, .html, .htm).
//
// Errors:
//   - Responds with 400 Bad Request if no file path is provided.
//   - Responds with 401 Unauthorized if access is denied.
//   - Responds with 204 No Content if the file is not found.
//
// Parameters:
//   - w: http.ResponseWriter to write the response.
//   - r: *http.Request representing the incoming HTTP request.
func ServeFileHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	//add prefix and clean
	rqstPath := path.Clean(r.URL.Path)

	if rqstPath == "/null" {
		http.Error(w, "No file path was given in the file url path!", http.StatusBadRequest)
	}

	// I will test if the requested path is in the reverse proxy list.
	// if it is the case I will redirect the request to the reverse proxy.
	for _, proxy := range globule.ReverseProxies {
		proxyPath := strings.TrimSpace(strings.Split(proxy.(string), "|")[1])
		proxyURLStr := strings.TrimSpace(strings.Split(proxy.(string), "|")[0])

		if strings.HasPrefix(rqstPath, proxyPath) {
			// Create a reverse proxy
			proxyURL, _ := url.Parse(proxyURLStr)

			// Connect to the proxy host
			hostURL, _ := url.Parse(proxyURL.Scheme + "://" + proxyURL.Host)

			reverseProxy := httputil.NewSingleHostReverseProxy(hostURL)

			// Update the request URL and headers
			r.URL, _ = url.Parse(proxyURLStr)

			// Update headers to reflect the forwarded host
			r.Host = proxyURL.Host

			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

			// Print request details
			//printRequestInfo(r)

			// Serve the request via the reverse proxy
			reverseProxy.ServeHTTP(w, r)
			return
		}
	}

	setupResponse(&w, r)
	dir := globule.webRoot

	// If a directory with the same name as the host in the request exist
	// it will be taken as root. Permission will be manage by the resource
	// manager and not simply the name of the directory. If you want to protect
	// a given you need to set permission on it.
	if Utility.Exists(dir + "/" + r.Host) {
		dir += "/" + r.Host
	}

	// Now I will test if a token is given in the header and manage it file access.
	application := r.Header.Get("application")
	token := r.Header.Get("token")

	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	if len(application) == 0 {
		// the token can be given by the url directly...
		application = r.URL.Query().Get("application")
	}

	// If the header dosent contain the required values i I will try to get it from the
	if token == "null" || token == "undefined" {
		token = ""
	}

	// If the path is '/' it mean's no application name was given and we are
	// at the root.
	if rqstPath == "/" {
		// if a default application is define in the globule i will use it.
		if len(globule.IndexApplication) > 0 {
			rqstPath += globule.IndexApplication
			application = globule.IndexApplication
		}

	} else if strings.Count(rqstPath, "/") == 1 {
		if strings.HasSuffix(rqstPath, ".js") ||
			strings.HasSuffix(rqstPath, ".json") ||
			strings.HasSuffix(rqstPath, ".css") ||
			strings.HasSuffix(rqstPath, ".htm") ||
			strings.HasSuffix(rqstPath, ".html") {
			if Utility.Exists(dir + "/" + rqstPath) {
				rqstPath = "/" + globule.IndexApplication + rqstPath
			}
		}
	}

	hasAccess := true
	var name string
	if strings.HasPrefix(rqstPath, "/users/") || strings.HasPrefix(rqstPath, "/applications/") || strings.HasPrefix(rqstPath, "/templates/") || strings.HasPrefix(rqstPath, "/projects/") {
		dir = globule.data + "/files"
		if !strings.Contains(rqstPath, "/.hidden/") {
			hasAccess = false
		}
	}

	// Now if the os is windows I will remove the leading /
	if len(rqstPath) > 3 {
		if runtime.GOOS == "windows" && rqstPath[0] == '/' && rqstPath[2] == ':' {
			rqstPath = rqstPath[1:]
		}
	}
	// path to file
	if !isPublic(rqstPath) {
		name = path.Join(dir, rqstPath)
	} else {
		name = rqstPath
		hasAccess = false // force validation (denied access...)
	}

	// stream, the validation is made on the directory containning the playlist...
	if strings.Contains(rqstPath, "/.hidden/") ||
		strings.HasSuffix(rqstPath, ".ts") ||
		strings.HasSuffix(rqstPath, "240p.m3u8") ||
		strings.HasSuffix(rqstPath, "360p.m3u8") ||
		strings.HasSuffix(rqstPath, "480p.m3u8") ||
		strings.HasSuffix(rqstPath, "720p.m3u8") ||
		strings.HasSuffix(rqstPath, "1080p.m3u8") ||
		strings.HasSuffix(rqstPath, "2160p.m3u8") {
		hasAccess = true
	}

	// this is the ca certificate use to sign client certificate.
	if rqstPath == "/ca.crt" {
		name = globule.creds + rqstPath
	}

	if strings.Contains(rqstPath, "pdf") {

		fmt.Println("validate access ", rqstPath)
	}
	hasAccessDenied := false
	var err error
	var userID string
	if len(token) != 0 && !hasAccess {
		var claims *security.Claims
		claims, err = security.ValidateToken(token)
		userID = claims.Id + "@" + claims.UserDomain
		if err == nil {
			hasAccess, hasAccessDenied, err = globule.validateAccess(userID, rbacpb.SubjectType_ACCOUNT, "read", rqstPath)
		} else {
			fmt.Println("fail to validate token with error: ", err)
		}
	}

	// Here I will validate applications...
	if isPublic(rqstPath) && !hasAccessDenied && !hasAccess {
		hasAccess = true
	} else if !hasAccess && !hasAccessDenied && len(application) != 0 {
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "read", rqstPath)
	}

	// validate ressource access...
	if !hasAccess || hasAccessDenied || err != nil {
		msg := "unable to read the file " + rqstPath + " Check your access privilege"
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	var code string
	// If the file is a javascript file...
	hasChange := false

	if !Utility.Exists(name) {
		name = "/" + rqstPath // try network path...
	}

	// Try to find the file in the hidden directory...
	if r.Method == http.MethodGet {
		if strings.HasSuffix(name, ".mkv") {
			streamHandler(name, w, r) // stream the video
			return
		}
	}

	// #nosec G304 -- File path is constructed from validated input and constant strings.
	f, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			name, err = findHashedFile(name)
			if err == nil {
				// #nosec G304 -- File path is constructed from validated input and constant strings.
				f, err = os.Open(name)
				if err != nil {
					http.Error(w, "File "+rqstPath+" not found!", http.StatusNoContent)
					return
				}
			} else {
				http.Error(w, "File "+rqstPath+" not found!", http.StatusNoContent)
				return
			}
		}
	}

	defer func() {
		err = f.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error closing file: %v", err), http.StatusInternalServerError)
		}
	}()

	if strings.HasSuffix(name, ".js") {
		w.Header().Add("Content-Type", "application/javascript")

		if err == nil {
			//hasChange = true
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "import") {
					if strings.Contains(line, `'@`) {
						filePath, err := resolveImportPath(rqstPath, line)
						if err == nil {
							line = line[0:strings.Index(line, `'@`)] + `'` + filePath + `'`
							hasChange = true
						}
					}
				}
				code += line + "\n"
			}
		}

	} else if strings.HasSuffix(name, ".css") {
		w.Header().Add("Content-Type", "text/css")
	} else if strings.HasSuffix(name, ".html") || strings.HasSuffix(name, ".htm") {
		w.Header().Add("Content-Type", "text/html")
	}

	// if the file has change...
	if !hasChange {
		//log.Println("server file ", name)
		http.ServeFile(w, r, name)
	} else {
		http.ServeContent(w, r, name, time.Now(), strings.NewReader(code))
	}
}

// GetIMDBPoster retrieves the highest resolution poster image URL for a given IMDb title ID.
// It performs the following steps:
//  1. Visits the IMDb title page to find the media viewer URL.
//  2. Extracts the image resource ID (rmID) from the media viewer URL.
//  3. Visits the media viewer page and parses the available images to find the poster
//     with the matching rmID, selecting the highest resolution image from the srcset attribute.
//  4. Returns the poster image URL or an error if the poster cannot be found.
//
// Parameters:
//
//	imdbID - The IMDb title identifier (e.g., "tt0111161").
//
// Returns:
//
//	string - The URL of the poster image.
//	error  - An error if the poster cannot be found or if any step fails.
func GetIMDBPoster(imdbID string) (string, error) {
	mainURL := "https://www.imdb.com/title/" + imdbID + "/"
	var posterViewerURL string
	var posterURL string

	c := colly.NewCollector()

	// Step 1: Find media viewer URL
	c.OnHTML("a.ipc-lockup-overlay", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "/mediaviewer/") && posterViewerURL == "" {
			posterViewerURL = "https://www.imdb.com" + href
		}
	})

	if err := c.Visit(mainURL); err != nil {
		return "", err
	}
	if posterViewerURL == "" {
		return "", fmt.Errorf("could not find media viewer URL")
	}

	// Step 2: Extract rmID from URL
	reRM := regexp.MustCompile(`/mediaviewer/(rm\d+)/`)
	match := reRM.FindStringSubmatch(posterViewerURL)
	if len(match) < 2 {
		return "", fmt.Errorf("could not extract rmID")
	}
	rmID := match[1] + "-curr"

	// Step 3: Visit media viewer and find correct image
	imgCollector := colly.NewCollector()

	imgCollector.OnHTML("img", func(e *colly.HTMLElement) {
		if e.Attr("data-image-id") == rmID {
			srcset := e.Attr("srcset")
			if srcset != "" {
				// Parse srcset and get highest resolution
				maxResURL := ""
				maxWidth := 0
				for _, part := range strings.Split(srcset, ",") {
					part = strings.TrimSpace(part)
					if items := strings.Split(part, " "); len(items) == 2 {
						url := items[0]
						widthStr := items[1]
						if strings.HasSuffix(widthStr, "w") {
							width, err := strconv.Atoi(strings.TrimSuffix(widthStr, "w"))
							if err == nil && width > maxWidth {
								maxWidth = width
								maxResURL = url
							}
						}
					}
				}
				if maxResURL != "" {
					posterURL = maxResURL
					return
				}
			}
			// fallback to src
			if posterURL == "" {
				posterURL = e.Attr("src")
			}
		}
	})

	if err := imgCollector.Visit(posterViewerURL); err != nil {
		return "", fmt.Errorf("failed to visit viewer page: %v", err)
	}
	if posterURL == "" {
		return "", fmt.Errorf("poster image not found")
	}
	return posterURL, nil
}

/**
 * Return a list of IMDB titles from a keyword...
 */
func getImdbTitlesHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// if the host is not the same...
	query := r.URL.Query().Get("query") // the csr in base64

	titles, err := imdb.SearchTitle(client, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(titles) == 0 {
		fmt.Fprintf(os.Stderr, "Not found.")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "No titles found", http.StatusBadRequest)
		}
		return
	}

	// now i will iterate over the titles and set the poster image...
	for i := range titles {
		title := titles[i]
		poster, err := GetIMDBPoster(titles[i].ID)
		if err == nil {
			// Now I will download the poster image...
			titles[i].Poster.ContentURL = poster
			titles[i].Poster.URL = poster
			titles[i].Poster.ID = title.ID
		}

	}

	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)

	jsonStr, err := json.MarshalIndent(titles, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, "fail to write json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Found ", len(titles), " titles.")

}

// ////////////////////////// imdb missing sesson and episode number info... /////////////////////////
// get the thumbnail fil with help of youtube dl...
func downloadThumbnail(videoID, videoURL, videoPath string) (string, error) {

	if len(videoID) == 0 {
		return "", errors.New("no video id was given")
	}
	if len(videoURL) == 0 {
		return "", errors.New("no video url was given")
	}
	if len(videoPath) == 0 {
		return "", errors.New("no video path was given")
	}

	lastIndex := -1
	if strings.Contains(videoPath, ".mp4") {
		lastIndex = strings.LastIndex(videoPath, ".")
	}

	// The hidden folder path...
	filePath := videoPath[0:strings.LastIndex(videoPath, "/")]

	name := videoPath[strings.LastIndex(videoPath, "/")+1:]
	if lastIndex != -1 {
		name = videoPath[strings.LastIndex(videoPath, "/")+1 : lastIndex]
	}

	thumbnailPath := filePath + "/.hidden/" + name + "/__thumbnail__"

	if Utility.Exists(thumbnailPath + "/" + "data_url.txt") {

		// #nosec G304 -- Path is anchored, allow-listed, symlinks resolved, and confined to base.
		thumbnail, err := os.ReadFile(thumbnailPath + "/" + "data_url.txt")
		if err != nil {
			return "", err
		}

		return string(thumbnail), nil
	}

	err := Utility.CreateDirIfNotExist(thumbnailPath)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("yt-dlp", videoURL, "-o", videoID, "--write-thumbnail", "--skip-download")
	cmd.Dir = thumbnailPath

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	files, err := Utility.ReadDir(thumbnailPath)
	if err != nil {
		return "", err
	}

	thumbnail, err := Utility.CreateThumbnail(filepath.Join(thumbnailPath, files[0].Name()), 300, 180)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(thumbnailPath+"/"+"data_url.txt", []byte(thumbnail), 0600)
	if err != nil {
		return "", err
	}

	// cointain the data url...
	return thumbnail, nil
}

// GetCoverDataURL handles HTTP requests to retrieve a cover image as a data URL for a video.
// It expects query parameters "id" (video ID), "url" (video URL), and "path" (video file path).
// The function attempts to download the thumbnail using these parameters and returns the data URL.
// On error, it responds with an appropriate HTTP error status and message.
func GetCoverDataURL(w http.ResponseWriter, r *http.Request) {
	// here in case of file uploaded from other website like pornhub...
	videoID := r.URL.Query().Get("id")
	videoURL := r.URL.Query().Get("url")
	videoPath := r.URL.Query().Get("path")

	dataURL, err := downloadThumbnail(videoID, videoURL, videoPath)
	if err != nil {
		http.Error(w, "fail to create data url with error'"+err.Error()+"'", http.StatusExpectationFailed)
		return
	}

	_, err = w.Write([]byte(dataURL))
	if err != nil {
		http.Error(w, "fail to write data url with error'"+err.Error()+"'", http.StatusInternalServerError)
		return
	}
}

func getSeasonAndEpisodeNumber(titleID string) (int, int, string, error) {

	resp, err := client.Get(`https://www.imdb.com/title/` + titleID)
	if err != nil {
		return -1, -1, "", err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fail to close response body with error: %v\n", err)
		}
	}()

	season := 0
	episode := 0
	serie := ""

	// The first regex to locate the information...
	regexSE := regexp.MustCompile(`>S\d{1,2}<!-- -->\.<!-- -->E\d{1,2}<`)
	page, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, "", err
	}

	matchsSE := regexSE.FindStringSubmatch(string(page))

	if len(matchsSE) > 0 {
		regexS := regexp.MustCompile(`S\d{1,2}`)
		matchsS := regexS.FindStringSubmatch(matchsSE[0])
		if len(matchsS) > 0 {
			season = Utility.ToInt(matchsS[0][1:])
		}

		regexE := regexp.MustCompile(`E\d{1,2}`)
		matchsE := regexE.FindStringSubmatch(matchsSE[0])
		if len(matchsE) > 0 {
			episode = Utility.ToInt(matchsE[0][1:])
		}
	}

	// Regex to find the series ID in the href attribute
	re := regexp.MustCompile(`data-testid="hero-title-block__series-link".*?href="/title/(tt\d{7,8})/`)

	// Extract the series ID
	matches := re.FindStringSubmatch(string(page))

	if len(matches) > 1 {
		serie = matches[1]
	}

	return season, episode, serie, nil
}

/**
 * Return a json object with the movie information from imdb.
 */
func getImdbTitleHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	id := r.URL.Query().Get("id") // the csr in base64

	title, err := imdb.NewTitle(client, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)

	titleMap, _ := Utility.ToMap(title)

	if title.Type == "TVEpisode" {
		s, e, t, err := getSeasonAndEpisodeNumber(id)
		if err == nil {
			titleMap["Season"] = s
			titleMap["Episode"] = e
			titleMap["Serie"] = t
		} else {
			fmt.Println("fail to find episode info ", err)
		}
	}

	// Now I will try to complete the casting informations...
	if titleMap["Actors"] != nil {
		for i := 0; i < len(titleMap["Actors"].([]interface{})); i++ {
			p := titleMap["Actors"].([]interface{})[i].(map[string]interface{})
			a, err := setPersonInformation(p)
			if err == nil {
				titleMap["Actors"].([]interface{})[i] = a
			}
		}
	}

	if titleMap["Writers"] != nil {
		for i := 0; i < len(titleMap["Writers"].([]interface{})); i++ {
			p := titleMap["Writers"].([]interface{})[i].(map[string]interface{})
			w, err := setPersonInformation(p)
			if err == nil {
				titleMap["Writers"].([]interface{})[i] = w
			}
		}
	}

	if titleMap["Directors"] != nil {
		for i := 0; i < len(titleMap["Directors"].([]interface{})); i++ {
			p := titleMap["Directors"].([]interface{})[i].(map[string]interface{})
			d, err := setPersonInformation(p)
			if err == nil {
				titleMap["Directors"].([]interface{})[i] = d
			}
		}
	}

	// now i will get the poster image...
	poster, err := GetIMDBPoster(title.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Now I will download the poster image...
	titleMap["Poster"] = map[string]interface{}{"URL": poster, "ContentURL": poster, "titleID": title.ID, "id": title.ID}

	jsonStr, err := json.MarshalIndent(titleMap, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Found ", len(titleMap), " titles.")

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, "fail to write json with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

func setPersonInformation(person map[string]interface{}) (map[string]interface{}, error) {
	movieCollector := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com", "imdb.com"),
	)

	// So here I will define collector's...
	biographySelector := `a[name="mini_bio"]`
	movieCollector.OnHTML(biographySelector, func(e *colly.HTMLElement) {

		// keep the text only...
		person["Biography"] = e.DOM.Next().Next().Text()
	})

	// The profile image.
	profilePictureSelector := `#main > div.article.listo > div.subpage_title_block.name-subpage-header-block > a > img`
	movieCollector.OnHTML(profilePictureSelector, func(e *colly.HTMLElement) {
		person["Picture"] = strings.TrimSpace(e.Attr("src"))
	})

	// The birtdate
	birthdateSelector := `#overviewTable > tbody > tr:nth-child(1) > td:nth-child(2) > time`
	movieCollector.OnHTML(birthdateSelector, func(e *colly.HTMLElement) {
		person["BirthDate"] = e.Attr("datetime")
	})

	// The birtplace.
	birthplaceSelector := `#overviewTable > tbody > tr:nth-child(1) > td:nth-child(2) > a`
	movieCollector.OnHTML(birthplaceSelector, func(e *colly.HTMLElement) {
		person["BirthPlace"] = e.Text
	})

	url := person["URL"].(string) + "/bio?ref_=nm_ov_bio_sm"

	err := movieCollector.Visit(url)
	if err != nil {
		return nil, err
	}

	return person, nil
}
