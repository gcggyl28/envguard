package reporter

import (
	"strings"
	"testing"

	"github.com/user/envguard/internal/auditor"
	"github.com/user/envguard/internal/validator"
)

func TestFprintMarkdown_NoIssues(t *testing.T) {
	var buf strings.Builder
	FprintMarkdown(&buf, nil, nil)
	out := buf.String()

	if !strings.Contains(out, "# envguard Report") {
		t.Error("expected report header")
	}
	if !strings.Contains(out, "No validation issues found.") {
		t.Error("expected no-findings message")
	}
	if !strings.Contains(out, "No audit notes.") {
		t.Error("expected no-notes message")
	}
	if !strings.Contains(out, "All checks passed") {
		t.Error("expected passing status")
	}
}

func TestFprintMarkdown_WithFindings(t *testing.T) {
	findings := []validator.Finding{
		{Key: "DB_HOST", Message: "required key is missing"},
	}
	var buf strings.Builder
	FprintMarkdown(&buf, findings, nil)
	out := buf.String()

	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected finding key in output")
	}
	if !strings.Contains(out, "required key is missing") {
		t.Error("expected finding message in output")
	}
	if !strings.Contains(out, "1 finding(s) require attention") {
		t.Error("expected failure status")
	}
}

func TestFprintMarkdown_WithNotes(t *testing.T) {
	notes := []auditor.Note{
		{Key: "LEGACY_KEY", Message: "undeclared key present in .env"},
	}
	var buf strings.Builder
	FprintMarkdown(&buf, nil, notes)
	out := buf.String()

	if !strings.Contains(out, "LEGACY_KEY") {
		t.Error("expected note key in output")
	}
	if !strings.Contains(out, "undeclared key present in .env") {
		t.Error("expected note message in output")
	}
	if !strings.Contains(out, "Total issues:** 1") {
		t.Error("expected total count of 1")
	}
}

func TestFprintMarkdown_TableFormat(t *testing.T) {
	findings := []validator.Finding{
		{Key: "API_KEY", Message: "pattern mismatch"},
	}
	var buf strings.Builder
	FprintMarkdown(&buf, findings, nil)
	out := buf.String()

	if !strings.Contains(out, "| Key | Message |") {
		t.Error("expected markdown table header for findings")
	}
	if !strings.Contains(out, "|-----|---------|" ) {
		t.Error("expected markdown table separator")
	}
}
