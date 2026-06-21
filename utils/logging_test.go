package utils

import (
	"testing"
)

func TestLogIfErrorNil(t *testing.T) {
	// Should not panic or log error for nil
	LogIfError(nil)
}

func TestLogIfErrorNotNil(t *testing.T) {
	// Should log error for non-nil (we can't easily verify logging output without capturing logs)
	err := "test error"
	LogIfError(err)
}

func TestLogStart(t *testing.T) {
	// Test that LogStart doesn't panic
	LogStart()
}

func TestLogEnd(t *testing.T) {
	// Test that LogEnd doesn't panic
	LogEnd()
}

func TestLogArgument(t *testing.T) {
	// Test that LogArgument doesn't panic with various types
	LogArgument("test", "string value")
	LogArgument("number", 42)
	LogArgument("float", 3.14)
}

func TestLogVariable(t *testing.T) {
	// Test that LogVariable doesn't panic with various types
	LogVariable("test", "string value")
	LogVariable("number", 42)
	LogVariable("float", 3.14)
	LogVariable("slice", []string{"a", "b", "c"})
}

func TestDebugf(t *testing.T) {
	// Test that Debugf doesn't panic
	Debugf("test message: %s %d", "hello", 42)
}

func TestInfof(t *testing.T) {
	// Test that Infof doesn't panic
	Infof("test message: %s %d", "hello", 42)
}

func TestErrorf(t *testing.T) {
	// Test that Errorf doesn't panic
	Errorf("test error message: %s", "error details")
}

func TestWarningf(t *testing.T) {
	// Test that Warningf doesn't panic
	Warningf("test warning message: %s", "warning details")
}
