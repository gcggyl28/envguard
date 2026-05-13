package reporter

import (
	"strings"
	"testing"

	"github.com/user/envguard/internal/sanitizer"
)

func TestFprintSanitize_NoChanges(t *testing.T) {
	result := sanitizer.Result{
		Env:     map[string]string{"FOO": "bar"},
		Changed: nil,
	}
	var buf strings.Builder
	FprintSanitize(&buf, result)
	if !strings.Contains(buf.String(), "No sanitization changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestFprintSanitize_WithChanges(t *testing.T) {
	result := sanitizer.Result{
		Env: map[string]string{"KEY": "hello"},
		Changed: []sanitizer.Change{
			{Key: "KEY", Before: "  hello  ", After: "hello", Reason: "trimmed whitespace"},
		},
	}
	var buf strings.Builder
	FprintSanitize(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "1 change") {
		t.Errorf("expected change count, got: %s", out)
	}
	if !strings.Contains(out, "trimmed whitespace") {
		t.Errorf("expected reason in output, got: %s", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestFprintSanitize_SortedOutput(t *testing.T) {
	result := sanitizer.Result{
		Env: map[string]string{"ZEBRA": "z", "ALPHA": "a"},
		Changed: []sanitizer.Change{
			{Key: "ZEBRA", Before: "z ", After: "z", Reason: "trimmed whitespace"},
			{Key: "ALPHA", Before: "a ", After: "a", Reason: "trimmed whitespace"},
		},
	}
	var buf strings.Builder
	FprintSanitize(&buf, result)
	out := buf.String()
	alpha := strings.Index(out, "ALPHA")
	zebra := strings.Index(out, "ZEBRA")
	if alpha == -1 || zebra == -1 || alpha > zebra {
		t.Errorf("expected ALPHA before ZEBRA in output: %s", out)
	}
}

func TestFprintSanitize_TotalKeyCount(t *testing.T) {
	result := sanitizer.Result{
		Env:     map[string]string{"A": "1", "B": "2", "C": "3"},
		Changed: []sanitizer.Change{{Key: "A", Before: " 1", After: "1", Reason: "trimmed whitespace"}},
	}
	var buf strings.Builder
	FprintSanitize(&buf, result)
	if !strings.Contains(buf.String(), "Total keys after sanitization: 3") {
		t.Errorf("expected total key count, got: %s", buf.String())
	}
}
