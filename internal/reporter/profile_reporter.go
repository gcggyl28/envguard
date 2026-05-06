package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envguard/internal/profiler"
)

// PrintProfile writes a human-readable profile summary to stdout.
func PrintProfile(p profiler.Profile) {
	FprintProfile(os.Stdout, p)
}

// FprintProfile writes a human-readable profile summary to the given writer.
func FprintProfile(w io.Writer, p profiler.Profile) {
	fmt.Fprintln(w, "=== .env Profile ===")
	fmt.Fprintf(w, "Total declared keys : %d\n", p.TotalKeys)
	fmt.Fprintf(w, "Required            : %d\n", len(p.RequiredKeys))
	fmt.Fprintf(w, "Optional            : %d\n", len(p.OptionalKeys))
	fmt.Fprintf(w, "With defaults       : %d\n", len(p.KeysWithDefault))
	fmt.Fprintf(w, "Sensitive           : %d\n", len(p.SensitiveKeys))
	fmt.Fprintf(w, "Undeclared          : %d\n", len(p.UndeclaredKeys))

	if len(p.SensitiveKeys) > 0 {
		sorted := sorted(p.SensitiveKeys)
		fmt.Fprintf(w, "\nSensitive keys      : %s\n", strings.Join(sorted, ", "))
	}

	if len(p.KeysWithDefault) > 0 {
		sorted := sorted(p.KeysWithDefault)
		fmt.Fprintf(w, "Keys with defaults  : %s\n", strings.Join(sorted, ", "))
	}

	if len(p.UndeclaredKeys) > 0 {
		sorted := sorted(p.UndeclaredKeys)
		fmt.Fprintf(w, "\n⚠  Undeclared keys  : %s\n", strings.Join(sorted, ", "))
	}
}

func sorted(keys []string) []string {
	copy_ := make([]string, len(keys))
	copy(copy_, keys)
	sort.Strings(copy_)
	return copy_
}
