package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourorg/envguard/internal/redactor"
)

// PrintRedacted writes the redacted env map to stdout.
func PrintRedacted(env map[string]string, useHeuristic bool, sensitiveKeys []string) {
	FprintRedacted(os.Stdout, env, useHeuristic, sensitiveKeys)
}

// FprintRedacted writes the redacted env map to w.
// When useHeuristic is true the redactor's pattern-based detection is applied;
// otherwise only keys in sensitiveKeys are masked.
func FprintRedacted(w io.Writer, env map[string]string, useHeuristic bool, sensitiveKeys []string) {
	var redacted map[string]string
	if useHeuristic {
		redacted = redactor.Redact(env)
	} else {
		redacted = redactor.RedactList(env, sensitiveKeys)
	}

	keys := make([]string, 0, len(redacted))
	for k := range redacted {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "=== Redacted Environment ===")
	for _, k := range keys {
		fmt.Fprintf(w, "  %-30s = %s\n", k, redacted[k])
	}
	fmt.Fprintf(w, "\n  Total keys: %d\n", len(keys))
}
