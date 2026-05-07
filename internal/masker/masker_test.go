package masker_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envguard/internal/masker"
)

func TestMask_Full(t *testing.T) {
	opts := masker.Options{Style: masker.StyleFull}
	got := masker.Mask("secret123", opts)
	if got != "*********" {
		t.Errorf("expected all asterisks, got %q", got)
	}
}

func TestMask_Partial_RevealsEnds(t *testing.T) {
	opts := masker.Options{Style: masker.StylePartial, RevealChars: 3}
	got := masker.Mask("abcXXXXXXdef", opts)
	if !strings.HasPrefix(got, "abc") {
		t.Errorf("expected prefix 'abc', got %q", got)
	}
	if !strings.HasSuffix(got, "def") {
		t.Errorf("expected suffix 'def', got %q", got)
	}
	if !strings.Contains(got, "*") {
		t.Errorf("expected middle to be masked, got %q", got)
	}
}

func TestMask_Partial_ShortValue(t *testing.T) {
	opts := masker.Options{Style: masker.StylePartial, RevealChars: 4}
	got := masker.Mask("hi", opts)
	if strings.ContainsAny(got, "hi") {
		t.Errorf("short value should be fully masked, got %q", got)
	}
}

func TestMask_Hash(t *testing.T) {
	opts := masker.Options{Style: masker.StyleHash}
	got := masker.Mask("hello", opts)
	if got != "[*****]" {
		t.Errorf("expected '[*****]', got %q", got)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	opts := masker.DefaultOptions()
	got := masker.Mask("", opts)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMaskMap_ReturnsNewMap(t *testing.T) {
	env := map[string]string{"API_KEY": "supersecret", "HOST": "localhost"}
	opts := masker.Options{Style: masker.StyleFull}
	out := masker.MaskMap(env, opts)

	if out["API_KEY"] == "supersecret" {
		t.Error("expected API_KEY to be masked")
	}
	if env["API_KEY"] != "supersecret" {
		t.Error("original map should not be mutated")
	}
	if len(out) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(out))
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := masker.DefaultOptions()
	if opts.Style != masker.StylePartial {
		t.Errorf("expected StylePartial, got %q", opts.Style)
	}
	if opts.RevealChars != 3 {
		t.Errorf("expected RevealChars=3, got %d", opts.RevealChars)
	}
}
