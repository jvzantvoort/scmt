package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/data"
	"github.com/jvzantvoort/scmt/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TemplateData represents the data structure available to templates
type TemplateData struct {
	Config    map[string]string `json:"config"`
	Roles     []string          `json:"roles"`
	Timestamp string            `json:"timestamp"`
	Engineer  string            `json:"engineer"`
}

// HasRole checks if a specific role exists
func (td TemplateData) HasRole(role string) bool {
	for _, r := range td.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// WriteCmd represents the write command
var WriteCmd = &cobra.Command{
	Use:   messages.GetUse("write"),
	Short: messages.GetShort("write"),
	Long:  messages.GetLong("write"),
	Args:  cobra.RangeArgs(1, 2),
	RunE:  handleWriteCmd,
}

// handleWriteCmd processes template files with server configuration data
func handleWriteCmd(cmd *cobra.Command, args []string) error {
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	templateFile := args[0]
	var outputFile string

	// Determine output destination
	if len(args) == 2 {
		outputFile = args[1]
	}

	// Load configuration
	cfg := config.New()
	d, err := data.New(*cfg)
	if err != nil {
		return fmt.Errorf("failed to create data: %w", err)
	}

	err = d.Open()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Prepare template data
	templateData, err := prepareTemplateData(d)
	if err != nil {
		return fmt.Errorf("failed to prepare template data: %w", err)
	}

	// Process template
	err = processTemplate(templateFile, outputFile, templateData)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	// Log the operation if output file is specified
	if outputFile != "" {
		if err := d.Log("TEMPLATE_WRITE", fmt.Sprintf("%s -> %s", templateFile, outputFile), Engineer, fmt.Sprintf("Template processing: %s", templateFile)); err != nil {
			log.Warnf("Failed to log template write: %v", err)
		}
	}

	return nil
}

// prepareTemplateData converts server data into template-friendly structure
func prepareTemplateData(d *data.Data) (*TemplateData, error) {
	// Extract configuration as key-value map
	configMap := make(map[string]string)
	for _, element := range d.Elements {
		configMap[element.Option] = element.Value.Value
	}

	// Get current timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	templateData := &TemplateData{
		Config:    configMap,
		Roles:     d.ListRoles(),
		Timestamp: timestamp,
		Engineer:  Engineer,
	}

	return templateData, nil
}

// processTemplate reads template file, processes it with data, and writes output
func processTemplate(templateFile, outputFile string, data *TemplateData) error {
	// Define template functions
	funcMap := template.FuncMap{
		"join":      strings.Join,
		"upper":     strings.ToUpper,
		"lower":     strings.ToLower,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"replace":   strings.Replace,
		"split":     strings.Split,
		"trim":      strings.TrimSpace,
	}

	// Read template file
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", templateFile, err)
	}

	// Parse template
	templateName := filepath.Base(templateFile)
	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateFile, err)
	}

	// Determine output destination
	var writer io.Writer
	var outputWriter *os.File

	if outputFile == "" {
		writer = os.Stdout
		log.Debugf("Writing template output to stdout")
	} else {
		// Create output directory if it doesn't exist
		outputDir := filepath.Dir(outputFile)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}

		// Create/open output file
		outputWriter, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
		}
		defer outputWriter.Close()
		writer = outputWriter
		log.Debugf("Writing template output to %s", outputFile)
	}

	// Execute template
	err = tmpl.Execute(writer, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Log success
	if outputFile != "" {
		log.Infof("Template %s processed successfully, output written to %s", templateFile, outputFile)
	} else {
		log.Debugf("Template %s processed successfully, output written to stdout", templateFile)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(WriteCmd)
}