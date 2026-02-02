package messages

import (
	"strings"
	"testing"
)

func TestGetContent(t *testing.T) {
	// Test getting version content
	version := GetContent("version", "content")
	if version == "" {
		t.Error("Expected non-empty version content")
	}
	if version == "undefined" {
		t.Error("Expected actual version content, got 'undefined'")
	}

	// Test getting non-existent content
	nonExistent := GetContent("nonexistent", "file")
	if nonExistent != "undefined" {
		t.Errorf("Expected 'undefined' for non-existent content, got '%s'", nonExistent)
	}
}

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("Expected non-empty version")
	}
	if version == "undefined" {
		t.Error("Expected actual version, got 'undefined'")
	}
}

func TestGetShort(t *testing.T) {
	// Test existing short messages
	rootShort := GetShort("root")
	if rootShort == "" {
		t.Error("Expected non-empty root short message")
	}
	if rootShort == "undefined" {
		t.Error("Expected actual root short message, got 'undefined'")
	}

	// Test role short message (added in our implementation)
	roleShort := GetShort("role")
	if roleShort == "" {
		t.Error("Expected non-empty role short message")
	}
	if roleShort == "undefined" {
		t.Error("Expected actual role short message, got 'undefined'")
	}
	if !strings.Contains(strings.ToLower(roleShort), "role") {
		t.Errorf("Expected role short message to contain 'role', got '%s'", roleShort)
	}

	// Test non-existent short message
	nonExistent := GetShort("nonexistent")
	if nonExistent != "undefined" {
		t.Errorf("Expected 'undefined' for non-existent short message, got '%s'", nonExistent)
	}
}

func TestGetUse(t *testing.T) {
	// Test existing use messages
	rootUse := GetUse("root")
	if rootUse == "" {
		t.Error("Expected non-empty root use message")
	}
	if rootUse == "undefined" {
		t.Error("Expected actual root use message, got 'undefined'")
	}

	// Test role use message (added in our implementation)
	roleUse := GetUse("role")
	if roleUse == "" {
		t.Error("Expected non-empty role use message")
	}
	if roleUse == "undefined" {
		t.Error("Expected actual role use message, got 'undefined'")
	}
	if roleUse != "role" {
		t.Errorf("Expected role use message to be 'role', got '%s'", roleUse)
	}

	// Test non-existent use message
	nonExistent := GetUse("nonexistent")
	if nonExistent != "undefined" {
		t.Errorf("Expected 'undefined' for non-existent use message, got '%s'", nonExistent)
	}
}

func TestGetLong(t *testing.T) {
	// Test existing long messages
	rootLong := GetLong("root")
	if rootLong == "" {
		t.Error("Expected non-empty root long message")
	}
	if rootLong == "undefined" {
		t.Error("Expected actual root long message, got 'undefined'")
	}

	// Test role long message (added in our implementation)
	roleLong := GetLong("role")
	if roleLong == "" {
		t.Error("Expected non-empty role long message")
	}
	if roleLong == "undefined" {
		t.Error("Expected actual role long message, got 'undefined'")
	}
	if !strings.Contains(strings.ToLower(roleLong), "role") {
		t.Errorf("Expected role long message to contain 'role', got '%s'", roleLong)
	}
	if !strings.Contains(roleLong, "add") || !strings.Contains(roleLong, "remove") || !strings.Contains(roleLong, "list") {
		t.Errorf("Expected role long message to contain 'add', 'remove', and 'list', got '%s'", roleLong)
	}

	// Test non-existent long message
	nonExistent := GetLong("nonexistent")
	if nonExistent != "undefined" {
		t.Errorf("Expected 'undefined' for non-existent long message, got '%s'", nonExistent)
	}
}

func TestContentConsistency(t *testing.T) {
	// Test that content returned is consistent (no trailing newlines)
	version1 := GetVersion()
	version2 := GetContent("version", "content")
	if version1 != version2 {
		t.Errorf("GetVersion() and GetContent('version', 'content') should return the same value")
	}

	// Test that content doesn't have trailing newlines
	roleShort := GetShort("role")
	if strings.HasSuffix(roleShort, "\n") {
		t.Error("GetShort should not return content with trailing newlines")
	}

	roleLong := GetLong("role")
	if strings.HasSuffix(roleLong, "\n") {
		t.Error("GetLong should not return content with trailing newlines")
	}

	roleUse := GetUse("role")
	if strings.HasSuffix(roleUse, "\n") {
		t.Error("GetUse should not return content with trailing newlines")
	}
}