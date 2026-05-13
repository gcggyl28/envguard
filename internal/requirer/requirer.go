// Package requirer checks that all required keys defined in a schema
// are present and non-empty in a given env map.
package requirer

import "github.com/user/envguard/internal/schema"

// Finding represents a missing or empty required key.
type Finding struct {
	Key     string
	Reason  string
}

// Result holds the outcome of a requirements check.
type Result struct {
	Findings []Finding
	Checked  int
}

// Check validates that every required key in the schema exists and is
// non-empty in env. Optional keys with no value are silently ignored.
func Check(env map[string]string, s *schema.Schema) Result {
	var findings []Finding
	checked := 0

	for _, field := range s.Fields {
		if !field.Required {
			continue
		}
		checked++
		val, exists := env[field.Key]
		if !exists {
			findings = append(findings, Finding{
				Key:    field.Key,
				Reason: "missing",
			})
			continue
		}
		if val == "" {
			findings = append(findings, Finding{
				Key:    field.Key,
				Reason: "empty value",
			})
		}
	}

	return Result{
		Findings: findings,
		Checked:  checked,
	}
}

// Passed returns true when no findings were recorded.
func (r Result) Passed() bool {
	return len(r.Findings) == 0
}
