package requirer_test

import (
	"testing"

	"github.com/user/envguard/internal/requirer"
	"github.com/user/envguard/internal/schema"
)

func baseSchema() *schema.Schema {
	return &schema.Schema{
		Fields: []schema.Field{
			{Key: "DATABASE_URL", Required: true},
			{Key: "API_KEY", Required: true},
			{Key: "LOG_LEVEL", Required: false, Default: "info"},
		},
	}
}

func TestCheck_AllPresent(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "secret",
	}
	res := requirer.Check(env, baseSchema())
	if !res.Passed() {
		t.Fatalf("expected no findings, got %v", res.Findings)
	}
	if res.Checked != 2 {
		t.Errorf("expected 2 checked, got %d", res.Checked)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	env := map[string]string{
		"API_KEY": "secret",
	}
	res := requirer.Check(env, baseSchema())
	if res.Passed() {
		t.Fatal("expected findings for missing DATABASE_URL")
	}
	if len(res.Findings) != 1 || res.Findings[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected findings: %v", res.Findings)
	}
	if res.Findings[0].Reason != "missing" {
		t.Errorf("expected reason 'missing', got %q", res.Findings[0].Reason)
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "",
		"API_KEY":      "secret",
	}
	res := requirer.Check(env, baseSchema())
	if res.Passed() {
		t.Fatal("expected findings for empty DATABASE_URL")
	}
	if res.Findings[0].Reason != "empty value" {
		t.Errorf("expected reason 'empty value', got %q", res.Findings[0].Reason)
	}
}

func TestCheck_OptionalNotChecked(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "secret",
		// LOG_LEVEL intentionally absent — it is optional
	}
	res := requirer.Check(env, baseSchema())
	if !res.Passed() {
		t.Fatalf("optional key absence should not produce findings: %v", res.Findings)
	}
}

func TestCheck_EmptySchema(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	res := requirer.Check(env, &schema.Schema{})
	if !res.Passed() {
		t.Fatalf("empty schema should produce no findings")
	}
	if res.Checked != 0 {
		t.Errorf("expected 0 checked, got %d", res.Checked)
	}
}
