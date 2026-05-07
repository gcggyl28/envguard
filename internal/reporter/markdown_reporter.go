package reporter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/envguard/internal/auditor"
	"github.com/user/envguard/internal/validator"
)

// PrintMarkdown writes a Markdown-formatted audit report to stdout.
func PrintMarkdown(findings []validator.Finding, notes []auditor.Note) {
	FprintMarkdown(os.Stdout, findings, notes)
}

// FprintMarkdown writes a Markdown-formatted audit report to the given writer.
func FprintMarkdown(w io.Writer, findings []validator.Finding, notes []auditor.Note) {
	fmt.Fprintf(w, "# envguard Report\n\n")
	fmt.Fprintf(w, "_Generated: %s_\n\n", time.Now().Format(time.RFC3339))

	fmt.Fprintf(w, "## Validation Findings\n\n")
	if len(findings) == 0 {
		fmt.Fprintf(w, "_No validation issues found._\n\n")
	} else {
		fmt.Fprintf(w, "| Key | Message |\n")
		fmt.Fprintf(w, "|-----|---------|\n")
		for _, f := range findings {
			fmt.Fprintf(w, "| `%s` | %s |\n", f.Key, f.Message)
		}
		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintf(w, "## Audit Notes\n\n")
	if len(notes) == 0 {
		fmt.Fprintf(w, "_No audit notes._\n\n")
	} else {
		fmt.Fprintf(w, "| Key | Note |\n")
		fmt.Fprintf(w, "|-----|------|\n")
		for _, n := range notes {
			fmt.Fprintf(w, "| `%s` | %s |\n", n.Key, n.Message)
		}
		fmt.Fprintf(w, "\n")
	}

	total := len(findings) + len(notes)
	status := summaryStatus(len(findings))
	fmt.Fprintf(w, "## Summary\n\n")
	fmt.Fprintf(w, "- **Status:** %s\n", status)
	fmt.Fprintf(w, "- **Total issues:** %d\n", total)
}

// summaryStatus returns a human-readable status string based on the number of findings.
func summaryStatus(findingCount int) string {
	if findingCount == 0 {
		return "✅ All checks passed"
	}
	return fmt.Sprintf("❌ %d finding(s) require attention", findingCount)
}
