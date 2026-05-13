package freezer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envguard/internal/freezer"
)

func tempFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "frozen-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestFreeze_AllKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	out := tempFile(t)
	res, err := freezer.Freeze(env, "test.env", out, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Frozen) != 2 {
		t.Errorf("expected 2 frozen, got %d", len(res.Frozen))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestFreeze_AllowKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := tempFile(t)
	res, err := freezer.Freeze(env, "test.env", out, []string{"A", "C"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Frozen) != 2 {
		t.Errorf("expected 2 frozen, got %d", len(res.Frozen))
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "B" {
		t.Errorf("expected B skipped, got %v", res.Skipped)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"X": "hello", "Y": "world"}
	out := tempFile(t)
	_, err := freezer.Freeze(env, "src.env", out, nil)
	if err != nil {
		t.Fatal(err)
	}
	fe, err := freezer.Load(out)
	if err != nil {
		t.Fatal(err)
	}
	if fe.Source != "src.env" {
		t.Errorf("expected source src.env, got %s", fe.Source)
	}
	if fe.Values["X"] != "hello" {
		t.Errorf("expected X=hello")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := freezer.Load(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestThaw_ReturnsMap(t *testing.T) {
	env := map[string]string{"K": "v"}
	out := tempFile(t)
	freezer.Freeze(env, "e", out, nil)
	fe, _ := freezer.Load(out)
	thawed := freezer.Thaw(fe)
	if thawed["K"] != "v" {
		t.Errorf("expected K=v after thaw")
	}
}

func TestFreeze_KeysAreSorted(t *testing.T) {
	env := map[string]string{"Z": "1", "A": "2", "M": "3"}
	out := tempFile(t)
	res, _ := freezer.Freeze(env, "", out, nil)
	if res.Frozen[0] != "A" || res.Frozen[1] != "M" || res.Frozen[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", res.Frozen)
	}
}
