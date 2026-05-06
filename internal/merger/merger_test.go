package merger_test

import (
	"testing"

	"github.com/yourusername/envguard/internal/merger"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"C": "3"}

	r := merger.Merge(base, override, merger.StrategyBase)

	if r.Merged["A"] != "1" || r.Merged["B"] != "2" || r.Merged["C"] != "3" {
		t.Errorf("unexpected merged map: %v", r.Merged)
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", r.Conflicts)
	}
	if len(r.Added) != 1 || r.Added[0] != "C" {
		t.Errorf("expected Added=[C], got %v", r.Added)
	}
}

func TestMerge_ConflictStrategyBase(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	override := map[string]string{"KEY": "new"}

	r := merger.Merge(base, override, merger.StrategyBase)

	if r.Merged["KEY"] != "original" {
		t.Errorf("expected base value to win, got %q", r.Merged["KEY"])
	}
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	if r.Conflicts[0].Resolved != "original" {
		t.Errorf("resolved should be base value")
	}
}

func TestMerge_ConflictStrategyOverride(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	override := map[string]string{"KEY": "new"}

	r := merger.Merge(base, override, merger.StrategyOverride)

	if r.Merged["KEY"] != "new" {
		t.Errorf("expected override value to win, got %q", r.Merged["KEY"])
	}
	if r.Conflicts[0].Resolved != "new" {
		t.Errorf("resolved should be override value")
	}
}

func TestMerge_AddedSorted(t *testing.T) {
	base := map[string]string{}
	override := map[string]string{"Z": "z", "A": "a", "M": "m"}

	r := merger.Merge(base, override, merger.StrategyBase)

	if len(r.Added) != 3 || r.Added[0] != "A" || r.Added[1] != "M" || r.Added[2] != "Z" {
		t.Errorf("expected sorted Added, got %v", r.Added)
	}
}

func TestStrategyFromString(t *testing.T) {
	s, err := merger.StrategyFromString("override")
	if err != nil || s != merger.StrategyOverride {
		t.Errorf("expected StrategyOverride, got %v %v", s, err)
	}

	s, err = merger.StrategyFromString("base")
	if err != nil || s != merger.StrategyBase {
		t.Errorf("expected StrategyBase, got %v %v", s, err)
	}

	_, err = merger.StrategyFromString("unknown")
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}
