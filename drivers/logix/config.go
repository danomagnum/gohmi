package logix

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/danomagnum/gologix"
)

type LogixConfig struct {
	DriverName string        `json:"drivername"`
	Rate       time.Duration `json:"rate"`
	IP         string        `json:"ip"`
	Path       string        `json:"path"`
}

func (drv *LogixDriver) Load(r io.Reader) error {
	d := json.NewDecoder(r)
	var conf LogixConfig
	err := d.Decode(&conf)
	if err != nil {
		log.Printf("Error loading logix driver: %v", err)
	}

	drv.IP = conf.IP
	drv.DriverName = conf.DriverName
	drv.on_scan = make(map[string]*readdata)
	drv.client = gologix.NewClient(conf.IP)
	drv.status = "Never Started"
	drv.Rate = conf.Rate
	drv.Path = conf.Path

	p, err := gologix.ParsePath(conf.Path)
	if err != nil {
		return fmt.Errorf("bad logix path: %v", err)
	}
	drv.client.Controller.Path = p

	return nil
}
func (drv *LogixDriver) Save(w io.Writer) error {
	e := json.NewEncoder(w)
	cfg := LogixConfig{
		DriverName: drv.DriverName,
		Rate:       drv.Rate,
		IP:         drv.IP,
		Path:       drv.Path,
	}
	err := e.Encode(cfg)
	if err != nil {
		return fmt.Errorf("error writing logix driver: %v", err)
	}
	return nil

}

func LoadAll(root string) map[string]*LogixDriver {
	drivers := make(map[string]*LogixDriver, 0)
	hmitag_configs, err := filepath.Glob(path.Join(root, "logix", "*.json"))
	if err != nil {
		log.Printf("problem opening logix configs: %v", err)
		return drivers
	}
	for _, fn := range hmitag_configs {
		log.Printf("loading %s", fn)

		f, err := os.Open(fn)
		if err != nil {
			log.Printf("problem opening %s", fn)
			continue
		}
		d := LogixDriver{}
		err = d.Load(f)
		if err != nil {
			log.Printf("problem loading logix tag driver: %v", err)
		}
		drivers[d.DriverName] = &d
	}
	return drivers
}
