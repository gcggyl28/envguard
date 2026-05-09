package flattener

import (
	"testing"
)

func TestFlatten_NoChanges(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	opts := DefaultOptions()
	res := Flatten(env, opts)

	if len(res.Flattened) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(res.Flattened))
	}
	if res.Flattened["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", res.Flattened["HOST"])
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", res.Renamed)
	}
}

func TestFlatten_Uppercase(t *testing.T) {
	env := map[string]string{
		"db__host": "localhost",
		"db__port": "5432",
	}
	opts := DefaultOptions()
	opts.Uppercase = true
	res := Flatten(env, opts)

	if _, ok := res.Flattened["DB__HOST"]; !ok {
		t.Errorf("expected DB__HOST in flattened, got %v", res.Flattened)
	}
	if _, ok := res.Flattened["DB__PORT"]; !ok {
		t.Errorf("expected DB__PORT in flattened, got %v", res.Flattened)
	}
	if len(res.Renamed) != 2 {
		t.Errorf("expected 2 renames, got %d", len(res.Renamed))
	}
}

func TestFlatten_PrefixFilter(t *testing.T) {
	env := map[string]string{
		"APP__HOST": "example.com",
		"APP__PORT": "8080",
		"DB__HOST":  "localhost",
	}
	opts := DefaultOptions()
	opts.Prefix = "APP"
	res := Flatten(env, opts)

	if len(res.Flattened) != 2 {
		t.Fatalf("expected 2 keys after prefix filter, got %d: %v", len(res.Flattened), res.Flattened)
	}
	if _, ok := res.Flattened["HOST"]; !ok {
		t.Errorf("expected HOST key after stripping prefix, got %v", res.Flattened)
	}
	if _, ok := res.Flattened["PORT"]; !ok {
		t.Errorf("expected PORT key after stripping prefix, got %v", res.Flattened)
	}
}

func TestFlatten_PrefixFilterUppercase(t *testing.T) {
	env := map[string]string{
		"app__debug": "true",
		"db__name":   "mydb",
	}
	opts := DefaultOptions()
	opts.Prefix = "app"
	opts.Uppercase = true
	res := Flatten(env, opts)

	if len(res.Flattened) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(res.Flattened), res.Flattened)
	}
	if res.Flattened["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %v", res.Flattened)
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	env := map[string]string{
		"APP.HOST": "example.com",
		"APP.PORT": "9000",
	}
	opts := Options{Separator: ".", Prefix: "APP"}
	res := Flatten(env, opts)

	if len(res.Flattened) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(res.Flattened), res.Flattened)
	}
	if res.Flattened["HOST"] != "example.com" {
		t.Errorf("expected HOST=example.com, got %v", res.Flattened)
	}
}

func TestFlatten_EmptyEnv(t *testing.T) {
	res := Flatten(map[string]string{}, DefaultOptions())
	if len(res.Flattened) != 0 {
		t.Errorf("expected empty result, got %v", res.Flattened)
	}
}
