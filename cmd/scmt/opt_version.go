// Package main is the entry point for the SCMT (Server Configuration Management Tool) CLI application.
package main

import (
	"fmt"

	"github.com/jvzantvoort/scmt/messages"
	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   messages.GetUse("version"),
	Short: messages.GetShort("version"),
	Long:  messages.GetLong("version"),
	Run:   handleVersionCmd,
}

// handleVersionCmd displays the application version.
func handleVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(messages.GetVersion())

}

func init() {
	rootCmd.AddCommand(VersionCmd)
}
