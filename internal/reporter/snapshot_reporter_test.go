package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envguard/internal/reporter"
)

func TestFprintSnapshot_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	reporter.FprintSnapshot(&buf, ".env", nil, nil, nil)
	out := buf.String()

	if !strings.Contains(out, "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestFprintSnapshot_WithAdded(t *testing.T) {
	var buf bytes.Buffer
	reporter.FprintSnapshot(&buf, ".env", []string{"NEW_KEY"}, nil, nil)
	out := buf.String()

	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected added key in output, got: %s", out)
	}
	if strings.Contains(out, "Removed") {
		t.Errorf("unexpected Removed section in output")
	}
}

func TestFprintSnapshot_WithRemoved(t *testing.T) {
	var buf bytes.Buffer
	reporter.FprintSnapshot(&buf, ".env", nil, []string{"OLD_KEY"}, nil)
	out := buf.String()

	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected removed key in output, got: %s", out)
	}
}

func TestFprintSnapshot_WithChanged(t *testing.T) {
	var buf bytes.Buffer
	reporter.FprintSnapshot(&buf, ".env", nil, nil, []string{"PORT"})
	out := buf.String()

	if !strings.Contains(out, "~ PORT") {
		t.Errorf("expected changed key in output, got: %s", out)
	}
}

func TestFprintSnapshot_SourceInHeader(t *testing.T) {
	var buf bytes.Buffer
	reporter.FprintSnapshot(&buf, "prod.env", nil, nil, nil)
	out := buf.String()

	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected source in header, got: %s", out)
	}
}
