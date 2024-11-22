package anomaly

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ConvertToKilobytes converts a human-readable size (e.g., "3.215MiB") into kilobytes.
func ConvertToKilobytes(input string) (float64, error) {
	// Regular expression to match numbers with units (e.g., "3.215MiB")
	re := regexp.MustCompile(`([0-9.]+)([a-zA-Z]+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid format for IO data: %s", input)
	}

	// Extract the numeric value and unit
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse IO value: %v", err)
	}
	unit := matches[2]

	// Convert to kilobytes based on the unit
	var multiplier float64
	switch unit {
	case "B":
		multiplier = 1.0 / 1000.0 // Bytes to kilobytes
	case "kB":
		multiplier = 1 // Already in kilobytes
	case "MB", "MiB":
		multiplier = 1024 // Megabytes to kilobytes
	case "GB", "GiB":
		multiplier = 1024 * 1024 // Gigabytes to kilobytes
	case "TB", "TiB":
		multiplier = 1024 * 1024 * 1024 // Terabytes to kilobytes
	default:
		return 0, fmt.Errorf("unsupported IO unit: %s", unit)
	}

	return value * multiplier, nil
}

// DetectAnomalies processes metrics and detects anomalies using a Python script.
func DetectAnomalies(cpu string, memory string, diskIO string, networkIO string) bool {
	log.Println("[INFO] Running anomaly detection...")

	// Extract and clean CPU usage
	cpuVal, err := CleanPercValue(cpu)
	if err != nil {
		log.Printf("[ERROR] Failed to clean CPU value: %v", err)
		return false
	}

	// Extract and clean memory usage
	memoryVal, err := CleanPercValue(memory)
	if err != nil {
		log.Printf("[ERROR] Failed to clean Memory value: %v", err)
		return false
	}

	// Convert DiskIO and NetworkIO values to kilobytes
	diskIOVal, err := ConvertToKilobytes(diskIO)
	if err != nil {
		log.Printf("[ERROR] Failed to parse Disk IO: %v", err)
		return false
	}

	networkIOVal, err := ConvertToKilobytes(networkIO)
	if err != nil {
		log.Printf("[ERROR] Failed to parse Network IO: %v", err)
		return false
	}

	// Log cleaned values for debugging
	log.Printf("[INFO] Cleaned values - CPU: %.2f%%, Memory: %.2f%%, DiskIO: %.2f KB, NetworkIO: %.2f KB",
		cpuVal, memoryVal, diskIOVal, networkIOVal)

	// Execute the Python script for anomaly detection
	cmd := exec.Command("python3", "./python/detect_anomaly.py",
		fmt.Sprintf("%.2f", cpuVal), fmt.Sprintf("%.2f", memoryVal),
		fmt.Sprintf("%.2f", diskIOVal), fmt.Sprintf("%.2f", networkIOVal))

	// Capture the output from the Python script
	output, err := cmd.Output()
	if err != nil {
		log.Printf("[ERROR] Failed to execute anomaly detection script: %v", err)
		return false
	}

	// Parse the result from the Python script
	result := strings.TrimSpace(string(output))
	log.Printf("[INFO] Anomaly Detection Result: %s", result)

	// Check if the result indicates an anomaly
	return result == "Anomaly Detected"
}

// CleanPercValue extracts the numeric percentage from the input string.
func CleanPercValue(value string) (float64, error) {
	// Remove the '%' symbol
	cleaned := strings.TrimRight(value, "%")
	// Convert the remaining string to a float
	cleanedValue, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse percentage value: %v", err)
	}
	return cleanedValue, nil
}
