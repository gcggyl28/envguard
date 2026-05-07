package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envguard/internal/scanner"
)

// PrintScan writes scan findings to stdout.
func PrintScan(findings []scanner.Finding) {
	FprintScan(os.Stdout, findings)
}

// FprintScan writes scan findings to the given writer.
func FprintScan(w io.Writer, findings []scanner.Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "✔  No scan issues found.")
		return
	}

	// Sort for deterministic output: severity then key.
	sorted := make([]scanner.Finding, len(findings))
	copy(sorted, findings)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Severity != sorted[j].Severity {
			return severityRank(sorted[i].Severity) < severityRank(sorted[j].Severity)
		}
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(w, "Scan Results (%d finding(s)):\n", len(sorted))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, f := range sorted {
		icon := scanIcon(f.Severity)
		fmt.Fprintf(w, "  %s [%s] %s — %s\n", icon, f.Severity, f.Key, f.Message)
	}

	counts := scanner.CountBySeverity(findings)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Summary: %d error(s), %d warning(s), %d info(s)\n",
		counts["error"], counts["warn"], counts["info"])
}

func severityRank(s string) int {
	switch s {
	case "error":
		return 0
	case "warn":
		return 1
	default:
		return 2
	}
}

func scanIcon(severity string) string {
	switch severity {
	case "error":
		return "✖"
	case "warn":
		return "⚠"
	default:
		return "ℹ"
	}
}
