package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envguard/internal/envparser"
	"github.com/yourusername/envguard/internal/reporter"
	"github.com/yourusername/envguard/internal/scoper"
)

func runScope(args []string) {
	fs := flag.NewFlagSet("scope", flag.ExitOnError)
	envFile := fs.String("env", ".env", "Path to the .env file")
	scopeName := fs.String("scope", "", "Scope name (e.g. production, staging)")
	prefixes := fs.String("prefixes", "", "Comma-separated list of key prefixes to include")
	strip := fs.Bool("strip", false, "Strip matched prefixes from output keys")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "error parsing flags:", err)
		os.Exit(1)
	}

	if *scopeName == "" {
		fmt.Fprintln(os.Stderr, "error: --scope is required")
		fs.Usage()
		os.Exit(1)
	}

	env, err := envparser.ParseFile(*envFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading env file:", err)
		os.Exit(1)
	}

	var pfxList []string
	if *prefixes != "" {
		for _, p := range strings.Split(*prefixes, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				pfxList = append(pfxList, p)
			}
		}
	}

	scope := scoper.Scope{Name: *scopeName, Prefixes: pfxList}
	result := scoper.Apply(env, scope)

	if *strip {
		stripped := scoper.Strip(result.Included, scope)
		fmt.Printf("# Scoped env: %s (prefixes stripped)\n", *scopeName)
		for k, v := range stripped {
			fmt.Printf("%s=%s\n", k, v)
		}
		return
	}

	reporter.PrintScope(result)
}
