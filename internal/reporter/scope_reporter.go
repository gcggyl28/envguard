package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/envguard/internal/scoper"
)

// PrintScope writes scope results to stdout.
func PrintScope(result scoper.Result) {
	FprintScope(os.Stdout, result)
}

// FprintScope writes scope results to the given writer.
func FprintScope(w io.Writer, result scoper.Result) {
	fmt.Fprintf(w, "Scope: %s\n", result.Scope)
	fmt.Fprintf(w, "  Included: %d key(s)\n", len(result.Included))

	if len(result.Included) > 0 {
		keys := sortedScopeKeys(result.Included)
		for _, k := range keys {
			fmt.Fprintf(w, "    + %s\n", k)
		}
	}

	fmt.Fprintf(w, "  Excluded: %d key(s)\n", len(result.Excluded))
	if len(result.Excluded) > 0 {
		keys := sortedScopeKeys(result.Excluded)
		for _, k := range keys {
			fmt.Fprintf(w, "    - %s\n", k)
		}
	}
}

func sortedScopeKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
