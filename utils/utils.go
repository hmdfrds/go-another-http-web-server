package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// HTTPDateFormat return RFC1123 formatted date (e.g. Thu, 05 Mar 2025 15:04:05 GMT)
func HTTPDateFormat(date time.Time) string {
	return date.Format(time.RFC1123)
}

// SafePath ensures that the provided path stay within the root directory.
func SafePath(root, path string) (string, error) {
	// Clean and make paths absolute
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute root path: %w", err)
	}

	absPath, err := filepath.Abs(filepath.Join(root, path))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if absPath is within absRoot
	if !strings.HasPrefix(absPath, absRoot) {
		return "", fmt.Errorf("access denied: %s is outside root directory", absPath)
	}
	return absPath, nil
}
