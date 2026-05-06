package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/envguard/internal/merger"
	"github.com/yourusername/envguard/internal/reporter"
)

func TestFprintMerge_NoConflictsNoAdded(t *testing.T) {
	r := merger.Result{
		Merged:    map[string]string{"A": "1"},
		Conflicts: nil,
		Added:     nil,
	}
	var buf bytes.Buffer
	reporter.FprintMerge(&buf, r, "base")
	out := buf.String()

	if !strings.Contains(out, "No conflicts detected") {
		t.Errorf("expected no-conflict message, got:\n%s", out)
	}
	if strings.Contains(out, "Added keys") {
		t.Errorf("should not print Added section when empty")
	}
}

func TestFprintMerge_WithAdded(t *testing.T) {
	r := merger.Result{
		Merged: map[string]string{"A": "1", "B": "2"},
		Added:  []string{"B"},
	}
	var buf bytes.Buffer
	reporter.FprintMerge(&buf, r, "override")
	out := buf.String()

	if !strings.Contains(out, "+ B") {
		t.Errorf("expected added key B in output, got:\n%s", out)
	}
	if !strings.Contains(out, "override") {
		t.Errorf("expected strategy in output")
	}
}

func TestFprintMerge_WithConflicts(t *testing.T) {
	r := merger.Result{
		Merged: map[string]string{"KEY": "original"},
		Conflicts: []merger.Conflict{
			{Key: "KEY", BaseValue: "original", OverrideValue: "new", Resolved: "original"},
		},
	}
	var buf bytes.Buffer
	reporter.FprintMerge(&buf, r, "base")
	out := buf.String()

	if !strings.Contains(out, "~ KEY") {
		t.Errorf("expected conflict marker for KEY, got:\n%s", out)
	}
	if !strings.Contains(out, "base:     original") {
		t.Errorf("expected base value in output")
	}
	if !strings.Contains(out, "override: new") {
		t.Errorf("expected override value in output")
	}
	if !strings.Contains(out, "resolved: original") {
		t.Errorf("expected resolved value in output")
	}
}

func TestMergedKeys_Sorted(t *testing.T) {
	m := map[string]string{"Z": "z", "A": "a", "M": "m"}
	keys := reporter.MergedKeys(m)
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}
