// Package main is the entry point for the SCMT (Server Configuration Management Tool) CLI application.
package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GetString retrieves a string flag value from a Cobra command.
// Logs a debug message about the retrieved value.
func GetString(cmd cobra.Command, name string) string {
	retv, _ := cmd.Flags().GetString(name)
	if len(retv) != 0 {
		log.Debugf("%s returned %s", name, retv)
	} else {
		log.Debugf("%s returned nothing", name)
	}
	return retv
}
