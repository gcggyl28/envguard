// Package interpolator resolves variable references within .env values.
// It supports ${VAR} and $VAR syntax, resolving references from the same
// env map or falling back to OS environment variables.
package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// ErrCircularReference is returned when a circular variable reference is detected.
type ErrCircularReference struct {
	Key string
}

func (e *ErrCircularReference) Error() string {
	return fmt.Sprintf("circular reference detected for key: %s", e.Key)
}

// Interpolate resolves all variable references in the provided env map.
// Values referencing undefined variables are left as-is.
func Interpolate(env map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := resolve(k, v, env, map[string]bool{k: true})
		if err != nil {
			return nil, err
		}
		result[k] = resolved
	}
	return result, nil
}

func resolve(key, value string, env map[string]string, visiting map[string]bool) (string, error) {
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		refKey := extractKey(match)
		if visiting[refKey] {
			resolveErr = &ErrCircularReference{Key: key}
			return match
		}
		if refVal, ok := env[refKey]; ok {
			nextVisiting := copyVisiting(visiting)
			nextVisiting[refKey] = true
			resolved, err := resolve(refKey, refVal, env, nextVisiting)
			if err != nil {
				resolveErr = err
				return match
			}
			return resolved
		}
		if osVal, ok := os.LookupEnv(refKey); ok {
			return osVal
		}
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

func extractKey(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}

func copyVisiting(m map[string]bool) map[string]bool {
	copy := make(map[string]bool, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}
