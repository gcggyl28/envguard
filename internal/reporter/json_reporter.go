package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envguard/internal/auditor"
	"github.com/yourorg/envguard/internal/validator"
)

// JSONReport is the structured output format for JSON reporting.
type JSONReport struct {
	Valid    bool              `json:"valid"`
	Findings []JSONFinding     `json:"findings"`
	Audit    []JSONAuditResult `json:"audit"`
	Summary  JSONSummary       `json:"summary"`
}

type JSONFinding struct {
	Key     string `json:"key"`
	Message string `json:"message"`
	Severity string `json:"severity"`
}

type JSONAuditResult struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type JSONSummary struct {
	TotalFindings int `json:"total_findings"`
	TotalAudit    int `json:"total_audit"`
}

// PrintJSON writes a JSON-formatted report to the given writer.
func PrintJSON(w io.Writer, findings []validator.Finding, auditResults []auditor.Result) {
	report := buildJSONReport(findings, auditResults)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON report: %v\n", err)
	}
}

func buildJSONReport(findings []validator.Finding, auditResults []auditor.Result) JSONReport {
	jsonFindings := make([]JSONFinding, 0, len(findings))
	for _, f := range findings {
		jsonFindings = append(jsonFindings, JSONFinding{
			Key:      f.Key,
			Message:  f.Message,
			Severity: "error",
		})
	}

	jsonAudit := make([]JSONAuditResult, 0, len(auditResults))
	for _, a := range auditResults {
		jsonAudit = append(jsonAudit, JSONAuditResult{
			Key:     a.Key,
			Message: a.Message,
		})
	}

	return JSONReport{
		Valid:    len(findings) == 0,
		Findings: jsonFindings,
		Audit:    jsonAudit,
		Summary: JSONSummary{
			TotalFindings: len(findings),
			TotalAudit:    len(auditResults),
		},
	}
}
