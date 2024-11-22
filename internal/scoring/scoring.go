package scoring

import (
	"github.com/ManiRzb/elixec/internal/image"
	"github.com/ManiRzb/elixec/internal/policy"
)

// AttackResult represents the result of a simulated attack
type AttackResult struct {
	Name        string `json:"name"`
	Severity    string `json:"severity"`
	Impact      int    `json:"impact"`
	Description string `json:"description"`
	Success     bool   `json:"success"`
	Output      string `json:"output"`
}

// RuntimeMetrics represents runtime resource usage of a container
type RuntimeMetrics struct {
	CPUUsage    string `json:"cpu_usage"`
	MemoryUsage string `json:"memory_usage"`
	DiskIO      string `json:"disk_io"`
	NetworkIO   string `json:"network_io"`
	IsAnomalous bool   `json:"isAnomalous"`
}

// Report represents the final summary of the security framework
type Report struct {
	TotalAttacks      int                   `json:"total_attacks"`
	SuccessfulAttacks int                   `json:"successful_attacks"`
	FailedAttacks     int                   `json:"failed_attacks"`
	FinalScore        int                   `json:"final_score"`
	Grade             string                `json:"grade"`
	ImageScore        int                   `json:"image_score"`
	Vulnerabilities   []image.Vulnerability `json:"vulnerabilities"`
	AttackResults     []AttackResult        `json:"attack_results"`
	RuntimeMetrics    RuntimeMetrics        `json:"runtime_metrics"`
	PolicyResults     []policy.PolicyResult `json:"policy_results"`
	Recommendations   []string              `json:"recommendations"`
}

// CalculateScore calculates the total score based on attack and policy results
func CalculateScore(attackResults []AttackResult, policyResults []policy.PolicyResult) int {
	score := 100

	// Subtract points for successful attacks
	for _, result := range attackResults {
		if result.Success {
			score += result.Impact
		}
	}

	// Subtract points for violated policies
	for _, policy := range policyResults {
		if policy.Violated {
			switch policy.Severity {
			case "Critical":
				score -= 10
			case "High":
				score -= 5
			case "Medium":
				score -= 3
			case "Low":
				score -= 1
			}
		}
	}

	if score < 0 {
		score = 0
	}
	return score
}

// AssignGrade maps the final score to a grade
func AssignGrade(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 70:
		return "Good"
	case score >= 50:
		return "Needs Improvement"
	default:
		return "Poor"
	}
}

// GenerateRecommendations provides recommendations based on successful attacks and policy violations
func GenerateRecommendations(
	attackResults []AttackResult,
	policyResults []policy.PolicyResult,
) []string {
	recommendations := []string{}

	for _, result := range attackResults {
		if result.Success {
			switch result.Severity {
			case "Critical":
				recommendations = append(recommendations, "Mitigate privilege escalation vulnerabilities by restricting user permissions.")
			case "High":
				recommendations = append(recommendations, "Secure sensitive files and enforce strict access controls.")
			case "Medium":
				recommendations = append(recommendations, "Implement resource quotas to prevent abuse.")
			case "Low":
				recommendations = append(recommendations, "Consider reviewing minor vulnerabilities for best practices.")
			}
		}
	}

	for _, policy := range policyResults {
		if policy.Violated {
			recommendations = append(recommendations, policy.Description)
		}
	}

	return recommendations
}

// GenerateFinalReport produces a comprehensive report
func GenerateFinalReport(
	attackResults []AttackResult,
	vulnerabilities []image.Vulnerability,
	imageScore int,
	runtimeMetrics RuntimeMetrics,
	policyResults []policy.PolicyResult,
) Report {
	finalScore := CalculateScore(attackResults, policyResults)
	grade := AssignGrade(finalScore)
	recommendations := GenerateRecommendations(attackResults, policyResults)

	return Report{
		TotalAttacks:      len(attackResults),
		SuccessfulAttacks: CountSuccessful(attackResults),
		FailedAttacks:     len(attackResults) - CountSuccessful(attackResults),
		FinalScore:        finalScore,
		Grade:             grade,
		ImageScore:        imageScore,
		Vulnerabilities:   vulnerabilities,
		AttackResults:     attackResults,
		RuntimeMetrics:    runtimeMetrics,
		PolicyResults:     policyResults,
		Recommendations:   recommendations,
	}
}

// CountSuccessful counts the number of successful attacks
func CountSuccessful(results []AttackResult) int {
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}
	return successCount
}
