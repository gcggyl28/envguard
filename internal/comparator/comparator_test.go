package comparator

import (
	"testing"
)

func TestCompare_Identical(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(a, b)
	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Errorf("expected no differences, got added=%v removed=%v changed=%v", r.Added, r.Removed, r.Changed)
	}
	if r.Summary() != "environments are identical" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestCompare_Added(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "NEW_KEY": "value"}
	r := Compare(a, b)
	if len(r.Added) != 1 || r.Added["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY in added, got %v", r.Added)
	}
	if len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Errorf("expected no removed/changed")
	}
}

func TestCompare_Removed(t *testing.T) {
	a := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	b := map[string]string{"FOO": "bar"}
	r := Compare(a, b)
	if len(r.Removed) != 1 || r.Removed["OLD_KEY"] != "gone" {
		t.Errorf("expected OLD_KEY in removed, got %v", r.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := Compare(a, b)
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(r.Changed))
	}
	ch := r.Changed["FOO"]
	if ch.Old != "old" || ch.New != "new" {
		t.Errorf("expected old=old new=new, got %+v", ch)
	}
}

func TestCompare_SortedAdded(t *testing.T) {
	a := map[string]string{}
	b := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := Compare(a, b)
	keys := r.SortedAdded()
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted added keys, got %v", keys)
	}
}

func TestCompare_SummaryDiffers(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"Y": "2"}
	r := Compare(a, b)
	if r.Summary() != "environments differ" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestCompare_Unchanged(t *testing.T) {
	a := map[string]string{"SAME": "val", "DIFF": "a"}
	b := map[string]string{"SAME": "val", "DIFF": "b"}
	r := Compare(a, b)
	if len(r.Unchanged) != 1 || r.Unchanged[0] != "SAME" {
		t.Errorf("expected SAME in unchanged, got %v", r.Unchanged)
	}
}
