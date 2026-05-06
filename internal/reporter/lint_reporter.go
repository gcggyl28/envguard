package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envguard/internal/linter"
)

// PrintLint writes lint hints to stdout in a human-readable format.
func PrintLint(hints []linter.Hint) {
	FprintLint(os.Stdout, hints)
}

// FprintLint writes lint hints to the provided writer.
func FprintLint(w io.Writer, hints []linter.Hint) {
	if len(hints) == 0 {
		fmt.Fprintln(w, "✔  No lint hints — your .env looks clean!")
		return
	}

	// Sort hints by key for deterministic output.
	sorted := make([]linter.Hint, len(hints))
	copy(sorted, hints)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Key == sorted[j].Key {
			return sorted[i].Message < sorted[j].Message
		}
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(w, "Lint hints (%d):\n", len(sorted))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, h := range sorted {
		fmt.Fprintf(w, "  %-30s  %s\n", h.Key, h.Message)
	}
}
