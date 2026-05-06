package sorter_test

import (
	"testing"

	"github.com/user/envguard/internal/sorter"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envguard",
	"APP_VERSION": "1.0.0",
	"LOG_LEVEL":   "info",
	"TIMEOUT":     "30",
}

func TestSort_Alpha(t *testing.T) {
	res := sorter.Sort(sampleEnv, sorter.GroupByAlpha)

	if len(res.Order) != len(sampleEnv) {
		t.Fatalf("expected %d keys, got %d", len(sampleEnv), len(res.Order))
	}

	for i := 1; i < len(res.Order); i++ {
		if res.Order[i] < res.Order[i-1] {
			t.Errorf("keys not sorted: %s before %s", res.Order[i-1], res.Order[i])
		}
	}

	if res.Groups != nil {
		t.Error("expected Groups to be nil for alpha sort")
	}
}

func TestSort_None(t *testing.T) {
	res := sorter.Sort(sampleEnv, sorter.GroupByNone)
	if len(res.Order) != len(sampleEnv) {
		t.Fatalf("expected %d keys, got %d", len(sampleEnv), len(res.Order))
	}
}

func TestSort_ByPrefix_GroupsCorrectly(t *testing.T) {
	res := sorter.Sort(sampleEnv, sorter.GroupByPrefix)

	if res.Groups == nil {
		t.Fatal("expected Groups to be populated for prefix sort")
	}

	appKeys, ok := res.Groups["APP"]
	if !ok {
		t.Fatal("expected group 'APP' to exist")
	}
	if len(appKeys) != 2 {
		t.Errorf("expected 2 APP keys, got %d", len(appKeys))
	}

	dbKeys := res.Groups["DB"]
	if len(dbKeys) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(dbKeys))
	}
}

func TestSort_ByPrefix_OrderIsGrouped(t *testing.T) {
	res := sorter.Sort(sampleEnv, sorter.GroupByPrefix)

	// APP group should come before DB group (alphabetically)
	appLast := -1
	dbFirst := len(res.Order)
	for i, k := range res.Order {
		if len(k) >= 3 && k[:3] == "APP" {
			appLast = i
		}
		if len(k) >= 2 && k[:2] == "DB" && i < dbFirst {
			dbFirst = i
		}
	}
	if appLast > dbFirst {
		t.Errorf("APP keys should appear before DB keys in prefix sort")
	}
}

func TestSort_NoUnderscore_KeyIsOwnPrefix(t *testing.T) {
	env := map[string]string{"TIMEOUT": "30", "DEBUG": "true"}
	res := sorter.Sort(env, sorter.GroupByPrefix)

	if _, ok := res.Groups["TIMEOUT"]; !ok {
		t.Error("expected TIMEOUT to be its own group prefix")
	}
	if _, ok := res.Groups["DEBUG"]; !ok {
		t.Error("expected DEBUG to be its own group prefix")
	}
}
