package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/envguard/internal/masker"
)

// PrintMasked writes the masked env map to stdout.
func PrintMasked(env map[string]string, opts masker.Options) {
	FprintMasked(os.Stdout, env, opts)
}

// FprintMasked writes the masked env map to w in a human-readable table.
func FprintMasked(w io.Writer, env map[string]string, opts masker.Options) {
	masked := masker.MaskMap(env, opts)

	keys := make([]string, 0, len(masked))
	for k := range masked {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "=== Masked Environment Variables ===")
	if len(keys) == 0 {
		fmt.Fprintln(w, "  (no variables)")
		return
	}

	// Compute column width for alignment.
	maxLen := 0
	for _, k := range keys {
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}

	for _, k := range keys {
		padding := maxLen - len(k)
		fmt.Fprintf(w, "  %s%s  =  %s\n", k, spaces(padding), masked[k])
	}
	fmt.Fprintf(w, "\nTotal: %d variable(s) | style: %s\n", len(keys), opts.Style)
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	s := make([]byte, n)
	for i := range s {
		s[i] = ' '
	}
	return string(s)
}
