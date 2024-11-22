package main

import (
	"log"
	"os"

	"github.com/ManiRzb/elixec/internal/container"
	"github.com/ManiRzb/elixec/internal/image"
	"github.com/ManiRzb/elixec/internal/policy"
	"github.com/ManiRzb/elixec/internal/report"
	"github.com/ManiRzb/elixec/internal/scoring"
)

func main() {
	// Read image name dynamically
	if len(os.Args) < 2 {
		log.Fatalf("[ERROR] Usage: %s <image-name>", os.Args[0])
	}
	imageName := os.Args[1]
	log.Printf("[INFO] Starting Container Security Framework for image: %s", imageName)

	// Step 1: Scan the image
	vulnerabilities, imageScore := image.ScanImage(imageName)

	// Step 2: Deploy the container
	containerID := container.DeployContainer(imageName)

	// Step 3: Validate policies
	policies := policy.LoadPolicies("configs/policies.yaml")
	containerConfig := container.GetContainerConfig(containerID)
	policyResults := policy.ValidatePolicies(containerConfig, policies)

	// Step 4: Simulate attacks
	attackResults := container.SimulateAttacks(containerID)

	// Step 5: Monitor runtime
	runtimeMetrics := container.MonitorContainer(containerID)

	// Step 6: Generate a comprehensive report
	finalReport := scoring.GenerateFinalReport(attackResults, vulnerabilities, imageScore, runtimeMetrics, policyResults)
	report.SaveFinalReport(finalReport)
	// Step 7: Cleanup
	container.CleanupContainer(containerID)

	log.Println("[INFO] Framework execution completed.")
}
