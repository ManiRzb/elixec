package image

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Vulnerability represents a single vulnerability in the Trivy scan result
type Vulnerability struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	InstalledVersion string `json:"InstalledVersion"`
	Severity         string `json:"Severity"`
	PrimaryURL       string `json:"PrimaryURL"`
}

// ScanResult represents the parsed JSON result from Trivy
type ScanResult struct {
	ArtifactName string           `json:"ArtifactName"`
	Results      []ScanResultItem `json:"Results"`
}

// ScanResultItem represents individual scan results
type ScanResultItem struct {
	Target          string          `json:"Target"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}

// ScanImage scans the Docker image for vulnerabilities, calculates a score, and prints a report
func ScanImage(imageName string) ([]Vulnerability, int) {
	log.Printf("[INFO] Scanning image %s for vulnerabilities...", imageName)

	// Run Trivy with JSON output
	outputFile := "trivy_result.json"
	cmd := exec.Command("trivy", "image", "--scanners", "vuln", "--format", "json", "--output", outputFile, imageName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("[ERROR] Failed to scan image: %v", err)
	}

	// Parse Trivy's JSON output
	file, err := os.Open(outputFile)
	if err != nil {
		log.Fatalf("[ERROR] Failed to open Trivy result file: %v", err)
	}
	defer file.Close()

	var result ScanResult
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		log.Fatalf("[ERROR] Failed to parse Trivy result: %v", err)
	}

	// Aggregate vulnerabilities and calculate score
	vulnerabilities := []Vulnerability{}
	score := calculateImageScore(result, &vulnerabilities)

	// Print vulnerability report
	// printVulnerabilityReport(imageName, vulnerabilities, score)

	log.Printf("[INFO] Image Security Score: %d", score)
	return vulnerabilities, score
}

// calculateImageScore calculates the image security score and aggregates vulnerabilities
func calculateImageScore(scanResult ScanResult, vulnerabilities *[]Vulnerability) int {
	score := 100
	severityWeights := map[string]int{
		"CRITICAL": -10,
		"HIGH":     -5,
		"MEDIUM":   -2,
		"LOW":      -1,
	}

	for _, item := range scanResult.Results {
		for _, vuln := range item.Vulnerabilities {
			*vulnerabilities = append(*vulnerabilities, vuln)
			if weight, exists := severityWeights[vuln.Severity]; exists {
				score += weight
			}
		}
	}

	if score < 0 {
		score = 0
	}
	return score
}

// printVulnerabilityReport prints a summary of vulnerabilities
func printVulnerabilityReport(imageName string, vulnerabilities []Vulnerability, score int) {
	fmt.Printf("\n[INFO] Vulnerability Report for Image: %s\n", imageName)
	fmt.Println("------------------------------------------------")
	fmt.Printf("Security Score: %d\n", score)
	fmt.Println("------------------------------------------------")
	fmt.Printf("%-15s %-20s %-10s %s\n", "Severity", "Package", "Version", "Vulnerability ID")
	fmt.Println("------------------------------------------------")
	for _, vuln := range vulnerabilities {
		fmt.Printf("%-15s %-20s %-10s %s\n", vuln.Severity, vuln.PkgName, vuln.InstalledVersion, vuln.VulnerabilityID)
	}
	fmt.Println("------------------------------------------------")
}
