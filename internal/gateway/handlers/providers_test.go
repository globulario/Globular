package handlers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadMinioCredentialsFilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "creds")
	if err := os.WriteFile(path, []byte("ak:sk"), 0o644); err != nil {
		t.Fatalf("write cred file: %v", err)
	}
	if _, _, err := readMinioCredentialsFile(path); !errors.Is(err, ErrObjectStoreUnavailable) {
		t.Fatalf("expected permission error, got %v", err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	ak, sk, err := readMinioCredentialsFile(path)
	if err != nil {
		t.Fatalf("expected success after chmod, got %v", err)
	}
	if ak != "ak" || sk != "sk" {
		t.Fatalf("unexpected credentials %q/%q", ak, sk)
	}
}

func TestReadMinioCredentialsFileMissing(t *testing.T) {
	if _, _, err := readMinioCredentialsFile("/does/not/exist"); !errors.Is(err, ErrObjectStoreUnavailable) {
		t.Fatalf("expected ErrObjectStoreUnavailable, got %v", err)
	}
}
