package flattener

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how nested key structures are flattened.
type Options struct {
	// Separator is the delimiter used to join nested key segments (default: "__").
	Separator string
	// Prefix is an optional prefix to prepend to all output keys.
	Prefix string
	// Uppercase forces all output keys to uppercase.
	Uppercase bool
}

// DefaultOptions returns sensible defaults for flattening.
func DefaultOptions() Options {
	return Options{
		Separator: "__",
	}
}

// Result holds the outcome of a flatten operation.
type Result struct {
	// Flattened is the resulting key=value map.
	Flattened map[string]string
	// Renamed tracks keys that were transformed (original -> new).
	Renamed map[string]string
}

// Flatten takes an env map whose keys may use a nested separator convention
// (e.g. "DB__HOST" or "APP_DB_HOST") and normalises them according to opts.
// Keys that are already in the desired form are kept as-is.
// When opts.Prefix is set, only keys that begin with that prefix are included
// and the prefix is stripped from the output key.
func Flatten(env map[string]string, opts Options) Result {
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	flattened := make(map[string]string, len(env))
	renamed := make(map[string]string)

	keys := sortedKeys(env)
	for _, k := range keys {
		v := env[k]

		outKey := k

		// Strip prefix if requested.
		if opts.Prefix != "" {
			pfx := opts.Prefix
			if opts.Uppercase {
				pfx = strings.ToUpper(pfx)
			}
			if !strings.HasPrefix(outKey, pfx) {
				continue
			}
			outKey = strings.TrimPrefix(outKey, pfx)
			outKey = strings.TrimPrefix(outKey, opts.Separator)
		}

		if opts.Uppercase {
			outKey = strings.ToUpper(outKey)
		}

		// Resolve collisions by appending an index suffix.
		if _, exists := flattened[outKey]; exists {
			outKey = fmt.Sprintf("%s%s1", outKey, opts.Separator)
		}

		if outKey != k {
			renamed[k] = outKey
		}
		flattened[outKey] = v
	}

	return Result{
		Flattened: flattened,
		Renamed:   renamed,
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
