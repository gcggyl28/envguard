package splitter_test

import (
	"testing"

	"github.com/user/envguard/internal/splitter"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"AWS_KEY":     "AKIA123",
	"AWS_SECRET":  "secret",
	"APP_NAME":    "envguard",
	"APP_VERSION": "1.0.0",
	"UNRELATED":   "value",
}

func TestSplit_NoRules(t *testing.T) {
	res := splitter.Split(baseEnv, nil)
	if len(res.Buckets) != 0 {
		t.Errorf("expected no buckets, got %d", len(res.Buckets))
	}
	if len(res.Leftover) != len(baseEnv) {
		t.Errorf("expected all keys in leftover, got %d", len(res.Leftover))
	}
}

func TestSplit_ByPrefix(t *testing.T) {
	rules := []splitter.Rule{
		{Name: "db", Prefixes: []string{"DB_"}},
		{Name: "aws", Prefixes: []string{"AWS_"}},
	}
	res := splitter.Split(baseEnv, rules)

	if len(res.Buckets["db"]) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(res.Buckets["db"]))
	}
	if len(res.Buckets["aws"]) != 2 {
		t.Errorf("expected 2 aws keys, got %d", len(res.Buckets["aws"]))
	}
	if len(res.Leftover) != 3 {
		t.Errorf("expected 3 leftover keys, got %d", len(res.Leftover))
	}
}

func TestSplit_CatchAllRule(t *testing.T) {
	rules := []splitter.Rule{
		{Name: "db", Prefixes: []string{"DB_"}},
		{Name: "rest", Prefixes: []string{}}, // catch-all
	}
	res := splitter.Split(baseEnv, rules)

	if len(res.Leftover) != 0 {
		t.Errorf("expected no leftover with catch-all, got %d", len(res.Leftover))
	}
	if len(res.Buckets["db"]) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(res.Buckets["db"]))
	}
	if len(res.Buckets["rest"]) != 5 {
		t.Errorf("expected 5 rest keys, got %d", len(res.Buckets["rest"]))
	}
}

func TestSplit_FirstRuleWins(t *testing.T) {
	rules := []splitter.Rule{
		{Name: "first", Prefixes: []string{"DB_"}},
		{Name: "second", Prefixes: []string{"DB_"}},
	}
	res := splitter.Split(baseEnv, rules)

	if len(res.Buckets["first"]) != 2 {
		t.Errorf("expected 2 keys in first, got %d", len(res.Buckets["first"]))
	}
	if len(res.Buckets["second"]) != 0 {
		t.Errorf("expected 0 keys in second (first rule wins), got %d", len(res.Buckets["second"]))
	}
}

func TestSplit_Summary(t *testing.T) {
	rules := []splitter.Rule{
		{Name: "db", Prefixes: []string{"DB_"}},
	}
	res := splitter.Split(baseEnv, rules)
	summary := res.Summary()

	if summary["db"] != 2 {
		t.Errorf("expected summary db=2, got %d", summary["db"])
	}
	if summary["(unmatched)"] != 5 {
		t.Errorf("expected summary unmatched=5, got %d", summary["(unmatched)"])
	}
}

func TestSplit_EmptyEnv(t *testing.T) {
	rules := []splitter.Rule{
		{Name: "db", Prefixes: []string{"DB_"}},
	}
	res := splitter.Split(map[string]string{}, rules)
	if len(res.Buckets["db"]) != 0 {
		t.Error("expected empty bucket for empty env")
	}
	if len(res.Leftover) != 0 {
		t.Error("expected empty leftover for empty env")
	}
}
