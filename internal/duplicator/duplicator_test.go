package duplicator

import (
	"testing"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"API_KEY":     "secret",
}

func TestDuplicate_Success(t *testing.T) {
	out, results := Duplicate(baseEnv, map[string]string{"DB_HOST": "DATABASE_HOST"}, true)
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Error("source key should still exist")
	}
	if len(results) != 1 || results[0].Skipped {
		t.Error("expected one successful duplication")
	}
}

func TestDuplicate_SkipsMissingSource(t *testing.T) {
	_, results := Duplicate(baseEnv, map[string]string{"MISSING_KEY": "NEW_KEY"}, true)
	if len(results) != 1 || !results[0].Skipped {
		t.Error("expected skip for missing source key")
	}
	if results[0].Reason != "source key not found" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestDuplicate_SkipsExistingDestNoOverwrite(t *testing.T) {
	env := map[string]string{"SRC": "val", "DST": "existing"}
	out, results := Duplicate(env, map[string]string{"SRC": "DST"}, false)
	if out["DST"] != "existing" {
		t.Error("destination should not be overwritten")
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Error("expected skip when destination exists and overwrite=false")
	}
}

func TestDuplicate_OverwritesExistingDest(t *testing.T) {
	env := map[string]string{"SRC": "new_val", "DST": "old_val"}
	out, results := Duplicate(env, map[string]string{"SRC": "DST"}, true)
	if out["DST"] != "new_val" {
		t.Errorf("expected DST=new_val, got %q", out["DST"])
	}
	if results[0].Skipped {
		t.Error("expected successful duplication with overwrite=true")
	}
}

func TestDuplicate_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"A": "1"}
	Duplicate(original, map[string]string{"A": "B"}, true)
	if _, ok := original["B"]; ok {
		t.Error("input map should not be mutated")
	}
}

func TestSummarize(t *testing.T) {
	results := []Result{
		{Key: "A", NewKey: "B", Skipped: false},
		{Key: "C", NewKey: "D", Skipped: true, Reason: "source key not found"},
		{Key: "E", NewKey: "F", Skipped: false},
	}
	s := Summarize(results)
	if s.Duplicated != 2 {
		t.Errorf("expected 2 duplicated, got %d", s.Duplicated)
	}
	if s.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", s.Skipped)
	}
}
