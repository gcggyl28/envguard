package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envguard/internal/deprecator"
)

// PrintDeprecation writes a deprecation report to stdout.
func PrintDeprecation(findings []deprecator.Finding) {
	FprintDeprecation(os.Stdout, findings)
}

// FprintDeprecation writes a deprecation report to the given writer.
func FprintDeprecation(w io.Writer, findings []deprecator.Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "✅ No deprecated keys found.")
		return
	}

	fmt.Fprintf(w, "⚠️  Deprecated keys found: %d\n\n", len(findings))

	for _, f := range findings {
		fmt.Fprintf(w, "  🔑 %s\n", f.Key)
		if f.Reason != "" {
			fmt.Fprintf(w, "     reason:      %s\n", f.Reason)
		}
		if f.Replacement != "" {
			fmt.Fprintf(w, "     replacement: %s\n", f.Replacement)
		}
	}
	fmt.Fprintln(w)
}
