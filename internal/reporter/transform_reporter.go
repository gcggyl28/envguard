package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/envguard/internal/transformer"
)

// PrintTransform writes a human-readable transformation summary to stdout.
func PrintTransform(result transformer.Result) {
	FprintTransform(os.Stdout, result)
}

// FprintTransform writes a human-readable transformation summary to w.
// Changes are sorted alphabetically by key (or OldKey for renames) for
// deterministic output. Each line is prefixed with the change type:
//   - [rename+value]: the key was renamed and its value was also changed
//   - [rename]:       only the key name changed
//   - [value]:        only the value changed
func FprintTransform(w io.Writer, result transformer.Result) {
	changes := result.Changes

	if len(changes) == 0 {
		fmt.Fprintln(w, "✔ No transformations applied.")
		return
	}

	// Sort for deterministic output.
	sort.Slice(changes, func(i, j int) bool {
		ki := changes[i].Key
		if changes[i].OldKey != "" {
			ki = changes[i].OldKey
		}
		kj := changes[j].Key
		if changes[j].OldKey != "" {
			kj = changes[j].OldKey
		}
		return ki < kj
	})

	fmt.Fprintf(w, "Transformations applied (%d):\n", len(changes))
	for _, c := range changes {
		switch {
		case c.OldKey != "" && c.OldValue != c.NewValue:
			fmt.Fprintf(w, "  [rename+value] %s -> %s  |  %q -> %q\n",
				c.OldKey, c.Key, c.OldValue, c.NewValue)
		case c.OldKey != "":
			fmt.Fprintf(w, "  [rename]       %s -> %s\n", c.OldKey, c.Key)
		default:
			fmt.Fprintf(w, "  [value]        %s: %q -> %q\n", c.Key, c.OldValue, c.NewValue)
		}
	}
}
