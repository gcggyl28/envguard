package envparser_test

import (
	"strings"
	"testing"

	"github.com/user/envguard/internal/envparser"
)

func TestParse_Basic(t *testing.T) {
	input := `
# comment
APP_ENV=production
PORT=8080
`
	env, err := envparser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env["APP_ENV"])
	}
	if env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", env["PORT"])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	input := `SECRET="my secret value"
TOKEN='abc123'`
	env, err := envparser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret value" {
		t.Errorf("expected stripped quotes, got %q", env["SECRET"])
	}
	if env["TOKEN"] != "abc123" {
		t.Errorf("expected stripped quotes, got %q", env["TOKEN"])
	}
}

func TestParse_InvalidLine(t *testing.T) {
	input := `INVALID_LINE`
	_, err := envparser.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}

func TestParse_EmptyKey(t *testing.T) {
	input := `=value`
	_, err := envparser.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestParse_SkipsBlanksAndComments(t *testing.T) {
	input := `
# this is a comment

KEY=val
`
	env, err := envparser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
}
