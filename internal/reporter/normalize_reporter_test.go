package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envguard/internal/normalizer"
)

func TestFprintNormalize_NoChanges(t *testing.T) {
	res := normalizer.Result{
		Normalized: map[string]string{"KEY": "value"},
		Changes:    nil,
	}
	var buf bytes.Buffer
	FprintNormalize(&buf, res)
	if !strings.Contains(buf.String(), "No normalization changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestFprintNormalize_WithKeyRename(t *testing.T) {
	res := normalizer.Result{
		Normalized: map[string]string{"APP_HOST": "localhost"},
		Changes: []normalizer.Change{
			{Key: "APP_HOST", OldKey: "app_host", OldValue: "localhost", NewValue: "localhost", Reason: "key renamed"},
		},
	}
	var buf bytes.Buffer
	FprintNormalize(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "app_host") {
		t.Error("expected old key in output")
	}
	if !strings.Contains(out, "APP_HOST") {
		t.Error("expected new key in output")
	}
	if !strings.Contains(out, "key renamed") {
		t.Error("expected reason in output")
	}
}

func TestFprintNormalize_WithValueChange(t *testing.T) {
	res := normalizer.Result{
		Normalized: map[string]string{"ENV": "production"},
		Changes: []normalizer.Change{
			{Key: "ENV", OldKey: "ENV", OldValue: "Production", NewValue: "production", Reason: "value changed"},
		},
	}
	var buf bytes.Buffer
	FprintNormalize(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "Production") {
		t.Error("expected old value in output")
	}
	if !strings.Contains(out, "production") {
		t.Error("expected new value in output")
	}
}

func TestFprintNormalize_CountInHeader(t *testing.T) {
	res := normalizer.Result{
		Normalized: map[string]string{},
		Changes: []normalizer.Change{
			{Key: "A", OldKey: "a", OldValue: "x", NewValue: "x", Reason: "key renamed"},
			{Key: "B", OldKey: "b", OldValue: "y", NewValue: "y", Reason: "key renamed"},
		},
	}
	var buf bytes.Buffer
	FprintNormalize(&buf, res)
	if !strings.Contains(buf.String(), "2 change(s)") {
		t.Errorf("expected count in header, got: %s", buf.String())
	}
}
