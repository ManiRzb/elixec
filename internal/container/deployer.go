package container

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

// DeployContainer deploys the container and returns its ID
func DeployContainer(imageName string) string {
	log.Println("[INFO] Deploying container...")
	cmd := exec.Command("docker", "run", "-d", "--name", "secure-container", imageName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("[ERROR] Failed to deploy container: %v", err)
	}
	containerID := strings.TrimSpace(out.String())
	log.Printf("[INFO] Container deployed: %s", containerID)
	return containerID
}

// CleanupContainer stops and removes the container
func CleanupContainer(containerID string) {
	log.Println("[INFO] Stopping and removing container...")
	cmd := exec.Command("docker", "rm", "-f", containerID)
	err := cmd.Run()
	if err != nil {
		log.Printf("[ERROR] Failed to cleanup container: %v", err)
	} else {
		log.Printf("[INFO] Container %s removed successfully.", containerID)
	}
}
