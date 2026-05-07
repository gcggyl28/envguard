// Package trimmer removes leading/trailing whitespace from env values
// and optionally strips surrounding quotes.
package trimmer

import "strings"

// Options controls trimmer behaviour.
type Options struct {
	// TrimValues removes leading and trailing whitespace from values.
	TrimValues bool
	// TrimKeys removes leading and trailing whitespace from keys.
	TrimKeys bool
	// StripQuotes removes a single layer of matching surrounding quotes.
	StripQuotes bool
}

// DefaultOptions returns sensible defaults: trim values and strip quotes.
func DefaultOptions() Options {
	return Options{
		TrimValues:  true,
		TrimKeys:    false,
		StripQuotes: true,
	}
}

// Result holds the outcome of a Trim operation.
type Result struct {
	// Trimmed is the cleaned env map.
	Trimmed map[string]string
	// Changes lists keys whose values (or keys) were modified.
	Changes []string
}

// Trim applies the given options to env and returns a Result.
// The original map is never mutated.
func Trim(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	changed := map[string]bool{}

	for k, v := range env {
		newKey := k
		if opts.TrimKeys {
			newKey = strings.TrimSpace(k)
			if newKey != k {
				changed[newKey] = true
			}
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}
		if opts.StripQuotes {
			newVal = stripQuotes(newVal)
		}
		if newVal != v {
			changed[newKey] = true
		}

		out[newKey] = newVal
	}

	keys := make([]string, 0, len(changed))
	for k := range changed {
		keys = append(keys, k)
	}
	sortStrings(keys)

	return Result{Trimmed: out, Changes: keys}
}

func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
