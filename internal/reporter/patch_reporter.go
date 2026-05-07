package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envguard/internal/patcher"
)

// PrintPatch writes a patch result summary to stdout.
func PrintPatch(result patcher.Result) {
	FprintPatch(os.Stdout, result)
}

// FprintPatch writes a patch result summary to the given writer.
func FprintPatch(w io.Writer, result patcher.Result) {
	applied := sorted(result.Applied)
	deleted := sorted(result.Deleted)
	skipped := sorted(result.Skipped)

	if len(applied) == 0 && len(deleted) == 0 && len(skipped) == 0 {
		fmt.Fprintln(w, "patch: no operations performed")
		return
	}

	if len(applied) > 0 {
		fmt.Fprintf(w, "patch: applied (%d)\n", len(applied))
		for _, k := range applied {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	if len(deleted) > 0 {
		fmt.Fprintf(w, "patch: deleted (%d)\n", len(deleted))
		for _, k := range deleted {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(skipped) > 0 {
		fmt.Fprintf(w, "patch: skipped — key not found (%d)\n", len(skipped))
		for _, k := range skipped {
			fmt.Fprintf(w, "  ? %s\n", k)
		}
	}
}

func sorted(keys []string) []string {
	if len(keys) == 0 {
		return keys
	}
	copy := append([]string(nil), keys...)
	sort.Strings(copy)
	return copy
}
