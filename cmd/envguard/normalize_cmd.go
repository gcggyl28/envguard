package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envguard/internal/envparser"
	"github.com/user/envguard/internal/normalizer"
	"github.com/user/envguard/internal/reporter"
)

func runNormalize(args []string) {
	fs := flag.NewFlagSet("normalize", flag.ExitOnError)
	upperKeys := fs.Bool("uppercase-keys", true, "Convert keys to uppercase")
	trimSpace := fs.Bool("trim", true, "Trim whitespace from keys and values")
	replaceHyphens := fs.Bool("replace-hyphens", true, "Replace hyphens in keys with underscores")
	lowerValues := fs.Bool("lowercase-values", false, "Convert values to lowercase")
	output := fs.String("output", "", "Write normalized env to file (default: stdout)")
	_ = fs.Parse(args)

	envFile := fs.Arg(0)
	if envFile == "" {
		fmt.Fprintln(os.Stderr, "usage: envguard normalize [options] <envfile>")
		os.Exit(1)
	}

	env, err := envparser.ParseFile(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading env file: %v\n", err)
		os.Exit(1)
	}

	opts := normalizer.Options{
		UppercaseKeys:   *upperKeys,
		TrimSpace:       *trimSpace,
		ReplaceHyphens:  *replaceHyphens,
		LowercaseValues: *lowerValues,
	}

	res := normalizer.Normalize(env, opts)
	reporter.PrintNormalize(res)

	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		for k, v := range res.Normalized {
			fmt.Fprintf(f, "%s=%s\n", k, v)
		}
	}
}
