package coercer

import (
	"testing"
)

var baseEnv = map[string]string{
	"ENABLED":  "yes",
	"PORT":     " 8080 ",
	"RATIO":    "0.75",
	"APP_NAME": "  myapp  ",
	"UNKNOWN":  "whatever",
}

func TestCoerce_BoolVariants(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"yes", "true"}, {"on", "true"}, {"1", "true"}, {"TRUE", "true"},
		{"no", "false"}, {"off", "false"}, {"0", "false"}, {"False", "false"},
	}
	for _, tt := range tests {
		env := map[string]string{"FLAG": tt.input}
		out, results := Coerce(env, []Rule{{Key: "FLAG", Target: TypeBool}})
		if out["FLAG"] != tt.expected {
			t.Errorf("input %q: got %q, want %q", tt.input, out["FLAG"], tt.expected)
		}
		if results[0].Error != "" {
			t.Errorf("unexpected error: %s", results[0].Error)
		}
	}
}

func TestCoerce_InvalidBool(t *testing.T) {
	env := map[string]string{"FLAG": "maybe"}
	out, results := Coerce(env, []Rule{{Key: "FLAG", Target: TypeBool}})
	if out["FLAG"] != "maybe" {
		t.Errorf("expected original value preserved, got %q", out["FLAG"])
	}
	if results[0].Error == "" {
		t.Error("expected error for invalid bool")
	}
}

func TestCoerce_IntTrimsSpace(t *testing.T) {
	env := map[string]string{"PORT": " 8080 "}
	out, results := Coerce(env, []Rule{{Key: "PORT", Target: TypeInt}})
	if out["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", out["PORT"])
	}
	if !results[0].Changed {
		t.Error("expected Changed=true")
	}
}

func TestCoerce_InvalidInt(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	_, results := Coerce(env, []Rule{{Key: "PORT", Target: TypeInt}})
	if results[0].Error == "" {
		t.Error("expected error for invalid int")
	}
}

func TestCoerce_Float(t *testing.T) {
	env := map[string]string{"RATIO": "0.75"}
	out, results := Coerce(env, []Rule{{Key: "RATIO", Target: TypeFloat}})
	if out["RATIO"] != "0.75" {
		t.Errorf("expected '0.75', got %q", out["RATIO"])
	}
	if results[0].Error != "" {
		t.Errorf("unexpected error: %s", results[0].Error)
	}
}

func TestCoerce_StringTrimsSpace(t *testing.T) {
	env := map[string]string{"APP_NAME": "  myapp  "}
	out, results := Coerce(env, []Rule{{Key: "APP_NAME", Target: TypeString}})
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected 'myapp', got %q", out["APP_NAME"])
	}
	if !results[0].Changed {
		t.Error("expected Changed=true")
	}
}

func TestCoerce_MissingKeySkipped(t *testing.T) {
	env := map[string]string{"A": "1"}
	_, results := Coerce(env, []Rule{{Key: "MISSING", Target: TypeBool}})
	if len(results) != 0 {
		t.Errorf("expected no results for missing key, got %d", len(results))
	}
}

func TestCoerce_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"FLAG": "yes"}
	Coerce(env, []Rule{{Key: "FLAG", Target: TypeBool}})
	if env["FLAG"] != "yes" {
		t.Error("input map was mutated")
	}
}
