package main

import (
	"os"
	"os/user"
	"strings"

	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Configdir  string
	Logfile    string
	Engineer   string
	Message    string
	OutputJSON bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   messages.GetUse("root"),
	Short: messages.GetShort("root"),
	Long:  messages.GetLong("root"),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return err
		}

		Configdir = viper.GetString("configdir")
		Logfile = viper.GetString("logfile")
		Engineer = viper.GetString("engineer")
		Message = viper.GetString("message")
		OutputJSON = viper.GetBool("json")

		setLogLevel(viper.GetString("loglevel"))
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	// Define flags and bind them to Viper keys
	rootCmd.PersistentFlags().StringP("configdir", "C", "", "Directory for config files")
	_ = viper.BindPFlag("configdir", rootCmd.PersistentFlags().Lookup("configdir"))

	rootCmd.PersistentFlags().StringP("engineer", "E", "", "Specify engineer name")
	_ = viper.BindPFlag("engineer", rootCmd.PersistentFlags().Lookup("engineer"))

	rootCmd.PersistentFlags().StringP("loglevel", "l", "", "Specify loglevel")
	_ = viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().StringP("logfile", "L", "", "Specify logfile")
	_ = viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))

	rootCmd.PersistentFlags().StringP("message", "M", "", "Specify message")
	_ = viper.BindPFlag("message", rootCmd.PersistentFlags().Lookup("message"))

	rootCmd.PersistentFlags().BoolP("json", "J", false, "JSON Output")
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))

}

func setLogLevel(loglevel string) {
	if len(loglevel) == 0 {
		log.Debugf("loglevel is empty")
		loglevel = "info"

	}

	switch loglevel {

	case "debug", "verbose":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error", "quiet":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.Debugf("log level set to %s", loglevel)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() error {
	// get current User
	user_obj, err := user.Current()
	cobra.CheckErr(err)

	viper.SetEnvPrefix("SCMT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("configdir", "/etc/scmt")
	viper.SetDefault("logfile", "/var/log/scmt.log")
	viper.SetDefault("engineer", user_obj.Username)
	viper.SetDefault("loglevel", "info")
	viper.SetDefault("message", "")

	// ---------------------------------------------------------------------
	// Config File Handling
	// ---------------------------------------------------------------------
	home, err := os.UserHomeDir() // Find home directory.
	cobra.CheckErr(err)

	// Search config in home directory with name ".scmt" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".scmt")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Debugf("Failed using config file: %v", viper.ConfigFileUsed())
	}
	return nil
}
