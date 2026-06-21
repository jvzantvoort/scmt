// Package data provides core data structures and operations for managing
// server configuration variables with change tracking and persistence.
package data

import (
	"fmt"
	"time"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/logger"
	log "github.com/sirupsen/logrus"
)

// defaultData contains the initial configuration values set during initialization
var (
	defaultData = map[string]string{
		"TYPE":           "server",
		"OWNER":          "Mad House",
		"COUNTRY_CODE":   "NL",
		"REGION_CODE":    "EU",
		"TIMEZONE":       "Europe/Amsterdam",
		"COMPUTE_ZONE":   "europe-west4-a",
		"COMPUTE_REGION": "europe-west3",
	}
)

// DataElementValue contains the metadata for a configuration value
type DataElementValue struct {
	Value    string    `json:"value"`    // The actual configuration value
	Engineer string    `json:"engineer"` // Name of engineer who set this value
	Message  string    `json:"message"`  // Reason for the change
	Changed  time.Time `json:"changed"`  // When the value was last changed
}

// DataElement represents a single configuration option with its current value and metadata
type DataElement struct {
	Option string           `json:"option"` // Configuration variable name
	Value  DataElementValue `json:"value"`  // Current value and metadata
}

// Data represents the complete configuration state
type Data struct {
	Config         config.Config `json:"-"`
	logger.Records `json:"-"`    // Embedded logger records for change tracking
	Elements       []DataElement `json:"elements"`
	Roles          []string      `json:"roles"`
}

// Get retrieves the value metadata for a given configuration option.
// Returns an error if the option is not found.
func (d Data) Get(option string) (*DataElementValue, error) {
	retv := &DataElementValue{}

	for _, row := range d.Elements {
		if row.Option != option {
			continue
		}
		// return value, engineer, changed and success
		return &row.Value, nil
	}
	return retv, fmt.Errorf("option %s not found", option)
}

// Log records a change to the audit log file.
func (d Data) Log(option, value, engineer, message string) error {
	logInstance, err := logger.New(d.Config.Logfile)
	if err != nil {
		return err
	}
	logInstance.Add(option, value, engineer, message)
	return logInstance.Save()
}

// Set updates a configuration option with a new value and metadata.
// Returns true if the value was changed, false if it was unchanged.
// If the option doesn't exist, it will be created.
func (d *Data) Set(option, value, engineer, message string) (bool, error) {
	log.Debugf("Set %s to %s, start", option, value)
	defer log.Debugf("Set %s to %s, end", option, value)
	log.Debugf("   By:     %s", engineer)
	log.Debugf("   Reason: %s", message)

	now := time.Now().UTC()
	changed := false
	found := false

	for i, element := range d.Elements {
		if element.Option == option {
			log.Debugf("found %s", option)
			orgval := d.Elements[i].Value.Value
			if orgval == value {
				log.Debugf("value is unchanged")
			} else {
				log.Debugf("value changed from %s to %s", orgval, value)
				if err := d.Log(option, value, engineer, message); err != nil {
					log.Warnf("Failed to log change: %v", err)
				}
				d.Elements[i].Value.Value = value
				d.Elements[i].Value.Engineer = engineer
				d.Elements[i].Value.Message = message
				d.Elements[i].Value.Changed = now
				changed = true
			}
			found = true
		}
	}

	// add if not found
	if !found {

		row := DataElement{}
		row.Option = option
		row.Value.Value = value
		row.Value.Engineer = engineer
		row.Value.Changed = now
		row.Value.Message = message
		if err := d.Log(option, value, engineer, message); err != nil {
			log.Warnf("Failed to log change: %v", err)
		}
		d.Elements = append(d.Elements, row)
		changed = true
	}

	return changed, nil
}

// SafeSet updates a configuration option and saves the data to disk if changed.
func (d *Data) SafeSet(option, value, engineer, message string) error {
	log.Debugf("Set %s to %s, start", option, value)
	defer log.Debugf("Set %s to %s, end", option, value)
	changed, err := d.Set(option, value, engineer, message)
	if err != nil {
		return err
	}
	if changed {
		return d.Save()
	}
	return nil
}

// Init initializes the Data structure with default configuration values.
func (d *Data) Init(engineer string) error {
	log.Debugf("Init data structure, start")
	defer log.Debugf("Init data structure, end")

	for option, val := range defaultData {
		_, err := d.Set(option, val, engineer, "Initialize")
		if err != nil {
			return err
		}
	}
	return nil
}

// AddRole adds a role to the roles list if it doesn't already exist
func (d *Data) AddRole(role, engineer, message string) (bool, error) {
	// Check if role already exists
	for _, r := range d.Roles {
		if r == role {
			return false, nil // Role already exists, no change
		}
	}

	// Add the role
	d.Roles = append(d.Roles, role)
	if err := d.Log("ROLE_ADD", role, engineer, message); err != nil {
		log.Warnf("Failed to log role addition: %v", err)
	}
	return true, nil
}

// RemoveRole removes a role from the roles list
func (d *Data) RemoveRole(role, engineer, message string) (bool, error) {
	for i, r := range d.Roles {
		if r == role {
			// Remove the role by slicing
			d.Roles = append(d.Roles[:i], d.Roles[i+1:]...)
			if err := d.Log("ROLE_REMOVE", role, engineer, message); err != nil {
				log.Warnf("Failed to log role removal: %v", err)
			}
			return true, nil
		}
	}
	return false, fmt.Errorf("role %s not found", role)
}

// ListRoles returns a copy of the roles list
func (d *Data) ListRoles() []string {
	roles := make([]string, len(d.Roles))
	copy(roles, d.Roles)
	return roles
}

// HasRole checks if a role exists in the roles list
func (d *Data) HasRole(role string) bool {
	for _, r := range d.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func New(cfg config.Config) (*Data, error) {
	retv := &Data{}
	retv.Config = cfg
	retv.Elements = make([]DataElement, 0)
	retv.Roles = make([]string, 0)
	return retv, nil

}
