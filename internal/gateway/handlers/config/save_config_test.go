package config_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "github.com/globulario/Globular/internal/gateway/handlers/config"
)

type okValidator struct{}

func (okValidator) Validate(string) error { return nil }

type failValidator struct{}

func (failValidator) Validate(string) error { return fmt.Errorf("bad token") }

type memSaver struct{ got map[string]any }

func (m *memSaver) Save(v map[string]any) error { m.got = v; return nil }

func TestSaveConfig_RequiresToken(t *testing.T) {
	h := cfg.NewSaveConfig(&memSaver{}, okValidator{})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/saveConfig", bytes.NewBufferString(`{}`))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
}

func TestSaveConfig_ValidTokenSaves204(t *testing.T) {
	s := &memSaver{}
	h := cfg.NewSaveConfig(s, okValidator{})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "https://globular.io/api/saveConfig", bytes.NewBufferString(`{"A":1}`))
	req.Header.Set("token", "abc")
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rr.Code)
	}
	if s.got["A"] != float64(1) { // JSON numbers decode to float64
		t.Fatalf("payload not saved, got %#v", s.got)
	}
}

func TestSaveConfig_InvalidToken401(t *testing.T) {
	s := &memSaver{}
	h := cfg.NewSaveConfig(s, failValidator{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "https://globular.io/api/saveConfig", bytes.NewBufferString(`{"A":1}`))
	req.Header.Set("token", "invalid")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 for invalid token, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	if s.got != nil {
		t.Fatalf("config should not be saved on invalid token, got %#v", s.got)
	}
}
