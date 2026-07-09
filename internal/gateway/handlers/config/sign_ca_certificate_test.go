package config

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeSigner struct {
	cert string
	err  error
}

func (f fakeSigner) SignCSR(_ []byte) (string, error) { return f.cert, f.err }

type allowAuthorizer struct{}

func (allowAuthorizer) AuthorizeSign(_ *http.Request) error { return nil }

type denyAuthorizer struct{}

func (denyAuthorizer) AuthorizeSign(_ *http.Request) error {
	return errors.New("requester is not a cluster member")
}

func signRequest() *http.Request {
	csr := base64.StdEncoding.EncodeToString([]byte("fake-csr"))
	return httptest.NewRequest(http.MethodGet, "/sign_ca_certificate?csr="+csr, nil)
}

// A nil authorizer must fail closed — an unwired gate must never fall back to
// open signing (a CA-signed cert is an etcd mTLS credential).
func TestSignCACertificate_NilAuthorizerFailsClosed(t *testing.T) {
	h := NewSignCACertificate(fakeSigner{cert: "PEM"}, nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, signRequest())
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("nil authorizer: got %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}
}

// A denying authorizer produces 403 and no certificate.
func TestSignCACertificate_DeniedRequestGets403(t *testing.T) {
	h := NewSignCACertificate(fakeSigner{cert: "PEM"}, denyAuthorizer{})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, signRequest())
	if rr.Code != http.StatusForbidden {
		t.Fatalf("denied request: got %d, want %d", rr.Code, http.StatusForbidden)
	}
	if body := rr.Body.String(); body == "PEM" {
		t.Fatal("denied request must not receive a signed certificate")
	}
}

// An authorized request is signed (legacy 201 behavior preserved).
func TestSignCACertificate_AuthorizedRequestSigned(t *testing.T) {
	h := NewSignCACertificate(fakeSigner{cert: "PEM"}, allowAuthorizer{})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, signRequest())
	if rr.Code != http.StatusCreated {
		t.Fatalf("authorized request: got %d, want %d", rr.Code, http.StatusCreated)
	}
	if rr.Body.String() != "PEM" {
		t.Fatalf("authorized request: got body %q, want signed cert", rr.Body.String())
	}
}
