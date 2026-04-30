package schema

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// VarType represents the expected type of an environment variable.
type VarType string

const (
	TypeString VarType = "string"
	TypeInt    VarType = "int"
	TypeBool   VarType = "bool"
	TypeURL    VarType = "url"
)

// VarSpec defines the schema for a single environment variable.
type VarSpec struct {
	Name     string
	Required bool
	Type     VarType
	Default  string
}

// Schema holds the full set of variable specifications.
type Schema struct {
	Vars []VarSpec
}

// LoadSchema parses a schema definition file.
// Each line format: VAR_NAME [required|optional] [string|int|bool|url] [default=VALUE]
func LoadSchema(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening schema file: %w", err)
	}
	defer f.Close()

	schema := &Schema{}
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: expected at least NAME and required/optional", lineNum)
		}

		spec := VarSpec{
			Name: parts[0],
			Type: TypeString,
		}

		switch strings.ToLower(parts[1]) {
		case "required":
			spec.Required = true
		case "optional":
			spec.Required = false
		default:
			return nil, fmt.Errorf("line %d: expected 'required' or 'optional', got %q", lineNum, parts[1])
		}

		if len(parts) >= 3 {
			spec.Type = VarType(strings.ToLower(parts[2]))
		}

		for _, part := range parts[3:] {
			if strings.HasPrefix(part, "default=") {
				spec.Default = strings.TrimPrefix(part, "default=")
			}
		}

		schema.Vars = append(schema.Vars, spec)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}

	return schema, nil
}
