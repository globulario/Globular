package config_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "github.com/globulario/Globular/internal/gateway/handlers/config"
)

type fakeSvcPerms struct {
	perms []any
	err   error
}

func (f fakeSvcPerms) LoadPermissions(serviceID string) ([]any, error) {
	return f.perms, f.err
}

func TestGetServicePermissions_OK_201(t *testing.T) {
	h := cfg.NewGetServicePermissions(fakeSvcPerms{
		perms: []any{"read", "write"},
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/api/get-service-permissions?id=434843b4-8512-321a-a54c-089bcd9a8cb7", nil)
	h.ServeHTTP(rr, req)

	t.Logf("response: %s", rr.Body.String())
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var got []any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 || got[0] != "read" || got[1] != "write" {
		t.Fatalf("unexpected perms: %#v", got)
	}
}

func TestGetServicePermissions_EmptyArrayWhenNil(t *testing.T) {
	h := cfg.NewGetServicePermissions(fakeSvcPerms{
		perms: nil, // provider returns nil → handler should return []
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/api/getServicePermissions?id=svc1", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var got any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// JSON [] unmarshals to []any
	if arr, ok := got.([]any); !ok || len(arr) != 0 {
		t.Fatalf("expected empty array, got %#v", got)
	}
}
