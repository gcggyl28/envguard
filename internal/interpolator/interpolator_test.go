package interpolator

import (
	"os"
	"testing"
)

func TestInterpolate_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	result, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["HOST"] != "localhost" || result["PORT"] != "8080" {
		t.Errorf("expected unchanged values, got %v", result)
	}
}

func TestInterpolate_BraceRef(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "443",
	}
	result, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["BASE_URL"] != "http://example.com:443" {
		t.Errorf("expected resolved URL, got %q", result["BASE_URL"])
	}
}

func TestInterpolate_DollarRef(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello $NAME",
		"NAME":     "World",
	}
	result, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["GREETING"] != "Hello World" {
		t.Errorf("got %q", result["GREETING"])
	}
}

func TestInterpolate_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{
		"VALUE": "${OS_VAR}-suffix",
	}
	result, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["VALUE"] != "from-os-suffix" {
		t.Errorf("got %q", result["VALUE"])
	}
}

func TestInterpolate_UndefinedRefPreserved(t *testing.T) {
	env := map[string]string{
		"VALUE": "${UNDEFINED_XYZ_123}",
	}
	os.Unsetenv("UNDEFINED_XYZ_123")
	result, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["VALUE"] != "${UNDEFINED_XYZ_123}" {
		t.Errorf("expected preserved ref, got %q", result["VALUE"])
	}
}

func TestInterpolate_CircularReference(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	_, err := Interpolate(env)
	if err == nil {
		t.Fatal("expected circular reference error")
	}
	if _, ok := err.(*ErrCircularReference); !ok {
		t.Errorf("expected ErrCircularReference, got %T", err)
	}
}

func TestInterpolate_ChainedRefs(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A}-world",
		"C": "${B}!",
	}
	result, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["C"] != "hello-world!" {
		t.Errorf("got %q", result["C"])
	}
}
