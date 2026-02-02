package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecord(t *testing.T) {
	// Test Record struct creation
	record := Record{
		Option:   "TEST_OPTION",
		Value:    "test_value",
		Engineer: "testuser",
		Message:  "test message",
		Changed:  time.Now(),
	}

	if record.Option != "TEST_OPTION" {
		t.Errorf("Expected Option 'TEST_OPTION', got '%s'", record.Option)
	}
	if record.Value != "test_value" {
		t.Errorf("Expected Value 'test_value', got '%s'", record.Value)
	}
	if record.Engineer != "testuser" {
		t.Errorf("Expected Engineer 'testuser', got '%s'", record.Engineer)
	}
	if record.Message != "test message" {
		t.Errorf("Expected Message 'test message', got '%s'", record.Message)
	}
}

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "test.log")

	logger, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if logger.Logfile != logfile {
		t.Errorf("Expected logfile '%s', got '%s'", logfile, logger.Logfile)
	}

	if logger.Records == nil {
		t.Error("Expected Records slice to be initialized")
	}
}

func TestLogger_Add(t *testing.T) {
	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "test.log")

	logger, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Add a record
	logger.Add("TEST_OPTION", "test_value", "testuser", "test message")

	if len(logger.Records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(logger.Records))
	}

	record := logger.Records[0]
	if record.Option != "TEST_OPTION" {
		t.Errorf("Expected Option 'TEST_OPTION', got '%s'", record.Option)
	}
	if record.Value != "test_value" {
		t.Errorf("Expected Value 'test_value', got '%s'", record.Value)
	}
	if record.Engineer != "testuser" {
		t.Errorf("Expected Engineer 'testuser', got '%s'", record.Engineer)
	}
	if record.Message != "test message" {
		t.Errorf("Expected Message 'test message', got '%s'", record.Message)
	}
}

func TestLogger_Writer(t *testing.T) {
	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "test.log")

	logger, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Add test records
	logger.Add("TEST1", "value1", "user1", "message1")
	logger.Add("TEST2", "value2", "user2", "message2")

	var buf bytes.Buffer
	err = logger.Writer(&buf)
	if err != nil {
		t.Fatalf("Failed to write logger: %v", err)
	}

	// Verify JSON output
	var result map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	records, ok := result["records"].([]interface{})
	if !ok || len(records) != 2 {
		t.Errorf("Expected 2 records in JSON output, got %v", records)
	}
}

func TestLogger_Reader(t *testing.T) {
	testData := `{
		"records": [
			{
				"option": "TEST_OPTION",
				"value": "test_value",
				"engineer": "testuser",
				"message": "test message",
				"changed": "2023-01-01T00:00:00Z"
			}
		]
	}`

	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "test.log")

	logger, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test reading from string reader
	reader := strings.NewReader(testData)
	err = logger.Reader(reader)
	if err != nil {
		t.Fatalf("Failed to read logger: %v", err)
	}

	// Verify parsed data
	if len(logger.Records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(logger.Records))
	}

	record := logger.Records[0]
	if record.Option != "TEST_OPTION" {
		t.Errorf("Expected Option 'TEST_OPTION', got '%s'", record.Option)
	}
	if record.Value != "test_value" {
		t.Errorf("Expected Value 'test_value', got '%s'", record.Value)
	}
}

func TestLogger_SaveAndOpen(t *testing.T) {
	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "test.log")

	// Create logger and add data
	logger1, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger1.Add("TEST1", "value1", "user1", "message1")
	logger1.Add("TEST2", "value2", "user2", "message2")

	// Save logger
	err = logger1.Save()
	if err != nil {
		t.Fatalf("Failed to save logger: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		t.Fatal("Log file was not created")
	}

	// Create new logger instance and load
	logger2, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create second logger: %v", err)
	}

	err = logger2.Open()
	if err != nil {
		t.Fatalf("Failed to open logger: %v", err)
	}

	// Verify loaded data
	if len(logger2.Records) != 2 {
		t.Errorf("Expected 2 loaded records, got %d", len(logger2.Records))
	}

	// Check first record
	record := logger2.Records[0]
	if record.Option != "TEST1" {
		t.Errorf("Expected first record Option 'TEST1', got '%s'", record.Option)
	}
	if record.Value != "value1" {
		t.Errorf("Expected first record Value 'value1', got '%s'", record.Value)
	}
}

func TestLogger_Dumper_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "test.log")

	logger, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Add("TEST1", "value1", "user1", "message1")
	logger.Add("TEST2", "value2", "user2", "message2")

	var buf bytes.Buffer
	err = logger.Dumper("TEST1", "json", &buf)
	if err != nil {
		t.Fatalf("Failed to dump logger as JSON: %v", err)
	}

	// Verify JSON output contains records
	output := buf.String()
	if !strings.Contains(output, "TEST1") {
		t.Errorf("Expected JSON output to contain 'TEST1', got: %s", output)
	}
}

func TestLogger_Dumper_Table(t *testing.T) {
	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "test.log")

	logger, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Add("TEST1", "value1", "user1", "message1")

	var buf bytes.Buffer
	err = logger.Dumper("TEST1", "table", &buf)
	if err != nil {
		t.Fatalf("Failed to dump logger as table: %v", err)
	}

	output := buf.String()
	
	// Check that table contains our data
	if !strings.Contains(output, "value1") {
		t.Errorf("Expected table output to contain 'value1', got: %s", output)
	}
	if !strings.Contains(output, "user1") {
		t.Errorf("Expected table output to contain 'user1', got: %s", output)
	}
	if !strings.Contains(output, "message1") {
		t.Errorf("Expected table output to contain 'message1', got: %s", output)
	}
}

func TestLogger_Open_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	logfile := filepath.Join(tmpDir, "nonexistent.log")

	logger, err := New(logfile)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test opening non-existent file should create empty logger
	err = logger.Open()
	if err != nil {
		t.Fatalf("Expected no error when opening non-existent file: %v", err)
	}

	// Should have empty records
	if len(logger.Records) != 0 {
		t.Errorf("Expected 0 records for non-existent file, got %d", len(logger.Records))
	}
}