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

	/*

		log.Debugf("Engineer: %s", Engineer)
		log.Debugf("Message: %s", Message)
		log.Debugf("Configdir: %s", Configdir)
		log.Debugf("Logfile: %s", Logfile)
		log.Debugf("OutputJSON: %v", OutputJSON)

			ProjectType := GetString(*cmd, "type")
			project_description := GetString(*cmd, "description")

			PrintFull, _ := cmd.Flags().GetBool("full")

			proj := project.NewProject(ProjectName)
			proj.SetDescription(project_description)
			err := proj.InitializeProject(ProjectType, true)
			if err != nil {
				utils.Fatalf("Encountered error: %q", err)
			} else {
				utils.Debugf("InitializeProject completed")

			}
	*/
}

func init() {
	rootCmd.AddCommand(InitCmd)
	InitCmd.Flags().StringP("type", "t", "default", "Type of installation")
	InitCmd.Flags().StringP("description", "d", "", "Description of the installation")
}
