package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envguard/internal/encoder"
)

// PrintEncoded writes the encoding result to stdout.
func PrintEncoded(res encoder.Result) {
	FprintEncoded(os.Stdout, res)
}

// FprintEncoded writes a human-readable encoding summary to w.
func FprintEncoded(w io.Writer, res encoder.Result) {
	fmt.Fprintf(w, "Encoded %d key(s) as %s\n", res.KeyCount, res.Format)
	fmt.Fprintf(w, "\n%s\n", res.Encoded)
}
