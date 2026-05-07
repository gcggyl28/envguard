// Package masker provides utilities for partially masking env variable values
// for safe display in logs, reports, and audit trails.
package masker

import "strings"

// Style controls how a value is masked.
type Style string

const (
	StyleFull    Style = "full"    // replace entire value with asterisks
	StylePartial Style = "partial" // reveal first/last N chars
	StyleHash    Style = "hash"    // show length hint only
)

// Options configures masking behaviour.
type Options struct {
	Style       Style
	RevealChars int // used by StylePartial
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Style:       StylePartial,
		RevealChars: 3,
	}
}

// Mask applies the given options to a single value string.
func Mask(value string, opts Options) string {
	if value == "" {
		return ""
	}
	switch opts.Style {
	case StyleFull:
		return strings.Repeat("*", len(value))
	case StyleHash:
		return maskHash(value)
	default:
		return maskPartial(value, opts.RevealChars)
	}
}

// MaskMap applies masking to every value in the provided map, returning a new map.
func MaskMap(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Mask(v, opts)
	}
	return out
}

func maskPartial(value string, reveal int) string {
	if reveal <= 0 || len(value) <= reveal*2 {
		return strings.Repeat("*", len(value))
	}
	prefix := value[:reveal]
	suffix := value[len(value)-reveal:]
	midLen := len(value) - reveal*2
	if midLen < 1 {
		midLen = 1
	}
	return prefix + strings.Repeat("*", midLen) + suffix
}

func maskHash(value string) string {
	return "[" + strings.Repeat("*", len(value)) + "]" // e.g. [*****]
}
