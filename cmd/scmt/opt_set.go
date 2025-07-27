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

// SetCmd represents the type command
var SetCmd = &cobra.Command{
	Use:   messages.GetUse("set"),
	Short: messages.GetShort("set"),
	Long:  messages.GetLong("set"),
	Run:   handleSetCmd,
}

// handleSetCmd handles the project create command
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
		scmto.Open()

		err := scmto.SafeSet(option_name, option_value, viper.GetString("engineer"), viper.GetString("message"))

		cobra.CheckErr(err)
	}
}

func init() {
	rootCmd.AddCommand(SetCmd)
}
