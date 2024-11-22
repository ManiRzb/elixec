package container

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/ManiRzb/elixec/internal/scoring"
	"gopkg.in/yaml.v2"
)

// Attack defines the structure of each attack
type Attack struct {
	Name        string `yaml:"name"`
	Command     string `yaml:"command"`
	Severity    string `yaml:"severity"`
	Impact      int    `yaml:"impact"`
	Description string `yaml:"description"`
}

// LoadAttacks loads attack definitions from a YAML file
func LoadAttacks(filePath string) []Attack {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read attacks YAML file: %v", err)
	}

	var attackList struct {
		Attacks []Attack `yaml:"attacks"`
	}
	err = yaml.Unmarshal(data, &attackList)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse YAML file: %v", err)
	}

	log.Printf("[INFO] Successfully loaded %d attacks from YAML file.", len(attackList.Attacks))
	return attackList.Attacks
}

// SimulateAttacks executes attacks and collects results
func SimulateAttacks(containerID string) []scoring.AttackResult {
	attacks := LoadAttacks("configs/attacks.yaml")
	results := []scoring.AttackResult{}

	for i, attack := range attacks {
		log.Printf("[INFO] Running attack %d: %s", i+1, attack.Name)
		cmd := exec.Command("docker", "exec", containerID, "sh", "-c", attack.Command)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		result := scoring.AttackResult{
			Name:        attack.Name,
			Severity:    attack.Severity,
			Impact:      attack.Impact,
			Description: attack.Description,
			Success:     err == nil,
			Output:      out.String(),
		}

		if err != nil {
			log.Printf("[WARNING] Attack failed: %v", err)
		} else {
			log.Printf("[INFO] Attack succeeded")
		}

		results = append(results, result)
	}
	return results
}
