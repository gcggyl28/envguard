// Package renamer provides utilities for renaming keys in a parsed .env map,
// supporting bulk renames via a mapping of old->new key names.
package renamer

import "fmt"

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	Renamed  map[string]string // oldKey -> newKey for successful renames
	Skipped  []string          // oldKeys that were not found in the env map
	Conflicts []string         // newKeys that already existed in the env map
}

// Rename applies a set of key renames to the given env map.
// It returns a new map with the renames applied and a RenameResult describing
// what happened. The original map is not mutated.
//
// Rules:
//   - If oldKey is not present, it is recorded in Skipped.
//   - If newKey already exists in the env map (and is not the same rename target),
//     the rename is skipped and recorded in Conflicts.
//   - Successful renames preserve the original value.
func Rename(env map[string]string, renames map[string]string) (map[string]string, RenameResult) {
	result := RenameResult{
		Renamed: make(map[string]string),
	}

	// Copy original map.
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for oldKey, newKey := range renames {
		if oldKey == newKey {
			continue
		}

		val, exists := out[oldKey]
		if !exists {
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}

		if _, conflict := out[newKey]; conflict {
			result.Conflicts = append(result.Conflicts, newKey)
			continue
		}

		delete(out, oldKey)
		out[newKey] = val
		result.Renamed[oldKey] = newKey
	}

	return out, result
}

// Summary returns a human-readable summary of a RenameResult.
func Summary(r RenameResult) string {
	return fmt.Sprintf(
		"renamed: %d, skipped: %d, conflicts: %d",
		len(r.Renamed), len(r.Skipped), len(r.Conflicts),
	)
}
