package security

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"errors"
	"path/filepath"
	"strings"
)

// ValidatePath validates that a file path is safe and doesn't contain path traversal attempts.
// It checks for path traversal sequences (..) to prevent reading files outside intended directories.
// For absolute safety, provide baseDir to restrict file access to a specific directory tree.
func ValidatePath(filePath string, baseDir string) error {
	// Check for path traversal sequences
	if strings.Contains(filePath, "..") {
		return errors.New("path contains path traversal sequence")
	}

	// If baseDir is provided, ensure the resolved path stays within it
	if baseDir != "" {
		absBase, err := filepath.Abs(baseDir)
		if err != nil {
			return err
		}

		absPath, err := filepath.Abs(filepath.Join(baseDir, filePath))
		if err != nil {
			return err
		}

		// Clean both paths to normalize them
		absBase = filepath.Clean(absBase)
		absPath = filepath.Clean(absPath)

		// Check if the resolved path is within the base directory
		if !strings.HasPrefix(absPath, absBase+string(filepath.Separator)) &&
			absPath != absBase {
			return errors.New("path escapes base directory")
		}
	}

	return nil
}
