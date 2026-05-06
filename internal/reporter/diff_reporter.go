package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourorg/envguard/internal/differ"
)

// PrintDiff writes a human-readable diff report to stdout.
func PrintDiff(d differ.DiffResult, labelA, labelB string) {
	FprintDiff(os.Stdout, d, labelA, labelB)
}

// FprintDiff writes a human-readable diff report to the given writer.
func FprintDiff(w io.Writer, d differ.DiffResult, labelA, labelB string) {
	fmt.Fprintf(w, "=== Env Diff: %s vs %s ===\n", labelA, labelB)

	if len(d.OnlyInA) > 0 {
		fmt.Fprintf(w, "\nOnly in %s:\n", labelA)
		for _, k := range d.OnlyInA {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(d.OnlyInB) > 0 {
		fmt.Fprintf(w, "\nOnly in %s:\n", labelB)
		for _, k := range d.OnlyInB {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(d.DiffValues) > 0 {
		fmt.Fprintln(w, "\nDiffering values:")
		keys := sortedKeys(d.DiffValues)
		for _, k := range keys {
			pair := d.DiffValues[k]
			fmt.Fprintf(w, "  ~ %s: [%s] -> [%s]\n", k, pair[0], pair[1])
		}
	}

	if len(d.OnlyInA) == 0 && len(d.OnlyInB) == 0 && len(d.DiffValues) == 0 {
		fmt.Fprintln(w, "No differences found.")
	}

	fmt.Fprintf(w, "\nSummary: %s\n", differ.Summary(d, labelA, labelB))
}

// PrintDiffStats writes only the numeric counts of differences to the given
// writer, useful for machine-readable or terse output modes.
func PrintDiffStats(w io.Writer, d differ.DiffResult, labelA, labelB string) {
	fmt.Fprintf(w, "only_in_%s=%d only_in_%s=%d differing=%d\n",
		labelA, len(d.OnlyInA),
		labelB, len(d.OnlyInB),
		len(d.DiffValues),
	)
}

func sortedKeys(m map[string][2]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
