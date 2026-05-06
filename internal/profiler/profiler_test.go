package profiler_test

import (
	"testing"

	"github.com/user/envguard/internal/profiler"
	"github.com/user/envguard/internal/schema"
)

func baseSchema() *schema.Schema {
	return &schema.Schema{
		Fields: []schema.Field{
			{Key: "APP_ENV", Required: true},
			{Key: "DB_PASSWORD", Required: true},
			{Key: "LOG_LEVEL", Required: false, Default: "info"},
			{Key: "API_KEY", Required: false},
			{Key: "PORT", Required: false, Default: "8080"},
		},
	}
}

func TestAnalyze_TotalKeys(t *testing.T) {
	s := baseSchema()
	env := map[string]string{"APP_ENV": "prod", "DB_PASSWORD": "secret"}
	p := profiler.Analyze(env, s)
	if p.TotalKeys != 5 {
		t.Errorf("expected TotalKeys=5, got %d", p.TotalKeys)
	}
}

func TestAnalyze_RequiredAndOptional(t *testing.T) {
	s := baseSchema()
	env := map[string]string{}
	p := profiler.Analyze(env, s)

	if len(p.RequiredKeys) != 2 {
		t.Errorf("expected 2 required keys, got %d", len(p.RequiredKeys))
	}
	if len(p.OptionalKeys) != 3 {
		t.Errorf("expected 3 optional keys, got %d", len(p.OptionalKeys))
	}
}

func TestAnalyze_SensitiveKeys(t *testing.T) {
	s := baseSchema()
	env := map[string]string{}
	p := profiler.Analyze(env, s)

	// DB_PASSWORD and API_KEY should be detected as sensitive
	if len(p.SensitiveKeys) != 2 {
		t.Errorf("expected 2 sensitive keys, got %d: %v", len(p.SensitiveKeys), p.SensitiveKeys)
	}
}

func TestAnalyze_KeysWithDefault(t *testing.T) {
	s := baseSchema()
	env := map[string]string{}
	p := profiler.Analyze(env, s)

	if len(p.KeysWithDefault) != 2 {
		t.Errorf("expected 2 keys with defaults, got %d", len(p.KeysWithDefault))
	}
}

func TestAnalyze_UndeclaredKeys(t *testing.T) {
	s := baseSchema()
	env := map[string]string{
		"APP_ENV":  "prod",
		"UNKNOWN1": "val1",
		"UNKNOWN2": "val2",
	}
	p := profiler.Analyze(env, s)

	if len(p.UndeclaredKeys) != 2 {
		t.Errorf("expected 2 undeclared keys, got %d: %v", len(p.UndeclaredKeys), p.UndeclaredKeys)
	}
}

func TestAnalyze_NoUndeclaredKeys(t *testing.T) {
	s := baseSchema()
	env := map[string]string{"APP_ENV": "prod", "PORT": "9000"}
	p := profiler.Analyze(env, s)

	if len(p.UndeclaredKeys) != 0 {
		t.Errorf("expected no undeclared keys, got %v", p.UndeclaredKeys)
	}
}
