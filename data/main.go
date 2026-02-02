package data

import (
	"fmt"
	"time"

	"github.com/jvzantvoort/scmt/config"
	"github.com/jvzantvoort/scmt/logger"
	log "github.com/sirupsen/logrus"
)

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

type DataElementValue struct {
	Value    string    `json:"value"`
	Engineer string    `json:"engineer"`
	Message  string    `json:"message"`
	Changed  time.Time `json:"changed"`
}

type DataElement struct {
	Option string           `json:"option"`
	Value  DataElementValue `json:"value"`
}

type Data struct {
	Config         config.Config `json:"-"`
	logger.Records `json:"-"`    // Embedded logger records for change tracking
	Elements       []DataElement `json:"elements"`
	Roles          []string      `json:"roles"`
}

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

func (d Data) Log(option, value, engineer, message string) error {
	logInstance, err := logger.New(d.Config.Logfile)
	if err != nil {
		return err
	}
	logInstance.Add(option, value, engineer, message)
	return logInstance.Save()
}

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
