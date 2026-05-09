// Package typecheck validates that .env values conform to declared types.
package typecheck

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Type represents a declared value type for a key.
type Type string

const (
	TypeString  Type = "string"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeBool    Type = "bool"
	TypeURL     Type = "url"
	TypeIP      Type = "ip"
	TypeEmail   Type = "email"
)

// Violation describes a type mismatch for a single key.
type Violation struct {
	Key      string
	Value    string
	Expected Type
	Reason   string
}

// Check validates each key in env against the provided type map.
// Keys not present in types are skipped.
func Check(env map[string]string, types map[string]Type) []Violation {
	var violations []Violation
	for key, expected := range types {
		value, ok := env[key]
		if !ok {
			continue
		}
		if reason, ok := validate(value, expected); !ok {
			violations = append(violations, Violation{
				Key:      key,
				Value:    value,
				Expected: expected,
				Reason:   reason,
			})
		}
	}
	return violations
}

func validate(value string, t Type) (string, bool) {
	switch t {
	case TypeString:
		return "", true
	case TypeInt:
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Sprintf("%q is not a valid integer", value), false
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Sprintf("%q is not a valid float", value), false
		}
	case TypeBool:
		lower := strings.ToLower(value)
		valid := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
		if !valid[lower] {
			return fmt.Sprintf("%q is not a valid boolean", value), false
		}
	case TypeURL:
		u, err := url.ParseRequestURI(value)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Sprintf("%q is not a valid URL", value), false
		}
	case TypeIP:
		if net.ParseIP(value) == nil {
			return fmt.Sprintf("%q is not a valid IP address", value), false
		}
	case TypeEmail:
		if !strings.Contains(value, "@") || strings.Count(value, "@") != 1 {
			return fmt.Sprintf("%q is not a valid email address", value), false
		}
		parts := strings.SplitN(value, "@", 2)
		if parts[0] == "" || !strings.Contains(parts[1], ".") {
			return fmt.Sprintf("%q is not a valid email address", value), false
		}
	default:
		return fmt.Sprintf("unknown type %q", t), false
	}
	return "", true
}
