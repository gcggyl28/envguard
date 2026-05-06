package validator_test

import (
	"testing"

	"github.com/user/envguard/internal/schema"
	"github.com/user/envguard/internal/validator"
)

func baseSchema() *schema.Schema {
	return &schema.Schema{
		Fields: []schema.Field{
			{Name: "APP_ENV", Required: true, AllowedValues: []string{"development", "staging", "production"}},
			{Name: "PORT", Required: true, Pattern: `^\d+$`},
			{Name: "LOG_LEVEL", Required: false, AllowedValues: []string{"debug", "info", "warn", "error"}},
		},
	}
}

func TestValidate_AllValid(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
		"LOG_LEVEL": "info",
	}
	report := validator.Validate(env, baseSchema())
	if !report.Valid {
		for _, r := range report.Results {
			if !r.Passed {
				t.Errorf("unexpected failure for %s: %s", r.Key, r.Message)
			}
		}
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
	}
	report := validator.Validate(env, baseSchema())
	if report.Valid {
		t.Fatal("expected report to be invalid")
	}
	found := false
	for _, r := range report.Results {
		if r.Key == "APP_ENV" && !r.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected failure for APP_ENV")
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "not-a-port",
	}
	report := validator.Validate(env, baseSchema())
	if report.Valid {
		t.Fatal("expected report to be invalid")
	}
}

func TestValidate_DisallowedValue(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "local",
		"PORT":    "3000",
	}
	report := validator.Validate(env, baseSchema())
	if report.Valid {
		t.Fatal("expected report to be invalid")
	}
}

func TestValidate_OptionalFieldAbsent(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "staging",
		"PORT":    "9000",
	}
	report := validator.Validate(env, baseSchema())
	if !report.Valid {
		t.Error("expected report to be valid when optional field is absent")
	}
}

func TestValidate_OptionalFieldInvalidValue(t *testing.T) {
	// An optional field that is present must still satisfy allowed values.
	env := map[string]string{
		"APP_ENV":   "production",
		"PORT":      "4000",
		"LOG_LEVEL": "verbose",
	}
	report := validator.Validate(env, baseSchema())
	if report.Valid {
		t.Fatal("expected report to be invalid when optional field has disallowed value")
	}
	found := false
	for _, r := range report.Results {
		if r.Key == "LOG_LEVEL" && !r.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected failure for LOG_LEVEL")
	}
}
