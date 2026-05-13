package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/envguard/internal/freezer"
)

// PrintFreeze writes a freeze result report to stdout.
func PrintFreeze(res freezer.FreezeResult) {
	FprintFreeze(os.Stdout, res)
}

// FprintFreeze writes a freeze result report to w.
func FprintFreeze(w io.Writer, res freezer.FreezeResult) {
	fmt.Fprintf(w, "❄️  Freeze Report\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 40))
	fmt.Fprintf(w, "  Output file : %s\n", res.File)
	fmt.Fprintf(w, "  Frozen keys : %d\n", len(res.Frozen))
	if len(res.Skipped) > 0 {
		fmt.Fprintf(w, "  Skipped keys: %d\n", len(res.Skipped))
		for _, k := range res.Skipped {
			fmt.Fprintf(w, "    - %s\n", k)
		}
	}
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 40))
	for _, k := range res.Frozen {
		fmt.Fprintf(w, "  ✔ %s\n", k)
	}
}

// PrintThaw writes a thaw (load) report to stdout.
func PrintThaw(fe freezer.FrozenEnv) {
	FprintThaw(os.Stdout, fe)
}

// FprintThaw writes a thaw report to w.
func FprintThaw(w io.Writer, fe freezer.FrozenEnv) {
	fmt.Fprintf(w, "🔥 Thaw Report\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 40))
	fmt.Fprintf(w, "  Source   : %s\n", fe.Source)
	fmt.Fprintf(w, "  Frozen at: %s\n", fe.FrozenAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "  Keys     : %d\n", len(fe.Keys))
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 40))
	for _, k := range fe.Keys {
		fmt.Fprintf(w, "  • %s\n", k)
	}
}
