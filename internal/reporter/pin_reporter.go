package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envguard/internal/pinner"
)

// PrintPin writes a drift report to stdout.
func PrintPin(p *pinner.PinnedEnv, result pinner.DriftResult) {
	FprintPin(os.Stdout, p, result)
}

// FprintPin writes a drift report to w.
func FprintPin(w io.Writer, p *pinner.PinnedEnv, result pinner.DriftResult) {
	fmt.Fprintf(w, "Pinned at: %s\n", p.PinnedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Pinned keys: %d\n", len(p.Keys))

	if len(result.Changed) == 0 && len(result.Removed) == 0 {
		fmt.Fprintln(w, "\n✔ No drift detected. All pinned values match.")
		return
	}

	fmt.Fprintln(w, "\n⚠ Drift detected:")

	if len(result.Changed) > 0 {
		fmt.Fprintln(w, "\n  Changed:")
		for _, e := range result.Changed {
			fmt.Fprintf(w, "    %-30s pinned=%q  current=%q\n", e.Key, e.Pinned, e.Current)
		}
	}

	if len(result.Removed) > 0 {
		fmt.Fprintln(w, "\n  Removed (key no longer present):")
		for _, k := range result.Removed {
			fmt.Fprintf(w, "    - %s\n", k)
		}
	}

	fmt.Fprintf(w, "\nSummary: %d changed, %d removed\n",
		len(result.Changed), len(result.Removed))
}
