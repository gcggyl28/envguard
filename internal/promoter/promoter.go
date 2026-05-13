// Package promoter handles promoting env values from one environment tier to another
// (e.g. staging → production), applying allow/deny rules and optional value transforms.
package promoter

import "strings"

// Options controls promotion behaviour.
type Options struct {
	// AllowKeys, if non-empty, restricts promotion to only these keys.
	AllowKeys []string
	// DenyKeys lists keys that must never be promoted.
	DenyKeys []string
	// Overwrite controls whether existing destination keys are replaced.
	Overwrite bool
}

// Result holds the outcome of a single key promotion.
type Result struct {
	Key       string
	Value     string
	Skipped   bool
	Reason    string
}

// Summary aggregates promotion statistics.
type Summary struct {
	Promoted int
	Skipped  int
}

// Promote copies keys from src into dst according to opts.
// It returns per-key results and an aggregate summary.
func Promote(src, dst map[string]string, opts Options) ([]Result, Summary) {
	allowSet := toSet(opts.AllowKeys)
	denySet := toSet(opts.DenyKeys)

	var results []Result
	var sum Summary

	for k, v := range src {
		r := Result{Key: k, Value: v}

		if len(allowSet) > 0 && !allowSet[strings.ToUpper(k)] {
			r.Skipped = true
			r.Reason = "not in allow list"
			results = append(results, r)
			sum.Skipped++
			continue
		}

		if denySet[strings.ToUpper(k)] {
			r.Skipped = true
			r.Reason = "in deny list"
			results = append(results, r)
			sum.Skipped++
			continue
		}

		if _, exists := dst[k]; exists && !opts.Overwrite {
			r.Skipped = true
			r.Reason = "already exists in destination"
			results = append(results, r)
			sum.Skipped++
			continue
		}

		dst[k] = v
		results = append(results, r)
		sum.Promoted++
	}

	return results, sum
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[strings.ToUpper(k)] = true
	}
	return s
}
