package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetString(cmd cobra.Command, name string) string {
	retv, _ := cmd.Flags().GetString(name)
	if len(retv) != 0 {
		log.Debugf("%s returned %s", name, retv)
	} else {
		log.Debugf("%s returned nothing", name)
	}
	return retv
}
