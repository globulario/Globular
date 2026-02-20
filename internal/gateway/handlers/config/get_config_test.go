package config_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "github.com/globulario/Globular/internal/gateway/handlers/config"
)

// ---- fakes ----

type fakeProvider struct {
	local    map[string]any
	services map[string]map[string]any
}

func (f fakeProvider) Address() (string, error)    { return "127.0.0.1:8080", nil }
func (f fakeProvider) MyIP() string                { return "1.2.3.4" }
func (f fakeProvider) LocalConfig() map[string]any { return shallowCopy(f.local) }
func (f fakeProvider) ServiceConfig(idOrName string) (map[string]any, error) {
	for key, svc := range f.services {
		if key == idOrName {
			return shallowCopy(svc), nil
		}
		if id, _ := svc["Id"].(string); id == idOrName {
			return shallowCopy(svc), nil
		}
		if name, _ := svc["Name"].(string); name == idOrName {
			return shallowCopy(svc), nil
		}
	}
	return nil, fmt.Errorf("service %s not found", idOrName)
}
func (f fakeProvider) RootDir() string      { return "/root" }
func (f fakeProvider) DataDir() string      { return "/data" }
func (f fakeProvider) ConfigDir() string    { return "/cfg" }
func (f fakeProvider) WebRootDir() string   { return "/web" }
func (f fakeProvider) PublicDirs() []string { return []string{"/public"} }

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// ---- tests ----

func TestGetConfig_LocalSnapshot_MasksSecret(t *testing.T) {
	local := map[string]any{
		"Services":           map[string]any{},
		"OAuth2ClientSecret": "should-be-masked",
	}
	h := cfg.NewGetConfig(fakeProvider{local: local})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/api/getConfig", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}

	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["OAuth2ClientSecret"] != "********" {
		t.Fatalf("secret not masked, got %v", got["OAuth2ClientSecret"])
	}
	if got["DataPath"] != "/data" {
		t.Fatalf("expected DataPath /data, got %v", got["DataPath"])
	}
	// Public is []string -> becomes []any after JSON decode
	if pubs, ok := got["Public"].([]any); !ok || len(pubs) == 0 || pubs[0] != "/public" {
		t.Fatalf("expected Public ['/public'], got %#v", got["Public"])
	}
}

func TestGetConfig_SelectsServiceByID(t *testing.T) {
	local := map[string]any{
		"Services": map[string]any{
			"svc1": map[string]any{"Id": "abc", "Name": "svc1"},
			"svc2": map[string]any{"Id": "def", "Name": "svc2"},
		},
	}
	p := fakeProvider{
		local: local,
		services: map[string]map[string]any{
			"svc1": {
				"Id":                       "abc",
				"Name":                     "svc1",
				"AutomaticVideoConversion": true,
				"Foo":                      "bar",
			},
		},
	}
	h := cfg.NewGetConfig(p)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/api/getConfig?id=abc", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}

	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["Id"] != "abc" || got["Foo"] != "bar" {
		t.Fatalf("expected service {Id:abc Foo:bar}, got %#v", got)
	}
	if got["AutomaticVideoConversion"] != true {
		t.Fatalf("expected AutomaticVideoConversion true, got %#v", got["AutomaticVideoConversion"])
	}
	// When a specific service is selected, the response should be that service only
	if _, exists := got["Services"]; exists {
		t.Fatalf("did not expect top-level Services in service view, got %#v", got["Services"])
	}
}
