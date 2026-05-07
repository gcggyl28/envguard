package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envguard/internal/comparator"
)

func makeResult(added, removed map[string]string, changed map[string]comparator.Change) comparator.Result {
	return comparator.Result{
		Added:   added,
		Removed: removed,
		Changed: changed,
	}
}

func TestFprintComparator_NoChanges(t *testing.T) {
	r := makeResult(nil, nil, nil)
	var buf bytes.Buffer
	FprintComparator(&buf, r, "base.env", "prod.env")
	out := buf.String()
	if !strings.Contains(out, "No differences found") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestFprintComparator_Added(t *testing.T) {
	r := makeResult(map[string]string{"NEW": "val"}, nil, nil)
	var buf bytes.Buffer
	FprintComparator(&buf, r, "a", "b")
	out := buf.String()
	if !strings.Contains(out, "Added") || !strings.Contains(out, "NEW=val") {
		t.Errorf("expected added section, got: %s", out)
	}
}

func TestFprintComparator_Removed(t *testing.T) {
	r := makeResult(nil, map[string]string{"OLD": "gone"}, nil)
	var buf bytes.Buffer
	FprintComparator(&buf, r, "a", "b")
	out := buf.String()
	if !strings.Contains(out, "Removed") || !strings.Contains(out, "OLD=gone") {
		t.Errorf("expected removed section, got: %s", out)
	}
}

func TestFprintComparator_Changed(t *testing.T) {
	changed := map[string]comparator.Change{
		"FOO": {Old: "old", New: "new"},
	}
	r := makeResult(nil, nil, changed)
	var buf bytes.Buffer
	FprintComparator(&buf, r, "a", "b")
	out := buf.String()
	if !strings.Contains(out, "Changed") || !strings.Contains(out, "FOO") {
		t.Errorf("expected changed section, got: %s", out)
	}
	if !strings.Contains(out, `"old"`) || !strings.Contains(out, `"new"`) {
		t.Errorf("expected old/new values quoted, got: %s", out)
	}
}

func TestFprintComparator_Header(t *testing.T) {
	r := makeResult(nil, nil, nil)
	var buf bytes.Buffer
	FprintComparator(&buf, r, "base.env", "prod.env")
	out := buf.String()
	if !strings.Contains(out, "base.env") || !strings.Contains(out, "prod.env") {
		t.Errorf("expected labels in header, got: %s", out)
	}
}
