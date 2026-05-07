package trimmer_test

import (
	"testing"

	"github.com/user/envguard/internal/trimmer"
)

func TestTrim_NoChanges(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	res := trimmer.Trim(env, trimmer.DefaultOptions())
	if len(res.Changes) != 0 {
		t.Fatalf("expected no changes, got %v", res.Changes)
	}
	if res.Trimmed["HOST"] != "localhost" {
		t.Errorf("unexpected value: %q", res.Trimmed["HOST"])
	}
}

func TestTrim_WhitespaceValues(t *testing.T) {
	env := map[string]string{"KEY": "  hello  "}
	res := trimmer.Trim(env, trimmer.Options{TrimValues: true})
	if res.Trimmed["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", res.Trimmed["KEY"])
	}
	if len(res.Changes) != 1 || res.Changes[0] != "KEY" {
		t.Errorf("expected KEY in changes, got %v", res.Changes)
	}
}

func TestTrim_StripDoubleQuotes(t *testing.T) {
	env := map[string]string{"DB": `"postgres"`}
	res := trimmer.Trim(env, trimmer.Options{StripQuotes: true})
	if res.Trimmed["DB"] != "postgres" {
		t.Errorf("expected 'postgres', got %q", res.Trimmed["DB"])
	}
}

func TestTrim_StripSingleQuotes(t *testing.T) {
	env := map[string]string{"TOKEN": "'abc123'"}
	res := trimmer.Trim(env, trimmer.Options{StripQuotes: true})
	if res.Trimmed["TOKEN"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", res.Trimmed["TOKEN"])
	}
}

func TestTrim_MismatchedQuotesUnchanged(t *testing.T) {
	env := map[string]string{"X": `"value'`}
	res := trimmer.Trim(env, trimmer.Options{StripQuotes: true})
	if res.Trimmed["X"] != `"value'` {
		t.Errorf("expected unchanged value, got %q", res.Trimmed["X"])
	}
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %v", res.Changes)
	}
}

func TestTrim_TrimKeys(t *testing.T) {
	env := map[string]string{" SPACED ": "value"}
	res := trimmer.Trim(env, trimmer.Options{TrimKeys: true})
	if _, ok := res.Trimmed["SPACED"]; !ok {
		t.Error("expected key 'SPACED' after trimming")
	}
	if len(res.Changes) != 1 || res.Changes[0] != "SPACED" {
		t.Errorf("expected SPACED in changes, got %v", res.Changes)
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"K": "  v  "}
	original := env["K"]
	trimmer.Trim(env, trimmer.DefaultOptions())
	if env["K"] != original {
		t.Error("original map was mutated")
	}
}

func TestTrim_ChangesAreSorted(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "  z  ",
		"ALPHA": "  a  ",
		"MANGO": "  m  ",
	}
	res := trimmer.Trim(env, trimmer.Options{TrimValues: true})
	for i := 1; i < len(res.Changes); i++ {
		if res.Changes[i] < res.Changes[i-1] {
			t.Errorf("changes not sorted: %v", res.Changes)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := trimmer.DefaultOptions()
	if !opts.TrimValues {
		t.Error("expected TrimValues=true")
	}
	if opts.TrimKeys {
		t.Error("expected TrimKeys=false")
	}
	if !opts.StripQuotes {
		t.Error("expected StripQuotes=true")
	}
}
