package sanitizer

import (
	"testing"
)

func TestSanitize_NoChanges(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Sanitize(env, DefaultOptions())
	if len(r.Changed) != 0 {
		t.Fatalf("expected no changes, got %d", len(r.Changed))
	}
	if r.Env["FOO"] != "bar" {
		t.Errorf("unexpected value for FOO: %q", r.Env["FOO"])
	}
}

func TestSanitize_TrimSpace(t *testing.T) {
	env := map[string]string{"KEY": "  hello  "}
	r := Sanitize(env, Options{TrimSpace: true})
	if r.Env["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", r.Env["KEY"])
	}
	if len(r.Changed) != 1 || r.Changed[0].Reason != "trimmed whitespace" {
		t.Errorf("expected trimmed whitespace change, got %+v", r.Changed)
	}
}

func TestSanitize_RemoveNewlines(t *testing.T) {
	env := map[string]string{"KEY": "line1\nline2"}
	r := Sanitize(env, Options{RemoveNewlines: true})
	if r.Env["KEY"] != "line1 line2" {
		t.Errorf("unexpected value: %q", r.Env["KEY"])
	}
	if len(r.Changed) == 0 {
		t.Error("expected a change record")
	}
}

func TestSanitize_NormalizeKeys(t *testing.T) {
	env := map[string]string{"my-key": "value"}
	r := Sanitize(env, Options{NormalizeKeys: true})
	if _, ok := r.Env["MY_KEY"]; !ok {
		t.Errorf("expected MY_KEY in output, got %v", r.Env)
	}
	if r.Env["MY_KEY"] != "value" {
		t.Errorf("unexpected value for MY_KEY: %q", r.Env["MY_KEY"])
	}
}

func TestSanitize_StripNonPrint(t *testing.T) {
	env := map[string]string{"KEY": "hello\x01world"}
	r := Sanitize(env, Options{StripNonPrint: true})
	if r.Env["KEY"] != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", r.Env["KEY"])
	}
}

func TestSanitize_MultipleOps(t *testing.T) {
	env := map[string]string{"db-host": "  localhost\n"}
	opts := Options{TrimSpace: true, RemoveNewlines: true, NormalizeKeys: true}
	r := Sanitize(env, opts)
	if v, ok := r.Env["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", r.Env)
	}
}

func TestSanitize_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": "  spaced  "}
	original := env["KEY"]
	Sanitize(env, DefaultOptions())
	if env["KEY"] != original {
		t.Error("input map was mutated")
	}
}
