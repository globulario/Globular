package files_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	files "github.com/globulario/Globular/internal/handlers/files"
)

type fakeUpload struct {
	dataRoot   string
	publicDirs []string
	allowWrite bool
}

func (f fakeUpload) DataRoot() string     { return f.dataRoot }
func (f fakeUpload) PublicDirs() []string { return f.publicDirs }
func (f fakeUpload) ParseUserID(tok string) (string, error) {
	if tok == "ok" {
		return "u@d", nil
	}
	return "", fmt.Errorf("bad token")
}
func (f fakeUpload) ValidateAccount(string, string, string) (bool, bool, error) {
	if f.allowWrite {
		return true, false, nil
	}
	return false, false, nil
}
func (f fakeUpload) ValidateApplication(string, string, string) (bool, bool, error) {
	return false, false, nil
}

func newMultipart(dir, filename, content string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	_ = w.WriteField("dir", dir)
	fw, _ := w.CreateFormFile("file", filename)
	_, _ = io.WriteString(fw, content)
	_ = w.Close()
	return body, w.FormDataContentType()
}

func TestUpload_DenyWithoutWrite_401(t *testing.T) {
	tmp := t.TempDir()
	p := fakeUpload{dataRoot: tmp, allowWrite: false}

	h := files.NewUploadFile(p)

	body, ctype := newMultipart("/users/alice", "note.txt", "hello")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/file-upload", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (body: %s)", rr.Code, rr.Body.String())
	}
}

func TestUpload_Success_201(t *testing.T) {
	tmp := t.TempDir()
	p := fakeUpload{dataRoot: tmp, allowWrite: true}

	h := files.NewUploadFile(p)

	body, ctype := newMultipart("/users/alice", "note.txt", "hello-world")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/file-upload", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body: %s)", rr.Code, rr.Body.String())
	}

	// file should be at <dataRoot>/files/users/alice/note.txt
	dst := filepath.Join(tmp, "files", "users", "alice", "note.txt")
	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("expected file written at %s: %v", dst, err)
	}
	if string(b) != "hello-world" {
		t.Fatalf("unexpected content: %q", string(b))
	}
}
