// Package templater renders .env files from a template with variable substitution.
package templater

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// RenderOptions controls template rendering behavior.
type RenderOptions struct {
	// Strict causes rendering to fail if any template variable is missing.
	Strict bool
}

// Result holds the output of a template render operation.
type Result struct {
	Output   string
	Missing  []string
	Rendered int
}

// Render processes a template string using the provided env map.
// Variables are referenced as {{.KEY}} in the template.
func Render(tmplContent string, env map[string]string, opts RenderOptions) (*Result, error) {
	result := &Result{}
	missing := map[string]struct{}{}

	funcMap := template.FuncMap{
		"env": func(key string) string {
			if val, ok := env[key]; ok {
				result.Rendered++
				return val
			}
			if val, ok := os.LookupEnv(key); ok {
				result.Rendered++
				return val
			}
			missing[key] = struct{}{}
			return ""
		},
	}

	tmpl, err := template.New("envfile").Option("missingkey=zero").Funcs(funcMap).Parse(tmplContent)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, env); err != nil {
		return nil, fmt.Errorf("template execute error: %w", err)
	}

	for k := range missing {
		result.Missing = append(result.Missing, k)
	}

	if opts.Strict && len(result.Missing) > 0 {
		return result, fmt.Errorf("strict mode: missing variables: %s", strings.Join(result.Missing, ", "))
	}

	result.Output = sb.String()
	return result, nil
}

// RenderFile reads a template file and renders it with the provided env map.
func RenderFile(path string, env map[string]string, opts RenderOptions) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template file: %w", err)
	}
	return Render(string(data), env, opts)
}
