package container

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/ManiRzb/elixec/internal/anomaly"
	"github.com/ManiRzb/elixec/internal/scoring"
)

// MonitorContainer monitors runtime performance and checks for anomalies
func MonitorContainer(containerID string) scoring.RuntimeMetrics {
	log.Println("[INFO] Monitoring container performance...")
	cmd := exec.Command("docker", "stats", containerID, "--no-stream", "--format", "{{json .}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("[ERROR] Failed to monitor container: %v", err)
		return scoring.RuntimeMetrics{}
	}

	// Parse the JSON output
	var usage struct {
		CPUPerc string `json:"CPUPerc"`
		MemPerc string `json:"MemPerc"`
		BlockIO string `json:"BlockIO"`
		NetIO   string `json:"NetIO"`
	}
	err = json.Unmarshal(out.Bytes(), &usage)
	if err != nil {
		log.Printf("[ERROR] Failed to parse monitoring output: %v", err)
		return scoring.RuntimeMetrics{}
	}

	// Detect anomalies
	isAnomalous := anomaly.DetectAnomalies(usage.CPUPerc, usage.MemPerc, usage.BlockIO, usage.NetIO)
	if isAnomalous {
		log.Println("[WARNING] Anomaly detected during runtime!")
	}

	return scoring.RuntimeMetrics{
		CPUUsage:    usage.CPUPerc,
		MemoryUsage: usage.MemPerc,
		DiskIO:      usage.BlockIO,
		NetworkIO:   usage.NetIO,
		IsAnomalous: isAnomalous,
	}
}
