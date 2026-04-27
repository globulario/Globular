package files

import "testing"

// TestBuildMinioObjectKey_UserPaths tests user path key construction
func TestBuildMinioObjectKey_UserPaths(t *testing.T) {
	cfg := &MinioProxyConfig{}
	tests := []struct {
		reqPath  string
		expected string
		wantErr  bool
	}{
		{"/users/alice/file.txt", "users/alice/file.txt", false},
		{"/users/bob/docs/report.pdf", "users/bob/docs/report.pdf", false},
		{"/users/../etc/passwd", "", true},
		{"/users/alice/../bob/file.txt", "", true},
		{"/users//alice/file", "users/alice/file", false},
	}

	for _, tt := range tests {
		got, err := buildMinioObjectKey(tt.reqPath, "host", cfg, true)
		if tt.wantErr {
			if err == nil {
				t.Errorf("expected error for %q, got nil", tt.reqPath)
			}
			continue
		}
		if err != nil {
			t.Fatalf("buildMinioObjectKey(%q) error = %v", tt.reqPath, err)
		}
		if got != tt.expected {
			t.Errorf("buildMinioObjectKey(%q) = %q, want %q", tt.reqPath, got, tt.expected)
		}
	}
}

// TestBuildMinioObjectKey_WebrootPaths tests webroot key construction
func TestBuildMinioObjectKey_WebrootPaths(t *testing.T) {
	cfg := &MinioProxyConfig{Domain: "globular.internal"}
	tests := []struct {
		reqPath  string
		host     string
		expected string
		wantErr  bool
	}{
		// External domain → domain-scoped prefix
		{"/index.html", "globular.io", "globular.io/webroot/index.html", false},
		// localhost → cluster default ("webroot")
		{"/index.html", "localhost", "webroot/index.html", false},
		// Internal cluster subdomain → cluster default
		{"/index.html", "globule-ryzen.globular.internal", "webroot/index.html", false},
		// Internal cluster domain itself → cluster default
		{"/index.html", "globular.internal", "webroot/index.html", false},
		// Subdirectory paths
		{"/apps/test/main.js", "localhost", "webroot/apps/test/main.js", false},
	}

	for _, tt := range tests {
		got, err := buildMinioObjectKey(tt.reqPath, tt.host, cfg, false)
		if tt.wantErr {
			if err == nil {
				t.Errorf("expected error for %q, got nil", tt.reqPath)
			}
			continue
		}
		if err != nil {
			t.Fatalf("buildMinioObjectKey(%q) error = %v", tt.reqPath, err)
		}
		if got != tt.expected {
			t.Errorf("buildMinioObjectKey(%q) = %q, want %q", tt.reqPath, got, tt.expected)
		}
	}
}

// TestBuildMinioObjectKey_Sanitization tests path sanitization
func TestBuildMinioObjectKey_Sanitization(t *testing.T) {
	cfg := &MinioProxyConfig{}
	maliciousPaths := []string{
		"/users/../../../etc/passwd",
		"/users/./../../sensitive",
		"/../etc/shadow",
		"/users/%2e%2e%2f%2e%2e%2fetc%2fpasswd",
	}
	for _, p := range maliciousPaths {
		if _, err := buildMinioObjectKey(p, "host", cfg, false); err == nil {
			t.Errorf("expected error for malicious path %q, got nil", p)
		}
	}
}
