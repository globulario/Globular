package config

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/globulario/services/golang/config"
)

type fakeDescriber struct {
	desc    config.ServiceDesc
	path    string
	timeout time.Duration
	err     error
}

func (f *fakeDescriber) DescribeService(_ string, timeout time.Duration) (config.ServiceDesc, string, error) {
	f.timeout = timeout
	if f.err != nil {
		return config.ServiceDesc{}, "", f.err
	}
	return f.desc, f.path, nil
}

func TestDescribeServiceMissingName(t *testing.T) {
	h := NewDescribeService(&fakeDescriber{})
	req := httptest.NewRequest(http.MethodGet, "/api/describe-service", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDescribeServiceError(t *testing.T) {
	provider := &fakeDescriber{err: errors.New("boom")}
	h := NewDescribeService(provider)

	req := httptest.NewRequest(http.MethodGet, "/api/describe-service?name=svc", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "boom") {
		t.Fatalf("body missing error, got %q", rec.Body.String())
	}
}

func TestDescribeServiceSuccess(t *testing.T) {
	desc := config.ServiceDesc{Name: "svc", Version: "1.2.3"}
	provider := &fakeDescriber{desc: desc, path: "/bin/svc"}
	h := NewDescribeService(provider)

	req := httptest.NewRequest(http.MethodGet, "/api/describe-service?name=svc&timeoutMs=1500", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
	}

	var out struct {
		config.ServiceDesc
		BinaryPath string `json:"binaryPath"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if out.Name != desc.Name {
		t.Fatalf("expected name %q, got %q", desc.Name, out.Name)
	}
	if out.BinaryPath != provider.path {
		t.Fatalf("expected binary %q, got %q", provider.path, out.BinaryPath)
	}
	if provider.timeout != 1500*time.Millisecond {
		t.Fatalf("timeout mismatch: got %v", provider.timeout)
	}
}
