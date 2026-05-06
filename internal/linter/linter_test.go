package linter

import (
	"testing"
)

func TestLint_NoHints(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	hints := Lint(env)
	// DATABASE_URL looks like a URL but is not quoted — expect 1 hint
	if len(hints) != 1 {
		t.Errorf("expected 1 hint (unquoted URL), got %d", len(hints))
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"appName": "myapp",
	}
	hints := Lint(env)
	if !containsMessage(hints, "uppercase") {
		t.Error("expected hint about uppercase key convention")
	}
}

func TestLint_KeyWithSpaces(t *testing.T) {
	env := map[string]string{
		"MY KEY": "value",
	}
	hints := Lint(env)
	if !containsMessage(hints, "spaces") {
		t.Error("expected hint about key containing spaces")
	}
}

func TestLint_ShortSecret(t *testing.T) {
	env := map[string]string{
		"API_SECRET": "abc123",
	}
	hints := Lint(env)
	if !containsMessage(hints, "short") {
		t.Error("expected hint about short secret value")
	}
}

func TestLint_LongSecretNoHint(t *testing.T) {
	env := map[string]string{
		"API_TOKEN": "supersecretlongvalue1234567890",
	}
	hints := Lint(env)
	for _, h := range hints {
		if h.Key == "API_TOKEN" {
			t.Errorf("unexpected hint for long secret: %s", h.Message)
		}
	}
}

func TestLint_UnquotedURL(t *testing.T) {
	env := map[string]string{
		"CALLBACK_URL": "https://example.com/callback",
	}
	hints := Lint(env)
	if !containsMessage(hints, "quoted") {
		t.Error("expected hint about unquoted URL")
	}
}

func TestLint_QuotedURLNoHint(t *testing.T) {
	env := map[string]string{
		"CALLBACK_URL": `"https://example.com/callback"`,
	}
	hints := Lint(env)
	for _, h := range hints {
		if h.Key == "CALLBACK_URL" {
			t.Errorf("unexpected hint for quoted URL: %s", h.Message)
		}
	}
}

func containsMessage(hints []Hint, substr string) bool {
	for _, h := range hints {
		if contains(h.Message, substr) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, sub string) bool {
	for i := range s {
		if i+len(sub) <= len(s) && s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
