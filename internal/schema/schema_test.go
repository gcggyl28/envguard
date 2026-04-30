package schema

import (
	"os"
	"testing"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "schema-*.txt")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadSchema_Valid(t *testing.T) {
	content := `
# This is a comment
DATABASE_URL required url
APP_PORT      optional int    default=8080
DEBUG         optional bool   default=false
APP_NAME      required string
`
	path := writeTempSchema(t, content)

	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.Vars) != 4 {
		t.Fatalf("expected 4 vars, got %d", len(s.Vars))
	}

	cases := []struct {
		name     string
		required bool
		varType  VarType
		defVal   string
	}{
		{"DATABASE_URL", true, TypeURL, ""},
		{"APP_PORT", false, TypeInt, "8080"},
		{"DEBUG", false, TypeBool, "false"},
		{"APP_NAME", true, TypeString, ""},
	}

	for i, c := range cases {
		v := s.Vars[i]
		if v.Name != c.name {
			t.Errorf("var %d: expected name %q, got %q", i, c.name, v.Name)
		}
		if v.Required != c.required {
			t.Errorf("var %d: expected required=%v, got %v", i, c.required, v.Required)
		}
		if v.Type != c.varType {
			t.Errorf("var %d: expected type %q, got %q", i, c.varType, v.Type)
		}
		if v.Default != c.defVal {
			t.Errorf("var %d: expected default %q, got %q", i, c.defVal, v.Default)
		}
	}
}

func TestLoadSchema_InvalidRequiredField(t *testing.T) {
	content := "MY_VAR unknown string\n"
	path := writeTempSchema(t, content)

	_, err := LoadSchema(path)
	if err == nil {
		t.Fatal("expected error for invalid required field, got nil")
	}
}

func TestLoadSchema_MissingFile(t *testing.T) {
	_, err := LoadSchema("/nonexistent/path/schema.txt")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadSchema_EmptyFile(t *testing.T) {
	path := writeTempSchema(t, "")
	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Vars) != 0 {
		t.Errorf("expected 0 vars, got %d", len(s.Vars))
	}
}
