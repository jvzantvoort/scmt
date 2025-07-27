package utils

import (
	"fmt"
	"os"
)

// MkdirAll creates a directory and all parent directories
func MkdirAll(targetpath string) error {
	if len(targetpath) == 0 {
		return fmt.Errorf("mkdir called with empty directory")
	}

	LogStart()
	defer LogEnd()

	if target, err := os.Stat(targetpath); !os.IsNotExist(err) {
		if !target.IsDir() {
			return fmt.Errorf("target exists but is not a directory")
		}
		Debugf("directory already exists")
		return nil
	}

	err := os.MkdirAll(targetpath, os.FileMode(int(0755)))

	if err != nil {
		return err
	}
	return nil

}
