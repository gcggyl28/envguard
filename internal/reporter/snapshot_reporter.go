package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// PrintSnapshot prints a snapshot diff summary to stdout.
func PrintSnapshot(source string, added, removed, changed []string) {
	FprintSnapshot(os.Stdout, source, added, removed, changed)
}

// FprintSnapshot writes a snapshot diff summary to the given writer.
func FprintSnapshot(w io.Writer, source string, added, removed, changed []string) {
	fmt.Fprintf(w, "Snapshot diff for: %s\n", source)

	if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
		fmt.Fprintln(w, "  No changes detected.")
		return
	}

	if len(added) > 0 {
		sort.Strings(added)
		fmt.Fprintln(w, "  Added:")
		for _, k := range added {
			fmt.Fprintf(w, "    + %s\n", k)
		}
	}

	if len(removed) > 0 {
		sort.Strings(removed)
		fmt.Fprintln(w, "  Removed:")
		for _, k := range removed {
			fmt.Fprintf(w, "    - %s\n", k)
		}
	}

	if len(changed) > 0 {
		sort.Strings(changed)
		fmt.Fprintln(w, "  Changed:")
		for _, k := range changed {
			fmt.Fprintf(w, "    ~ %s\n", k)
		}
	}
}
