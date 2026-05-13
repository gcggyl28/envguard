package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envguard/internal/scorer"
)

// PrintScore writes the score report to stdout.
func PrintScore(r scorer.Result) {
	FprintScore(os.Stdout, r)
}

// FprintScore writes a formatted score report to w.
func FprintScore(w io.Writer, r scorer.Result) {
	bar := scoreBar(r.Score)

	fmt.Fprintf(w, "\n── Env Health Score ──────────────────────────\n")
	fmt.Fprintf(w, "  Score : %d / 100  [%s]\n", r.Score, bar)
	fmt.Fprintf(w, "  Grade : %s\n", gradeLabel(r.Grade))
	fmt.Fprintf(w, "──────────────────────────────────────────────\n")

	if r.Total == 0 {
		fmt.Fprintf(w, "  ✓ No issues found — perfect score!\n")
	} else {
		fmt.Fprintf(w, "  Deductions:\n")
		if r.ValidationLoss > 0 {
			fmt.Fprintf(w, "    ✗ Validation errors : -%d pts\n", r.ValidationLoss)
		}
		if r.AuditLoss > 0 {
			fmt.Fprintf(w, "    ✗ Audit findings    : -%d pts\n", r.AuditLoss)
		}
		if r.LintLoss > 0 {
			fmt.Fprintf(w, "    ✗ Lint hints        : -%d pts\n", r.LintLoss)
		}
	}
	fmt.Fprintln(w)
}

func gradeLabel(g scorer.Grade) string {
	switch g {
	case scorer.GradeA:
		return "A  (Excellent)"
	case scorer.GradeB:
		return "B  (Good)"
	case scorer.GradeC:
		return "C  (Fair)"
	case scorer.GradeD:
		return "D  (Poor)"
	default:
		return "F  (Critical)"
	}
}

func scoreBar(score int) string {
	const width = 20
	filled := score * width / 100
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
