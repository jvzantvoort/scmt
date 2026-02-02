package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestNew(t *testing.T) {
	// Set up test values in viper
	viper.Set("configdir", "/test/config")
	viper.Set("logfile", "/test/log/test.log")
	viper.Set("json", true)

	cfg := New()

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}

	if cfg.Configdir != "/test/config" {
		t.Errorf("Expected configdir '/test/config', got '%s'", cfg.Configdir)
	}

	if cfg.Logfile != "/test/log/test.log" {
		t.Errorf("Expected logfile '/test/log/test.log', got '%s'", cfg.Logfile)
	}

	if !cfg.OutputJSON {
		t.Error("Expected OutputJSON to be true")
	}

	expectedDataFile := "/test/config/data.json"
	if cfg.ConfigDatafile != expectedDataFile {
		t.Errorf("Expected ConfigDatafile '%s', got '%s'", expectedDataFile, cfg.ConfigDatafile)
	}

	// Clean up
	viper.Reset()
}

func TestNew_DefaultValues(t *testing.T) {
	// Reset viper to ensure clean state
	viper.Reset()

	// Set minimal required values
	viper.Set("configdir", "/etc/scmt")
	viper.Set("logfile", "/var/log/scmt.log")
	viper.Set("json", false)

	cfg := New()

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}

	if cfg.Configdir != "/etc/scmt" {
		t.Errorf("Expected configdir '/etc/scmt', got '%s'", cfg.Configdir)
	}

	if cfg.OutputJSON {
		t.Error("Expected OutputJSON to be false")
	}

	expectedDataFile := "/etc/scmt/data.json"
	if cfg.ConfigDatafile != expectedDataFile {
		t.Errorf("Expected ConfigDatafile '%s', got '%s'", expectedDataFile, cfg.ConfigDatafile)
	}

	// Clean up
	viper.Reset()
}