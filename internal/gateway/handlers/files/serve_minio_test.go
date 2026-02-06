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
		{"/users/alice/file.txt", "localhost/users/alice/file.txt", false},
		{"/users/bob/docs/report.pdf", "localhost/users/bob/docs/report.pdf", false},
		{"/users/../etc/passwd", "", true},
		{"/users/alice/../bob/file.txt", "", true},
		{"/users//alice/file", "localhost/users/alice/file", false},
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
	cfg := &MinioProxyConfig{}
	tests := []struct {
		reqPath    string
		host       string
		useHostPfx bool
		expected   string
		wantErr    bool
	}{
		{"/index.html", "globular.io", true, "globular.io/webroot/index.html", false},
		{"/index.html", "globular.io", false, "globular.io/webroot/index.html", false},
		{"/apps/test/main.js", "localhost", true, "localhost/webroot/apps/test/main.js", false},
	}

	for _, tt := range tests {
		cfg.UseHostPrefix = tt.useHostPfx
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
