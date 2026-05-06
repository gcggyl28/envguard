package linter

import (
	"fmt"
	"strings"
)

// Hint represents a style or best-practice suggestion for a .env file entry.
type Hint struct {
	Key     string
	Message string
}

// Lint inspects the provided key-value pairs and returns style hints.
// It checks for common issues such as lowercase keys, keys with spaces,
// values that look like unquoted URLs, and suspiciously short secret values.
func Lint(env map[string]string) []Hint {
	var hints []Hint

	for key, value := range env {
		if key != strings.ToUpper(key) {
			hints = append(hints, Hint{
				Key:     key,
				Message: "key should be uppercase by convention",
			})
		}

		if strings.Contains(key, " ") {
			hints = append(hints, Hint{
				Key:     key,
				Message: "key contains spaces, which may cause parsing issues",
			})
		}

		if isSecretKey(key) && len(value) > 0 && len(value) < 16 {
			hints = append(hints, Hint{
				Key:     key,
				Message: fmt.Sprintf("secret value appears short (%d chars); consider a stronger value", len(value)),
			})
		}

		if looksLikeURL(value) && !isQuoted(value) {
			hints = append(hints, Hint{
				Key:     key,
				Message: "URL value should be quoted to avoid shell interpretation issues",
			})
		}
	}

	return hints
}

func isSecretKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, word := range []string{"SECRET", "PASSWORD", "TOKEN", "API_KEY", "PRIVATE"} {
		if strings.Contains(upper, word) {
			return true
		}
	}
	return false
}

func looksLikeURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}

func isQuoted(value string) bool {
	return (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'"))
}
