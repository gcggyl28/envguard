package reporter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/envguard/internal/pinner"
	"github.com/user/envguard/internal/reporter"
)

func makePinned(keys map[string]string) *pinner.PinnedEnv {
	return &pinner.PinnedEnv{
		PinnedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Keys:     keys,
	}
}

func TestFprintPin_NoDrift(t *testing.T) {
	p := makePinned(map[string]string{"A": "1"})
	result := pinner.DriftResult{}
	var buf strings.Builder
	reporter.FprintPin(&buf, p, result)
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got:\n%s", out)
	}
	if !strings.Contains(out, "Pinned keys: 1") {
		t.Errorf("expected pinned key count, got:\n%s", out)
	}
}

func TestFprintPin_WithChanged(t *testing.T) {
	p := makePinned(map[string]string{"DB": "old"})
	result := pinner.DriftResult{
		Changed: []pinner.DriftEntry{
			{Key: "DB", Pinned: "old", Current: "new"},
		},
	}
	var buf strings.Builder
	reporter.FprintPin(&buf, p, result)
	out := buf.String()
	if !strings.Contains(out, "Drift detected") {
		t.Errorf("expected drift header, got:\n%s", out)
	}
	if !strings.Contains(out, "DB") {
		t.Errorf("expected key DB in output, got:\n%s", out)
	}
	if !strings.Contains(out, `"old"`) {
		t.Errorf("expected pinned value in output, got:\n%s", out)
	}
}

func TestFprintPin_WithRemoved(t *testing.T) {
	p := makePinned(map[string]string{"GONE": "val"})
	result := pinner.DriftResult{
		Removed: []string{"GONE"},
	}
	var buf strings.Builder
	reporter.FprintPin(&buf, p, result)
	out := buf.String()
	if !strings.Contains(out, "Removed") {
		t.Errorf("expected removed section, got:\n%s", out)
	}
	if !strings.Contains(out, "GONE") {
		t.Errorf("expected GONE in output, got:\n%s", out)
	}
}

func TestFprintPin_SummaryLine(t *testing.T) {
	p := makePinned(map[string]string{"X": "a", "Y": "b"})
	result := pinner.DriftResult{
		Changed: []pinner.DriftEntry{{Key: "X", Pinned: "a", Current: "z"}},
		Removed: []string{"Y"},
	}
	var buf strings.Builder
	reporter.FprintPin(&buf, p, result)
	out := buf.String()
	if !strings.Contains(out, "1 changed, 1 removed") {
		t.Errorf("expected summary line, got:\n%s", out)
	}
}
