// Package logger provides audit logging functionality for tracking all
// configuration changes with metadata including engineer, timestamp, and reason.
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/jvzantvoort/scmt/utils"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// Record represents a single audit log entry for a configuration change.
type Record struct {
	Option   string    `json:"option"`   // Option name or key
	Value    string    `json:"value"`    // Value associated with the option
	Engineer string    `json:"engineer"` // Engineer who made the change
	Message  string    `json:"message"`  // Description or message
	Changed  time.Time `json:"changed"`  // Timestamp of change
}

// Records is a slice of Record entries for iteration over the audit log.
type Records []Record

// Logger manages audit log entries and provides file I/O operations.
type Logger struct {
	Logfile string   `json:"-"`    // Path to the audit log file
	Records []Record `json:"records"` // Slice of all audit log entries
}

// Writer writes the Logger as indented JSON to the provided io.Writer.
func (rec Logger) Writer(writer io.Writer) error {
	utils.LogStart()
	defer utils.LogEnd()

	content, err := json.MarshalIndent(rec, "", "  ")
	if err == nil {
		_, err := fmt.Fprintf(writer, "%s\n", string(content))
		if err != nil {
			return err
		}
	}
	return err
}

// Reader loads Logger from a JSON-encoded io.Reader.
func (rec *Logger) Reader(reader io.Reader) error {
	utils.LogStart()
	defer utils.LogEnd()

	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &rec)
	if err != nil {
		return err
	}
	return nil
}

// Select returns all log records for a specific configuration option,
// sorted by timestamp (most recent first).
func (rec Logger) Select(option string) []Record {
	utils.LogStart()
	defer utils.LogEnd()

	retv := []Record{}

	for _, row := range rec.Records {
		if option == row.Option {
			retv = append(retv, row)
		}
	}
	sort.Slice(retv, func(i, j int) bool {
		return retv[i].Changed.After(retv[j].Changed)
	})
	return retv
}

// JSONDumper outputs selected records as JSON to the provided writer.
func (rec Logger) JSONDumper(option string, writer io.Writer) error {
	dataset := rec.Select(option)

	if content, err := json.MarshalIndent(dataset, "", "  "); err == nil {
		_, nerr := fmt.Fprintf(writer, "%s\n", string(content))
		if nerr != nil {
			return nerr
		}
	}
	return nil
}

// TableDumper outputs selected records as a formatted table to the provided writer.
func (rec Logger) TableDumper(option string, writer io.Writer) error {
	dataset := rec.Select(option)
	table := tablewriter.NewWriter(writer)
	table.Header([]string{"Value", "Engineer", "Changed", "Message"})
	tabledata := [][]string{}

	for _, element := range dataset {
		cols := []string{}
		cols = append(cols, element.Value)
		cols = append(cols, element.Engineer)
		cols = append(cols, element.Changed.Format("2006-01-02 15:04"))
		cols = append(cols, element.Message)
		tabledata = append(tabledata, cols)
	}
	if err := table.Bulk(tabledata); err != nil {
		return err
	}
	return table.Render()
}

// Dumper outputs selected Logger in either JSON or table format to the writer.
func (rec Logger) Dumper(option, outputtype string, writer io.Writer) error {
	utils.LogStart()
	defer utils.LogEnd()

	if outputtype == "json" {
		if err := rec.JSONDumper(option, writer); err != nil {
			return err
		}
	}

	if outputtype == "table" {
		if err := rec.TableDumper(option, writer); err != nil {
			return err
		}
	}

	return nil
}

// Add appends a new audit log record with the current timestamp.
func (rec *Logger) Add(option, value, engineer, message string) {
	utils.LogStart()
	defer utils.LogEnd()

	row := Record{
		Option:   option,
		Value:    value,
		Engineer: engineer,
		Message:  message,
		Changed:  time.Now(),
	}

	rec.Records = append(rec.Records, row)
}

// Log adds a record to the audit log and saves it to disk immediately.
func (rec *Logger) Log(option, value, engineer, message string) error {
	utils.LogStart()
	defer utils.LogEnd()
	log.Infof("%s %s %s", option, value, engineer)
	rec.Add(option, value, engineer, message)
	return rec.Save()
}

// Open loads Logger from the specified logfile.
func (rec *Logger) Open() error {
	utils.LogStart()
	defer utils.LogEnd()

	// target doesn't exist
	if _, err := os.Stat(rec.Logfile); os.IsNotExist(err) {
		return nil
	}

	filehandle, err := os.Open(rec.Logfile)
	if err != nil {
		return err
	}
	defer func() { _ = filehandle.Close() }()

	return rec.Reader(filehandle)
}

// Save writes Logger to the specified logfile.
func (rec Logger) Save() error {
	utils.LogStart()
	defer utils.LogEnd()

	filehandle, err := os.OpenFile(rec.Logfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = filehandle.Close() }()
	return rec.Writer(filehandle)
}

// New creates a new Logger instance and loads existing audit logs from the file.
func New(logfile string) (*Logger, error) {
	utils.LogStart()
	defer utils.LogEnd()

	rec := &Logger{}
	rec.Logfile = logfile
	rec.Records = []Record{}
	err := rec.Open()
	return rec, err
}
