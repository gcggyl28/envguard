package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/envguard/internal/merger"
)

// PrintMerge writes a human-readable merge report to stdout.
func PrintMerge(result merger.Result, strategy string) {
	FprintMerge(os.Stdout, result, strategy)
}

// FprintMerge writes a human-readable merge report to w.
func FprintMerge(w io.Writer, result merger.Result, strategy string) {
	fmt.Fprintf(w, "Merge strategy: %s\n", strategy)
	fmt.Fprintf(w, "Total keys in result: %d\n", len(result.Merged))

	if len(result.Added) > 0 {
		fmt.Fprintf(w, "\nAdded keys (%d):\n", len(result.Added))
		for _, k := range result.Added {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(result.Conflicts) == 0 {
		fmt.Fprintln(w, "\nNo conflicts detected.")
		return
	}

	fmt.Fprintf(w, "\nConflicts (%d):\n", len(result.Conflicts))
	for _, c := range result.Conflicts {
		fmt.Fprintf(w, "  ~ %s\n", c.Key)
		fmt.Fprintf(w, "      base:     %s\n", c.BaseValue)
		fmt.Fprintf(w, "      override: %s\n", c.OverrideValue)
		fmt.Fprintf(w, "      resolved: %s\n", c.Resolved)
	}
}

// MergedKeys returns the sorted keys of the merged map (helper for reporters).
func MergedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
