package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envguard/internal/validator"
)

const (
	passIcon  = "✔"
	failIcon  = "✘"
	borderLen = 50
)

// Print writes a human-readable validation report to stdout.
func Print(report *validator.Report) {
	Fprint(os.Stdout, report)
}

// Fprint writes a human-readable validation report to the given writer.
func Fprint(w io.Writer, report *validator.Report) {
	border := strings.Repeat("-", borderLen)

	fmt.Fprintln(w, border)
	fmt.Fprintln(w, "  envguard — Validation Report")
	fmt.Fprintln(w, border)

	for _, r := range report.Results {
		icon := passIcon
		if !r.Passed {
			icon = failIcon
		}
		fmt.Fprintf(w, "  %s  %-30s %s\n", icon, r.Key, r.Message)
	}

	fmt.Fprintln(w, border)

	if report.Valid {
		fmt.Fprintln(w, "  Result: PASSED")
	} else {
		fmt.Fprintln(w, "  Result: FAILED")
	}

	fmt.Fprintln(w, border)
}

// ExitCode returns 0 if the report is valid, 1 otherwise.
func ExitCode(report *validator.Report) int {
	if report.Valid {
		return 0
	}
	return 1
}

// Summary returns a short one-line summary of the report, e.g.:
//
//	"3 checks passed, 1 failed"
func Summary(report *validator.Report) string {
	passed, failed := 0, 0
	for _, r := range report.Results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}
	return fmt.Sprintf("%d checks passed, %d failed", passed, failed)
}
