package reporter

import (
	"strings"
	"testing"

	"github.com/user/envguard/internal/grouper"
)

func TestFprintGroup_NoGroups(t *testing.T) {
	result := grouper.GroupResult{}
	env := map[string]string{}
	var buf strings.Builder
	FprintGroup(&buf, result, env)
	out := buf.String()
	if !strings.Contains(out, "0 total") {
		t.Errorf("expected '0 total' in output, got:\n%s", out)
	}
}

func TestFprintGroup_WithGroups(t *testing.T) {
	result := grouper.GroupResult{
		Groups: []grouper.Group{
			{Name: "DB", Keys: []string{"DB_HOST", "DB_PORT"}},
			{Name: "APP", Keys: []string{"APP_NAME"}},
		},
	}
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "envguard",
	}
	var buf strings.Builder
	FprintGroup(&buf, result, env)
	out := buf.String()
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected '[DB]' group header in output")
	}
	if !strings.Contains(out, "[APP]") {
		t.Errorf("expected '[APP]' group header in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST key in output")
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected value 'localhost' in output")
	}
	if !strings.Contains(out, "3 total") {
		t.Errorf("expected '3 total' in output, got:\n%s", out)
	}
}

func TestFprintGroup_WithUngrouped(t *testing.T) {
	result := grouper.GroupResult{
		Groups:    []grouper.Group{{Name: "DB", Keys: []string{"DB_HOST"}}},
		Ungrouped: []string{"PORT"},
	}
	env := map[string]string{
		"DB_HOST": "localhost",
		"PORT":    "8080",
	}
	var buf strings.Builder
	FprintGroup(&buf, result, env)
	out := buf.String()
	if !strings.Contains(out, "[ungrouped]") {
		t.Errorf("expected '[ungrouped]' section in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in ungrouped section")
	}
	if !strings.Contains(out, "Ungrouped: 1") {
		t.Errorf("expected 'Ungrouped: 1' in summary line, got:\n%s", out)
	}
}

func TestFprintGroup_EmptyGroupSkipped(t *testing.T) {
	result := grouper.GroupResult{
		Groups: []grouper.Group{
			{Name: "EMPTY", Keys: []string{}},
			{Name: "APP", Keys: []string{"APP_NAME"}},
		},
	}
	env := map[string]string{"APP_NAME": "envguard"}
	var buf strings.Builder
	FprintGroup(&buf, result, env)
	out := buf.String()
	if strings.Contains(out, "[EMPTY]") {
		t.Errorf("empty group should be skipped in output")
	}
	if !strings.Contains(out, "[APP]") {
		t.Errorf("expected '[APP]' in output")
	}
}
