package files

import (
	"testing"
)

func TestUsersObjectKeyMapping(t *testing.T) {
	cfg := &MinioProxyConfig{Domain: "example.com"}
	key, err := usersObjectKey(cfg, "/users/a/b.txt")
	if err != nil {
		t.Fatalf("usersObjectKey error: %v", err)
	}
	if key != "users/a/b.txt" {
		t.Fatalf("unexpected key %q", key)
	}
}

func TestWebrootObjectKeyMapping(t *testing.T) {
	cfg := &MinioProxyConfig{Domain: "example.com"}
	key, err := webrootObjectKey(cfg, "globular.io", "/index.html")
	if err != nil {
		t.Fatalf("webrootObjectKey error: %v", err)
	}
	if key != "webroot/index.html" {
		t.Fatalf("unexpected key %q", key)
	}
}

func TestPathSanitizationRejectsTraversal(t *testing.T) {
	cfg := &MinioProxyConfig{Domain: "example.com"}
	if _, err := usersObjectKey(cfg, "/../secret"); err == nil {
		t.Fatalf("expected traversal error")
	}
	if key, err := webrootObjectKey(cfg, "globular.io", "//double//slash"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if key != "webroot/double/slash" {
		t.Fatalf("unexpected normalized key %q", key)
	}
}
