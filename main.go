package main

import (
	"flag"
	"gohmi/drivers/hmitags"
	"gohmi/drivers/logix"
	"log/slog"
	"time"

	"github.com/danomagnum/admin"
)

var drivers = map[string]Driver{}

var configDir = flag.String("configdir", "./config", "directory where config files are located")

func main() {

	h := hmitags.LoadAll(*configDir)
	for k, v := range h {
		drivers[k] = v
	}

	l := logix.LoadAll(*configDir)
	for k, v := range l {
		drivers[k] = v
	}

	for k, v := range drivers {
		err := v.Start()
		if err != nil {
			slog.Error("failed to start %s: %w", k, err)
		}
	}
	Admin = admin.NewAdmin(admin.SetDurationTimebase(time.Millisecond))
	go web_startup()

	Admin.RegisterFunc("New Logix Driver", func() {
		Admin.RegisterStruct("New Logix Driver", logix.NewLogixDriver("New Logix Driver", "0.0.0.0", "1,0", time.Second))
	})
	Admin.RegisterFunc("New HMI Tag Driver", func() {
		Admin.RegisterStruct("New HMI Tag Driver", hmitags.NewTagStore("New HMI Tag Driver"))
	})

	select {}
}
