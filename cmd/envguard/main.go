package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envguard/internal/auditor"
	"github.com/yourorg/envguard/internal/envparser"
	"github.com/yourorg/envguard/internal/reporter"
	"github.com/yourorg/envguard/internal/schema"
	"github.com/yourorg/envguard/internal/validator"
)

func main() {
	envFile := flag.String("env", ".env", "Path to the .env file")
	schemaFile := flag.String("schema", ".env.schema.json", "Path to the schema file")
	jsonOutput := flag.Bool("json", false, "Output results as JSON")
	quiet := flag.Bool("quiet", false, "Suppress output, only use exit code")
	flag.Parse()

	sc, err := schema.LoadSchema(*schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading schema: %v\n", err)
		os.Exit(2)
	}

	env, err := envparser.ParseFile(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env file: %v\n", err)
		os.Exit(2)
	}

	findings := validator.Validate(env, sc)
	auditResults := auditor.Audit(env, sc)

	if !*quiet {
		if *jsonOutput {
			reporter.PrintJSON(os.Stdout, findings, auditResults)
		} else {
			reporter.Print(findings, auditResults)
		}
	}

	os.Exit(reporter.ExitCode(findings, auditResults))
}
