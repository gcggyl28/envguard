// Package deprecator identifies deprecated or sunset keys in an env map
// based on a deprecation manifest and produces structured findings.
package deprecator

import "strings"

// Finding represents a single deprecated key found in the env.
type Finding struct {
	Key        string
	Replacement string // empty if no replacement suggested
	Reason     string
}

// Rule defines a deprecation rule for a single key.
type Rule struct {
	Key         string `json:"key"`
	Replacement string `json:"replacement,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// Deprecate checks the provided env map against a list of deprecation rules
// and returns findings for any deprecated keys that are present.
func Deprecate(env map[string]string, rules []Rule) []Finding {
	if len(rules) == 0 || len(env) == 0 {
		return nil
	}

	ruleIndex := make(map[string]Rule, len(rules))
	for _, r := range rules {
		ruleIndex[strings.ToUpper(r.Key)] = r
	}

	var findings []Finding
	for k := range env {
		if rule, ok := ruleIndex[strings.ToUpper(k)]; ok {
			findings = append(findings, Finding{
				Key:         k,
				Replacement: rule.Replacement,
				Reason:      rule.Reason,
			})
		}
	}

	// stable sort by key
	for i := 1; i < len(findings); i++ {
		for j := i; j > 0 && findings[j].Key < findings[j-1].Key; j-- {
			findings[j], findings[j-1] = findings[j-1], findings[j]
		}
	}
	return findings
}

// Summary returns a brief human-readable summary string.
func Summary(findings []Finding) string {
	if len(findings) == 0 {
		return "no deprecated keys found"
	}
	if len(findings) == 1 {
		return "1 deprecated key found"
	}
	return strings.Join([]string{
		strconv(len(findings)), " deprecated keys found",
	}, "")
}

func strconv(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = digits[n%10]
		n /= 10
	}
	return string(buf[pos:])
}
