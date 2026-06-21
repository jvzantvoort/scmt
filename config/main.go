// Package config handles configuration management for SCMT,
// including loading from environment variables, config files, and flags.
package config

import (
	"path"

	"github.com/spf13/viper"
)

// Config represents the runtime configuration for SCMT.
type Config struct {
	Configdir      string // Directory where configuration data is stored
	ConfigDatafile string // Full path to the data.json configuration file
	Logfile        string // Path to the audit log file
	OutputJSON     bool   // Whether to output in JSON format
}

// New creates a new Config instance by reading values from Viper.
// It loads configuration from environment variables, config files, and command-line flags.
func New() *Config {
	retv := &Config{}

	retv.Configdir = viper.GetString("configdir")
	retv.Logfile = viper.GetString("logfile")
	retv.OutputJSON = viper.GetBool("json")
	retv.ConfigDatafile = path.Join(retv.Configdir, "data.json")

	return retv
}
