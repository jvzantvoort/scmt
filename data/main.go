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
	log, err := logger.New(d.Config.Logfile)
	if err != nil {
		return err
	}
	log.Add(option, value, engineer, message)
	return log.Save()
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
				d.Log(option, value, engineer, message)
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
		d.Log(option, value, engineer, message)
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

func New(cfg config.Config) (*Data, error) {
	retv := &Data{}
	retv.Config = cfg
	return retv, nil

}
