package policy

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Policy represents a security policy definition
type Policy struct {
	Name        string `yaml:"name"`
	Condition   string `yaml:"condition"`
	Severity    string `yaml:"severity"`
	Action      string `yaml:"action"`
	Description string `yaml:"description"`
}

// PolicyResult represents the result of a policy validation
type PolicyResult struct {
	PolicyName  string `json:"policy_name"`
	Severity    string `json:"severity"`
	Violated    bool   `json:"violated"`
	Description string `json:"description"`
}

// LoadPolicies loads security policies from a YAML file
func LoadPolicies(filePath string) []Policy {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to load policies: %v", err)
	}

	var policyList struct {
		Policies []Policy `yaml:"policies"`
	}
	if err := yaml.Unmarshal(file, &policyList); err != nil {
		log.Fatalf("[ERROR] Failed to parse policies: %v", err)
	}
	return policyList.Policies
}

// ValidatePolicies validates a container configuration against loaded policies
func ValidatePolicies(containerConfig map[string]string, policies []Policy) []PolicyResult {
	results := []PolicyResult{}

	for _, policy := range policies {
		violated := false
		description := "Policy condition satisfied"

		// Example checks for policy conditions
		switch policy.Condition {
		case "container.User != 'root'":
			if containerConfig["User"] == "root" {
				violated = true
				description = "Container runs as root user."
			}
		case "container.CPULimit > 0":
			if containerConfig["CPULimit"] == "0" {
				violated = true
				description = "Container has no CPU limit set."
			}
		case "container.MemoryLimit > 0":
			if containerConfig["MemoryLimit"] == "0" {
				violated = true
				description = "Container has no memory limit set."
			}
		case "container.ImageTag != 'latest'":
			if containerConfig["ImageTag"] == "latest" {
				violated = true
				description = "Container uses an image with the 'latest' tag."
			}
		case "container.Privileged == false":
			if containerConfig["Privileged"] == "true" {
				violated = true
				description = "Container runs in privileged mode."
			}
		}

		results = append(results, PolicyResult{
			PolicyName:  policy.Name,
			Severity:    policy.Severity,
			Violated:    violated,
			Description: description,
		})
	}
	return results
}
