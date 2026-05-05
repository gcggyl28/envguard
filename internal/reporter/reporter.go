package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envguard/internal/validator"
)

const (
	passIcon = "✔"
	failIcon = "✘"
)

// Print writes a human-readable validation report to stdout.
func Print(report *validator.Report) {
	Fprint(os.Stdout, report)
}

// Fprint writes a human-readable validation report to the given writer.
func Fprint(w io.Writer, report *validator.Report) {
	fmt.Fprintln(w, strings.Repeat("-", 50))
	fmt.Fprintln(w, "  envguard — Validation Report")
	fmt.Fprintln(w, strings.Repeat("-", 50))

	for _, r := range report.Results {
		icon := passIcon
		if !r.Passed {
			icon = failIcon
		}
		fmt.Fprintf(w, "  %s  %-30s %s\n", icon, r.Key, r.Message)
	}

	fmt.Fprintln(w, strings.Repeat("-", 50))

	if report.Valid {
		fmt.Fprintln(w, "  Result: PASSED")
	} else {
		fmt.Fprintln(w, "  Result: FAILED")
	}

	fmt.Fprintln(w, strings.Repeat("-", 50))
}

// ExitCode returns 0 if the report is valid, 1 otherwise.
func ExitCode(report *validator.Report) int {
	if report.Valid {
		return 0
	}
	return 1
}
