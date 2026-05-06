// Package differ compares two .env files and reports differences
// relative to a schema, useful for auditing env drift between environments.
package differ

import (
	"fmt"
	"sort"

	"github.com/yourorg/envguard/internal/schema"
)

// DiffResult holds the result of comparing two env maps.
type DiffResult struct {
	OnlyInA    []string          // keys present only in the first env
	OnlyInB    []string          // keys present only in the second env
	DiffValues map[string][2]string // keys with differing values: key -> [valueA, valueB]
	Common     []string          // keys present in both with equal values
}

// Diff compares two parsed env maps (key->value) and returns a DiffResult.
func Diff(envA, envB map[string]string) DiffResult {
	result := DiffResult{
		DiffValues: make(map[string][2]string),
	}

	seen := make(map[string]bool)

	for k, vA := range envA {
		seen[k] = true
		if vB, ok := envB[k]; ok {
			if vA == vB {
				result.Common = append(result.Common, k)
			} else {
				result.DiffValues[k] = [2]string{vA, vB}
			}
		} else {
			result.OnlyInA = append(result.OnlyInA, k)
		}
	}

	for k := range envB {
		if !seen[k] {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Common)

	return result
}

// DiffAgainstSchema filters a DiffResult to only include keys declared in the schema.
func DiffAgainstSchema(envA, envB map[string]string, s *schema.Schema) DiffResult {
	filteredA := make(map[string]string)
	filteredB := make(map[string]string)

	for _, field := range s.Fields {
		key := field.Key
		if v, ok := envA[key]; ok {
			filteredA[key] = v
		}
		if v, ok := envB[key]; ok {
			filteredB[key] = v
		}
	}

	return Diff(filteredA, filteredB)
}

// Summary returns a human-readable summary string of a DiffResult.
func Summary(d DiffResult, labelA, labelB string) string {
	return fmt.Sprintf(
		"Only in %s: %d | Only in %s: %d | Differing values: %d | Common: %d",
		labelA, len(d.OnlyInA),
		labelB, len(d.OnlyInB),
		len(d.DiffValues),
		len(d.Common),
	)
}
