package data

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jvzantvoort/scmt/config"
)

func TestData_ConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Test with non-existent file
	configFile, found := d.ConfigFile()
	if found {
		t.Error("Expected ConfigFile to return false for non-existent file")
	}
	expectedPath := filepath.Join(tmpDir, "data.json")
	if configFile != expectedPath {
		t.Errorf("Expected config file path '%s', got '%s'", expectedPath, configFile)
	}

	// Create the file
	file, err := os.Create(cfg.ConfigDatafile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Test with existing file
	configFile, found = d.ConfigFile()
	if !found {
		t.Error("Expected ConfigFile to return true for existing file")
	}
	if configFile != expectedPath {
		t.Errorf("Expected config file path '%s', got '%s'", expectedPath, configFile)
	}
}

func TestData_ConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	configDir := d.ConfigDir()
	if configDir != tmpDir {
		t.Errorf("Expected config dir '%s', got '%s'", tmpDir, configDir)
	}
}

func TestData_Writer(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Add some test data
	d.Set("TEST", "value", "testuser", "test")
	d.AddRole("web-server", "testuser", "test role")

	var buf bytes.Buffer
	err = d.Writer(&buf)
	if err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}

	// Verify JSON output
	var result map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Check elements
	elements, ok := result["elements"].([]interface{})
	if !ok || len(elements) != 1 {
		t.Errorf("Expected 1 element in JSON output, got %v", elements)
	}

	// Check roles
	roles, ok := result["roles"].([]interface{})
	if !ok || len(roles) != 1 {
		t.Errorf("Expected 1 role in JSON output, got %v", roles)
	}
	if roles[0] != "web-server" {
		t.Errorf("Expected role 'web-server', got %v", roles[0])
	}
}

func TestData_Dumper_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Add test data
	d.Set("TEST1", "value1", "testuser", "test")
	d.Set("TEST2", "value2", "testuser", "test")

	var buf bytes.Buffer
	err = d.Dumper("json", &buf)
	if err != nil {
		t.Fatalf("Failed to dump data as JSON: %v", err)
	}

	// Verify JSON output contains only the key-value pairs
	var result map[string]string
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["TEST1"] != "value1" {
		t.Errorf("Expected TEST1 to be 'value1', got '%s'", result["TEST1"])
	}
	if result["TEST2"] != "value2" {
		t.Errorf("Expected TEST2 to be 'value2', got '%s'", result["TEST2"])
	}
}

func TestData_Dumper_Table(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Add test data
	d.Set("TEST1", "value1", "testuser", "test message")

	var buf bytes.Buffer
	err = d.Dumper("table", &buf)
	if err != nil {
		t.Fatalf("Failed to dump data as table: %v", err)
	}

	output := buf.String()
	
	// Check that table contains our data
	if !strings.Contains(output, "TEST1") {
		t.Errorf("Expected table output to contain 'TEST1', got: %s", output)
	}
	if !strings.Contains(output, "value1") {
		t.Errorf("Expected table output to contain 'value1', got: %s", output)
	}
	if !strings.Contains(output, "testuser") {
		t.Errorf("Expected table output to contain 'testuser', got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected table output to contain 'test message', got: %s", output)
	}
}

func TestData_Reader(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	// Create test JSON data
	testData := `{
		"elements": [
			{
				"option": "TEST_OPTION",
				"value": {
					"value": "test_value",
					"engineer": "testuser",
					"message": "test message",
					"changed": "2023-01-01T00:00:00Z"
				}
			}
		],
		"roles": ["web-server", "database"]
	}`

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Test reading from string reader
	reader := strings.NewReader(testData)
	err = d.Reader(reader)
	if err != nil {
		t.Fatalf("Failed to read data: %v", err)
	}

	// Verify parsed data
	value, err := d.Get("TEST_OPTION")
	if err != nil {
		t.Fatalf("Failed to get parsed option: %v", err)
	}
	if value.Value != "test_value" {
		t.Errorf("Expected value 'test_value', got '%s'", value.Value)
	}

	// Verify parsed roles
	if len(d.Roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(d.Roles))
	}
	if !d.HasRole("web-server") || !d.HasRole("database") {
		t.Errorf("Expected roles to contain 'web-server' and 'database', got %v", d.Roles)
	}
}

func TestData_Open_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "nonexistent.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Test opening non-existent file
	err = d.Open()
	if err == nil {
		t.Error("Expected error when opening non-existent file")
	}
	if !strings.Contains(err.Error(), "configfile not found") {
		t.Errorf("Expected 'configfile not found' error, got '%s'", err.Error())
	}
}

func TestData_Reader_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Test reading invalid JSON
	invalidJSON := `{"invalid": json}`
	reader := strings.NewReader(invalidJSON)
	err = d.Reader(reader)
	if err == nil {
		t.Error("Expected error when reading invalid JSON")
	}
}