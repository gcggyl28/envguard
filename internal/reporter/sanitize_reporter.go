package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envguard/internal/sanitizer"
)

// PrintSanitize writes the sanitization report to stdout.
func PrintSanitize(result sanitizer.Result) {
	FprintSanitize(os.Stdout, result)
}

// FprintSanitize writes the sanitization report to w.
func FprintSanitize(w io.Writer, result sanitizer.Result) {
	if len(result.Changed) == 0 {
		fmt.Fprintln(w, "✔  No sanitization changes.")
		return
	}

	fmt.Fprintf(w, "Sanitization report — %d change(s):\n", len(result.Changed))
	fmt.Fprintln(w, "")

	// Sort changes by key for deterministic output.
	sorted := make([]sanitizer.Change, len(result.Changed))
	copy(sorted, result.Changed)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, c := range sorted {
		fmt.Fprintf(w, "  %-24s  [%s]\n", c.Key, c.Reason)
		fmt.Fprintf(w, "    before: %q\n", c.Before)
		fmt.Fprintf(w, "    after:  %q\n", c.After)
	}

	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "Total keys after sanitization: %d\n", len(result.Env))
}
