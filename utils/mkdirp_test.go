package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirAllEmptyPath(t *testing.T) {
	err := MkdirAll("")
	if err == nil {
		t.Errorf("expected error for empty path, got nil")
	}
}

func TestMkdirAllCreateDirectory(t *testing.T) {
	tmpdir := t.TempDir()
	testdir := filepath.Join(tmpdir, "test_scmt")

	err := MkdirAll(testdir)
	if err != nil {
		t.Errorf("expected no error creating directory, got %v", err)
	}

	// Check if directory was created
	info, err := os.Stat(testdir)
	if err != nil {
		t.Errorf("expected directory to exist, got error: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("expected path to be a directory")
	}
}

func TestMkdirAllNestedDirectories(t *testing.T) {
	tmpdir := t.TempDir()
	testdir := filepath.Join(tmpdir, "test", "nested", "dir", "scmt")

	err := MkdirAll(testdir)
	if err != nil {
		t.Errorf("expected no error creating nested directories, got %v", err)
	}

	// Check if directory was created
	info, err := os.Stat(testdir)
	if err != nil {
		t.Errorf("expected nested directory to exist, got error: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("expected path to be a directory")
	}
}

func TestMkdirAllExistingDirectory(t *testing.T) {
	tmpdir := t.TempDir()

	err := MkdirAll(tmpdir)
	if err != nil {
		t.Errorf("expected no error for existing directory, got %v", err)
	}
}

func TestMkdirAllFileExists(t *testing.T) {
	tmpdir := t.TempDir()
	testfile := filepath.Join(tmpdir, "testfile")

	// Create a file
	f, err := os.Create(testfile)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	f.Close()

	// Try to create a directory with the same name as the file
	err = MkdirAll(testfile)
	if err == nil {
		t.Errorf("expected error when path is a file, got nil")
	}

	if !os.IsExist(err) && err.Error() != "target exists but is not a directory" {
		// Some systems may return different errors, so we check for both
		t.Logf("expected IsExist error or 'not a directory' message, got: %v", err)
	}
}
