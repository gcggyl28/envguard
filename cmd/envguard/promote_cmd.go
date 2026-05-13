package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envguard/internal/envparser"
	"github.com/user/envguard/internal/promoter"
	"github.com/user/envguard/internal/reporter"
)

func runPromote(args []string) {
	fs := flag.NewFlagSet("promote", flag.ExitOnError)
	srcFile := fs.String("src", "", "source .env file (required)")
	dstFile := fs.String("dst", "", "destination .env file (required)")
	allow := fs.String("allow", "", "comma-separated list of keys to promote")
	deny := fs.String("deny", "", "comma-separated list of keys to block")
	overwrite := fs.Bool("overwrite", false, "overwrite existing destination keys")
	fs.Parse(args) //nolint:errcheck

	if *srcFile == "" || *dstFile == "" {
		fmt.Fprintln(os.Stderr, "error: --src and --dst are required")
		os.Exit(1)
	}

	src, err := envparser.ParseFile(*srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading src: %v\n", err)
		os.Exit(1)
	}

	dst, err := envparser.ParseFile(*dstFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error reading dst: %v\n", err)
		os.Exit(1)
	}
	if dst == nil {
		dst = map[string]string{}
	}

	opts := promoter.Options{
		Overwrite:  *overwrite,
		AllowKeys:  splitCSV(*allow),
		DenyKeys:   splitCSV(*deny),
	}

	results, sum := promoter.Promote(src, dst, opts)
	reporter.PrintPromote(results, sum)

	if sum.Promoted > 0 {
		if err := writeEnvFile(*dstFile, dst); err != nil {
			fmt.Fprintf(os.Stderr, "error writing dst: %v\n", err)
			os.Exit(1)
		}
	}
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func writeEnvFile(path string, env map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for k, v := range env {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}
	return nil
}
