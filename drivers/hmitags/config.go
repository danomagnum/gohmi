package hmitags

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

type TagConfig struct {
	DriverName string         `json:"drivername"`
	Tags       map[string]any `json:"tags"`
}

func (drv *TagStore) Load(r io.Reader) error {
	d := json.NewDecoder(r)
	var conf TagConfig
	err := d.Decode(&conf)
	if err != nil {
		log.Printf("Error loading hmi tag driver: %v", err)
	}

	drv.DriverName = conf.DriverName
	if drv.data != nil {
		drv.data = conf.Tags
	} else {
		drv.data = make(map[string]any)
	}

	return nil
}
func (drv *TagStore) Save(w io.Writer) error {
	e := json.NewEncoder(w)
	cfg := TagConfig{
		DriverName: drv.DriverName,
		Tags:       drv.data,
	}
	err := e.Encode(cfg)
	if err != nil {
		log.Printf("Error writing hmi tag driver: %v", err)
	}
	return nil

}

func LoadAll(root string) map[string]*TagStore {
	drivers := make(map[string]*TagStore, 0)
	hmitag_configs, err := filepath.Glob(path.Join(root, "hmitags", "*.json"))
	if err != nil {
		log.Printf("problem opening hmitag configs: %v", err)
		return drivers
	}
	for _, fn := range hmitag_configs {
		log.Printf("loading %s", fn)

		f, err := os.Open(fn)
		if err != nil {
			log.Printf("problem opening %s", fn)
			continue
		}
		d := TagStore{}
		err = d.Load(f)
		if err != nil {
			log.Printf("problem loading hmi tag driver: %v", err)
		}
		drivers[d.DriverName] = &d
	}
	return drivers
}
