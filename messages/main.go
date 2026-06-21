// Package messages provides embedded help text, usage information, and version
// details for CLI commands using Go's embed functionality.
package messages

import (
	"embed"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Content contains all embedded message files for help text and documentation.
//
//go:embed long/* use/* short/* version/*
var Content embed.FS

// GetContent retrieves embedded content from the specified folder and filename.
// Returns "undefined" if the file is not found.
func GetContent(folder, name string) string {
	filename := fmt.Sprintf("%s/%s", folder, name)

	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Errorf("%s", err)
		msgstr = []byte("undefined")
	}
	return strings.TrimSuffix(string(msgstr), "\n")

}

// GetVersion retrieves the version string from embedded content.
func GetVersion() string {
	return GetContent("version", "content")

}

// GetShort retrieves the short description for a command.
func GetShort(name string) string {
	return GetContent("short", name)
}

// GetUse retrieves the usage string for a command.
func GetUse(name string) string {
	return GetContent("use", name)
}

// GetLong retrieves the long description for a command.
func GetLong(name string) string {
	return GetContent("long", name)
}
