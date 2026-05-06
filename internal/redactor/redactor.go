// Package redactor masks sensitive values in env maps before output.
package redactor

import "strings"

const masked = "********"

// sensitivePatterns holds substrings that indicate a key holds a secret.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "APIKEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "SEED", "SALT",
}

// IsSensitive reports whether the given key name looks like it holds a secret.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Redact returns a copy of env where values whose keys are sensitive are
// replaced with the masked placeholder.
func Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitive(k) {
			out[k] = masked
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactList returns a copy of env where only keys present in sensitiveKeys
// are masked.  This allows callers to supply an explicit allowlist derived
// from a schema rather than relying on heuristics.
func RedactList(env map[string]string, sensitiveKeys []string) map[string]string {
	set := make(map[string]struct{}, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		set[k] = struct{}{}
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if _, ok := set[k]; ok {
			out[k] = masked
		} else {
			out[k] = v
		}
	}
	return out
}
