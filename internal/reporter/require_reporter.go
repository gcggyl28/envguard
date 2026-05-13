package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envguard/internal/requirer"
)

// PrintRequire writes a requirements-check report to stdout.
func PrintRequire(res requirer.Result) {
	FprintRequire(os.Stdout, res)
}

// FprintRequire writes a requirements-check report to w.
func FprintRequire(w io.Writer, res requirer.Result) {
	fmt.Fprintf(w, "Requirements check — %d required key(s) checked\n", res.Checked)

	if res.Passed() {
		fmt.Fprintln(w, "✔ All required keys are present and non-empty.")
		return
	}

	fmt.Fprintf(w, "✖ %d required key(s) failed:\n", len(res.Findings))
	for _, f := range res.Findings {
		fmt.Fprintf(w, "  • %-30s %s\n", f.Key, f.Reason)
	}
}
