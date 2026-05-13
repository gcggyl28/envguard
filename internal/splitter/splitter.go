// Package splitter splits an env map into multiple files based on a key predicate or prefix rules.
package splitter

import "sort"

// Rule defines how keys are assigned to a named output bucket.
type Rule struct {
	Name    string
	Prefixes []string
}

// Result holds the split output.
type Result struct {
	Buckets  map[string]map[string]string // name -> env map
	Leftover map[string]string            // keys that matched no rule
}

// Summary returns a human-readable count per bucket.
func (r Result) Summary() map[string]int {
	out := make(map[string]int, len(r.Buckets)+1)
	for name, m := range r.Buckets {
		out[name] = len(m)
	}
	if len(r.Leftover) > 0 {
		out["(unmatched)"] = len(r.Leftover)
	}
	return out
}

// Split distributes env keys into named buckets according to the provided rules.
// Each key is assigned to the first matching rule. Unmatched keys go to Leftover.
// Rules with no Prefixes act as a catch-all for that bucket.
func Split(env map[string]string, rules []Rule) Result {
	buckets := make(map[string]map[string]string, len(rules))
	for _, r := range rules {
		buckets[r.Name] = make(map[string]string)
	}
	leftover := make(map[string]string)

	keys := sortedKeys(env)
	for _, k := range keys {
		matched := false
		for _, r := range rules {
			if len(r.Prefixes) == 0 {
				// catch-all rule
				buckets[r.Name][k] = env[k]
				matched = true
				break
			}
			for _, p := range r.Prefixes {
				if len(k) >= len(p) && k[:len(p)] == p {
					buckets[r.Name][k] = env[k]
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
		if !matched {
			leftover[k] = env[k]
		}
	}
	return Result{Buckets: buckets, Leftover: leftover}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
