package pinner_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envguard/internal/pinner"
)

func TestPin_WritesFile(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}
	path := filepath.Join(t.TempDir(), "pinned.json")
	p, err := pinner.Pin(env, []string{"DB_HOST", "API_KEY"}, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Keys["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", p.Keys["DB_HOST"])
	}
	data, _ := os.ReadFile(path)
	var loaded pinner.PinnedEnv
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("file not valid JSON: %v", err)
	}
}

func TestPin_MissingKey(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	path := filepath.Join(t.TempDir(), "pinned.json")
	_, err := pinner.Pin(env, []string{"MISSING_KEY"}, path)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"PORT": "8080", "ENV": "production"}
	path := filepath.Join(t.TempDir(), "pinned.json")
	_, err := pinner.Pin(env, []string{"PORT", "ENV"}, path)
	if err != nil {
		t.Fatalf("pin: %v", err)
	}
	loaded, err := pinner.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Keys["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", loaded.Keys["PORT"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := pinner.Load("/nonexistent/pinned.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheckDrift_NoChanges(t *testing.T) {
	pinned := &pinner.PinnedEnv{Keys: map[string]string{"A": "1", "B": "2"}}
	env := map[string]string{"A": "1", "B": "2"}
	result := pinner.CheckDrift(pinned, env)
	if len(result.Changed) != 0 || len(result.Removed) != 0 {
		t.Errorf("expected no drift, got changed=%v removed=%v", result.Changed, result.Removed)
	}
}

func TestCheckDrift_Changed(t *testing.T) {
	pinned := &pinner.PinnedEnv{Keys: map[string]string{"DB": "old"}}
	env := map[string]string{"DB": "new"}
	result := pinner.CheckDrift(pinned, env)
	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(result.Changed))
	}
	if result.Changed[0].Pinned != "old" || result.Changed[0].Current != "new" {
		t.Errorf("unexpected drift entry: %+v", result.Changed[0])
	}
}

func TestCheckDrift_Removed(t *testing.T) {
	pinned := &pinner.PinnedEnv{Keys: map[string]string{"GONE": "val"}}
	env := map[string]string{}
	result := pinner.CheckDrift(pinned, env)
	if len(result.Removed) != 1 || result.Removed[0] != "GONE" {
		t.Errorf("expected GONE in removed, got %v", result.Removed)
	}
}
