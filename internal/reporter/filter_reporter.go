package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourorg/envguard/internal/filter"
)

// PrintFilter writes the filter results to stdout.
func PrintFilter(env map[string]string, s filter.Summary) {
	FprintFilter(os.Stdout, env, s)
}

// FprintFilter writes the filter results to w.
func FprintFilter(w io.Writer, env map[string]string, s filter.Summary) {
	fmt.Fprintf(w, "Filter Results: %d included / %d excluded (total: %d)\n",
		s.Included, s.Excluded, s.Total)

	if len(env) == 0 {
		fmt.Fprintln(w, "  (no keys matched)")
		return
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "  ✔ %s=%s\n", k, env[k])
	}
}
