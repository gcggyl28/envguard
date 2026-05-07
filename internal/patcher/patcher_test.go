package patcher

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"LOG_LEVEL": "info",
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	env := baseEnv()
	ops := []Op{{Key: "NEW_KEY", Value: "hello"}}
	out, result := Patch(env, ops)
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
	if len(result.Applied) != 1 || result.Applied[0] != "NEW_KEY" {
		t.Errorf("expected Applied=[NEW_KEY], got %v", result.Applied)
	}
}

func TestPatch_OverwriteExistingKey(t *testing.T) {
	env := baseEnv()
	ops := []Op{{Key: "APP_ENV", Value: "staging"}}
	out, result := Patch(env, ops)
	if out["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", out["APP_ENV"])
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(result.Applied))
	}
}

func TestPatch_DeleteExistingKey(t *testing.T) {
	env := baseEnv()
	ops := []Op{{Key: "DB_PORT", Delete: true}}
	out, result := Patch(env, ops)
	if _, ok := out["DB_PORT"]; ok {
		t.Error("expected DB_PORT to be deleted")
	}
	if len(result.Deleted) != 1 || result.Deleted[0] != "DB_PORT" {
		t.Errorf("expected Deleted=[DB_PORT], got %v", result.Deleted)
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	env := baseEnv()
	ops := []Op{{Key: "GHOST_KEY", Delete: true}}
	_, result := Patch(env, ops)
	if len(result.Skipped) != 1 || result.Skipped[0] != "GHOST_KEY" {
		t.Errorf("expected Skipped=[GHOST_KEY], got %v", result.Skipped)
	}
}

func TestPatch_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	ops := []Op{{Key: "APP_ENV", Value: "test"}}
	Patch(env, ops)
	if env["APP_ENV"] != "production" {
		t.Error("original env map was mutated")
	}
}

func TestParseOps_Valid(t *testing.T) {
	lines := []string{"APP_ENV=staging", "!DB_PORT", "# comment", "", "LOG_LEVEL=debug"}
	ops, err := ParseOps(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(ops))
	}
	if ops[1].Delete != true || ops[1].Key != "DB_PORT" {
		t.Errorf("expected delete op for DB_PORT, got %+v", ops[1])
	}
}

func TestParseOps_InvalidLine(t *testing.T) {
	_, err := ParseOps([]string{"NOEQUALS"})
	if err == nil {
		t.Error("expected error for invalid op line")
	}
}

func TestParseOps_EmptyDeleteKey(t *testing.T) {
	_, err := ParseOps([]string{"!"})
	if err == nil {
		t.Error("expected error for empty delete key")
	}
}
