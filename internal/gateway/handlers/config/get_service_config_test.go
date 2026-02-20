package config_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "github.com/globulario/Globular/internal/gateway/handlers/config"
)

func TestGetServiceConfig_ByName(t *testing.T) {
	local := map[string]any{
		"Services": map[string]any{
			"media": map[string]any{"Id": "svc-media", "Name": "media.MediaService"},
		},
	}

	h := cfg.NewGetServiceConfig(fakeProvider{
		local: local,
		services: map[string]map[string]any{
			"media.MediaService": {
				"Id":                           "svc-media",
				"Name":                         "media.MediaService",
				"AutomaticVideoConversion":     true,
				"AutomaticStreamConversion":    false,
				"StartVideoConversionHour":     "03:00",
				"MaximumVideoConversionDelay":  "4h",
				"Foo":                          "bar",
				"AutomaticAudioTranscription":  true,
				"MaximumTranscriptionParallel": 2,
			},
		},
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/config/media.MediaService", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}

	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["Name"] != "media.MediaService" || got["Foo"] != "bar" {
		t.Fatalf("expected media service config, got %#v", got)
	}
	if got["MaximumVideoConversionDelay"] != "4h" {
		t.Fatalf("expected detailed property, got %#v", got["MaximumVideoConversionDelay"])
	}
}

func TestGetServiceConfig_NotFound(t *testing.T) {
	local := map[string]any{
		"Services": map[string]any{},
	}

	h := cfg.NewGetServiceConfig(fakeProvider{local: local})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/config/missing.Service", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when service is missing, got %d", rr.Code)
	}
}
