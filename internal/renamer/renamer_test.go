package renamer_test

import (
	"testing"

	"github.com/yourusername/envguard/internal/renamer"
)

var baseEnv = map[string]string{
	"DB_HOST":  "localhost",
	"DB_PORT":  "5432",
	"APP_NAME": "myapp",
}

func TestRename_Success(t *testing.T) {
	out, result := renamer.Rename(baseEnv, map[string]string{
		"DB_HOST": "DATABASE_HOST",
	})

	if _, ok := out["DATABASE_HOST"]; !ok {
		t.Error("expected DATABASE_HOST to exist in output")
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed from output")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected value 'localhost', got %q", out["DATABASE_HOST"])
	}
	if len(result.Renamed) != 1 {
		t.Errorf("expected 1 renamed, got %d", len(result.Renamed))
	}
}

func TestRename_SkipsMissingKey(t *testing.T) {
	_, result := renamer.Rename(baseEnv, map[string]string{
		"MISSING_KEY": "NEW_KEY",
	})

	if len(result.Skipped) != 1 || result.Skipped[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in Skipped, got %v", result.Skipped)
	}
}

func TestRename_ConflictWithExistingKey(t *testing.T) {
	_, result := renamer.Rename(baseEnv, map[string]string{
		"DB_HOST": "DB_PORT", // DB_PORT already exists
	})

	if len(result.Conflicts) != 1 || result.Conflicts[0] != "DB_PORT" {
		t.Errorf("expected DB_PORT in Conflicts, got %v", result.Conflicts)
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"FOO": "bar"}
	renamer.Rename(original, map[string]string{"FOO": "BAZ"})

	if _, ok := original["FOO"]; !ok {
		t.Error("original map was mutated")
	}
}

func TestRename_SameKeyIsNoop(t *testing.T) {
	out, result := renamer.Rename(baseEnv, map[string]string{
		"DB_HOST": "DB_HOST",
	})

	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST to remain, got %q", out["DB_HOST"])
	}
	if len(result.Renamed) != 0 {
		t.Errorf("expected 0 renames for same-key, got %d", len(result.Renamed))
	}
}

func TestSummary(t *testing.T) {
	_, result := renamer.Rename(baseEnv, map[string]string{
		"DB_HOST":    "DATABASE_HOST",
		"MISSING":    "WHATEVER",
		"DB_PORT":    "APP_NAME", // conflict
	})

	s := renamer.Summary(result)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
