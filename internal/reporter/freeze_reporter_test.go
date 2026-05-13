package reporter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/envguard/internal/freezer"
	"github.com/yourusername/envguard/internal/reporter"
)

func TestFprintFreeze_NoSkipped(t *testing.T) {
	res := freezer.FreezeResult{
		Frozen:  []string{"A", "B"},
		Skipped: nil,
		File:    "frozen.json",
	}
	var sb strings.Builder
	reporter.FprintFreeze(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "frozen.json") {
		t.Error("expected file name in output")
	}
	if !strings.Contains(out, "Frozen keys : 2") {
		t.Error("expected frozen count")
	}
	if strings.Contains(out, "Skipped") {
		t.Error("should not show skipped section when none")
	}
}

func TestFprintFreeze_WithSkipped(t *testing.T) {
	res := freezer.FreezeResult{
		Frozen:  []string{"A"},
		Skipped: []string{"B", "C"},
		File:    "out.json",
	}
	var sb strings.Builder
	reporter.FprintFreeze(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "Skipped keys: 2") {
		t.Error("expected skipped count")
	}
	if !strings.Contains(out, "- B") {
		t.Error("expected skipped key B listed")
	}
}

func TestFprintThaw_Basic(t *testing.T) {
	fe := freezer.FrozenEnv{
		FrozenAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Source:   "prod.env",
		Keys:     []string{"DB_HOST", "DB_PORT"},
		Values:   map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"},
	}
	var sb strings.Builder
	reporter.FprintThaw(&sb, fe)
	out := sb.String()
	if !strings.Contains(out, "prod.env") {
		t.Error("expected source in output")
	}
	if !strings.Contains(out, "2024-01-15") {
		t.Error("expected frozen date in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key DB_HOST listed")
	}
}

func TestFprintThaw_KeyCount(t *testing.T) {
	fe := freezer.FrozenEnv{
		FrozenAt: time.Now(),
		Source:   "s",
		Keys:     []string{"A", "B", "C"},
		Values:   map[string]string{},
	}
	var sb strings.Builder
	reporter.FprintThaw(&sb, fe)
	if !strings.Contains(sb.String(), "Keys     : 3") {
		t.Error("expected key count 3")
	}
}
