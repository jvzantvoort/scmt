// Package data provides I/O operations for persisting and retrieving
// configuration data in JSON format, with support for both file-based and stream operations.
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/jvzantvoort/scmt/utils"
)

// ConfigFile returns the path to the configuration data file and whether it exists.
func (data Data) ConfigFile() (string, bool) {
	configfile := data.Config.ConfigDatafile
	found := true

	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		found = false
	}

	return configfile, found

}

// ConfigDir returns the configuration directory path.
func (data Data) ConfigDir() string {
	return data.Config.Configdir

}

// Writer serializes the Data structure as JSON and writes it to the provided writer.
func (d Data) Writer(writer io.Writer) error {
	utils.LogStart()
	defer utils.LogEnd()

	content, err := json.MarshalIndent(d, "", "  ")
	if err == nil {
		_, err := fmt.Fprintf(writer, "%s\n", string(content))
		if err != nil {
			return err
		}
	}
	return err

}

// Dumper outputs configuration data in the specified format (json or table) to the provided writer.
func (d Data) Dumper(outputtype string, writer io.Writer) error {
	utils.LogStart()
	defer utils.LogEnd()

	mdata := map[string]string{}
	for indx, element := range d.Elements {
		utils.LogVariable(indx, element)
		mdata[element.Option] = element.Value.Value
	}

	if outputtype == "json" {
		content, err := json.MarshalIndent(mdata, "", "  ")
		if err == nil {
			_, err := fmt.Fprintf(writer, "%s\n", string(content))
			if err != nil {
				return err
			}
		}
	} else if outputtype == "table" {
		table := tablewriter.NewWriter(writer)
		table.Header([]string{"Name", "Value", "Engineer", "Changed", "Message"})
		tabledata := [][]string{}

		for _, element := range d.Elements {
			cols := []string{}
			cols = append(cols, element.Option)
			cols = append(cols, element.Value.Value)
			cols = append(cols, element.Value.Engineer)
			cols = append(cols, element.Value.Changed.Format("2006-01-02 15:04"))
			cols = append(cols, element.Value.Message)
			tabledata = append(tabledata, cols)
		}
		if err := table.Bulk(tabledata); err != nil {
			return err
		}
		return table.Render()

	}

	return nil
}

// Reader deserializes JSON data from the provided reader into the Data structure.
func (data *Data) Reader(reader io.Reader) error {
	utils.LogStart()
	defer utils.LogEnd()

	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	return nil
}

// Open loads configuration data from the data file into the Data structure.
func (data *Data) Open() error {
	utils.LogStart()
	defer utils.LogEnd()

	configfile, found := data.ConfigFile()

	if found {

		filehandle, err := os.Open(configfile)
		if err != nil {
			utils.Errorf("cannot open config file for reading: %s", err)
			return err
		}

		return data.Reader(filehandle)
	}
	return fmt.Errorf("configfile not found")

}

// Save persists the Data structure to the configuration data file.
// It creates a backup of the existing file before writing the new data.
func (data Data) Save() error {
	utils.LogStart()
	defer utils.LogEnd()

	err := utils.MkdirAll(data.ConfigDir())
	if err != nil {
		return err
	}

	configfile, _ := data.ConfigFile()
	utils.Debugf("project file: %s", configfile)
	_ = os.Rename(configfile, configfile+".bck") // Ignore error if backup fails

	filehandle, err := os.OpenFile(configfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Errorf("cannot open project file for writing: %s", err)
		return err
	}
	defer filehandle.Close()
	return data.Writer(filehandle)
}
