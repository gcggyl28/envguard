// Package exporter provides functionality to export validated env variables
// to various output formats (e.g., shell export statements, Docker env-file).
package exporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents the output format for exported env variables.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatDotenv Format = "dotenv"
)

// Export writes the env map to stdout in the given format.
func Export(env map[string]string, format Format) error {
	return Fexport(os.Stdout, env, format)
}

// Fexport writes the env map to the given writer in the given format.
func Fexport(w io.Writer, env map[string]string, format Format) error {
	keys := sortedKeys(env)

	switch format {
	case FormatShell:
		for _, k := range keys {
			_, err := fmt.Fprintf(w, "export %s=%q\n", k, env[k])
			if err != nil {
				return err
			}
		}
	case FormatDocker:
		for _, k := range keys {
			_, err := fmt.Fprintf(w, "%s=%s\n", k, env[k])
			if err != nil {
				return err
			}
		}
	case FormatDotenv:
		for _, k := range keys {
			v := env[k]
			if strings.ContainsAny(v, " \t#") {
				v = fmt.Sprintf("%q", v)
			}
			_, err := fmt.Fprintf(w, "%s=%s\n", k, v)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}

	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
