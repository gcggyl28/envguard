package coercer

import (
	"fmt"
	"strconv"
	"strings"
)

// TargetType represents the type to coerce a value into.
type TargetType string

const (
	TypeBool   TargetType = "bool"
	TypeInt    TargetType = "int"
	TypeFloat  TargetType = "float"
	TypeString TargetType = "string"
)

// Rule defines a coercion rule for a specific key.
type Rule struct {
	Key    string
	Target TargetType
}

// Result holds the outcome of a single coercion attempt.
type Result struct {
	Key      string
	Original string
	Coerced  string
	Type     TargetType
	Changed  bool
	Error    string
}

// Coerce applies coercion rules to the provided env map, returning a new map
// with normalised values and a slice of Results describing each operation.
func Coerce(env map[string]string, rules []Rule) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	results := make([]Result, 0, len(rules))

	for _, rule := range rules {
		original, exists := env[rule.Key]
		if !exists {
			continue
		}

		coerced, err := coerceValue(original, rule.Target)
		r := Result{
			Key:      rule.Key,
			Original: original,
			Type:     rule.Target,
		}
		if err != nil {
			r.Error = err.Error()
			r.Coerced = original
		} else {
			r.Coerced = coerced
			r.Changed = coerced != original
			out[rule.Key] = coerced
		}
		results = append(results, r)
	}

	return out, results
}

func coerceValue(value string, target TargetType) (string, error) {
	switch target {
	case TypeBool:
		norm := strings.ToLower(strings.TrimSpace(value))
		switch norm {
		case "1", "yes", "on", "true":
			return "true", nil
		case "0", "no", "off", "false":
			return "false", nil
		default:
			return "", fmt.Errorf("cannot coerce %q to bool", value)
		}
	case TypeInt:
		trimmed := strings.TrimSpace(value)
		if _, err := strconv.ParseInt(trimmed, 10, 64); err != nil {
			return "", fmt.Errorf("cannot coerce %q to int", value)
		}
		return trimmed, nil
	case TypeFloat:
		trimmed := strings.TrimSpace(value)
		f, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return "", fmt.Errorf("cannot coerce %q to float", value)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case TypeString:
		return strings.TrimSpace(value), nil
	default:
		return "", fmt.Errorf("unknown target type %q", target)
	}
}
