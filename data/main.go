package data

import (
	"fmt"
	"time"

	"github.com/jvzantvoort/scmt/config"
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
	Changed  time.Time `json:"changed"`
}

type DataElement struct {
	Option string           `json:"option"`
	Value  DataElementValue `json:"value"`
}

type Data struct {
	Config   config.Config `json:"-"`
	Elements []DataElement `json:"elements"`
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
	return retv, fmt.Errorf("Option %s not found", option)
}

func (d *Data) Set(option, value, engineer string) error {
	log.Debugf("Set %s to %s, start", option, value)
	defer log.Debugf("Set %s to %s, end", option, value)

	now := time.Now().UTC()
	changed := false

	for i, element := range d.Elements {
		if element.Option == option {
			d.Elements[i].Value.Value = value
			d.Elements[i].Value.Engineer = engineer
			d.Elements[i].Value.Changed = now
			changed = true
		}
	}
	if changed {
		return nil
	}

	row := DataElement{}
	row.Option = option
	row.Value.Value = value
	row.Value.Engineer = engineer
	row.Value.Changed = now
	d.Elements = append(d.Elements, row)

	return nil
}

func (d *Data) Init(engineer string) error {
	log.Debugf("Init data structure, start")
	defer log.Debugf("Init data structure, end")
	fmt.Printf("%#v\n", defaultData)

	for option, val := range defaultData {
		if err := d.Set(option, val, engineer); err != nil {
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
