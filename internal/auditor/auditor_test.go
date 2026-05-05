package auditor_test

import (
	"testing"

	"github.com/user/envguard/internal/auditor"
	"github.com/user/envguard/internal/schema"
)

var baseSchema = &schema.Schema{
	Vars: []schema.VarSpec{
		{Name: "APP_ENV", Required: true},
		{Name: "PORT", Required: false, Default: "8080"},
		{Name: "DB_URL", Required: true},
	},
}

func TestAudit_NoFindings(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "3000",
		"DB_URL":  "postgres://localhost/mydb",
	}
	result := auditor.Audit(env, baseSchema)
	if len(result.Findings) != 0 {
		t.Errorf("expected no findings, got %d: %s", len(result.Findings), result.Summary())
	}
	if result.HasWarnings() {
		t.Error("expected HasWarnings to be false")
	}
}

func TestAudit_UndeclaredKey(t *testing.T) {
	env := map[string]string{
		"APP_ENV":    "production",
		"PORT":       "3000",
		"DB_URL":     "postgres://localhost/mydb",
		"SECRET_KEY": "abc123",
	}
	result := auditor.Audit(env, baseSchema)
	if !result.HasWarnings() {
		t.Error("expected warnings for undeclared key")
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
	f := result.Findings[0]
	if f.Key != "SECRET_KEY" {
		t.Errorf("expected finding for SECRET_KEY, got %s", f.Key)
	}
	if f.Severity != "warn" {
		t.Errorf("expected severity warn, got %s", f.Severity)
	}
}

func TestAudit_EmptyValueWithDefault(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "staging",
		"PORT":    "",
		"DB_URL":  "postgres://localhost/mydb",
	}
	result := auditor.Audit(env, baseSchema)
	if result.HasWarnings() {
		t.Error("expected no warnings")
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 info finding, got %d", len(result.Findings))
	}
	f := result.Findings[0]
	if f.Key != "PORT" {
		t.Errorf("expected finding for PORT, got %s", f.Key)
	}
	if f.Severity != "info" {
		t.Errorf("expected severity info, got %s", f.Severity)
	}
}

func TestAudit_SummaryEmpty(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "test",
		"PORT":    "9000",
		"DB_URL":  "sqlite://:memory:",
	}
	result := auditor.Audit(env, baseSchema)
	if result.Summary() != "audit passed: no issues found" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestAudit_SummaryContainsKeys(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "dev",
		"PORT":    "8080",
		"DB_URL":  "postgres://localhost/mydb",
		"UNKNOWN": "value",
	}
	result := auditor.Audit(env, baseSchema)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	if !containsStr(summary, "UNKNOWN") {
		t.Errorf("expected summary to mention UNKNOWN, got: %s", summary)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
