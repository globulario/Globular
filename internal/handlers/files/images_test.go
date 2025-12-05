package files_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	files "github.com/globulario/Globular/internal/handlers/files"
)

type fakeLister struct {
	out []string
	err error
}

func (f fakeLister) ListImages(dir string) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.out, nil
}

func TestGetImages_MissingDir_400(t *testing.T) {
	h := files.NewGetImages(fakeLister{})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/get-images", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetImages_DirRejected_400(t *testing.T) {
	h := files.NewGetImages(fakeLister{err: errors.New("dir not allowed")})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/api/get-images?dir=/secret", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (body: %s)", rr.Code, rr.Body.String())
	}
}

func TestGetImages_OK_200(t *testing.T) {
	h := files.NewGetImages(fakeLister{out: []string{"a.jpg", "b.png"}})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://globular.io/api/get-images?dir=/public", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	var got struct {
		Files []string `json:"files"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got.Files) != 2 || got.Files[0] != "a.jpg" || got.Files[1] != "b.png" {
		t.Fatalf("unexpected files: %#v", got.Files)
	}
}
