package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envguard/internal/deprecator"
	"github.com/user/envguard/internal/reporter"
)

func TestFprintDeprecation_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	reporter.FprintDeprecation(&buf, nil)
	if !strings.Contains(buf.String(), "No deprecated keys found") {
		t.Errorf("expected no-findings message, got: %q", buf.String())
	}
}

func TestFprintDeprecation_WithFindings(t *testing.T) {
	findings := []deprecator.Finding{
		{Key: "OLD_API_KEY", Replacement: "API_KEY", Reason: "renamed in v2"},
	}
	var buf bytes.Buffer
	reporter.FprintDeprecation(&buf, findings)
	out := buf.String()
	if !strings.Contains(out, "OLD_API_KEY") {
		t.Errorf("expected key in output, got: %q", out)
	}
	if !strings.Contains(out, "renamed in v2") {
		t.Errorf("expected reason in output, got: %q", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected replacement in output, got: %q", out)
	}
}

func TestFprintDeprecation_NoReplacement(t *testing.T) {
	findings := []deprecator.Finding{
		{Key: "DISABLE_FEATURE_X", Reason: "feature removed"},
	}
	var buf bytes.Buffer
	reporter.FprintDeprecation(&buf, findings)
	out := buf.String()
	if strings.Contains(out, "replacement:") {
		t.Errorf("did not expect replacement line when empty, got: %q", out)
	}
	if !strings.Contains(out, "feature removed") {
		t.Errorf("expected reason in output")
	}
}

func TestFprintDeprecation_CountInHeader(t *testing.T) {
	findings := []deprecator.Finding{
		{Key: "A"},
		{Key: "B"},
		{Key: "C"},
	}
	var buf bytes.Buffer
	reporter.FprintDeprecation(&buf, findings)
	if !strings.Contains(buf.String(), "3") {
		t.Errorf("expected count 3 in header, got: %q", buf.String())
	}
}
