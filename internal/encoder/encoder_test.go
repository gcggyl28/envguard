package encoder

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

var baseEnv = map[string]string{
	"APP_ENV": "production",
	"DB_HOST": "localhost",
	"SECRET":  "s3cr3t",
}

func TestEncode_Base64(t *testing.T) {
	res, err := Encode(baseEnv, FormatBase64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeyCount != 3 {
		t.Errorf("expected 3 keys, got %d", res.KeyCount)
	}
	decoded, err := base64.StdEncoding.DecodeString(res.Encoded)
	if err != nil {
		t.Fatalf("invalid base64: %v", err)
	}
	if !strings.Contains(string(decoded), "APP_ENV=production") {
		t.Errorf("decoded output missing expected key=value pair")
	}
}

func TestEncode_JSON(t *testing.T) {
	res, err := Encode(baseEnv, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(res.Encoded), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestEncode_Query(t *testing.T) {
	res, err := Encode(baseEnv, FormatQuery)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Encoded, "APP_ENV=production") {
		t.Errorf("query string missing APP_ENV=production, got: %s", res.Encoded)
	}
}

func TestEncode_UnsupportedFormat(t *testing.T) {
	_, err := Encode(baseEnv, Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestEncode_EmptyEnv(t *testing.T) {
	res, err := Encode(map[string]string{}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeyCount != 0 {
		t.Errorf("expected 0 keys, got %d", res.KeyCount)
	}
}

func TestFormatFromString_Valid(t *testing.T) {
	for _, s := range []string{"base64", "json", "query", "BASE64", "JSON"} {
		if _, err := FormatFromString(s); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", s, err)
		}
	}
}

func TestFormatFromString_Invalid(t *testing.T) {
	if _, err := FormatFromString("toml"); err == nil {
		t.Error("expected error for unknown format")
	}
}
