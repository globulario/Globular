package files

import "strings"

// AuthDecision represents the result of an authorization check.
type AuthDecision int

const (
	// AuthDeny means access is explicitly denied
	AuthDeny AuthDecision = iota
	// AuthAllow means access is explicitly allowed
	AuthAllow
	// AuthRequireCheck means this rule doesn't apply, continue checking other rules
	AuthRequireCheck
)

// String returns a string representation of the decision.
func (d AuthDecision) String() string {
	switch d {
	case AuthDeny:
		return "DENY"
	case AuthAllow:
		return "ALLOW"
	case AuthRequireCheck:
		return "REQUIRE_CHECK"
	default:
		return "UNKNOWN"
	}
}

// AuthRule is a function that evaluates whether to allow/deny access to a path.
// Rules return:
//   - AuthAllow: Access is explicitly granted, stop checking other rules
//   - AuthDeny: Access is explicitly denied, stop checking other rules
//   - AuthRequireCheck: This rule doesn't apply, continue to next rule
type AuthRule func(path string) AuthDecision

// AuthorizationEngine evaluates a chain of authorization rules to determine access.
// Rules are evaluated in order. The first rule that returns Allow or Deny wins.
// If all rules return RequireCheck, access is denied by default.
type AuthorizationEngine struct {
	rules []AuthRule
}

// NewAuthorizationEngine creates an authorization engine with the given rules.
// Rules are evaluated in the order provided.
func NewAuthorizationEngine(rules ...AuthRule) *AuthorizationEngine {
	return &AuthorizationEngine{rules: rules}
}

// Decide evaluates all rules in order and returns the first Allow or Deny decision.
// If all rules return RequireCheck, returns Deny (default deny).
func (e *AuthorizationEngine) Decide(path string) AuthDecision {
	for _, rule := range e.rules {
		decision := rule(path)
		if decision != AuthRequireCheck {
			return decision
		}
	}
	// Default deny: if no rule explicitly allows, deny access
	return AuthDeny
}

// AuthRuleHiddenDirectories allows access to files in /.hidden/ directories.
// This is used for HLS streaming segments that should be accessible even in
// protected areas.
func AuthRuleHiddenDirectories(path string) AuthDecision {
	if strings.Contains(path, "/.hidden/") {
		return AuthAllow
	}
	return AuthRequireCheck
}

// AuthRuleHLSFiles allows access to HLS streaming files (.ts, resolution-specific m3u8).
// These files need to be accessible for video playback even in protected areas.
func AuthRuleHLSFiles(path string) AuthDecision {
	if isHLSFile(path) {
		return AuthAllow
	}
	return AuthRequireCheck
}

// AuthRuleUserDirectories denies access to /users/ paths unless explicitly allowed
// by subsequent rules (like token validation).
func AuthRuleUserDirectories(path string) AuthDecision {
	if strings.HasPrefix(path, "/users/") {
		return AuthRequireCheck // Requires additional checks (token, etc.)
	}
	return AuthRequireCheck
}

// AuthRulePublicDirectories returns a rule that allows access to specified public directories.
func AuthRulePublicDirectories(publicDirs []string) AuthRule {
	return func(path string) AuthDecision {
		if isPublicLike(path, publicDirs) {
			return AuthAllow
		}
		return AuthRequireCheck
	}
}

// AuthRuleTokenValidation returns a rule that validates access using a token.
// If a valid token is provided, it checks with the provider's ValidateAccount method.
func AuthRuleTokenValidation(provider ServeProvider, token string) AuthRule {
	return func(path string) AuthDecision {
		// No token = cannot validate
		if token == "" {
			return AuthRequireCheck
		}

		// Parse user ID from token
		uid, err := provider.ParseUserID(token)
		if err != nil || uid == "" {
			return AuthRequireCheck // Invalid token, continue to other rules
		}

		// Validate account access
		hasAccess, hasDenied, err := provider.ValidateAccount(uid, "read", path)
		if err != nil {
			return AuthRequireCheck // Error in validation, continue to other rules
		}

		if hasDenied {
			return AuthDeny // Explicitly denied
		}
		if hasAccess {
			return AuthAllow // Explicitly allowed
		}

		return AuthRequireCheck // No decision from this rule
	}
}

// AuthRuleApplicationValidation returns a rule that validates access using an application identifier.
func AuthRuleApplicationValidation(provider ServeProvider, app string) AuthRule {
	return func(path string) AuthDecision {
		// No app identifier = cannot validate
		if app == "" {
			return AuthRequireCheck
		}

		// Validate application access
		hasAccess, hasDenied, err := provider.ValidateApplication(app, "read", path)
		if err != nil {
			return AuthRequireCheck // Error in validation, continue to other rules
		}

		if hasDenied {
			return AuthDeny // Explicitly denied
		}
		if hasAccess {
			return AuthAllow // Explicitly allowed
		}

		return AuthRequireCheck // No decision from this rule
	}
}

// AuthRuleWebrootDefault allows access to non-protected paths (anything not under /users/).
// This implements the default-allow behavior for webroot files.
func AuthRuleWebrootDefault(path string) AuthDecision {
	// Protected paths require explicit authorization
	if strings.HasPrefix(path, "/users/") {
		return AuthRequireCheck
	}
	// Everything else (webroot, applications, templates, etc.) is allowed by default
	return AuthAllow
}

// BuildAuthorizationEngine creates the standard authorization engine for file serving.
// Rules are evaluated in this order:
//  1. Hidden directories (.hidden/) - ALLOW
//  2. HLS streaming files - ALLOW
//  3. Public directories - ALLOW (if in publicDirs)
//  4. Token validation - ALLOW/DENY based on account access
//  5. Application validation - ALLOW/DENY based on app access
//  6. Webroot default - ALLOW (for non-protected paths)
//  7. Default - DENY (if no rule explicitly allows)
func BuildAuthorizationEngine(provider ServeProvider, publicDirs []string, token, app string) *AuthorizationEngine {
	return NewAuthorizationEngine(
		AuthRuleHiddenDirectories,
		AuthRuleHLSFiles,
		AuthRulePublicDirectories(publicDirs),
		AuthRuleTokenValidation(provider, token),
		AuthRuleApplicationValidation(provider, app),
		AuthRuleWebrootDefault,
	)
}
