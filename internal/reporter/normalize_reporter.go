package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envguard/internal/normalizer"
)

// PrintNormalize writes normalization results to stdout.
func PrintNormalize(res normalizer.Result) {
	FprintNormalize(os.Stdout, res)
}

// FprintNormalize writes normalization results to w.
func FprintNormalize(w io.Writer, res normalizer.Result) {
	if len(res.Changes) == 0 {
		fmt.Fprintln(w, "✔ No normalization changes.")
		return
	}

	fmt.Fprintf(w, "Normalization: %d change(s)\n", len(res.Changes))
	fmt.Fprintln(w, strings.Repeat("-", 48))

	sorted := make([]normalizer.Change, len(res.Changes))
	copy(sorted, res.Changes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, c := range sorted {
		if c.OldKey != c.Key {
			fmt.Fprintf(w, "  ~ key:   %s → %s\n", c.OldKey, c.Key)
		}
		if c.OldValue != c.NewValue {
			fmt.Fprintf(w, "  ~ value: %s → %s\n", c.OldValue, c.NewValue)
		}
		if c.Reason != "" {
			fmt.Fprintf(w, "    reason: %s\n", c.Reason)
		}
	}
}
