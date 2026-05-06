package differ_test

import (
	"testing"

	"github.com/yourorg/envguard/internal/differ"
	"github.com/yourorg/envguard/internal/schema"
)

func TestDiff_AllCommon(t *testing.T) {
	envA := map[string]string{"FOO": "bar", "BAZ": "qux"}
	envB := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := differ.Diff(envA, envB)

	if len(result.Common) != 2 {
		t.Errorf("expected 2 common keys, got %d", len(result.Common))
	}
	if len(result.OnlyInA) != 0 || len(result.OnlyInB) != 0 || len(result.DiffValues) != 0 {
		t.Error("expected no differences")
	}
}

func TestDiff_OnlyInA(t *testing.T) {
	envA := map[string]string{"FOO": "bar", "EXTRA": "only"}
	envB := map[string]string{"FOO": "bar"}

	result := differ.Diff(envA, envB)

	if len(result.OnlyInA) != 1 || result.OnlyInA[0] != "EXTRA" {
		t.Errorf("expected EXTRA only in A, got %v", result.OnlyInA)
	}
}

func TestDiff_OnlyInB(t *testing.T) {
	envA := map[string]string{"FOO": "bar"}
	envB := map[string]string{"FOO": "bar", "NEW_KEY": "val"}

	result := differ.Diff(envA, envB)

	if len(result.OnlyInB) != 1 || result.OnlyInB[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY only in B, got %v", result.OnlyInB)
	}
}

func TestDiff_DifferingValues(t *testing.T) {
	envA := map[string]string{"DB_HOST": "localhost"}
	envB := map[string]string{"DB_HOST": "prod.db.example.com"}

	result := differ.Diff(envA, envB)

	pair, ok := result.DiffValues["DB_HOST"]
	if !ok {
		t.Fatal("expected DB_HOST to appear in DiffValues")
	}
	if pair[0] != "localhost" || pair[1] != "prod.db.example.com" {
		t.Errorf("unexpected values: %v", pair)
	}
}

func TestDiffAgainstSchema_FiltersKeys(t *testing.T) {
	s := &schema.Schema{
		Fields: []schema.Field{
			{Key: "APP_ENV", Required: true},
		},
	}

	envA := map[string]string{"APP_ENV": "staging", "SECRET": "abc"}
	envB := map[string]string{"APP_ENV": "production", "SECRET": "xyz"}

	result := differ.DiffAgainstSchema(envA, envB, s)

	if _, ok := result.DiffValues["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in DiffValues")
	}
	if _, ok := result.DiffValues["SECRET"]; ok {
		t.Error("SECRET should be filtered out by schema")
	}
}

func TestSummary(t *testing.T) {
	d := differ.DiffResult{
		OnlyInA:    []string{"A"},
		OnlyInB:    []string{"B", "C"},
		DiffValues: map[string][2]string{"X": {"1", "2"}},
		Common:     []string{"Y"},
	}

	summary := differ.Summary(d, "staging", "production")
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
