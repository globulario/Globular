package files_test

import (
	"bytes"
	"context"
	"encoding/json"
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
	minioCfg   *files.MinioProxyConfig
}

func (f fakeUpload) DataRoot() string { return f.dataRoot }

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
func (f fakeUpload) AddResourceOwner(path, resourceType, owner string) error { return nil }
func (f fakeUpload) FileServiceMinioConfig() (*files.MinioProxyConfig, bool) {
	if f.minioCfg == nil {
		return nil, false
	}
	return f.minioCfg, true
}

func newMultipart(dir, filename, content string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	_ = w.WriteField("dir", dir)
	fw, _ := w.CreateFormFile("multiplefiles", filename)
	_, _ = io.WriteString(fw, content)
	_ = w.Close()
	return body, w.FormDataContentType()
}

func TestUpload_DenyWithoutWrite_401(t *testing.T) {

	p := fakeUpload{allowWrite: false}

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
	p := fakeUpload{allowWrite: true}

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
	dst := filepath.Join("files", "users", "alice", "note.txt")
	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("expected file written at %s: %v", dst, err)
	}
	if string(b) != "hello-world" {
		t.Fatalf("unexpected content: %q", string(b))
	}
}

func TestUpload_MinioUsers(t *testing.T) {
	uploaded := struct {
		bucket string
		key    string
		data   string
	}{}
	cfg := &files.MinioProxyConfig{
		Bucket: "bucket",
		Prefix: "files",
		Put: func(ctx context.Context, bucket, key string, src io.Reader, size int64, contentType string) error {
			b, err := io.ReadAll(src)
			if err != nil {
				return err
			}
			uploaded.bucket = bucket
			uploaded.key = key
			uploaded.data = string(b)
			return nil
		},
	}

	tmp := t.TempDir()
	p := fakeUpload{allowWrite: true, minioCfg: cfg}

	h := files.NewUploadFile(p)

	body, ctype := newMultipart("/users/alice", "note.txt", "cloud-data")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/file-upload", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("token", "ok")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	if uploaded.bucket != "bucket" || uploaded.key != "files/users/alice/note.txt" {
		t.Fatalf("unexpected object uploaded: %#v", uploaded)
	}
	if uploaded.data != "cloud-data" {
		t.Fatalf("unexpected object data: %q", uploaded.data)
	}
	if _, err := os.Stat(filepath.Join(tmp, "files", "users", "alice", "note.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected no local file, stat err=%v", err)
	}
	var resp struct {
		Paths []string `json:"paths"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Paths) != 1 || resp.Paths[0] != "/users/alice/note.txt" {
		t.Fatalf("unexpected response paths: %#v", resp.Paths)
	}
}
