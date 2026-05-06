package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// InterpolationChange records a key whose value changed after interpolation.
type InterpolationChange struct {
	Key      string
	Original string
	Resolved string
}

// PrintInterpolation writes interpolation changes to stdout.
func PrintInterpolation(changes []InterpolationChange) {
	FprintInterpolation(os.Stdout, changes)
}

// FprintInterpolation writes interpolation changes to the provided writer.
func FprintInterpolation(w io.Writer, changes []InterpolationChange) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No interpolations applied.")
		return
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	fmt.Fprintf(w, "Interpolation results (%d change(s)):\n", len(changes))
	fmt.Fprintln(w, "--------------------------------------")
	for _, c := range changes {
		fmt.Fprintf(w, "  %-24s %s -> %s\n", c.Key, c.Original, c.Resolved)
	}
}

// BuildInterpolationChanges compares original and resolved env maps,
// returning a slice of keys whose values differ.
func BuildInterpolationChanges(original, resolved map[string]string) []InterpolationChange {
	var changes []InterpolationChange
	for k, orig := range original {
		if res, ok := resolved[k]; ok && res != orig {
			changes = append(changes, InterpolationChange{
				Key:      k,
				Original: orig,
				Resolved: res,
			})
		}
	}
	return changes
}
