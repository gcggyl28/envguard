// Package sorter provides utilities for sorting and grouping .env file keys.
package sorter

import (
	"sort"
	"strings"
)

// GroupBy defines how keys should be grouped when sorting.
type GroupBy string

const (
	GroupByPrefix GroupBy = "prefix"
	GroupByAlpha  GroupBy = "alpha"
	GroupByNone   GroupBy = "none"
)

// Result holds the sorted key-value pairs and any grouping metadata.
type Result struct {
	Sorted map[string]string
	Order  []string
	Groups map[string][]string // prefix -> keys (only populated for GroupByPrefix)
}

// Sort returns a Result with keys ordered according to the given GroupBy strategy.
func Sort(env map[string]string, by GroupBy) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch by {
	case GroupByPrefix:
		return sortByPrefix(env, keys)
	case GroupByAlpha, GroupByNone:
		fallthrough
	default:
		sort.Strings(keys)
		return Result{Sorted: env, Order: keys, Groups: nil}
	}
}

func sortByPrefix(env map[string]string, keys []string) Result {
	groups := make(map[string][]string)

	for _, k := range keys {
		prefix := extractPrefix(k)
		groups[prefix] = append(groups[prefix], k)
	}

	// Sort within each group
	for p := range groups {
		sort.Strings(groups[p])
	}

	// Sort group prefixes
	prefixes := make([]string, 0, len(groups))
	for p := range groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	ordered := make([]string, 0, len(keys))
	for _, p := range prefixes {
		ordered = append(ordered, groups[p]...)
	}

	return Result{Sorted: env, Order: ordered, Groups: groups}
}

// extractPrefix returns the portion of the key before the first underscore,
// or the full key if no underscore is present.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
