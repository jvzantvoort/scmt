package config

import (
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	Configdir      string
	ConfigDatafile string
	Logfile        string
	OutputJSON     bool
}

func New() *Config {
	retv := &Config{}

	retv.Configdir = viper.GetString("configdir")
	retv.Logfile = viper.GetString("logfile")
	retv.OutputJSON = viper.GetBool("json")
	retv.ConfigDatafile = path.Join(retv.Configdir, "data.json")

	return retv
}
