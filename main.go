package main

import (
	"gohmi/drivers/hmitags"
	"gohmi/drivers/logix"
	"log/slog"
	"time"
)

var drivers = map[string]Driver{}

func main() {
	go web_startup()
	parse_screens()

	lgx := logix.NewLogixDriver("GaragePLC", "192.168.2.241", time.Second)
	drivers[lgx.Name()] = lgx

	hmitag := hmitags.NewTagStore("builtin")
	drivers[hmitag.Name()] = hmitag

	for k, v := range drivers {
		err := v.Start()
		if err != nil {
			slog.Error("failed to start %s: %w", k, err)
		}
	}

	select {}
}
