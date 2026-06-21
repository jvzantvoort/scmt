// Package main is the entry point for the SCMT (Server Configuration Management Tool) CLI application.
package main

import (
	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command for initializing a new configuration
var InitCmd = &cobra.Command{
	Use:   messages.GetUse("init"),
	Short: messages.GetShort("init"),
	Long:  messages.GetLong("init"),
	Run:   handleInitCmd,
}

// handleInitCmd initializes a new configuration with default values and saves it to disk.
func handleInitCmd(cmd *cobra.Command, args []string) {
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)
	cfg := config.New()

	if scmto, err := data.New(*cfg); err == nil {
		if err := scmto.Init(Engineer); err != nil {
			log.Errorf("Failed to initialize: %v", err)
			return
		}
		if err := scmto.Save(); err != nil {
			log.Errorf("Failed to save: %v", err)
		}
		return
	}

}

func init() {
	rootCmd.AddCommand(InitCmd)
	InitCmd.Flags().StringP("type", "t", "default", "Type of installation")
	InitCmd.Flags().StringP("description", "d", "", "Description of the installation")
}
