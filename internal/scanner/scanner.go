package scanner

import (
	"fmt"
	"strings"
)

// Finding represents a single scan result for a key-value pair.
type Finding struct {
	Key      string
	Value    string
	Severity string // "error", "warn", "info"
	Message  string
}

// Options controls which checks the scanner performs.
type Options struct {
	CheckEmptyValues   bool
	CheckDuplicateKeys bool
	CheckLongValues    bool
	MaxValueLength     int
}

// DefaultOptions returns a sensible default set of scan options.
func DefaultOptions() Options {
	return Options{
		CheckEmptyValues:   true,
		CheckDuplicateKeys: true,
		CheckLongValues:    true,
		MaxValueLength:     256,
	}
}

// Scan inspects the provided env map for common issues and returns findings.
func Scan(env map[string]string, opts Options) []Finding {
	var findings []Finding
	seen := make(map[string]int)

	for key, value := range env {
		seen[key]++

		if opts.CheckEmptyValues && strings.TrimSpace(value) == "" {
			findings = append(findings, Finding{
				Key:      key,
				Value:    value,
				Severity: "warn",
				Message:  "value is empty or whitespace-only",
			})
		}

		if opts.CheckLongValues && len(value) > opts.MaxValueLength {
			findings = append(findings, Finding{
				Key:      key,
				Value:    value[:32] + "...",
				Severity: "info",
				Message:  fmt.Sprintf("value exceeds %d characters (%d)", opts.MaxValueLength, len(value)),
			})
		}
	}

	if opts.CheckDuplicateKeys {
		for key, count := range seen {
			if count > 1 {
				findings = append(findings, Finding{
					Key:      key,
					Severity: "error",
					Message:  fmt.Sprintf("key appears %d times", count),
				})
			}
		}
	}

	return findings
}

// CountBySeverity returns a map of severity -> count from a findings slice.
func CountBySeverity(findings []Finding) map[string]int {
	counts := map[string]int{"error": 0, "warn": 0, "info": 0}
	for _, f := range findings {
		counts[f.Severity]++
	}
	return counts
}
