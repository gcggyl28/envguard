package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envguard/internal/exporter"
)

var sampleEnv = map[string]string{
	"APP_ENV":    "production",
	"DB_URL":     "postgres://localhost:5432/mydb",
	"SECRET_KEY": "s3cr3t value with spaces",
}

func TestFexport_Shell(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Fexport(&buf, sampleEnv, exporter.FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=\"") {
		t.Errorf("expected shell export syntax, got:\n%s", out)
	}
	if !strings.Contains(out, "export SECRET_KEY=\"") {
		t.Errorf("expected SECRET_KEY in output, got:\n%s", out)
	}
}

func TestFexport_Docker(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Fexport(&buf, sampleEnv, exporter.FormatDocker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected plain key=value, got:\n%s", out)
	}
}

func TestFexport_Dotenv_QuotesSpaces(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Fexport(&buf, map[string]string{"KEY": "value with spaces"}, exporter.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `KEY="value with spaces"`) {
		t.Errorf("expected quoted value for spaces, got: %s", out)
	}
}

func TestFexport_Dotenv_NoQuotesPlain(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Fexport(&buf, map[string]string{"APP_ENV": "production"}, exporter.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "APP_ENV=production" {
		t.Errorf("expected unquoted plain value, got: %s", out)
	}
}

func TestFexport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Fexport(&buf, sampleEnv, exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestFexport_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	var buf bytes.Buffer
	_ = exporter.Fexport(&buf, env, exporter.FormatDocker)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
