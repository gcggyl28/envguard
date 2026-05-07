package grouper

import (
	"testing"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envguard",
	"APP_VERSION": "1.0",
	"LOG_LEVEL":   "info",
	"PORT":        "8080",
}

func TestGroup_Default(t *testing.T) {
	result := Group(sampleEnv, Options{})
	if len(result.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result.Groups))
	}
	if result.Groups[0].Name != "all" {
		t.Errorf("expected group name 'all', got %q", result.Groups[0].Name)
	}
	if len(result.Groups[0].Keys) != len(sampleEnv) {
		t.Errorf("expected %d keys, got %d", len(sampleEnv), len(result.Groups[0].Keys))
	}
}

func TestGroup_ByPrefix(t *testing.T) {
	result := Group(sampleEnv, Options{ByPrefix: true})
	groupNames := map[string]bool{}
	for _, g := range result.Groups {
		groupNames[g.Name] = true
	}
	for _, expected := range []string{"DB", "APP", "LOG"} {
		if !groupNames[expected] {
			t.Errorf("expected group %q not found", expected)
		}
	}
	if len(result.Ungrouped) != 1 || result.Ungrouped[0] != "PORT" {
		t.Errorf("expected PORT in ungrouped, got %v", result.Ungrouped)
	}
}

func TestGroup_ByPrefix_KeysAreSorted(t *testing.T) {
	result := Group(sampleEnv, Options{ByPrefix: true})
	for _, g := range result.Groups {
		for i := 1; i < len(g.Keys); i++ {
			if g.Keys[i-1] > g.Keys[i] {
				t.Errorf("keys not sorted in group %q: %v", g.Name, g.Keys)
			}
		}
	}
}

func TestGroup_CustomGroups(t *testing.T) {
	opts := Options{
		CustomGroups: map[string][]string{
			"database": {"DB"},
			"app":      {"APP"},
		},
	}
	result := Group(sampleEnv, opts)
	groupMap := map[string][]string{}
	for _, g := range result.Groups {
		groupMap[g.Name] = g.Keys
	}
	if len(groupMap["database"]) != 2 {
		t.Errorf("expected 2 database keys, got %v", groupMap["database"])
	}
	if len(groupMap["app"]) != 2 {
		t.Errorf("expected 2 app keys, got %v", groupMap["app"])
	}
	// LOG_LEVEL and PORT should be ungrouped
	if len(result.Ungrouped) != 2 {
		t.Errorf("expected 2 ungrouped keys, got %v", result.Ungrouped)
	}
}

func TestGroup_CustomGroups_ExactMatch(t *testing.T) {
	env := map[string]string{
		"PORT":    "8080",
		"HOST":    "localhost",
		"TIMEOUT": "30",
	}
	opts := Options{
		CustomGroups: map[string][]string{
			"network": {"PORT", "HOST"},
		},
	}
	result := Group(env, opts)
	if len(result.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result.Groups))
	}
	if len(result.Groups[0].Keys) != 2 {
		t.Errorf("expected 2 keys in network group, got %v", result.Groups[0].Keys)
	}
	if len(result.Ungrouped) != 1 || result.Ungrouped[0] != "TIMEOUT" {
		t.Errorf("expected TIMEOUT ungrouped, got %v", result.Ungrouped)
	}
}
