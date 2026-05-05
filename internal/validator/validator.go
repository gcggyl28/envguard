package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envguard/internal/schema"
)

// Result holds the outcome of a single variable validation.
type Result struct {
	Key     string
	Passed  bool
	Message string
}

// Report aggregates all validation results.
type Report struct {
	Results []Result
	Valid   bool
}

// Validate checks the provided env map against the given schema.
func Validate(env map[string]string, s *schema.Schema) *Report {
	report := &Report{Valid: true}

	for _, field := range s.Fields {
		val, exists := env[field.Name]

		if field.Required && !exists {
			report.add(field.Name, false, "required variable is missing")
			continue
		}

		if !exists {
			continue
		}

		if field.Pattern != "" {
			matched, err := regexp.MatchString(field.Pattern, val)
			if err != nil {
				report.add(field.Name, false, fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err))
				continue
			}
			if !matched {
				report.add(field.Name, false, fmt.Sprintf("value does not match pattern %q", field.Pattern))
				continue
			}
		}

		if len(field.AllowedValues) > 0 {
			if !contains(field.AllowedValues, val) {
				report.add(field.Name, false, fmt.Sprintf("value %q is not in allowed values [%s]", val, strings.Join(field.AllowedValues, ", ")))
				continue
			}
		}

		report.add(field.Name, true, "ok")
	}

	return report
}

func (r *Report) add(key string, passed bool, message string) {
	if !passed {
		r.Valid = false
	}
	r.Results = append(r.Results, Result{Key: key, Passed: passed, Message: message})
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
