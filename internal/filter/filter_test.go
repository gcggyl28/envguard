package filter_test

import (
	"testing"

	"github.com/yourorg/envguard/internal/filter"
)

var baseEnv = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"LOG_LEVEL":   "info",
}

func TestFilter_NoOptions(t *testing.T) {
	out, err := filter.Filter(baseEnv, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(baseEnv) {
		t.Errorf("expected %d keys, got %d", len(baseEnv), len(out))
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	out, err := filter.Filter(baseEnv, filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestFilter_BySuffix(t *testing.T) {
	out, err := filter.Filter(baseEnv, filter.Options{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_ByContains(t *testing.T) {
	out, err := filter.Filter(baseEnv, filter.Options{Contains: "DB"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_ByPattern(t *testing.T) {
	out, err := filter.Filter(baseEnv, filter.Options{Pattern: "^(APP|LOG)_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := filter.Filter(baseEnv, filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestFilter_Exclude(t *testing.T) {
	out, err := filter.Filter(baseEnv, filter.Options{Exclude: []string{"DB_PASSWORD", "LOG_LEVEL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}

func TestSummarize(t *testing.T) {
	out, _ := filter.Filter(baseEnv, filter.Options{Prefix: "APP_"})
	s := filter.Summarize(baseEnv, out)
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Included != 2 {
		t.Errorf("expected Included=2, got %d", s.Included)
	}
	if s.Excluded != 3 {
		t.Errorf("expected Excluded=3, got %d", s.Excluded)
	}
}
