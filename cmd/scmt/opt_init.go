package main

import (
	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// InitCmd represents the type command
var InitCmd = &cobra.Command{
	Use:   messages.GetUse("init"),
	Short: messages.GetShort("init"),
	Long:  messages.GetLong("init"),
	Run:   handleInitCmd,
}

// handleInitCmd handles the project create command
func handleInitCmd(cmd *cobra.Command, args []string) {
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)
	cfg := config.New()

	if scmto, err := data.New(*cfg); err == nil {
		scmto.Init(Engineer)
		scmto.Save()
		return
	}

}

func init() {
	rootCmd.AddCommand(InitCmd)
	InitCmd.Flags().StringP("type", "t", "default", "Type of installation")
	InitCmd.Flags().StringP("description", "d", "", "Description of the installation")
}
