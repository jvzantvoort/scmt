// Package main is the entry point for the SCMT (Server Configuration Management Tool) CLI application.
package main

const (
	ApplicationName string = "scmt" // Name of the application
	ConfigDirName   string = "scmt" // Default name of the configuration directory
	// ConstDirMode specifies permissions on newly created configuration directories
	ConstDirMode int = 0755
	// ConstFileMode specifies permissions on newly created configuration files
	ConstFileMode int = 0644
)
