package scorer_test

import (
	"testing"

	"github.com/user/envguard/internal/auditor"
	"github.com/user/envguard/internal/linter"
	"github.com/user/envguard/internal/scorer"
	"github.com/user/envguard/internal/validator"
)

func TestScore_Perfect(t *testing.T) {
	r := scorer.Score(nil, nil, nil)
	if r.Score != 100 {
		t.Errorf("expected 100, got %d", r.Score)
	}
	if r.Grade != scorer.GradeA {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
}

func TestScore_ValidationErrors(t *testing.T) {
	errs := []validator.Error{
		{Key: "FOO", Message: "missing"},
		{Key: "BAR", Message: "pattern"},
	}
	r := scorer.Score(errs, nil, nil)
	// 100 - 2*15 = 70
	if r.Score != 70 {
		t.Errorf("expected 70, got %d", r.Score)
	}
	if r.Grade != scorer.GradeC {
		t.Errorf("expected grade C, got %s", r.Grade)
	}
	if r.ValidationLoss != 30 {
		t.Errorf("expected validation loss 30, got %d", r.ValidationLoss)
	}
}

func TestScore_AuditFindings(t *testing.T) {
	findings := []auditor.Finding{
		{Key: "SECRET", Message: "undeclared"},
	}
	r := scorer.Score(nil, findings, nil)
	// 100 - 1*10 = 90
	if r.Score != 90 {
		t.Errorf("expected 90, got %d", r.Score)
	}
	if r.Grade != scorer.GradeA {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
}

func TestScore_LintHints(t *testing.T) {
	hints := []linter.Hint{
		{Key: "db_host", Message: "key should be uppercase"},
		{Key: "API KEY", Message: "key contains spaces"},
		{Key: "TOKEN", Message: "short secret"},
	}
	r := scorer.Score(nil, nil, hints)
	// 100 - 3*5 = 85
	if r.Score != 85 {
		t.Errorf("expected 85, got %d", r.Score)
	}
	if r.Grade != scorer.GradeB {
		t.Errorf("expected grade B, got %s", r.Grade)
	}
}

func TestScore_ClampedAtZero(t *testing.T) {
	var errs []validator.Error
	for i := 0; i < 10; i++ {
		errs = append(errs, validator.Error{Key: "K", Message: "err"})
	}
	r := scorer.Score(errs, nil, nil)
	if r.Score != 0 {
		t.Errorf("expected 0, got %d", r.Score)
	}
	if r.Grade != scorer.GradeF {
		t.Errorf("expected grade F, got %s", r.Grade)
	}
}

func TestScore_GradeD(t *testing.T) {
	// 100 - 4*15 = 40 => D
	var errs []validator.Error
	for i := 0; i < 4; i++ {
		errs = append(errs, validator.Error{Key: "K", Message: "err"})
	}
	r := scorer.Score(errs, nil, nil)
	if r.Grade != scorer.GradeD {
		t.Errorf("expected grade D, got %s", r.Grade)
	}
}
