package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ManiRzb/elixec/internal/image"
	"github.com/ManiRzb/elixec/internal/policy"
	"github.com/ManiRzb/elixec/internal/scoring"
)

// SaveFinalReport generates and saves summarized reports.
func SaveFinalReport(report scoring.Report) {
	saveAsSummaryJSON(report)
	saveAsPlainText(report)
}

// Summarized JSON report for integration with tools
func saveAsSummaryJSON(report scoring.Report) {
	file, err := os.Create("summary_report.json")
	if err != nil {
		log.Fatalf("[ERROR] Failed to create JSON report: %v", err)
	}
	defer file.Close()

	// Summarized structure
	summary := map[string]interface{}{
		"vulnerabilities": map[string]int{
			"critical": len(filterVulnerabilities(report.Vulnerabilities, "CRITICAL")),
			"high":     len(filterVulnerabilities(report.Vulnerabilities, "HIGH")),
			"medium":   len(filterVulnerabilities(report.Vulnerabilities, "MEDIUM")),
			"low":      len(filterVulnerabilities(report.Vulnerabilities, "LOW")),
		},
		"successful_attacks": report.SuccessfulAttacks,
		"policy_violations":  len(filterViolatedPolicies(report.PolicyResults)),
		"anomaly_detected":   report.RuntimeMetrics.IsAnomalous,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		log.Fatalf("[ERROR] Failed to write JSON report: %v", err)
	}

	log.Println("[INFO] Summary report saved to summary_report.json")
}

// Helper to filter vulnerabilities by severity
func filterVulnerabilities(vulnerabilities []image.Vulnerability, severity string) []image.Vulnerability {
	var filtered []image.Vulnerability
	for _, v := range vulnerabilities {
		if v.Severity == severity {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// Helper to filter violated policies
func filterViolatedPolicies(policies []policy.PolicyResult) []policy.PolicyResult {
	var violated []policy.PolicyResult
	for _, p := range policies {
		if p.Violated {
			violated = append(violated, p)
		}
	}
	return violated
}

// Plain-text report for human readability
func saveAsPlainText(report scoring.Report) {
	file, err := os.Create("final_report.txt")
	if err != nil {
		log.Fatalf("[ERROR] Failed to create Plain Text report: %v", err)
	}
	defer file.Close()

	content := fmt.Sprintf(`
Container Security Report
===================================
Vulnerabilities:
- Critical: %d
- High: %d
- Medium: %d
- Low: %d

Policy Violations:
- Total Violated Policies: %d

Attack Results:
- Successful Attacks: %d / %d

Runtime Metrics:
- CPU Usage: %s
- Memory Usage: %s
- Disk I/O: %s
- Network I/O: %s

Recommendations:
`, len(filterVulnerabilities(report.Vulnerabilities, "CRITICAL")),
		len(filterVulnerabilities(report.Vulnerabilities, "HIGH")),
		len(filterVulnerabilities(report.Vulnerabilities, "MEDIUM")),
		len(filterVulnerabilities(report.Vulnerabilities, "LOW")),
		len(filterViolatedPolicies(report.PolicyResults)),
		report.SuccessfulAttacks,
		report.TotalAttacks,
		report.RuntimeMetrics.CPUUsage,
		report.RuntimeMetrics.MemoryUsage,
		report.RuntimeMetrics.DiskIO,
		report.RuntimeMetrics.NetworkIO)

	for _, rec := range report.Recommendations {
		content += fmt.Sprintf("- %s\n", rec)
	}

	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("[ERROR] Failed to write Plain Text report: %v", err)
	}

	log.Println("[INFO] Plain Text report saved to final_report.txt")
}
