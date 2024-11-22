package container

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// ContainerConfig represents the configuration of a running container
type ContainerConfig struct {
	User        string `json:"User"`
	CPULimit    int    `json:"CpuShares"` // CPU limit is an integer
	MemoryLimit int64  `json:"Memory"`    // Memory limit is an integer
	ImageTag    string `json:"ImageTag"`
	Privileged  bool   `json:"Privileged"` // Privileged mode is a boolean
}

// GetContainerConfig retrieves the container configuration for policy validation
func GetContainerConfig(containerID string) map[string]string {
	log.Printf("[INFO] Fetching configuration for container: %s", containerID)

	// Get container details
	cmd := exec.Command("docker", "inspect", containerID)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("[ERROR] Failed to fetch container configuration: %v", err)
	}

	// Log the raw output for debugging purposes
	// log.Printf("[INFO] Raw output from docker inspect: %s", out.String())

	// Parse the JSON output
	var containerInfo []struct {
		Config struct {
			User       string `json:"User"`
			Privileged bool   `json:"Privileged"`
		} `json:"Config"`
		HostConfig struct {
			CpuShares int   `json:"CpuShares"`
			Memory    int64 `json:"Memory"`
		} `json:"HostConfig"`
		RepoTags []string `json:"RepoTags"`
	}

	// Handle case where docker inspect returns no data
	if len(out.Bytes()) == 0 {
		log.Fatalf("[ERROR] No data returned from docker inspect for container: %s", containerID)
	}

	err = json.Unmarshal(out.Bytes(), &containerInfo)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse container configuration: %v", err)
	}

	// Check if containerInfo is empty, which would indicate an invalid container ID or empty response
	if len(containerInfo) == 0 {
		log.Fatalf("[ERROR] No configuration found for container ID: %s", containerID)
	}

	// Extract relevant configuration details with safe checks
	config := map[string]string{}

	// Safe check for 'User'
	if containerInfo[0].Config.User != "" {
		config["User"] = containerInfo[0].Config.User
	} else {
		config["User"] = "default"
	}

	// Safe check for 'Privileged'
	if containerInfo[0].Config.Privileged {
		config["Privileged"] = "true"
	} else {
		config["Privileged"] = "false"
	}

	// Safe check for 'CpuShares'
	if containerInfo[0].HostConfig.CpuShares != 0 {
		config["CPULimit"] = fmt.Sprintf("%d", containerInfo[0].HostConfig.CpuShares)
	} else {
		config["CPULimit"] = "0"
	}

	// Safe check for 'Memory'
	if containerInfo[0].HostConfig.Memory != 0 {
		config["MemoryLimit"] = fmt.Sprintf("%d", containerInfo[0].HostConfig.Memory)
	} else {
		config["MemoryLimit"] = "0"
	}

	// Safe check for 'ImageTag'
	if len(containerInfo[0].RepoTags) > 0 {
		config["ImageTag"] = containerInfo[0].RepoTags[0]
	} else {
		config["ImageTag"] = "unknown"
	}

	return config
}
