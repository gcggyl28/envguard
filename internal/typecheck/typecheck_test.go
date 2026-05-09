package typecheck_test

import (
	"testing"

	"github.com/yourusername/envguard/internal/typecheck"
)

var baseEnv = map[string]string{
	"PORT":       "8080",
	"RATIO":      "0.75",
	"DEBUG":      "true",
	"API_URL":    "https://api.example.com",
	"SERVER_IP":  "192.168.1.1",
	"ADMIN_EMAIL": "admin@example.com",
	"APP_NAME":   "envguard",
}

func TestCheck_NoViolations(t *testing.T) {
	types := map[string]typecheck.Type{
		"PORT":        typecheck.TypeInt,
		"RATIO":       typecheck.TypeFloat,
		"DEBUG":       typecheck.TypeBool,
		"API_URL":     typecheck.TypeURL,
		"SERVER_IP":   typecheck.TypeIP,
		"ADMIN_EMAIL": typecheck.TypeEmail,
		"APP_NAME":    typecheck.TypeString,
	}
	violations := typecheck.Check(baseEnv, types)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestCheck_InvalidInt(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	types := map[string]typecheck.Type{"PORT": typecheck.TypeInt}
	violations := typecheck.Check(env, types)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected key PORT, got %s", violations[0].Key)
	}
}

func TestCheck_InvalidBool(t *testing.T) {
	env := map[string]string{"DEBUG": "maybe"}
	types := map[string]typecheck.Type{"DEBUG": typecheck.TypeBool}
	violations := typecheck.Check(env, types)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_ValidBoolVariants(t *testing.T) {
	for _, val := range []string{"true", "false", "1", "0", "yes", "no", "TRUE", "FALSE"} {
		env := map[string]string{"FLAG": val}
		types := map[string]typecheck.Type{"FLAG": typecheck.TypeBool}
		violations := typecheck.Check(env, types)
		if len(violations) != 0 {
			t.Errorf("expected %q to be valid bool, got violation: %s", val, violations[0].Reason)
		}
	}
}

func TestCheck_InvalidURL(t *testing.T) {
	env := map[string]string{"API_URL": "not-a-url"}
	types := map[string]typecheck.Type{"API_URL": typecheck.TypeURL}
	violations := typecheck.Check(env, types)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_InvalidIP(t *testing.T) {
	env := map[string]string{"SERVER_IP": "999.999.999.999"}
	types := map[string]typecheck.Type{"SERVER_IP": typecheck.TypeIP}
	violations := typecheck.Check(env, types)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_InvalidEmail(t *testing.T) {
	env := map[string]string{"ADMIN_EMAIL": "notanemail"}
	types := map[string]typecheck.Type{"ADMIN_EMAIL": typecheck.TypeEmail}
	violations := typecheck.Check(env, types)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_SkipsMissingKeys(t *testing.T) {
	env := map[string]string{"APP_NAME": "envguard"}
	types := map[string]typecheck.Type{
		"MISSING_KEY": typecheck.TypeInt,
		"APP_NAME":    typecheck.TypeString,
	}
	violations := typecheck.Check(env, types)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for missing key, got %d", len(violations))
	}
}
