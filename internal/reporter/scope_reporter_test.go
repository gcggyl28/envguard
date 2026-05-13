package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/envguard/internal/reporter"
	"github.com/yourusername/envguard/internal/scoper"
)

func TestFprintScope_NoKeys(t *testing.T) {
	result := scoper.Result{
		Scope:    "production",
		Included: map[string]string{},
		Excluded: map[string]string{},
	}
	var buf bytes.Buffer
	reporter.FprintScope(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Error("expected scope name in output")
	}
	if !strings.Contains(out, "Included: 0") {
		t.Error("expected 0 included keys")
	}
}

func TestFprintScope_WithIncluded(t *testing.T) {
	result := scoper.Result{
		Scope: "staging",
		Included: map[string]string{
			"STAGING_DB_HOST": "db.staging",
			"STAGING_API_KEY": "key123",
		},
		Excluded: map[string]string{},
	}
	var buf bytes.Buffer
	reporter.FprintScope(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "+ STAGING_API_KEY") {
		t.Error("expected STAGING_API_KEY in included output")
	}
	if !strings.Contains(out, "Included: 2") {
		t.Error("expected 2 included keys")
	}
}

func TestFprintScope_WithExcluded(t *testing.T) {
	result := scoper.Result{
		Scope:    "production",
		Included: map[string]string{"PROD_KEY": "val"},
		Excluded: map[string]string{"STAGING_KEY": "other"},
	}
	var buf bytes.Buffer
	reporter.FprintScope(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "- STAGING_KEY") {
		t.Error("expected STAGING_KEY in excluded output")
	}
	if !strings.Contains(out, "Excluded: 1") {
		t.Error("expected 1 excluded key")
	}
}

func TestFprintScope_SortedOutput(t *testing.T) {
	result := scoper.Result{
		Scope: "prod",
		Included: map[string]string{
			"PROD_Z": "z",
			"PROD_A": "a",
			"PROD_M": "m",
		},
		Excluded: map[string]string{},
	}
	var buf bytes.Buffer
	reporter.FprintScope(&buf, result)
	out := buf.String()
	idxA := strings.Index(out, "PROD_A")
	idxM := strings.Index(out, "PROD_M")
	idxZ := strings.Index(out, "PROD_Z")
	if !(idxA < idxM && idxM < idxZ) {
		t.Error("expected keys in sorted order")
	}
}
