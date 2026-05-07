package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envguard/internal/patcher"
)

func TestFprintPatch_NoOps(t *testing.T) {
	var buf bytes.Buffer
	FprintPatch(&buf, patcher.Result{})
	if !strings.Contains(buf.String(), "no operations") {
		t.Errorf("expected 'no operations' message, got: %q", buf.String())
	}
}

func TestFprintPatch_Applied(t *testing.T) {
	var buf bytes.Buffer
	result := patcher.Result{Applied: []string{"APP_ENV", "LOG_LEVEL"}}
	FprintPatch(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "applied (2)") {
		t.Errorf("expected 'applied (2)', got: %q", out)
	}
	if !strings.Contains(out, "~ APP_ENV") {
		t.Errorf("expected '~ APP_ENV' in output, got: %q", out)
	}
}

func TestFprintPatch_Deleted(t *testing.T) {
	var buf bytes.Buffer
	result := patcher.Result{Deleted: []string{"DB_PORT"}}
	FprintPatch(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "deleted (1)") {
		t.Errorf("expected 'deleted (1)', got: %q", out)
	}
	if !strings.Contains(out, "- DB_PORT") {
		t.Errorf("expected '- DB_PORT' in output, got: %q", out)
	}
}

func TestFprintPatch_Skipped(t *testing.T) {
	var buf bytes.Buffer
	result := patcher.Result{Skipped: []string{"GHOST"}}
	FprintPatch(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' in output, got: %q", out)
	}
	if !strings.Contains(out, "? GHOST") {
		t.Errorf("expected '? GHOST' in output, got: %q", out)
	}
}

func TestFprintPatch_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	result := patcher.Result{Applied: []string{"Z_KEY", "A_KEY", "M_KEY"}}
	FprintPatch(&buf, result)
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	mIdx := strings.Index(out, "M_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Errorf("expected sorted output A < M < Z, positions: %d %d %d", aIdx, mIdx, zIdx)
	}
}
