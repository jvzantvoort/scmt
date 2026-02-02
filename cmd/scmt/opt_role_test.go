package main

import (
	"path/filepath"
	"testing"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/spf13/viper"
)

// setupTestEnvironment creates a temporary directory and sets up test configuration
func setupTestEnvironment(t *testing.T) string {
	tmpDir := t.TempDir()

	// Set viper values for testing
	viper.Set("configdir", tmpDir)
	viper.Set("logfile", filepath.Join(tmpDir, "test.log"))
	viper.Set("engineer", "testuser")
	viper.Set("message", "test message")
	viper.Set("json", false)

	// Update global variables
	Configdir = tmpDir
	Logfile = filepath.Join(tmpDir, "test.log")
	Engineer = "testuser"
	Message = "test message"
	OutputJSON = false

	return tmpDir
}

// initializeTestData creates and initializes test data
func initializeTestData(t *testing.T) {
	cfg := config.New()
	d, err := data.New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	err = d.Init("testuser")
	if err != nil {
		t.Fatalf("Failed to initialize data: %v", err)
	}

	err = d.Save()
	if err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}
}

func TestRoleAddCommand_Integration(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Test adding a role using the command function directly
	err := roleAddCmd.RunE(roleAddCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to add role: %v", err)
	}

	// Verify role was actually added by checking the data file directly
	cfg := config.New()
	d, err := data.New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	err = d.Open()
	if err != nil {
		t.Fatalf("Failed to open data: %v", err)
	}

	if !d.HasRole("web-server") {
		t.Error("Expected role 'web-server' to be added")
	}
}

func TestRoleRemoveCommand_Integration(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Add a role first
	err := roleAddCmd.RunE(roleAddCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to add role: %v", err)
	}

	// Remove the role
	err = roleRemoveCmd.RunE(roleRemoveCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to remove role: %v", err)
	}

	// Verify role was removed
	cfg := config.New()
	d, err := data.New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	err = d.Open()
	if err != nil {
		t.Fatalf("Failed to open data: %v", err)
	}

	if d.HasRole("web-server") {
		t.Error("Expected role 'web-server' to be removed")
	}
}

func TestRoleListCommand_Integration(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Add multiple roles
	var err error
	err = roleAddCmd.RunE(roleAddCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to add web-server role: %v", err)
	}
	err = roleAddCmd.RunE(roleAddCmd, []string{"database"})
	if err != nil {
		t.Fatalf("Failed to add database role: %v", err)
	}

	// Verify roles exist in data
	cfg := config.New()
	d, err := data.New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	err = d.Open()
	if err != nil {
		t.Fatalf("Failed to open data: %v", err)
	}

	roles := d.ListRoles()
	if len(roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(roles))
	}
	if !d.HasRole("web-server") || !d.HasRole("database") {
		t.Errorf("Expected both 'web-server' and 'database' roles, got %v", roles)
	}

	// Test list command doesn't error
	err = roleListCmd.RunE(roleListCmd, []string{})
	if err != nil {
		t.Fatalf("Failed to list roles: %v", err)
	}
}

func TestRoleAddCommand_Duplicate_Integration(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Add a role first
	err := roleAddCmd.RunE(roleAddCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to add role: %v", err)
	}

	// Try adding the same role again - should not error
	err = roleAddCmd.RunE(roleAddCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to handle duplicate role: %v", err)
	}

	// Verify still only one role
	cfg := config.New()
	d, err := data.New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	err = d.Open()
	if err != nil {
		t.Fatalf("Failed to open data: %v", err)
	}

	roles := d.ListRoles()
	if len(roles) != 1 {
		t.Errorf("Expected 1 role after duplicate add, got %d", len(roles))
	}
}

func TestRoleRemoveCommand_NonExistent_Integration(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Try removing non-existent role - should not error (prints message instead)
	err := roleRemoveCmd.RunE(roleRemoveCmd, []string{"nonexistent"})
	if err != nil {
		t.Fatalf("Command should handle non-existent role gracefully: %v", err)
	}
}

func TestRoleCommands_JSONOutput_Integration(t *testing.T) {
	setupTestEnvironment(t)
	initializeTestData(t)

	// Enable JSON output
	viper.Set("json", true)
	OutputJSON = true

	// Test JSON output commands don't error
	err := roleAddCmd.RunE(roleAddCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to add role with JSON output: %v", err)
	}

	err = roleListCmd.RunE(roleListCmd, []string{})
	if err != nil {
		t.Fatalf("Failed to list roles with JSON output: %v", err)
	}

	err = roleRemoveCmd.RunE(roleRemoveCmd, []string{"web-server"})
	if err != nil {
		t.Fatalf("Failed to remove role with JSON output: %v", err)
	}
}
