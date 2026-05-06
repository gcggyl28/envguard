// Package transformer provides utilities to apply key/value transformations
// to a parsed env map, such as uppercasing keys, trimming whitespace, or
// applying a prefix/suffix to all keys.
package transformer

import (
	"strings"
)

// Options controls which transformations are applied.
type Options struct {
	UppercaseKeys  bool
	TrimValues     bool
	KeyPrefix      string
	KeySuffix      string
	ReplaceKeys    map[string]string // old key -> new key
}

// Result holds the transformed map and a log of changes made.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Change records a single transformation applied to a key or value.
type Change struct {
	Key      string
	OldKey   string // non-empty when the key itself was renamed
	OldValue string
	NewValue string
}

// Transform applies the given Options to env and returns a Result.
// The original map is never mutated.
func Transform(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var changes []Change

	for k, v := range env {
		newKey := k
		newVal := v

		// Key rename table takes priority over other key transforms.
		if mapped, ok := opts.ReplaceKeys[k]; ok {
			newKey = mapped
		}

		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}

		if opts.KeyPrefix != "" {
			newKey = opts.KeyPrefix + newKey
		}

		if opts.KeySuffix != "" {
			newKey = newKey + opts.KeySuffix
		}

		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}

		if newKey != k || newVal != v {
			c := Change{Key: newKey, OldValue: v, NewValue: newVal}
			if newKey != k {
				c.OldKey = k
			}
			changes = append(changes, c)
		}

		out[newKey] = newVal
	}

	return Result{Env: out, Changes: changes}
}
