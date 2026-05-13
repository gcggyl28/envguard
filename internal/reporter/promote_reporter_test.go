package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envguard/internal/promoter"
	"github.com/user/envguard/internal/reporter"
)

func TestFprintPromote_NoResults(t *testing.T) {
	var buf bytes.Buffer
	reporter.FprintPromote(&buf, nil, promoter.Summary{})
	if !strings.Contains(buf.String(), "promoted: 0") {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestFprintPromote_AllPromoted(t *testing.T) {
	results := []promoter.Result{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "db.prod"},
	}
	sum := promoter.Summary{Promoted: 2}
	var buf bytes.Buffer
	reporter.FprintPromote(&buf, results, sum)
	out := buf.String()
	if !strings.Contains(out, "promoted: 2") {
		t.Fatalf("expected promoted count in output: %s", out)
	}
	if !strings.Contains(out, "APP_ENV") || !strings.Contains(out, "DB_HOST") {
		t.Fatalf("expected keys in output: %s", out)
	}
}

func TestFprintPromote_WithSkipped(t *testing.T) {
	results := []promoter.Result{
		{Key: "SECRET_KEY", Skipped: true, Reason: "in deny list"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	sum := promoter.Summary{Promoted: 1, Skipped: 1}
	var buf bytes.Buffer
	reporter.FprintPromote(&buf, results, sum)
	out := buf.String()
	if !strings.Contains(out, "skipped: 1") {
		t.Fatalf("expected skipped count: %s", out)
	}
	if !strings.Contains(out, "in deny list") {
		t.Fatalf("expected skip reason: %s", out)
	}
}

func TestFprintPromote_SortedOutput(t *testing.T) {
	results := []promoter.Result{
		{Key: "Z_KEY", Value: "z"},
		{Key: "A_KEY", Value: "a"},
	}
	sum := promoter.Summary{Promoted: 2}
	var buf bytes.Buffer
	reporter.FprintPromote(&buf, results, sum)
	out := buf.String()
	if strings.Index(out, "A_KEY") > strings.Index(out, "Z_KEY") {
		t.Fatal("output should be sorted alphabetically")
	}
}
