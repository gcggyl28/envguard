// Package encoder converts an env map into various encoded string formats
// such as base64, JSON key=value pairs, and URL query strings.
package encoder

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Format represents a supported encoding format.
type Format string

const (
	FormatBase64 Format = "base64"
	FormatJSON   Format = "json"
	FormatQuery  Format = "query"
)

// Result holds the encoded output and the format used.
type Result struct {
	Format  Format
	Encoded string
	KeyCount int
}

// Encode encodes the provided env map into the requested format.
func Encode(env map[string]string, format Format) (Result, error) {
	keys := sortedKeys(env)

	switch format {
	case FormatBase64:
		var sb strings.Builder
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(sb.String()))
		return Result{Format: format, Encoded: encoded, KeyCount: len(keys)}, nil

	case FormatJSON:
		b, err := json.Marshal(env)
		if err != nil {
			return Result{}, fmt.Errorf("encoder: json marshal failed: %w", err)
		}
		return Result{Format: format, Encoded: string(b), KeyCount: len(keys)}, nil

	case FormatQuery:
		vals := url.Values{}
		for _, k := range keys {
			vals.Set(k, env[k])
		}
		return Result{Format: format, Encoded: vals.Encode(), KeyCount: len(keys)}, nil

	default:
		return Result{}, fmt.Errorf("encoder: unsupported format %q", format)
	}
}

// FormatFromString parses a format string, returning an error if unrecognised.
func FormatFromString(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatBase64, FormatJSON, FormatQuery:
		return Format(strings.ToLower(s)), nil
	default:
		return "", fmt.Errorf("encoder: unknown format %q (choose base64, json, query)", s)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
