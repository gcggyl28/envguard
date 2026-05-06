package templater_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envguard/internal/templater"
)

func TestRender_BasicSubstitution(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080", "APP_HOST": "localhost"}
	tmpl := "HOST={{.APP_HOST}}\nPORT={{.APP_PORT}}\n"

	res, err := templater.Render(tmpl, env, templater.RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "HOST=localhost") {
		t.Errorf("expected HOST=localhost in output, got: %s", res.Output)
	}
	if !strings.Contains(res.Output, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got: %s", res.Output)
	}
}

func TestRender_MissingKeyNonStrict(t *testing.T) {
	env := map[string]string{"APP_PORT": "9000"}
	tmpl := "PORT={{.APP_PORT}}\nSECRET={{.APP_SECRET}}\n"

	res, err := templater.Render(tmpl, env, templater.RenderOptions{Strict: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "SECRET=\n") {
		t.Errorf("expected empty SECRET, got: %s", res.Output)
	}
}

func TestRender_MissingKeyStrict(t *testing.T) {
	env := map[string]string{}
	tmpl := "DB={{.DB_URL}}\n"

	_, err := templater.Render(tmpl, env, templater.RenderOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode for missing key")
	}
}

func TestRender_EnvFuncSubstitution(t *testing.T) {
	env := map[string]string{"SERVICE": "api"}
	tmpl := `NAME={{ env "SERVICE" }}`

	res, err := templater.Render(tmpl, env, templater.RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "NAME=api") {
		t.Errorf("expected NAME=api, got: %s", res.Output)
	}
	if res.Rendered != 1 {
		t.Errorf("expected 1 rendered substitution, got %d", res.Rendered)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	env := map[string]string{}
	tmpl := "BROKEN={{.unclosed"

	_, err := templater.Render(tmpl, env, templater.RenderOptions{})
	if err == nil {
		t.Fatal("expected parse error for invalid template")
	}
}

func TestRenderFile_ReadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "template.env.tmpl")
	content := "APP_ENV={{.APP_ENV}}\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	env := map[string]string{"APP_ENV": "production"}
	res, err := templater.RenderFile(path, env, templater.RenderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production, got: %s", res.Output)
	}
}

func TestRenderFile_MissingFile(t *testing.T) {
	_, err := templater.RenderFile("/nonexistent/path.tmpl", map[string]string{}, templater.RenderOptions{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
