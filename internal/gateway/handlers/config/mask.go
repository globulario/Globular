package config

import (
	"strings"
)

// secretKeys lists config field names that contain credentials.
// Matching is case-insensitive and checked with HasSuffix/Contains.
var secretKeys = []string{
	"SecretKey", "MinioSecretKey", "ResticPassword",
	"OAuth2ClientSecret", "Password", "Secret",
	"ApiSecret", "ApiToken", "AccessToken",
}

// secretOptionKeys are keys inside Options/Credentials maps that are secrets.
var secretOptionKeys = []string{
	"secret_key", "password", "api_secret", "api_token",
	"secret_access_key", "access_token",
}

const maskPrefix = "****"

// IsSecretField returns true if a config field name likely holds a secret value.
func IsSecretField(key string) bool {
	lower := strings.ToLower(key)
	for _, sk := range secretKeys {
		if strings.ToLower(sk) == lower {
			return true
		}
	}
	// Fallback heuristic: field name ends with "secret", "password", or "token"
	// but NOT "endpoint", "path", "repo", "location".
	if strings.HasSuffix(lower, "secret") || strings.HasSuffix(lower, "password") {
		return true
	}
	if strings.HasSuffix(lower, "secretkey") {
		return true
	}
	return false
}

// isSecretOptionKey returns true for keys inside Options/credentials maps.
func isSecretOptionKey(key string) bool {
	lower := strings.ToLower(key)
	for _, sk := range secretOptionKeys {
		if lower == sk {
			return true
		}
	}
	return strings.Contains(lower, "secret") || strings.Contains(lower, "password")
}

// MaskValue fully masks a secret value, revealing nothing about its contents
// or length. Returns "" for empty input so callers can distinguish "unset".
//
// SECURITY: this function is used ONLY to mask secret fields (MaskConfigSecrets /
// maskDestination — secretKeys like OAuth2ClientSecret, Password, SecretKey,
// ApiToken, …). It previously returned maskPrefix+v[:2]+"..." which leaked the
// first TWO characters of every secret through the config API — exactly what the
// MasksSecret test was defending against. A partial "hint" mask is never safe for
// a credential, so it now emits a constant full mask. The result still begins with
// maskPrefix so IsMaskedValue (used by UnmaskPatch to avoid overwriting secrets)
// keeps working.
func MaskValue(v string) string {
	if v == "" {
		return ""
	}
	return "********"
}

// IsMaskedValue returns true if the value looks like it was masked by us.
func IsMaskedValue(v string) bool {
	return strings.HasPrefix(v, maskPrefix)
}

// MaskConfigSecrets returns a shallow copy of the config with secret values masked.
// It handles top-level string fields, and also walks Destinations[].Options maps.
func MaskConfigSecrets(cfg map[string]any) map[string]any {
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		switch val := v.(type) {
		case string:
			if IsSecretField(k) && val != "" {
				out[k] = MaskValue(val)
			} else {
				out[k] = v
			}
		case []any:
			// Handle Destinations array — each element may have an Options map
			masked := make([]any, len(val))
			for i, item := range val {
				if m, ok := item.(map[string]any); ok {
					masked[i] = maskDestination(m)
				} else {
					masked[i] = item
				}
			}
			out[k] = masked
		default:
			out[k] = v
		}
	}
	return out
}

// maskDestination masks secret fields inside a destination's Options map.
func maskDestination(dest map[string]any) map[string]any {
	out := make(map[string]any, len(dest))
	for k, v := range dest {
		if k == "Options" {
			if opts, ok := v.(map[string]any); ok {
				maskedOpts := make(map[string]any, len(opts))
				for ok2, ov := range opts {
					if s, isStr := ov.(string); isStr && isSecretOptionKey(ok2) && s != "" {
						maskedOpts[ok2] = MaskValue(s)
					} else {
						maskedOpts[ok2] = ov
					}
				}
				out[k] = maskedOpts
				continue
			}
		}
		out[k] = v
	}
	return out
}

// UnmaskPatch takes a patch (from a save request) and a current config,
// and replaces any masked values in the patch with the real values from current.
// This prevents the UI from accidentally overwriting secrets with "****..." strings.
func UnmaskPatch(patch, current map[string]any) {
	if current == nil {
		return
	}
	for k, v := range patch {
		if s, ok := v.(string); ok && IsMaskedValue(s) {
			if orig, exists := current[k]; exists {
				patch[k] = orig
			}
		}
	}
}
