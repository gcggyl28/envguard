package scoper_test

import (
	"testing"

	"github.com/yourusername/envguard/internal/scoper"
)

var baseEnv = map[string]string{
	"PROD_DB_HOST":     "prod.db.local",
	"PROD_API_KEY":     "secret-prod",
	"STAGING_DB_HOST":  "staging.db.local",
	"STAGING_API_KEY":  "secret-staging",
	"SHARED_LOG_LEVEL": "info",
}

func TestApply_NoPrefix_IncludesAll(t *testing.T) {
	scope := scoper.Scope{Name: "all"}
	result := scoper.Apply(baseEnv, scope)
	if len(result.Included) != len(baseEnv) {
		t.Errorf("expected %d included, got %d", len(baseEnv), len(result.Included))
	}
	if len(result.Excluded) != 0 {
		t.Errorf("expected 0 excluded, got %d", len(result.Excluded))
	}
}

func TestApply_ProdScope(t *testing.T) {
	scope := scoper.Scope{Name: "production", Prefixes: []string{"PROD_"}}
	result := scoper.Apply(baseEnv, scope)
	if len(result.Included) != 2 {
		t.Errorf("expected 2 included, got %d", len(result.Included))
	}
	if _, ok := result.Included["PROD_DB_HOST"]; !ok {
		t.Error("expected PROD_DB_HOST in included")
	}
	if _, ok := result.Excluded["STAGING_DB_HOST"]; !ok {
		t.Error("expected STAGING_DB_HOST in excluded")
	}
}

func TestApply_ScopeNamePreserved(t *testing.T) {
	scope := scoper.Scope{Name: "staging", Prefixes: []string{"STAGING_"}}
	result := scoper.Apply(baseEnv, scope)
	if result.Scope != "staging" {
		t.Errorf("expected scope name 'staging', got %q", result.Scope)
	}
}

func TestApply_CaseInsensitivePrefix(t *testing.T) {
	env := map[string]string{"prod_secret": "val"}
	scope := scoper.Scope{Name: "prod", Prefixes: []string{"PROD_"}}
	result := scoper.Apply(env, scope)
	if len(result.Included) != 1 {
		t.Errorf("expected 1 included (case-insensitive), got %d", len(result.Included))
	}
}

func TestStrip_RemovesPrefix(t *testing.T) {
	env := map[string]string{
		"PROD_DB_HOST": "prod.db.local",
		"PROD_API_KEY": "secret",
	}
	scope := scoper.Scope{Name: "production", Prefixes: []string{"PROD_"}}
	stripped := scoper.Strip(env, scope)
	if _, ok := stripped["DB_HOST"]; !ok {
		t.Error("expected DB_HOST after stripping PROD_ prefix")
	}
	if _, ok := stripped["API_KEY"]; !ok {
		t.Error("expected API_KEY after stripping PROD_ prefix")
	}
}

func TestStrip_NoMatchRetainsKey(t *testing.T) {
	env := map[string]string{"SHARED_LOG": "info"}
	scope := scoper.Scope{Name: "prod", Prefixes: []string{"PROD_"}}
	stripped := scoper.Strip(env, scope)
	if _, ok := stripped["SHARED_LOG"]; !ok {
		t.Error("expected SHARED_LOG to be retained unchanged")
	}
}
