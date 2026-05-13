package normalizer

import (
	"testing"
)

func TestNormalize_NoChanges(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "PORT": "8080"}
	opts := Options{}
	res := Normalize(env, opts)
	if len(res.Changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(res.Changes))
	}
	if res.Normalized["APP_HOST"] != "localhost" {
		t.Errorf("unexpected value for APP_HOST")
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"app_host": "localhost"}
	opts := Options{UppercaseKeys: true}
	res := Normalize(env, opts)
	if _, ok := res.Normalized["APP_HOST"]; !ok {
		t.Error("expected APP_HOST to exist after uppercasing")
	}
	if len(res.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changes))
	}
	if res.Changes[0].OldKey != "app_host" {
		t.Errorf("unexpected OldKey: %s", res.Changes[0].OldKey)
	}
}

func TestNormalize_ReplaceHyphens(t *testing.T) {
	env := map[string]string{"my-key": "value"}
	opts := Options{ReplaceHyphens: true}
	res := Normalize(env, opts)
	if _, ok := res.Normalized["my_key"]; !ok {
		t.Error("expected my_key after hyphen replacement")
	}
}

func TestNormalize_TrimSpace(t *testing.T) {
	env := map[string]string{" KEY ": "  val  "}
	opts := Options{TrimSpace: true}
	res := Normalize(env, opts)
	if v, ok := res.Normalized["KEY"]; !ok || v != "val" {
		t.Errorf("expected KEY=val, got %q=%q", "KEY", v)
	}
}

func TestNormalize_LowercaseValues(t *testing.T) {
	env := map[string]string{"ENV": "Production"}
	opts := Options{LowercaseValues: true}
	res := Normalize(env, opts)
	if res.Normalized["ENV"] != "production" {
		t.Errorf("expected 'production', got %q", res.Normalized["ENV"])
	}
}

func TestNormalize_DefaultOptions(t *testing.T) {
	env := map[string]string{"my-service-url": " https://example.com "}
	opts := DefaultOptions()
	res := Normalize(env, opts)
	if _, ok := res.Normalized["MY_SERVICE_URL"]; !ok {
		t.Error("expected MY_SERVICE_URL after default normalization")
	}
	if res.Normalized["MY_SERVICE_URL"] != "https://example.com" {
		t.Errorf("expected trimmed value, got %q", res.Normalized["MY_SERVICE_URL"])
	}
}

func TestNormalize_ReasonContainsKeyRenamed(t *testing.T) {
	env := map[string]string{"lower": "value"}
	opts := Options{UppercaseKeys: true}
	res := Normalize(env, opts)
	if len(res.Changes) == 0 {
		t.Fatal("expected a change")
	}
	if res.Changes[0].Reason == "" {
		t.Error("expected non-empty reason")
	}
}
