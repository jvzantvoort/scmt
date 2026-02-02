package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirAll(t *testing.T) {
	tmpDir := t.TempDir()

	// Test creating a new directory
	newDir := filepath.Join(tmpDir, "test", "nested", "directory")
	err := MkdirAll(newDir)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}

	// Test creating directory that already exists (should not error)
	err = MkdirAll(newDir)
	if err != nil {
		t.Errorf("Expected no error when directory already exists, got: %v", err)
	}

	// Test with empty path
	err = MkdirAll("")
	if err == nil {
		t.Error("Expected error when path is empty")
	}
	if err.Error() != "mkdir called with empty directory" {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}

	// Test when target exists but is not a directory
	testFile := filepath.Join(tmpDir, "testfile")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	err = MkdirAll(testFile)
	if err == nil {
		t.Error("Expected error when target exists but is not a directory")
	}
	if err.Error() != "target exists but is not a directory" {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestLogIfError(t *testing.T) {
	// Test with nil (should not panic)
	LogIfError(nil)

	// Test with non-nil error (should not panic)
	LogIfError("test error")

	// Since this function only logs, we can't easily test the output
	// but we can ensure it doesn't panic
}

func TestLogStart(t *testing.T) {
	// Test that LogStart doesn't panic
	LogStart()
}

func TestLogEnd(t *testing.T) {
	// Test that LogEnd doesn't panic
	LogEnd()
}
