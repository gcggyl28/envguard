package redactor_test

import (
	"testing"

	"github.com/yourorg/envguard/internal/redactor"
)

func TestIsSensitive_Matches(t *testing.T) {
	cases := []string{
		"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY",
		"AWS_SECRET_ACCESS_KEY", "JWT_SECRET", "APP_CREDENTIAL",
	}
	for _, key := range cases {
		if !redactor.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_NoMatch(t *testing.T) {
	cases := []string{"APP_ENV", "PORT", "LOG_LEVEL", "DATABASE_HOST"}
	for _, key := range cases {
		if redactor.IsSensitive(key) {
			t.Errorf("expected %q NOT to be sensitive", key)
		}
	}
}

func TestRedact_MasksSensitiveValues(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_ENV":     "production",
		"API_KEY":     "abc123",
		"PORT":        "8080",
	}
	result := redactor.Redact(env)

	if result["DB_PASSWORD"] != "********" {
		t.Errorf("DB_PASSWORD should be masked, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "********" {
		t.Errorf("API_KEY should be masked, got %q", result["API_KEY"])
	}
	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged, got %q", result["APP_ENV"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", result["PORT"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret"}
	_ = redactor.Redact(env)
	if env["DB_PASSWORD"] != "secret" {
		t.Error("Redact must not mutate the input map")
	}
}

func TestRedactList_ExplicitKeys(t *testing.T) {
	env := map[string]string{
		"CUSTOM_TOKEN": "tok",
		"APP_NAME":     "envguard",
		"SOME_VALUE":   "hello",
	}
	result := redactor.RedactList(env, []string{"CUSTOM_TOKEN", "SOME_VALUE"})

	if result["CUSTOM_TOKEN"] != "********" {
		t.Errorf("CUSTOM_TOKEN should be masked")
	}
	if result["SOME_VALUE"] != "********" {
		t.Errorf("SOME_VALUE should be masked")
	}
	if result["APP_NAME"] != "envguard" {
		t.Errorf("APP_NAME should be unchanged")
	}
}

func TestRedactList_EmptyList(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret", "PORT": "3000"}
	result := redactor.RedactList(env, nil)
	for k, v := range env {
		if result[k] != v {
			t.Errorf("key %q should be unchanged when list is empty", k)
		}
	}
}
