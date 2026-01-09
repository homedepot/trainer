package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "path-validation-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		filePath string
		baseDir  string
		wantErr  bool
	}{
		{
			name:     "valid relative path",
			filePath: "data/config.yml",
			baseDir:  tmpDir,
			wantErr:  false,
		},
		{
			name:     "path traversal with ../",
			filePath: "../../../etc/passwd",
			baseDir:  tmpDir,
			wantErr:  true,
		},
		{
			name:     "path traversal with ..",
			filePath: "data/../../secrets.txt",
			baseDir:  tmpDir,
			wantErr:  true,
		},
		{
			name:     "valid nested path",
			filePath: "plans/test/config.yml",
			baseDir:  tmpDir,
			wantErr:  false,
		},
		{
			name:     "empty base dir allows any path without ..",
			filePath: "data/config.yml",
			baseDir:  "",
			wantErr:  false,
		},
		{
			name:     "empty base dir rejects ..",
			filePath: "../data/config.yml",
			baseDir:  "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.filePath, tt.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePath_ActualFileSystem(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "path-validation-fs-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Test valid subdirectory path
	err = ValidatePath("subdir/file.txt", tmpDir)
	if err != nil {
		t.Errorf("ValidatePath() rejected valid subdirectory path: %v", err)
	}

	// Test that .. cannot escape
	err = ValidatePath("subdir/../../../etc/passwd", tmpDir)
	if err == nil {
		t.Error("ValidatePath() allowed path traversal escape")
	}
}
