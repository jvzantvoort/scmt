// Package main is the entry point for the SCMT (Server Configuration Management Tool) CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/logger"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// LogCmd represents the log command for displaying audit log history
var LogCmd = &cobra.Command{
	Use:   messages.GetUse("log"),
	Short: messages.GetShort("log"),
	Long:  messages.GetLong("log"),
	Run:   handleLogCmd,
}

// handleLogCmd displays the audit log history for a specific configuration variable.
// Usage: log <option_name>
func handleLogCmd(cmd *cobra.Command, args []string) {
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 1 {
		fmt.Printf("USAGE:\n\n\t%s %s <name>\n\n", os.Args[0], cmd.Use)
		return
	}

	option_name := args[0]

	cfg := config.New()
	logh, _ := logger.New(cfg.Logfile)
	if err := logh.TableDumper(option_name, os.Stdout); err != nil {
		cobra.CheckErr(err)
	}

}

func init() {
	rootCmd.AddCommand(LogCmd)
}
