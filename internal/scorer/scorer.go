// Package scorer computes a health score for an .env file based on
// validation, lint, and audit findings.
package scorer

import (
	"github.com/user/envguard/internal/auditor"
	"github.com/user/envguard/internal/linter"
	"github.com/user/envguard/internal/validator"
)

// Grade represents a letter grade for the env file health.
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

// Result holds the computed score and contributing details.
type Result struct {
	Score          int
	Grade          Grade
	ValidationLoss int
	AuditLoss      int
	LintLoss       int
	Total          int
}

// Score computes a 0-100 health score for the env file.
// Deductions:
//   - Each validation error:  -15 points
//   - Each audit finding:     -10 points
//   - Each lint hint:         -5  points
func Score(
	validationErrs []validator.Error,
	auditFindings []auditor.Finding,
	lintHints []linter.Hint,
) Result {
	const base = 100

	validationLoss := len(validationErrs) * 15
	auditLoss := len(auditFindings) * 10
	lintLoss := len(lintHints) * 5

	total := validationLoss + auditLoss + lintLoss
	score := base - total
	if score < 0 {
		score = 0
	}

	return Result{
		Score:          score,
		Grade:          grade(score),
		ValidationLoss: validationLoss,
		AuditLoss:      auditLoss,
		LintLoss:       lintLoss,
		Total:          total,
	}
}

func grade(score int) Grade {
	switch {
	case score >= 90:
		return GradeA
	case score >= 75:
		return GradeB
	case score >= 60:
		return GradeC
	case score >= 40:
		return GradeD
	default:
		return GradeF
	}
}
