package sanitizer

import (
	"strings"
	"unicode"
)

// Options controls which sanitization passes are applied.
type Options struct {
	TrimSpace      bool
	RemoveNewlines bool
	NormalizeKeys  bool // uppercase + replace hyphens with underscores
	StripNonPrint  bool
}

// DefaultOptions returns sensible sanitization defaults.
func DefaultOptions() Options {
	return Options{
		TrimSpace:      true,
		RemoveNewlines: true,
		NormalizeKeys:  false,
		StripNonPrint:  true,
	}
}

// Result holds the sanitized env map and a record of what changed.
type Result struct {
	Env     map[string]string
	Changed []Change
}

// Change records a single sanitization mutation.
type Change struct {
	Key    string
	Before string
	After  string
	Reason string
}

// Sanitize applies the given options to env and returns a Result.
func Sanitize(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var changes []Change

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.NormalizeKeys {
			nk := strings.ToUpper(strings.ReplaceAll(k, "-", "_"))
			if nk != k {
				newKey = nk
			}
		}

		if opts.TrimSpace {
			t := strings.TrimSpace(newVal)
			if t != newVal {
				changes = append(changes, Change{Key: newKey, Before: newVal, After: t, Reason: "trimmed whitespace"})
				newVal = t
			}
		}

		if opts.RemoveNewlines {
			r := strings.NewReplacer("\n", " ", "\r", "")
			t := r.Replace(newVal)
			if t != newVal {
				changes = append(changes, Change{Key: newKey, Before: newVal, After: t, Reason: "removed newlines"})
				newVal = t
			}
		}

		if opts.StripNonPrint {
			t := strings.Map(func(r rune) rune {
				if unicode.IsPrint(r) || r == '\t' {
					return r
				}
				return -1
			}, newVal)
			if t != newVal {
				changes = append(changes, Change{Key: newKey, Before: newVal, After: t, Reason: "stripped non-printable"})
				newVal = t
			}
		}

		if newKey != k {
			changes = append(changes, Change{Key: k, Before: k, After: newKey, Reason: "normalized key"})
		}

		out[newKey] = newVal
	}

	return Result{Env: out, Changed: changes}
}
