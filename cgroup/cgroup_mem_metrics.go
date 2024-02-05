package cgroup

import (
	"os"
	"strconv"
	"strings"
)

// readMetricValue reads a memory metric from the cgroup memory subsystem.
func readMetricValue(filePath string) (float64, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Convert content to string, trim whitespace, and parse to float64
	valueStr := strings.TrimSpace(string(content))
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
