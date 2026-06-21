// Package main is the entry point for the SCMT (Server Configuration Management Tool) CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SetCmd represents the set command for updating configuration variables
var SetCmd = &cobra.Command{
	Use:   messages.GetUse("set"),
	Short: messages.GetShort("set"),
	Long:  messages.GetLong("set"),
	Run:   handleSetCmd,
}

// handleSetCmd updates a configuration variable with the provided value and metadata.
// Usage: set <option_name> <option_value>
func handleSetCmd(cmd *cobra.Command, args []string) {
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 2 {
		fmt.Printf("USAGE:\n\n\t%s <name> <value>\n\n", os.Args[0])
		fmt.Printf("%#v\n", os.Args)
		return
	}

	option_name := args[0]
	option_value := args[1]

	cfg := config.New()

	if scmto, err := data.New(*cfg); err == nil {
		if err := scmto.Open(); err != nil {
			cobra.CheckErr(err)
			return
		}

		err := scmto.SafeSet(option_name, option_value, viper.GetString("engineer"), viper.GetString("message"))

		cobra.CheckErr(err)
	}
}

func init() {
	rootCmd.AddCommand(SetCmd)
}
