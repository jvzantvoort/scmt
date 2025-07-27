package main

import (
	"fmt"

	"github.com/jvzantvoort/scmt/messages"
	"github.com/spf13/cobra"
)

// VersionCmd represents the type command
var VersionCmd = &cobra.Command{
	Use:   messages.GetUse("version"),
	Short: messages.GetShort("version"),
	Long:  messages.GetLong("version"),
	Run:   handleVersionCmd,
}

// handleVersionCmd handles the project create command
func handleVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(messages.GetVersion())

}

func init() {
	rootCmd.AddCommand(VersionCmd)
}
