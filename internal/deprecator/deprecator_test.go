package deprecator_test

import (
	"testing"

	"github.com/user/envguard/internal/deprecator"
)

var baseRules = []deprecator.Rule{
	{Key: "OLD_API_KEY", Replacement: "API_KEY", Reason: "renamed in v2"},
	{Key: "LEGACY_DB_URL", Replacement: "DATABASE_URL", Reason: "standardised name"},
	{Key: "DISABLE_FEATURE_X", Reason: "feature removed"},
}

func TestDeprecate_NoFindings(t *testing.T) {
	env := map[string]string{"API_KEY": "abc", "DATABASE_URL": "postgres://"}
	findings := deprecator.Deprecate(env, baseRules)
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestDeprecate_SingleDeprecatedKey(t *testing.T) {
	env := map[string]string{"OLD_API_KEY": "secret", "OTHER": "val"}
	findings := deprecator.Deprecate(env, baseRules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "OLD_API_KEY" {
		t.Errorf("unexpected key %q", findings[0].Key)
	}
	if findings[0].Replacement != "API_KEY" {
		t.Errorf("unexpected replacement %q", findings[0].Replacement)
	}
}

func TestDeprecate_MultipleDeprecatedKeys(t *testing.T) {
	env := map[string]string{
		"OLD_API_KEY":    "s",
		"LEGACY_DB_URL":  "pg",
		"DISABLE_FEATURE_X": "true",
	}
	findings := deprecator.Deprecate(env, baseRules)
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings, got %d", len(findings))
	}
	// sorted by key
	if findings[0].Key != "DISABLE_FEATURE_X" {
		t.Errorf("expected DISABLE_FEATURE_X first, got %q", findings[0].Key)
	}
}

func TestDeprecate_CaseInsensitiveMatch(t *testing.T) {
	env := map[string]string{"old_api_key": "secret"}
	findings := deprecator.Deprecate(env, baseRules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestDeprecate_EmptyEnv(t *testing.T) {
	findings := deprecator.Deprecate(map[string]string{}, baseRules)
	if findings != nil {
		t.Errorf("expected nil findings for empty env")
	}
}

func TestDeprecate_EmptyRules(t *testing.T) {
	env := map[string]string{"OLD_API_KEY": "x"}
	findings := deprecator.Deprecate(env, nil)
	if findings != nil {
		t.Errorf("expected nil findings for empty rules")
	}
}

func TestDeprecate_NoReplacementField(t *testing.T) {
	env := map[string]string{"DISABLE_FEATURE_X": "1"}
	findings := deprecator.Deprecate(env, baseRules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding")
	}
	if findings[0].Replacement != "" {
		t.Errorf("expected empty replacement, got %q", findings[0].Replacement)
	}
	if findings[0].Reason != "feature removed" {
		t.Errorf("unexpected reason %q", findings[0].Reason)
	}
}
