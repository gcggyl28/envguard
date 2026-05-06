// Package profiler analyzes .env files and produces a usage profile,
// summarizing key statistics such as total keys, required vs optional,
// sensitive keys, and keys with defaults.
package profiler

import "github.com/user/envguard/internal/schema"

// Profile holds aggregated statistics about an .env file relative to a schema.
type Profile struct {
	TotalKeys      int
	RequiredKeys   []string
	OptionalKeys   []string
	SensitiveKeys  []string
	KeysWithDefault []string
	UndeclaredKeys []string
}

// sensitivePatterns mirrors the logic used in the redactor package.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "PWD",
}

// Analyze builds a Profile from the provided env map and schema.
func Analyze(env map[string]string, s *schema.Schema) Profile {
	p := Profile{}

	declared := make(map[string]bool)
	for _, field := range s.Fields {
		declared[field.Key] = true
		p.TotalKeys++

		if field.Required {
			p.RequiredKeys = append(p.RequiredKeys, field.Key)
		} else {
			p.OptionalKeys = append(p.OptionalKeys, field.Key)
		}

		if field.Default != "" {
			p.KeysWithDefault = append(p.KeysWithDefault, field.Key)
		}

		if isSensitive(field.Key) {
			p.SensitiveKeys = append(p.SensitiveKeys, field.Key)
		}
	}

	for k := range env {
		if !declared[k] {
			p.UndeclaredKeys = append(p.UndeclaredKeys, k)
		}
	}

	return p
}

func isSensitive(key string) bool {
	for _, pat := range sensitivePatterns {
		if containsIgnoreCase(key, pat) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, sub string) bool {
	if len(sub) > len(s) {
		return false
	}
	sU := toUpper(s)
	return len(sU) >= len(sub) && (sU == sub || len(sU) > 0 && containsStr(sU, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func toUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		b[i] = c
	}
	return string(b)
}
