package messages

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Content missing godoc.
//
//go:embed long/* use/* short/* config/* version/*
var Content embed.FS

func GetContent(folder, name string) string {
	filename := fmt.Sprintf("%s/%s", folder, name)

	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Errorf("%s", err)
		msgstr = []byte("undefined")
	}
	return strings.TrimSuffix(string(msgstr), "\n")

}

func GetVersion() string {
	return GetContent("version", "content")

}

func GetShort(name string) string {
	return GetContent("short", name)
}

func GetUse(name string) string {
	return GetContent("use", name)
}

func GetLong(name string) string {
	return GetContent("long", name)
}

func GetConfig(name string) string {
	return GetContent("config", name)
}

func Copy(srcFile, destFile string, mode fs.FileMode) error {

	// Create the directory if it doesn't exist
	destDir := filepath.Dir(destFile)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory: %s", err)
		}
	}

	content := GetConfig(srcFile)

	filehandle, err := os.OpenFile(destFile, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		log.Errorf("cannot open %s file for writing: %s", destFile, err)
		return err
	}
	_, err = filehandle.WriteString(content)
	filehandle.Close()
	return err

}
