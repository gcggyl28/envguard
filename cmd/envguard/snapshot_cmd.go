package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envguard/internal/envparser"
	"github.com/user/envguard/internal/reporter"
	"github.com/user/envguard/internal/snapshotter"
)

func runSnapshot(args []string) {
	fs := flag.NewFlagSet("snapshot", flag.ExitOnError)
	save := fs.String("save", "", "Save a new snapshot to this file")
	compare := fs.String("compare", "", "Compare current env file against this snapshot")
	envFile := fs.String("env", ".env", "Path to the .env file")
	fs.Parse(args)

	if *save == "" && *compare == "" {
		fmt.Fprintln(os.Stderr, "error: provide --save or --compare")
		fs.Usage()
		os.Exit(1)
	}

	env, err := envparser.ParseFile(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not parse env file: %v\n", err)
		os.Exit(1)
	}

	if *save != "" {
		if err := snapshotter.Save(env, *envFile, *save); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Snapshot saved to %s\n", *save)
		return
	}

	snap, err := snapshotter.Load(*compare)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not load snapshot: %v\n", err)
		os.Exit(1)
	}

	current := &snapshotter.Snapshot{Source: *envFile, Env: env}
	added, removed, changed := snapshotter.Compare(snap, current)
	reporter.PrintSnapshot(*envFile, added, removed, changed)

	if len(added)+len(removed)+len(changed) > 0 {
		os.Exit(1)
	}
}
