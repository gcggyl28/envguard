package deduplicator

import (
	"testing"
)

func pairs(kv ...string) [][2]string {
	if len(kv)%2 != 0 {
		panic("pairs: odd number of arguments")
	}
	out := make([][2]string, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		out = append(out, [2]string{kv[i], kv[i+1]})
	}
	return out
}

func TestDeduplicate_NoDuplicates(t *testing.T) {
	r := Deduplicate(pairs("A", "1", "B", "2", "C", "3"), KeepLast)
	if len(r.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", r.Duplicates)
	}
	if r.Unique["A"] != "1" || r.Unique["B"] != "2" || r.Unique["C"] != "3" {
		t.Errorf("unexpected unique map: %v", r.Unique)
	}
}

func TestDeduplicate_KeepLast(t *testing.T) {
	r := Deduplicate(pairs("HOST", "localhost", "PORT", "5432", "HOST", "prod.example.com"), KeepLast)
	if r.Unique["HOST"] != "prod.example.com" {
		t.Errorf("KeepLast: expected prod.example.com, got %q", r.Unique["HOST"])
	}
	if _, ok := r.Duplicates["HOST"]; !ok {
		t.Error("expected HOST in duplicates map")
	}
}

func TestDeduplicate_KeepFirst(t *testing.T) {
	r := Deduplicate(pairs("HOST", "localhost", "HOST", "prod.example.com"), KeepFirst)
	if r.Unique["HOST"] != "localhost" {
		t.Errorf("KeepFirst: expected localhost, got %q", r.Unique["HOST"])
	}
}

func TestDeduplicate_DuplicateHistoryOrdered(t *testing.T) {
	r := Deduplicate(pairs("KEY", "v1", "KEY", "v2", "KEY", "v3"), KeepLast)
	hist := r.Duplicates["KEY"]
	if len(hist) != 3 {
		t.Fatalf("expected 3 history entries, got %d: %v", len(hist), hist)
	}
	if hist[0] != "v1" || hist[1] != "v2" || hist[2] != "v3" {
		t.Errorf("unexpected history order: %v", hist)
	}
}

func TestDeduplicate_MultipleKeysWithDuplicates(t *testing.T) {
	r := Deduplicate(pairs("A", "1", "B", "x", "A", "2", "B", "y"), KeepLast)
	if len(r.Duplicates) != 2 {
		t.Errorf("expected 2 duplicate keys, got %d", len(r.Duplicates))
	}
	if r.Unique["A"] != "2" {
		t.Errorf("expected A=2, got %q", r.Unique["A"])
	}
	if r.Unique["B"] != "y" {
		t.Errorf("expected B=y, got %q", r.Unique["B"])
	}
}

func TestStrategyFromString(t *testing.T) {
	if StrategyFromString("first") != KeepFirst {
		t.Error("expected KeepFirst for 'first'")
	}
	if StrategyFromString("last") != KeepLast {
		t.Error("expected KeepLast for 'last'")
	}
	if StrategyFromString("") != KeepLast {
		t.Error("expected KeepLast for empty string")
	}
}
