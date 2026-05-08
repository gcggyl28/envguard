package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envguard/internal/tagger"
)

// PrintTag writes tagging results to stdout.
func PrintTag(tags []tagger.Tag) {
	FprintTag(os.Stdout, tags)
}

// FprintTag writes tagging results to w.
func FprintTag(w io.Writer, tags []tagger.Tag) {
	tagged := 0
	for _, t := range tags {
		if len(t.Tags) > 0 {
			tagged++
		}
	}

	fmt.Fprintf(w, "Tags: %d keys scanned, %d tagged\n\n", len(tags), tagged)

	if len(tags) == 0 {
		fmt.Fprintln(w, "  (no keys)")
		return
	}

	for _, t := range tags {
		if len(t.Tags) == 0 {
			fmt.Fprintf(w, "  %-30s  (untagged)\n", t.Key)
		} else {
			fmt.Fprintf(w, "  %-30s  [%s]\n", t.Key, strings.Join(t.Tags, ", "))
		}
	}
}
