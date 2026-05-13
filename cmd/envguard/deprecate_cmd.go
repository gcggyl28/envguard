package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/user/envguard/internal/deprecator"
	"github.com/user/envguard/internal/envparser"
	"github.com/user/envguard/internal/reporter"
)

func runDeprecate(args []string) {
	fs := flag.NewFlagSet("deprecate", flag.ExitOnError)
	envFile := fs.String("env", ".env", "path to the .env file")
	rulesFile := fs.String("rules", ".deprecations.json", "path to deprecation rules JSON file")
	fs.Parse(args) //nolint:errcheck

	env, err := envparser.ParseFile(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading env file: %v\n", err)
		os.Exit(1)
	}

	rules, err := loadDeprecationRules(*rulesFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading rules file: %v\n", err)
		os.Exit(1)
	}

	findings := deprecator.Deprecate(env, rules)
	reporter.PrintDeprecation(findings)

	if len(findings) > 0 {
		os.Exit(1)
	}
}

func loadDeprecationRules(path string) ([]deprecator.Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var rules []deprecator.Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return rules, nil
}
