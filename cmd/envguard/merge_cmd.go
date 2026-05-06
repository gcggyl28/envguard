package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/envguard/internal/envparser"
	"github.com/yourusername/envguard/internal/exporter"
	"github.com/yourusername/envguard/internal/merger"
	"github.com/yourusername/envguard/internal/reporter"
)

// runMerge handles the `merge` subcommand.
// Usage: envguard merge --base .env --override .env.local [--strategy override] [--out merged.env]
func runMerge(args []string) {
	fs := flag.NewFlagSet("merge", flag.ExitOnError)
	baseFile := fs.String("base", ".env", "base .env file")
	overrideFile := fs.String("override", "", "override .env file (required)")
	strategyStr := fs.String("strategy", "base", "conflict resolution strategy: base|override")
	outFile := fs.String("out", "", "write merged env to file (default: stdout summary only)")
	_ = fs.Parse(args)

	if *overrideFile == "" {
		fmt.Fprintln(os.Stderr, "error: --override is required")
		fs.Usage()
		os.Exit(1)
	}

	strategy, err := merger.StrategyFromString(*strategyStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	baseEnv, err := envparser.ParseFile(*baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file: %v\n", err)
		os.Exit(1)
	}

	overrideEnv, err := envparser.ParseFile(*overrideFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading override file: %v\n", err)
		os.Exit(1)
	}

	result := merger.Merge(baseEnv, overrideEnv, strategy)
	reporter.PrintMerge(result, *strategyStr)

	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := exporter.Fexport(f, result.Merged, "dotenv"); err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Merged env written to %s\n", *outFile)
	}
}
