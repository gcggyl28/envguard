package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envguard/internal/promoter"
)

// PrintPromote writes promotion results to stdout.
func PrintPromote(results []promoter.Result, sum promoter.Summary) {
	FprintPromote(os.Stdout, results, sum)
}

// FprintPromote writes promotion results to w.
func FprintPromote(w io.Writer, results []promoter.Result, sum promoter.Summary) {
	fmt.Fprintf(w, "Promotion report — promoted: %d, skipped: %d\n", sum.Promoted, sum.Skipped)

	// Sort for deterministic output.
	sorted := make([]promoter.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		if r.Skipped {
			fmt.Fprintf(w, "  ⏭  %-30s skipped  (%s)\n", r.Key, r.Reason)
		} else {
			fmt.Fprintf(w, "  ✔  %-30s promoted\n", r.Key)
		}
	}
}
