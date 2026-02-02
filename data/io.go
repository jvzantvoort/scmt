package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/jvzantvoort/scmt/utils"
)

func (data Data) ConfigFile() (string, bool) {
	configfile := data.Config.ConfigDatafile
	found := true

	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		found = false
	}

	return configfile, found

}

func (data Data) ConfigDir() string {
	return data.Config.Configdir

}

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

// Read session content from a [io.Reader] object.
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

// Write session configuration to a projectfile
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
