package deduplicator

// Result holds the outcome of a deduplication pass.
type Result struct {
	// Unique is the deduplicated env map (last-write-wins by default).
	Unique map[string]string
	// Duplicates maps each key that appeared more than once to all of its
	// observed values in the order they were encountered.
	Duplicates map[string][]string
}

// Strategy controls which occurrence is kept when a duplicate is found.
type Strategy int

const (
	// KeepFirst retains the first occurrence of a duplicated key.
	KeepFirst Strategy = iota
	// KeepLast retains the last occurrence of a duplicated key (default).
	KeepLast
)

// StrategyFromString parses a strategy name, defaulting to KeepLast.
func StrategyFromString(s string) Strategy {
	switch s {
	case "first":
		return KeepFirst
	default:
		return KeepLast
	}
}

// Deduplicate scans an ordered list of key=value pairs (as would be produced
// by reading a raw .env file line by line) and returns a Result.
//
// pairs is a slice of [2]string{key, value} tuples preserving source order.
func Deduplicate(pairs [][2]string, strategy Strategy) Result {
	seen := make(map[string]bool)
	duplicates := make(map[string][]string)
	unique := make(map[string]string)

	for _, pair := range pairs {
		key, val := pair[0], pair[1]
		if seen[key] {
			duplicates[key] = append(duplicates[key], val)
			if strategy == KeepLast {
				unique[key] = val
			}
		} else {
			seen[key] = true
			// Record the first value in the duplicates list too so the
			// full history is always available once a second occurrence
			// is found later.
			unique[key] = val
		}
	}

	// Prepend the first-seen value into each duplicate slice so callers
	// always get the full ordered history.
	for key := range duplicates {
		first := unique[key]
		if strategy == KeepFirst {
			// unique already holds the first value; nothing to change.
			_ = first
		}
		// Build full history: first value + subsequent values.
		// duplicates[key] currently contains values from the 2nd occurrence
		// onwards; prepend the original.
		originalFirst := pairs[firstIndex(pairs, key)][1]
		duplicates[key] = append([][2]string{{key, originalFirst}}[0][1:0:0], append([]string{originalFirst}, duplicates[key]...)...)
	}

	return Result{
		Unique:     unique,
		Duplicates: duplicates,
	}
}

// firstIndex returns the slice index of the first pair with the given key.
func firstIndex(pairs [][2]string, key string) int {
	for i, p := range pairs {
		if p[0] == key {
			return i
		}
	}
	return 0
}
