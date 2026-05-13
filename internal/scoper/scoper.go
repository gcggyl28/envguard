// Package scoper filters env variables to a named deployment scope (e.g. "production", "staging").
package scoper

import "strings"

// Scope represents a named deployment environment scope.
type Scope struct {
	Name     string
	Prefixes []string // e.g. ["PROD_", "PRODUCTION_"]
}

// Result holds the outcome of scoping an env map.
type Result struct {
	Scope    string
	Included map[string]string // keys that matched the scope
	Excluded map[string]string // keys that did not match
}

// Apply filters env to keys that match any of the scope's prefixes.
// If no prefixes are defined, all keys are included.
func Apply(env map[string]string, scope Scope) Result {
	result := Result{
		Scope:    scope.Name,
		Included: make(map[string]string),
		Excluded: make(map[string]string),
	}

	if len(scope.Prefixes) == 0 {
		for k, v := range env {
			result.Included[k] = v
		}
		return result
	}

	for k, v := range env {
		if matchesAnyPrefix(k, scope.Prefixes) {
			result.Included[k] = v
		} else {
			result.Excluded[k] = v
		}
	}
	return result
}

// Strip returns a new env map with scope prefixes removed from matched keys.
func Strip(env map[string]string, scope Scope) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		stripped := stripPrefix(k, scope.Prefixes)
		out[stripped] = v
	}
	return out
}

func matchesAnyPrefix(key string, prefixes []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range prefixes {
		if strings.HasPrefix(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func stripPrefix(key string, prefixes []string) string {
	upper := strings.ToUpper(key)
	for _, p := range prefixes {
		pu := strings.ToUpper(p)
		if strings.HasPrefix(upper, pu) {
			return key[len(p):]
		}
	}
	return key
}
