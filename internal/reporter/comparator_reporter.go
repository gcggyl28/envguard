package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envguard/internal/comparator"
)

// PrintComparator writes a comparison report to stdout.
func PrintComparator(r comparator.Result, baseLabel, targetLabel string) {
	FprintComparator(os.Stdout, r, baseLabel, targetLabel)
}

// FprintComparator writes a comparison report to the given writer.
func FprintComparator(w io.Writer, r comparator.Result, baseLabel, targetLabel string) {
	fmt.Fprintf(w, "Comparing: %s → %s\n", baseLabel, targetLabel)
	fmt.Fprintln(w, "")

	if len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0 {
		fmt.Fprintln(w, "✅ No differences found.")
		return
	}

	if len(r.Added) > 0 {
		fmt.Fprintf(w, "➕ Added (%d):\n", len(r.Added))
		for _, k := range r.SortedAdded() {
			fmt.Fprintf(w, "   + %s=%s\n", k, r.Added[k])
		}
		fmt.Fprintln(w, "")
	}

	if len(r.Removed) > 0 {
		fmt.Fprintf(w, "➖ Removed (%d):\n", len(r.Removed))
		for _, k := range r.SortedRemoved() {
			fmt.Fprintf(w, "   - %s=%s\n", k, r.Removed[k])
		}
		fmt.Fprintln(w, "")
	}

	if len(r.Changed) > 0 {
		fmt.Fprintf(w, "✏️  Changed (%d):\n", len(r.Changed))
		for _, k := range r.SortedChanged() {
			ch := r.Changed[k]
			fmt.Fprintf(w, "   ~ %s: %q → %q\n", k, ch.Old, ch.New)
		}
		fmt.Fprintln(w, "")
	}

	fmt.Fprintf(w, "Summary: %s\n", r.Summary())
}
