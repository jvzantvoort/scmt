package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
)

func TestWriteCommand_Integration(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Create a simple template
	tmpDir := t.TempDir()
	templateFile := filepath.Join(tmpDir, "test.template")
	templateContent := `Server: {{.Config.TYPE}}
Owner: {{.Config.OWNER}}
Roles: {{range .Roles}}{{.}} {{end}}
HasWebServer: {{.HasRole "web-server"}}
Timestamp: {{.Timestamp}}
Engineer: {{.Engineer}}`

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Add test data
	testErr := roleAddCmd.RunE(roleAddCmd, []string{"web-server"})
	if testErr != nil {
		t.Fatalf("Failed to add web-server role: %v", testErr)
	}

	// Test write to stdout (no output file)
	err = WriteCmd.RunE(WriteCmd, []string{templateFile})
	if err != nil {
		t.Fatalf("Failed to process template to stdout: %v", err)
	}

	// Test write to file
	outputFile := filepath.Join(tmpDir, "output.txt")
	err = WriteCmd.RunE(WriteCmd, []string{templateFile, outputFile})
	if err != nil {
		t.Fatalf("Failed to process template to file: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Read and verify output content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "Server: server") {
		t.Error("Expected server type in output")
	}
	if !strings.Contains(output, "Owner: Mad House") {
		t.Error("Expected owner in output")
	}
	if !strings.Contains(output, "web-server") {
		t.Error("Expected web-server role in output")
	}
	if !strings.Contains(output, "HasWebServer: true") {
		t.Error("Expected HasRole function to work")
	}
	if !strings.Contains(output, "Engineer: testuser") {
		t.Error("Expected engineer name in output")
	}
}

func TestWriteCommand_NonExistentTemplate(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Test with non-existent template file
	err := WriteCmd.RunE(WriteCmd, []string{"/nonexistent/template.txt"})
	if err == nil {
		t.Error("Expected error for non-existent template file")
	}
	if !strings.Contains(err.Error(), "failed to read template file") {
		t.Errorf("Expected template read error, got: %v", err)
	}
}

func TestWriteCommand_InvalidTemplate(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Create invalid template
	tmpDir := t.TempDir()
	templateFile := filepath.Join(tmpDir, "invalid.template")
	invalidTemplate := `{{.InvalidField {{.AnotherInvalid}}`

	err := os.WriteFile(templateFile, []byte(invalidTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid template file: %v", err)
	}

	// Test with invalid template
	err = WriteCmd.RunE(WriteCmd, []string{templateFile})
	if err == nil {
		t.Error("Expected error for invalid template")
	}
	if !strings.Contains(err.Error(), "failed to parse template") {
		t.Errorf("Expected template parse error, got: %v", err)
	}
}

func TestWriteCommand_OutputDirectory(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Create simple template
	tmpDir := t.TempDir()
	templateFile := filepath.Join(tmpDir, "test.template")
	templateContent := `Type: {{.Config.TYPE}}`

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Test write to file in non-existent directory
	outputFile := filepath.Join(tmpDir, "nested", "dir", "output.txt")
	err = WriteCmd.RunE(WriteCmd, []string{templateFile, outputFile})
	if err != nil {
		t.Fatalf("Failed to process template with nested output dir: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created in nested directory")
	}
}

func TestPrepareTemplateData(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := data.New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Add test data
	_, err = d.Set("TEST_KEY", "test_value", "testuser", "test")
	if err != nil {
		t.Fatalf("Failed to set test data: %v", err)
	}
	_, err = d.AddRole("test-role", "testuser", "test")
	if err != nil {
		t.Fatalf("Failed to add test role: %v", err)
	}

	// Set global Engineer for testing
	originalEngineer := Engineer
	Engineer = "test-engineer"
	defer func() { Engineer = originalEngineer }()

	// Test prepareTemplateData
	templateData, err := prepareTemplateData(d)
	if err != nil {
		t.Fatalf("Failed to prepare template data: %v", err)
	}

	// Verify template data structure
	if templateData.Config["TEST_KEY"] != "test_value" {
		t.Errorf("Expected TEST_KEY='test_value', got %s", templateData.Config["TEST_KEY"])
	}

	if len(templateData.Roles) != 1 || templateData.Roles[0] != "test-role" {
		t.Errorf("Expected roles=['test-role'], got %v", templateData.Roles)
	}

	if templateData.Engineer != "test-engineer" {
		t.Errorf("Expected engineer='test-engineer', got %s", templateData.Engineer)
	}

	if templateData.Timestamp == "" {
		t.Error("Expected non-empty timestamp")
	}

	// Test HasRole function
	if !templateData.HasRole("test-role") {
		t.Error("Expected HasRole('test-role') to return true")
	}
	if templateData.HasRole("nonexistent-role") {
		t.Error("Expected HasRole('nonexistent-role') to return false")
	}
}

func TestTemplateData_HasRole(t *testing.T) {
	td := TemplateData{
		Roles: []string{"web-server", "database", "monitoring"},
	}

	// Test existing roles
	if !td.HasRole("web-server") {
		t.Error("Expected HasRole('web-server') to return true")
	}
	if !td.HasRole("database") {
		t.Error("Expected HasRole('database') to return true")
	}
	if !td.HasRole("monitoring") {
		t.Error("Expected HasRole('monitoring') to return true")
	}

	// Test non-existing role
	if td.HasRole("nonexistent") {
		t.Error("Expected HasRole('nonexistent') to return false")
	}

	// Test empty roles
	emptyTd := TemplateData{Roles: []string{}}
	if emptyTd.HasRole("any-role") {
		t.Error("Expected HasRole to return false for empty roles")
	}
}