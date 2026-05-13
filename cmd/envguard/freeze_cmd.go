package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envguard/internal/envparser"
	"github.com/yourusername/envguard/internal/freezer"
	"github.com/yourusername/envguard/internal/reporter"
)

func runFreeze(args []string) {
	fs := flag.NewFlagSet("freeze", flag.ExitOnError)
	output := fs.String("out", "frozen.json", "output file for frozen env")
	keys := fs.String("keys", "", "comma-separated list of keys to freeze (empty = all)")
	thaw := fs.Bool("thaw", false, "load and display a frozen env file instead of freezing")
	fs.Parse(args)

	if *thaw {
		runThaw(*output)
		return
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: envguard freeze [--out frozen.json] [--keys A,B] <file.env>")
		os.Exit(1)
	}

	envFile := fs.Arg(0)
	env, err := envparser.ParseFile(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "freeze: parse error: %v\n", err)
		os.Exit(1)
	}

	var allowKeys []string
	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				allowKeys = append(allowKeys, k)
			}
		}
	}

	res, err := freezer.Freeze(env, envFile, *output, allowKeys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "freeze: %v\n", err)
		os.Exit(1)
	}

	reporter.PrintFreeze(res)
}

func runThaw(path string) {
	fe, err := freezer.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "thaw: %v\n", err)
		os.Exit(1)
	}
	reporter.PrintThaw(fe)
}
