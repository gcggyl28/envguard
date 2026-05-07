package scanner

import (
	"testing"
)

func TestScan_NoFindings(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "envguard",
		"PORT":     "8080",
	}
	findings := Scan(env, DefaultOptions())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestScan_EmptyValue(t *testing.T) {
	env := map[string]string{
		"EMPTY_KEY": "",
		"SPACE_KEY": "   ",
	}
	opts := DefaultOptions()
	findings := Scan(env, opts)
	if len(findings) != 2 {
		t.Errorf("expected 2 findings for empty values, got %d", len(findings))
	}
	for _, f := range findings {
		if f.Severity != "warn" {
			t.Errorf("expected severity 'warn', got %q", f.Severity)
		}
	}
}

func TestScan_LongValue(t *testing.T) {
	env := map[string]string{
		"BIG_KEY": string(make([]byte, 300)),
	}
	opts := DefaultOptions()
	findings := Scan(env, opts)
	found := false
	for _, f := range findings {
		if f.Key == "BIG_KEY" && f.Severity == "info" {
			found = true
		}
	}
	if !found {
		t.Error("expected an info finding for long value")
	}
}

func TestScan_DisableEmptyCheck(t *testing.T) {
	env := map[string]string{"KEY": ""}
	opts := DefaultOptions()
	opts.CheckEmptyValues = false
	findings := Scan(env, opts)
	if len(findings) != 0 {
		t.Errorf("expected no findings when empty check disabled, got %d", len(findings))
	}
}

func TestCountBySeverity(t *testing.T) {
	findings := []Finding{
		{Severity: "error"},
		{Severity: "error"},
		{Severity: "warn"},
		{Severity: "info"},
	}
	counts := CountBySeverity(findings)
	if counts["error"] != 2 {
		t.Errorf("expected 2 errors, got %d", counts["error"])
	}
	if counts["warn"] != 1 {
		t.Errorf("expected 1 warn, got %d", counts["warn"])
	}
	if counts["info"] != 1 {
		t.Errorf("expected 1 info, got %d", counts["info"])
	}
}

func TestCountBySeverity_Empty(t *testing.T) {
	counts := CountBySeverity(nil)
	if counts["error"] != 0 || counts["warn"] != 0 || counts["info"] != 0 {
		t.Error("expected all zero counts for empty findings")
	}
}
