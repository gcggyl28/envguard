package reporter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yourorg/envguard/internal/auditor"
	"github.com/yourorg/envguard/internal/validator"
)

func TestPrintJSON_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	PrintJSON(&buf, nil, nil)

	var report JSONReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}
	if !report.Valid {
		t.Errorf("expected Valid=true when no findings")
	}
	if report.Summary.TotalFindings != 0 {
		t.Errorf("expected 0 findings, got %d", report.Summary.TotalFindings)
	}
}

func TestPrintJSON_WithFindings(t *testing.T) {
	findings := []validator.Finding{
		{Key: "DB_HOST", Message: "required key is missing"},
		{Key: "PORT", Message: "value does not match pattern"},
	}
	auditResults := []auditor.Result{
		{Key: "OLD_KEY", Message: "undeclared key found in env file"},
	}

	var buf bytes.Buffer
	PrintJSON(&buf, findings, auditResults)

	var report JSONReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if report.Valid {
		t.Errorf("expected Valid=false when findings exist")
	}
	if report.Summary.TotalFindings != 2 {
		t.Errorf("expected 2 findings, got %d", report.Summary.TotalFindings)
	}
	if report.Summary.TotalAudit != 1 {
		t.Errorf("expected 1 audit result, got %d", report.Summary.TotalAudit)
	}
	if report.Findings[0].Key != "DB_HOST" {
		t.Errorf("expected first finding key DB_HOST, got %s", report.Findings[0].Key)
	}
	if report.Findings[0].Severity != "error" {
		t.Errorf("expected severity 'error', got %s", report.Findings[0].Severity)
	}
	if report.Audit[0].Key != "OLD_KEY" {
		t.Errorf("expected audit key OLD_KEY, got %s", report.Audit[0].Key)
	}
}

func TestBuildJSONReport_EmptySlices(t *testing.T) {
	report := buildJSONReport([]validator.Finding{}, []auditor.Result{})
	if !report.Valid {
		t.Errorf("expected Valid=true for empty findings")
	}
	if len(report.Findings) != 0 {
		t.Errorf("expected empty findings slice")
	}
	if len(report.Audit) != 0 {
		t.Errorf("expected empty audit slice")
	}
}
