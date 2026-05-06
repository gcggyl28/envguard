package transformer_test

import (
	"testing"

	"github.com/yourusername/envguard/internal/transformer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"db_host": "localhost",
		"db_port": "  5432  ",
		"api_key": "secret",
	}
}

func TestTransform_NoOps(t *testing.T) {
	env := baseEnv()
	res := transformer.Transform(env, transformer.Options{})
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
	if res.Env["db_host"] != "localhost" {
		t.Errorf("unexpected value: %s", res.Env["db_host"])
	}
}

func TestTransform_UppercaseKeys(t *testing.T) {
	res := transformer.Transform(baseEnv(), transformer.Options{UppercaseKeys: true})
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to exist")
	}
	if _, ok := res.Env["db_host"]; ok {
		t.Error("expected lowercase key to be gone")
	}
	if len(res.Changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(res.Changes))
	}
}

func TestTransform_TrimValues(t *testing.T) {
	res := transformer.Transform(baseEnv(), transformer.Options{TrimValues: true})
	if res.Env["db_port"] != "5432" {
		t.Errorf("expected trimmed value, got %q", res.Env["db_port"])
	}
	if len(res.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changes))
	}
}

func TestTransform_KeyPrefix(t *testing.T) {
	res := transformer.Transform(baseEnv(), transformer.Options{KeyPrefix: "APP_"})
	if _, ok := res.Env["APP_db_host"]; !ok {
		t.Error("expected APP_db_host")
	}
}

func TestTransform_KeySuffix(t *testing.T) {
	res := transformer.Transform(baseEnv(), transformer.Options{KeySuffix: "_V2"})
	if _, ok := res.Env["db_host_V2"]; !ok {
		t.Error("expected db_host_V2")
	}
}

func TestTransform_ReplaceKeys(t *testing.T) {
	opts := transformer.Options{
		ReplaceKeys: map[string]string{"db_host": "DATABASE_HOST"},
	}
	res := transformer.Transform(baseEnv(), opts)
	if _, ok := res.Env["DATABASE_HOST"]; !ok {
		t.Error("expected DATABASE_HOST after rename")
	}
	if _, ok := res.Env["db_host"]; ok {
		t.Error("old key should not exist")
	}
	var found bool
	for _, c := range res.Changes {
		if c.OldKey == "db_host" && c.Key == "DATABASE_HOST" {
			found = true
		}
	}
	if !found {
		t.Error("expected change record for key rename")
	}
}

func TestTransform_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	transformer.Transform(env, transformer.Options{UppercaseKeys: true, TrimValues: true})
	if _, ok := env["db_host"]; !ok {
		t.Error("original map was mutated")
	}
}
