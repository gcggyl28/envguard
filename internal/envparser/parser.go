package envparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ParseFile reads a .env file and returns a map of key-value pairs.
// Lines starting with '#' and blank lines are ignored.
// Lines must follow KEY=VALUE format.
func ParseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envparser: cannot open file %q: %w", path, err)
	}
	defer f.Close()
	return Parse(f)
}

// Parse reads env entries from an io.Reader.
func Parse(r io.Reader) (map[string]string, error) {
	env := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("envparser: invalid syntax on line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		if key == "" {
			return nil, fmt.Errorf("envparser: empty key on line %d", lineNum)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envparser: read error: %w", err)
	}

	return env, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
