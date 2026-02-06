package files

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

// mockServeProvider implements ServeProvider for testing
type mockServeProvider struct {
	validateAccountFunc     func(userID, action, reqPath string) (bool, bool, error)
	validateApplicationFunc func(app, action, reqPath string) (bool, bool, error)
	parseUserIDFunc         func(token string) (string, error)
}

func (m *mockServeProvider) ParseUserID(token string) (string, error) {
	if m.parseUserIDFunc != nil {
		return m.parseUserIDFunc(token)
	}
	return "user123", nil
}

func (m *mockServeProvider) ValidateAccount(userID, action, reqPath string) (bool, bool, error) {
	if m.validateAccountFunc != nil {
		return m.validateAccountFunc(userID, action, reqPath)
	}
	return true, false, nil
}

func (m *mockServeProvider) ValidateApplication(app, action, reqPath string) (bool, bool, error) {
	if m.validateApplicationFunc != nil {
		return m.validateApplicationFunc(app, action, reqPath)
	}
	return true, false, nil
}

// Stub implementations for other ServeProvider methods
func (m *mockServeProvider) WebRoot() string                                    { return "/webroot" }
func (m *mockServeProvider) DataRoot() string                                   { return "/data" }
func (m *mockServeProvider) CredsDir() string                                   { return "/creds" }
func (m *mockServeProvider) PublicDirs() []string                               { return []string{"/public"} }
func (m *mockServeProvider) IndexApplication() string                           { return "" }
func (m *mockServeProvider) Exists(p string) bool                               { return false }
func (m *mockServeProvider) FindHashedFile(p string) (string, error)            { return "", nil }
func (m *mockServeProvider) FileServiceMinioConfig() (*MinioProxyConfig, error) { return nil, nil }
func (m *mockServeProvider) FileServiceMinioConfigStrict(ctx context.Context) (*MinioProxyConfig, error) {
	return nil, nil
}
func (m *mockServeProvider) Mode() string { return "" }
func (m *mockServeProvider) ResolveImportPath(basePath string, importLine string) (string, error) {
	return "", nil
}
func (m *mockServeProvider) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	return false
}
func (m *mockServeProvider) ResolveProxy(reqPath string) (string, bool) { return "", false }

// TestAuthDecision_String tests the String method
func TestAuthDecision_String(t *testing.T) {
	tests := []struct {
		decision AuthDecision
		expected string
	}{
		{AuthDeny, "DENY"},
		{AuthAllow, "ALLOW"},
		{AuthRequireCheck, "REQUIRE_CHECK"},
	}

	for _, tt := range tests {
		if got := tt.decision.String(); got != tt.expected {
			t.Errorf("AuthDecision(%d).String() = %s, want %s", tt.decision, got, tt.expected)
		}
	}
}

// TestAuthRuleHiddenDirectories tests the hidden directory rule
func TestAuthRuleHiddenDirectories(t *testing.T) {
	tests := []struct {
		path     string
		expected AuthDecision
	}{
		{"/users/alice/.hidden/video.mp4", AuthAllow},
		{"/data/.hidden/file.txt", AuthAllow},
		{"/.hidden/secret", AuthAllow},
		{"/users/alice/video.mp4", AuthRequireCheck},
		{"/public/file.txt", AuthRequireCheck},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := AuthRuleHiddenDirectories(tt.path)
			if got != tt.expected {
				t.Errorf("AuthRuleHiddenDirectories(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestAuthRuleHLSFiles tests the HLS file rule
func TestAuthRuleHLSFiles(t *testing.T) {
	tests := []struct {
		path     string
		expected AuthDecision
	}{
		{"/videos/segment.ts", AuthAllow},
		{"/videos/720p.m3u8", AuthAllow},
		{"/videos/1080p.m3u8", AuthAllow},
		{"/videos/video.mp4", AuthRequireCheck},
		{"/documents/file.txt", AuthRequireCheck},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := AuthRuleHLSFiles(tt.path)
			if got != tt.expected {
				t.Errorf("AuthRuleHLSFiles(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestAuthRuleUserDirectories tests the user directory rule
func TestAuthRuleUserDirectories(t *testing.T) {
	tests := []struct {
		path     string
		expected AuthDecision
	}{
		{"/users/alice/file.txt", AuthRequireCheck},
		{"/users/bob/video.mp4", AuthRequireCheck},
		{"/public/file.txt", AuthRequireCheck},
		{"/webroot/index.html", AuthRequireCheck},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := AuthRuleUserDirectories(tt.path)
			if got != tt.expected {
				t.Errorf("AuthRuleUserDirectories(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestAuthRulePublicDirectories tests the public directory rule
func TestAuthRulePublicDirectories(t *testing.T) {
	publicDirs := []string{"/public", "/shared"}
	rule := AuthRulePublicDirectories(publicDirs)

	tests := []struct {
		path     string
		expected AuthDecision
	}{
		{"/public/file.txt", AuthAllow},
		{"/shared/document.pdf", AuthAllow},
		{"/users/alice/file.txt", AuthRequireCheck},
		{"/private/secret.txt", AuthRequireCheck},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := rule(tt.path)
			if got != tt.expected {
				t.Errorf("AuthRulePublicDirectories(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestAuthRuleTokenValidation tests the token validation rule
func TestAuthRuleTokenValidation(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		path     string
		provider *mockServeProvider
		expected AuthDecision
	}{
		{
			name:  "no_token",
			token: "",
			path:  "/users/alice/file.txt",
			provider: &mockServeProvider{
				parseUserIDFunc: func(token string) (string, error) {
					return "", nil
				},
			},
			expected: AuthRequireCheck,
		},
		{
			name:  "valid_token_allow",
			token: "valid_token",
			path:  "/users/alice/file.txt",
			provider: &mockServeProvider{
				parseUserIDFunc: func(token string) (string, error) {
					return "alice", nil
				},
				validateAccountFunc: func(userID, action, reqPath string) (bool, bool, error) {
					return true, false, nil // has access, not denied
				},
			},
			expected: AuthAllow,
		},
		{
			name:  "valid_token_deny",
			token: "valid_token",
			path:  "/users/bob/file.txt",
			provider: &mockServeProvider{
				parseUserIDFunc: func(token string) (string, error) {
					return "alice", nil
				},
				validateAccountFunc: func(userID, action, reqPath string) (bool, bool, error) {
					return false, true, nil // no access, explicitly denied
				},
			},
			expected: AuthDeny,
		},
		{
			name:  "valid_token_no_decision",
			token: "valid_token",
			path:  "/users/charlie/file.txt",
			provider: &mockServeProvider{
				parseUserIDFunc: func(token string) (string, error) {
					return "alice", nil
				},
				validateAccountFunc: func(userID, action, reqPath string) (bool, bool, error) {
					return false, false, nil // no access, not denied (neutral)
				},
			},
			expected: AuthRequireCheck,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := AuthRuleTokenValidation(tt.provider, tt.token)
			got := rule(tt.path)
			if got != tt.expected {
				t.Errorf("AuthRuleTokenValidation(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestAuthRuleWebrootDefault tests the webroot default allow rule
func TestAuthRuleWebrootDefault(t *testing.T) {
	tests := []struct {
		path     string
		expected AuthDecision
	}{
		// Protected paths - require check
		{"/users/alice/file.txt", AuthRequireCheck},
		{"/users/bob/video.mp4", AuthRequireCheck},

		// Webroot paths - allow by default
		{"/index.html", AuthAllow},
		{"/styles/main.css", AuthAllow},
		{"/applications/app/index.html", AuthAllow},
		{"/templates/page.html", AuthAllow},
		{"/projects/proj/file.txt", AuthAllow},
		{"/public/file.txt", AuthAllow},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := AuthRuleWebrootDefault(tt.path)
			if got != tt.expected {
				t.Errorf("AuthRuleWebrootDefault(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestAuthRuleApplicationValidation tests the application validation rule
func TestAuthRuleApplicationValidation(t *testing.T) {
	tests := []struct {
		name     string
		app      string
		path     string
		provider *mockServeProvider
		expected AuthDecision
	}{
		{
			name: "no_app",
			app:  "",
			path: "/apps/test/file.txt",
			provider: &mockServeProvider{
				validateApplicationFunc: func(app, action, reqPath string) (bool, bool, error) {
					return false, false, nil
				},
			},
			expected: AuthRequireCheck,
		},
		{
			name: "valid_app_allow",
			app:  "test_app",
			path: "/apps/test/file.txt",
			provider: &mockServeProvider{
				validateApplicationFunc: func(app, action, reqPath string) (bool, bool, error) {
					return true, false, nil // has access
				},
			},
			expected: AuthAllow,
		},
		{
			name: "valid_app_deny",
			app:  "test_app",
			path: "/apps/restricted/file.txt",
			provider: &mockServeProvider{
				validateApplicationFunc: func(app, action, reqPath string) (bool, bool, error) {
					return false, true, nil // explicitly denied
				},
			},
			expected: AuthDeny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := AuthRuleApplicationValidation(tt.provider, tt.app)
			got := rule(tt.path)
			if got != tt.expected {
				t.Errorf("AuthRuleApplicationValidation(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestAuthorizationEngine_Decide tests the engine's decision logic
func TestAuthorizationEngine_Decide(t *testing.T) {
	tests := []struct {
		name     string
		rules    []AuthRule
		path     string
		expected AuthDecision
	}{
		{
			name: "first_rule_allows",
			rules: []AuthRule{
				func(path string) AuthDecision { return AuthAllow },
				func(path string) AuthDecision { return AuthDeny },
			},
			path:     "/test",
			expected: AuthAllow,
		},
		{
			name: "first_rule_denies",
			rules: []AuthRule{
				func(path string) AuthDecision { return AuthDeny },
				func(path string) AuthDecision { return AuthAllow },
			},
			path:     "/test",
			expected: AuthDeny,
		},
		{
			name: "first_rule_requires_check_second_allows",
			rules: []AuthRule{
				func(path string) AuthDecision { return AuthRequireCheck },
				func(path string) AuthDecision { return AuthAllow },
			},
			path:     "/test",
			expected: AuthAllow,
		},
		{
			name: "all_require_check_default_deny",
			rules: []AuthRule{
				func(path string) AuthDecision { return AuthRequireCheck },
				func(path string) AuthDecision { return AuthRequireCheck },
			},
			path:     "/test",
			expected: AuthDeny,
		},
		{
			name:     "no_rules_default_deny",
			rules:    []AuthRule{},
			path:     "/test",
			expected: AuthDeny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewAuthorizationEngine(tt.rules...)
			got := engine.Decide(tt.path)
			if got != tt.expected {
				t.Errorf("engine.Decide(%q) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestBuildAuthorizationEngine tests the standard engine builder
func TestBuildAuthorizationEngine(t *testing.T) {
	provider := &mockServeProvider{
		parseUserIDFunc: func(token string) (string, error) {
			if token == "valid_token" {
				return "user123", nil
			}
			return "", nil
		},
		validateAccountFunc: func(userID, action, reqPath string) (bool, bool, error) {
			// Allow access to /users/user123/
			if userID == "user123" && strings.HasPrefix(reqPath, "/users/user123/") {
				return true, false, nil
			}
			return false, false, nil
		},
	}

	tests := []struct {
		name       string
		path       string
		token      string
		app        string
		publicDirs []string
		expected   AuthDecision
		desc       string
	}{
		{
			name:       "hidden_directory",
			path:       "/users/alice/.hidden/video.mp4",
			token:      "",
			app:        "",
			publicDirs: []string{},
			expected:   AuthAllow,
			desc:       "Hidden directories should be allowed",
		},
		{
			name:       "hls_file",
			path:       "/videos/720p.m3u8",
			token:      "",
			app:        "",
			publicDirs: []string{},
			expected:   AuthAllow,
			desc:       "HLS files should be allowed",
		},
		{
			name:       "public_directory",
			path:       "/public/file.txt",
			token:      "",
			app:        "",
			publicDirs: []string{"/public"},
			expected:   AuthAllow,
			desc:       "Public directories should be allowed",
		},
		{
			name:       "valid_token_own_file",
			path:       "/users/user123/document.pdf",
			token:      "valid_token",
			app:        "",
			publicDirs: []string{},
			expected:   AuthAllow,
			desc:       "User should access their own files",
		},
		{
			name:       "no_token_private_file",
			path:       "/users/alice/private.txt",
			token:      "",
			app:        "",
			publicDirs: []string{},
			expected:   AuthDeny,
			desc:       "Private files without token should be denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := BuildAuthorizationEngine(provider, tt.publicDirs, tt.token, tt.app)
			got := engine.Decide(tt.path)
			if got != tt.expected {
				t.Errorf("%s: engine.Decide(%q) = %s, want %s", tt.desc, tt.path, got, tt.expected)
			} else {
				t.Logf("✓ %s", tt.desc)
			}
		})
	}
}
