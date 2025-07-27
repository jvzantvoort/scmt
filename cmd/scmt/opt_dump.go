package main

import (
	"os"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DumpCmd represents the type command
var DumpCmd = &cobra.Command{
	Use:   messages.GetUse("dump"),
	Short: messages.GetShort("dump"),
	Long:  messages.GetLong("dump"),
	Run:   handleDumpCmd,
}

// handleDumpCmd handles the project create command
func handleDumpCmd(cmd *cobra.Command, args []string) {
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)
	cfg := config.New()

	if scmto, err := data.New(*cfg); err == nil {
		scmto.Open()
		if cfg.OutputJSON {
			scmto.Dumper("json", os.Stdout)

		} else {
			scmto.Dumper("table", os.Stdout)
		}
	}
}

func init() {
	rootCmd.AddCommand(DumpCmd)
}
