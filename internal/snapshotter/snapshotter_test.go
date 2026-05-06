package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envguard/internal/snapshotter"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	tmp := filepath.Join(t.TempDir(), "snap.json")

	if err := snapshotter.Save(env, ".env", tmp); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := snapshotter.Load(tmp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if snap.Source != ".env" {
		t.Errorf("expected source .env, got %s", snap.Source)
	}
	if snap.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", snap.Env["APP_ENV"])
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshotter.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not json"), 0600)

	_, err := snapshotter.Load(tmp)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCompare_AddedRemovedChanged(t *testing.T) {
	base := &snapshotter.Snapshot{
		Timestamp: time.Now(),
		Source:    ".env",
		Env:       map[string]string{"A": "1", "B": "2", "C": "3"},
	}
	current := &snapshotter.Snapshot{
		Timestamp: time.Now(),
		Source:    ".env",
		Env:       map[string]string{"A": "1", "B": "changed", "D": "4"},
	}

	added, removed, changed := snapshotter.Compare(base, current)

	if len(added) != 1 || added[0] != "D" {
		t.Errorf("expected added=[D], got %v", added)
	}
	if len(removed) != 1 || removed[0] != "C" {
		t.Errorf("expected removed=[C], got %v", removed)
	}
	if len(changed) != 1 || changed[0] != "B" {
		t.Errorf("expected changed=[B], got %v", changed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	base := &snapshotter.Snapshot{Env: env}
	current := &snapshotter.Snapshot{Env: map[string]string{"X": "1", "Y": "2"}}

	added, removed, changed := snapshotter.Compare(base, current)
	if len(added)+len(removed)+len(changed) != 0 {
		t.Errorf("expected no changes, got added=%v removed=%v changed=%v", added, removed, changed)
	}
}
