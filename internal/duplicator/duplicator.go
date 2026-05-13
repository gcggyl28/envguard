package duplicator

// Result holds the outcome of a key duplication operation.
type Result struct {
	Key      string
	NewKey   string
	Skipped  bool
	Reason   string
}

// Summary returns counts of successful duplications and skips.
type Summary struct {
	Duplicated int
	Skipped    int
}

// Duplicate copies values from srcKey to dstKey in the provided env map.
// Each entry in pairs maps a source key to a destination key.
// If overwrite is false, existing destination keys are skipped.
func Duplicate(env map[string]string, pairs map[string]string, overwrite bool) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	results := make([]Result, 0, len(pairs))

	// Iterate in a stable order for deterministic output.
	srcKeys := sortedKeys(pairs)
	for _, src := range srcKeys {
		dst := pairs[src]
		srcVal, srcExists := out[src]
		if !srcExists {
			results = append(results, Result{
				Key:     src,
				NewKey:  dst,
				Skipped: true,
				Reason:  "source key not found",
			})
			continue
		}
		if _, dstExists := out[dst]; dstExists && !overwrite {
			results = append(results, Result{
				Key:     src,
				NewKey:  dst,
				Skipped: true,
				Reason:  "destination key already exists",
			})
			continue
		}
		out[dst] = srcVal
		results = append(results, Result{
			Key:    src,
			NewKey: dst,
		})
	}
	return out, results
}

// Summarize aggregates a slice of Results into a Summary.
func Summarize(results []Result) Summary {
	s := Summary{}
	for _, r := range results {
		if r.Skipped {
			s.Skipped++
		} else {
			s.Duplicated++
		}
	}
	return s
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

import "sort"
