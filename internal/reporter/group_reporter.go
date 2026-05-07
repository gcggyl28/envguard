package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envguard/internal/grouper"
)

// PrintGroup writes a grouped env report to stdout.
func PrintGroup(result grouper.GroupResult, env map[string]string) {
	FprintGroup(os.Stdout, result, env)
}

// FprintGroup writes a grouped env report to w.
func FprintGroup(w io.Writer, result grouper.GroupResult, env map[string]string) {
	total := 0
	for _, g := range result.Groups {
		total += len(g.Keys)
	}
	total += len(result.Ungrouped)

	fmt.Fprintf(w, "Grouped Environment Keys (%d total)\n", total)
	fmt.Fprintln(w, strings.Repeat("─", 40))

	for _, g := range result.Groups {
		if len(g.Keys) == 0 {
			continue
		}
		fmt.Fprintf(w, "\n[%s] (%d keys)\n", g.Name, len(g.Keys))
		for _, k := range g.Keys {
			fmt.Fprintf(w, "  %-30s = %s\n", k, env[k])
		}
	}

	if len(result.Ungrouped) > 0 {
		fmt.Fprintf(w, "\n[ungrouped] (%d keys)\n", len(result.Ungrouped))
		for _, k := range result.Ungrouped {
			fmt.Fprintf(w, "  %-30s = %s\n", k, env[k])
		}
	}

	fmt.Fprintln(w, strings.Repeat("─", 40))
	fmt.Fprintf(w, "Groups: %d  |  Ungrouped: %d\n", len(result.Groups), len(result.Ungrouped))
}
