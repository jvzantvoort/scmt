package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jvzantvoort/scmt/config"
)

func TestData_AddRole(t *testing.T) {
	// Create a test config
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

	// Test adding a new role
	changed, err := d.AddRole("web-server", "testuser", "test message")
	if err != nil {
		t.Fatalf("Failed to add role: %v", err)
	}
	if !changed {
		t.Error("Expected AddRole to return true for new role")
	}
	if len(d.Roles) != 1 || d.Roles[0] != "web-server" {
		t.Errorf("Expected roles to contain 'web-server', got %v", d.Roles)
	}

	// Test adding duplicate role
	changed, err = d.AddRole("web-server", "testuser", "test message")
	if err != nil {
		t.Fatalf("Failed to add duplicate role: %v", err)
	}
	if changed {
		t.Error("Expected AddRole to return false for duplicate role")
	}
	if len(d.Roles) != 1 {
		t.Errorf("Expected only 1 role after duplicate add, got %d", len(d.Roles))
	}

	// Test adding second unique role
	changed, err = d.AddRole("database", "testuser", "test message")
	if err != nil {
		t.Fatalf("Failed to add second role: %v", err)
	}
	if !changed {
		t.Error("Expected AddRole to return true for second unique role")
	}
	if len(d.Roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(d.Roles))
	}
}

func TestData_RemoveRole(t *testing.T) {
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

	// Add some roles first
	d.Roles = []string{"web-server", "database", "cache"}

	// Test removing existing role
	changed, err := d.RemoveRole("database", "testuser", "test message")
	if err != nil {
		t.Fatalf("Failed to remove role: %v", err)
	}
	if !changed {
		t.Error("Expected RemoveRole to return true for existing role")
	}
	if len(d.Roles) != 2 {
		t.Errorf("Expected 2 roles after removal, got %d", len(d.Roles))
	}
	for _, role := range d.Roles {
		if role == "database" {
			t.Error("Role 'database' should have been removed")
		}
	}

	// Test removing non-existent role
	changed, err = d.RemoveRole("nonexistent", "testuser", "test message")
	if err == nil {
		t.Error("Expected error when removing non-existent role")
	}
	if changed {
		t.Error("Expected RemoveRole to return false for non-existent role")
	}
	if len(d.Roles) != 2 {
		t.Errorf("Expected roles count to remain unchanged, got %d", len(d.Roles))
	}

	// Test removing first role
	changed, err = d.RemoveRole("web-server", "testuser", "test message")
	if err != nil {
		t.Fatalf("Failed to remove first role: %v", err)
	}
	if !changed {
		t.Error("Expected RemoveRole to return true")
	}
	if len(d.Roles) != 1 || d.Roles[0] != "cache" {
		t.Errorf("Expected only 'cache' role remaining, got %v", d.Roles)
	}
}

func TestData_ListRoles(t *testing.T) {
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

	// Test empty roles list
	roles := d.ListRoles()
	if len(roles) != 0 {
		t.Errorf("Expected empty roles list, got %v", roles)
	}

	// Add some roles
	d.Roles = []string{"web-server", "database"}

	// Test listing roles
	roles = d.ListRoles()
	if len(roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(roles))
	}

	// Verify it's a copy (modifying returned slice shouldn't affect original)
	roles[0] = "modified"
	if d.Roles[0] == "modified" {
		t.Error("ListRoles should return a copy, not reference to original slice")
	}
}

func TestData_HasRole(t *testing.T) {
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

	// Test with empty roles
	if d.HasRole("web-server") {
		t.Error("Expected HasRole to return false for empty roles list")
	}

	// Add some roles
	d.Roles = []string{"web-server", "database"}

	// Test existing roles
	if !d.HasRole("web-server") {
		t.Error("Expected HasRole to return true for existing role 'web-server'")
	}
	if !d.HasRole("database") {
		t.Error("Expected HasRole to return true for existing role 'database'")
	}

	// Test non-existing role
	if d.HasRole("cache") {
		t.Error("Expected HasRole to return false for non-existing role 'cache'")
	}
}

func TestData_SetAndGet(t *testing.T) {
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

	// Test setting a new option
	changed, err := d.Set("TEST_OPTION", "test_value", "testuser", "test message")
	if err != nil {
		t.Fatalf("Failed to set option: %v", err)
	}
	if !changed {
		t.Error("Expected Set to return true for new option")
	}

	// Test getting the option
	value, err := d.Get("TEST_OPTION")
	if err != nil {
		t.Fatalf("Failed to get option: %v", err)
	}
	if value.Value != "test_value" {
		t.Errorf("Expected value 'test_value', got '%s'", value.Value)
	}
	if value.Engineer != "testuser" {
		t.Errorf("Expected engineer 'testuser', got '%s'", value.Engineer)
	}
	if value.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", value.Message)
	}

	// Test setting same value (should not change)
	changed, err = d.Set("TEST_OPTION", "test_value", "testuser", "same value")
	if err != nil {
		t.Fatalf("Failed to set same value: %v", err)
	}
	if changed {
		t.Error("Expected Set to return false when value unchanged")
	}

	// Test updating value
	changed, err = d.Set("TEST_OPTION", "new_value", "testuser2", "updated value")
	if err != nil {
		t.Fatalf("Failed to update option: %v", err)
	}
	if !changed {
		t.Error("Expected Set to return true when value changed")
	}

	// Verify updated value
	value, err = d.Get("TEST_OPTION")
	if err != nil {
		t.Fatalf("Failed to get updated option: %v", err)
	}
	if value.Value != "new_value" {
		t.Errorf("Expected updated value 'new_value', got '%s'", value.Value)
	}
	if value.Engineer != "testuser2" {
		t.Errorf("Expected engineer 'testuser2', got '%s'", value.Engineer)
	}

	// Test getting non-existent option
	_, err = d.Get("NONEXISTENT")
	if err == nil {
		t.Error("Expected error when getting non-existent option")
	}
}

func TestData_Init(t *testing.T) {
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

	// Test initialization
	err = d.Init("testuser")
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Check that default values are set
	expectedDefaults := map[string]string{
		"TYPE":           "server",
		"OWNER":          "Mad House",
		"COUNTRY_CODE":   "NL",
		"REGION_CODE":    "EU",
		"TIMEZONE":       "Europe/Amsterdam",
		"COMPUTE_ZONE":   "europe-west4-a",
		"COMPUTE_REGION": "europe-west3",
	}

	for key, expectedValue := range expectedDefaults {
		value, err := d.Get(key)
		if err != nil {
			t.Errorf("Failed to get default option %s: %v", key, err)
			continue
		}
		if value.Value != expectedValue {
			t.Errorf("Expected %s to be '%s', got '%s'", key, expectedValue, value.Value)
		}
		if value.Engineer != "testuser" {
			t.Errorf("Expected engineer 'testuser' for %s, got '%s'", key, value.Engineer)
		}
		if value.Message != "Initialize" {
			t.Errorf("Expected message 'Initialize' for %s, got '%s'", key, value.Message)
		}
	}
}

func TestData_SaveAndOpen(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	// Create and populate data
	d1, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	// Add some data
	d1.Set("TEST_OPTION", "test_value", "testuser", "test message")
	d1.AddRole("web-server", "testuser", "test role")
	d1.AddRole("database", "testuser", "another role")

	// Save data
	err = d1.Save()
	if err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(cfg.ConfigDatafile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Create new data instance and load
	d2, err := New(*cfg)
	if err != nil {
		t.Fatalf("Failed to create second data instance: %v", err)
	}

	err = d2.Open()
	if err != nil {
		t.Fatalf("Failed to open data: %v", err)
	}

	// Verify loaded data
	value, err := d2.Get("TEST_OPTION")
	if err != nil {
		t.Fatalf("Failed to get loaded option: %v", err)
	}
	if value.Value != "test_value" {
		t.Errorf("Expected loaded value 'test_value', got '%s'", value.Value)
	}

	// Verify loaded roles
	if len(d2.Roles) != 2 {
		t.Errorf("Expected 2 loaded roles, got %d", len(d2.Roles))
	}
	if !d2.HasRole("web-server") || !d2.HasRole("database") {
		t.Errorf("Expected loaded roles to contain 'web-server' and 'database', got %v", d2.Roles)
	}
}

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := config.Config{
		Configdir:      tmpDir,
		ConfigDatafile: filepath.Join(tmpDir, "data.json"),
		Logfile:        filepath.Join(tmpDir, "test.log"),
	}

	d, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create data: %v", err)
	}

	if d == nil {
		t.Fatal("Expected non-nil data instance")
	}

	if d.Config.Configdir != tmpDir {
		t.Errorf("Expected configdir '%s', got '%s'", tmpDir, d.Config.Configdir)
	}

	if d.Elements == nil {
		t.Error("Expected Elements slice to be initialized")
	}

	if d.Roles == nil {
		t.Error("Expected Roles slice to be initialized")
	}
}