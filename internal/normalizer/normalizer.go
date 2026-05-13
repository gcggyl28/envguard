package normalizer

import (
	"strings"
)

// Options controls normalization behavior.
type Options struct {
	UppercaseKeys   bool
	LowercaseValues bool
	TrimSpace       bool
	ReplaceHyphens  bool // replace hyphens in keys with underscores
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:  true,
		TrimSpace:      true,
		ReplaceHyphens: true,
	}
}

// Result holds the outcome of a normalization pass.
type Result struct {
	Normalized map[string]string
	Changes    []Change
}

// Change describes a single key or value that was altered.
type Change struct {
	Key      string
	OldKey   string
	OldValue string
	NewValue string
	Reason   string
}

// Normalize applies the given options to env and returns a Result.
func Normalize(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var changes []Change

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.TrimSpace {
			newKey = strings.TrimSpace(newKey)
			newVal = strings.TrimSpace(newVal)
		}
		if opts.ReplaceHyphens {
			newKey = strings.ReplaceAll(newKey, "-", "_")
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.LowercaseValues {
			newVal = strings.ToLower(newVal)
		}

		if newKey != k || newVal != v {
			reason := buildReason(k, newKey, v, newVal)
			changes = append(changes, Change{
				Key:      newKey,
				OldKey:   k,
				OldValue: v,
				NewValue: newVal,
				Reason:   reason,
			})
		}
		out[newKey] = newVal
	}

	return Result{Normalized: out, Changes: changes}
}

func buildReason(oldKey, newKey, oldVal, newVal string) string {
	var parts []string
	if oldKey != newKey {
		parts = append(parts, "key renamed")
	}
	if oldVal != newVal {
		parts = append(parts, "value changed")
	}
	return strings.Join(parts, ", ")
}
