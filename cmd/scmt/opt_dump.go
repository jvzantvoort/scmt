// Package main is the entry point for the SCMT (Server Configuration Management Tool) CLI application.
package main

import (
	"os"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DumpCmd represents the dump command for displaying all configuration variables
var DumpCmd = &cobra.Command{
	Use:   messages.GetUse("dump"),
	Short: messages.GetShort("dump"),
	Long:  messages.GetLong("dump"),
	Run:   handleDumpCmd,
}

// handleDumpCmd displays all configuration variables in the specified output format.
func handleDumpCmd(cmd *cobra.Command, args []string) {
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)
	cfg := config.New()

	if scmto, err := data.New(*cfg); err == nil {
		if err := scmto.Open(); err != nil {
			log.Errorf("Failed to open data: %v", err)
			return
		}
		if cfg.OutputJSON {
			if err := scmto.Dumper("json", os.Stdout); err != nil {
				log.Errorf("Failed to dump JSON: %v", err)
			}
		} else {
			if err := scmto.Dumper("table", os.Stdout); err != nil {
				log.Errorf("Failed to dump table: %v", err)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(DumpCmd)
}
