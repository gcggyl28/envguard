package auditor

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envguard/internal/schema"
)

// Finding represents a single audit observation about an env variable.
type Finding struct {
	Key      string
	Severity string // "warn" or "info"
	Message  string
}

// AuditResult holds all findings from an audit run.
type AuditResult struct {
	Findings []Finding
}

// HasWarnings returns true if any finding has severity "warn".
func (r *AuditResult) HasWarnings() bool {
	for _, f := range r.Findings {
		if f.Severity == "warn" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r *AuditResult) Summary() string {
	if len(r.Findings) == 0 {
		return "audit passed: no issues found"
	}
	var sb strings.Builder
	for _, f := range r.Findings {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", strings.ToUpper(f.Severity), f.Key, f.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Audit inspects the parsed env map against the schema and returns an AuditResult.
// It checks for:
//   - keys present in env but not declared in the schema (undeclared keys)
//   - keys with empty values that have a default defined in the schema
func Audit(env map[string]string, s *schema.Schema) *AuditResult {
	result := &AuditResult{}

	// Build a set of declared keys for quick lookup.
	declared := make(map[string]bool, len(s.Vars))
	defaults := make(map[string]string, len(s.Vars))
	for _, v := range s.Vars {
		declared[v.Name] = true
		if v.Default != "" {
			defaults[v.Name] = v.Default
		}
	}

	// Collect keys in deterministic order.
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if !declared[k] {
			result.Findings = append(result.Findings, Finding{
				Key:      k,
				Severity: "warn",
				Message:  "key is not declared in schema",
			})
			continue
		}
		if env[k] == "" {
			if def, ok := defaults[k]; ok {
				result.Findings = append(result.Findings, Finding{
					Key:      k,
					Severity: "info",
					Message:  fmt.Sprintf("value is empty; schema default is %q", def),
				})
			}
		}
	}

	return result
}
