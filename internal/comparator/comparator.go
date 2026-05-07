package comparator

import "sort"

// Result holds the comparison outcome between two env maps.
type Result struct {
	Added    map[string]string // keys present in B but not A
	Removed  map[string]string // keys present in A but not B
	Changed  map[string]Change // keys present in both but with different values
	Unchanged []string          // keys present in both with identical values
}

// Change captures the before and after value for a modified key.
type Change struct {
	Old string
	New string
}

// Summary returns a human-readable one-liner for the comparison result.
func (r Result) Summary() string {
	if len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0 {
		return "environments are identical"
	}
	return "environments differ"
}

// SortedAdded returns added keys in alphabetical order.
func (r Result) SortedAdded() []string {
	return sortedKeys(r.Added)
}

// SortedRemoved returns removed keys in alphabetical order.
func (r Result) SortedRemoved() []string {
	return sortedKeys(r.Removed)
}

// SortedChanged returns changed keys in alphabetical order.
func (r Result) SortedChanged() []string {
	keys := make([]string, 0, len(r.Changed))
	for k := range r.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Compare compares two env maps (base vs target) and returns a Result.
func Compare(base, target map[string]string) Result {
	result := Result{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]Change),
	}

	for k, v := range target {
		if old, ok := base[k]; !ok {
			result.Added[k] = v
		} else if old != v {
			result.Changed[k] = Change{Old: old, New: v}
		} else {
			result.Unchanged = append(result.Unchanged, k)
		}
	}

	for k, v := range base {
		if _, ok := target[k]; !ok {
			result.Removed[k] = v
		}
	}

	sort.Strings(result.Unchanged)
	return result
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
