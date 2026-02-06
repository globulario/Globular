package files

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubServeProvider struct {
	publicDirs       []string
	parseUserIDFn    func(token string) (string, error)
	validateAcctFn   func(userID, action, reqPath string) (bool, bool, error)
	validateAppFn    func(app, action, reqPath string) (bool, bool, error)
	resolveImportFn  func(basePath, importLine string) (string, error)
	maybeStreamFn    func(name string, w http.ResponseWriter, r *http.Request) bool
	resolveProxyFn   func(reqPath string) (string, bool)
	fileServiceCfgFn func() (*MinioProxyConfig, error)
}

func (s *stubServeProvider) WebRoot() string                       { return "" }
func (s *stubServeProvider) DataRoot() string                      { return "" }
func (s *stubServeProvider) CredsDir() string                      { return "" }
func (s *stubServeProvider) IndexApplication() string              { return "" }
func (s *stubServeProvider) PublicDirs() []string                  { return s.publicDirs }
func (s *stubServeProvider) Exists(string) bool                    { return false }
func (s *stubServeProvider) FindHashedFile(string) (string, error) { return "", nil }
func (s *stubServeProvider) FileServiceMinioConfig() (*MinioProxyConfig, error) {
	if s.fileServiceCfgFn != nil {
		return s.fileServiceCfgFn()
	}
	return nil, nil
}
func (s *stubServeProvider) FileServiceMinioConfigStrict(ctx context.Context) (*MinioProxyConfig, error) {
	return s.FileServiceMinioConfig()
}
func (s *stubServeProvider) Mode() string { return "" }
func (s *stubServeProvider) ParseUserID(token string) (string, error) {
	if s.parseUserIDFn != nil {
		return s.parseUserIDFn(token)
	}
	return "", nil
}
func (s *stubServeProvider) ValidateAccount(userID, action, reqPath string) (bool, bool, error) {
	if s.validateAcctFn != nil {
		return s.validateAcctFn(userID, action, reqPath)
	}
	return false, false, nil
}
func (s *stubServeProvider) ValidateApplication(app, action, reqPath string) (bool, bool, error) {
	if s.validateAppFn != nil {
		return s.validateAppFn(app, action, reqPath)
	}
	return false, false, nil
}
func (s *stubServeProvider) ResolveImportPath(basePath, importLine string) (string, error) {
	if s.resolveImportFn != nil {
		return s.resolveImportFn(basePath, importLine)
	}
	return "", nil
}
func (s *stubServeProvider) MaybeStream(name string, w http.ResponseWriter, r *http.Request) bool {
	if s.maybeStreamFn != nil {
		return s.maybeStreamFn(name, w, r)
	}
	return false
}
func (s *stubServeProvider) ResolveProxy(reqPath string) (string, bool) {
	if s.resolveProxyFn != nil {
		return s.resolveProxyFn(reqPath)
	}
	return "", false
}

func TestAuthorizationHandler_AllowsWhenEngineAllows(t *testing.T) {
	h := &AuthorizationHandler{
		EngineFactory: func() *AuthorizationEngine {
			return NewAuthorizationEngine(func(string) AuthDecision { return AuthAllow })
		},
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/files", nil)
	allowed := h.Authorize(rr, req, &pathInfo{reqPath: "/files"})
	if !allowed {
		t.Fatalf("expected allow, got deny")
	}
	if rr.Code != 0 && rr.Code != http.StatusOK {
		t.Fatalf("expected no response written, got status %d", rr.Code)
	}
}

func TestAuthorizationHandler_DeniesWithResponse(t *testing.T) {
	h := &AuthorizationHandler{
		EngineFactory: func() *AuthorizationEngine {
			return NewAuthorizationEngine(func(string) AuthDecision { return AuthDeny })
		},
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	allowed := h.Authorize(rr, req, &pathInfo{reqPath: "/secret"})
	if allowed {
		t.Fatalf("expected deny, got allow")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthorizationHandler_UsesTokenRule(t *testing.T) {
	sp := &stubServeProvider{
		parseUserIDFn: func(token string) (string, error) {
			if token == "valid" {
				return "user@domain", nil
			}
			return "", nil
		},
		validateAcctFn: func(userID, action, reqPath string) (bool, bool, error) {
			if userID == "user@domain" && action == "read" && reqPath == "/users/alice/file.txt" {
				return true, false, nil
			}
			return false, false, nil
		},
	}

	h := &AuthorizationHandler{
		Provider:   sp,
		PublicDirs: []string{},
		Token:      "valid",
		App:        "",
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/alice/file.txt", nil)
	allowed := h.Authorize(rr, req, &pathInfo{reqPath: "/users/alice/file.txt"})
	if !allowed {
		t.Fatalf("expected allow based on token validation")
	}
}
