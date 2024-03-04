package logix

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/danomagnum/admin"
	"github.com/danomagnum/gologix"
)

type LogixDriver struct {
	DriverName   string
	m            sync.RWMutex
	client       *gologix.Client
	on_scan      map[string]*readdata
	close_signal context.CancelFunc
	status       string
	Rate         time.Duration
	IP           string
	Path         string
}

func NewLogixDriver(name, ip, path string, rate time.Duration) *LogixDriver {

	drv := &LogixDriver{
		IP:         ip,
		DriverName: name,
		on_scan:    make(map[string]*readdata),
		client:     gologix.NewClient(ip),
		status:     "Never Started",
		Rate:       rate,
		Path:       path,
	}
	p, err := gologix.ParsePath(path)
	if err != nil {
		log.Panicf("bad path: %v", err)
	}
	drv.client.Path = p
	return drv
}

type readdata struct {
	lastreq time.Time
	value   any
}

func (drv *LogixDriver) Read(key string) (any, error) {
	drv.m.RLock()
	dat, ok := drv.on_scan[key]
	drv.m.RUnlock()

	if !ok {
		// new tag we need to read the value and add then to the scan list.

		dat = &readdata{}
		val, err := drv.client.Read_single(key, gologix.CIPTypeUnknown, 1)
		if err != nil {
			return nil, err
		}
		dat.value = val
		dat.lastreq = time.Now()
		drv.m.Lock()
		drv.on_scan[key] = dat
		drv.m.Unlock()
		return dat.value, nil
	}
	dat.lastreq = time.Now()
	return dat.value, nil
}

func (drv *LogixDriver) Write(key string, value any) error {
	return drv.client.Write(key, value)
}

func (drv *LogixDriver) Start() error {
	drv.status = "starting"
	ctx, cancel := context.WithCancel(context.Background())
	drv.close_signal = cancel
	go drv.run(ctx)
	return nil
}

func (drv *LogixDriver) Stop() error {
	drv.status = "stopping"
	drv.close_signal()
	return nil
}

func (drv *LogixDriver) Status() string {
	return drv.status
}

func (drv *LogixDriver) Name() string {
	return drv.DriverName
}

func (drv *LogixDriver) run(ctx context.Context) {
	drv.client.Connect()
	if drv.client.Connected {
		drv.status = "running. connected."
	} else {
		drv.status = "running. no connection."
	}
	defer func() {
		drv.client.Disconnect()
		drv.status = "stopped"
	}()

	t := time.NewTicker(drv.Rate)

	for {
		select {
		case <-t.C:
			// update all on-scan tags and purge old ones
			drv.m.Lock()
			tags := make([]string, 0, len(drv.on_scan))
			types := make([]gologix.CIPType, 0, len(drv.on_scan))
			elements := make([]int, 0, len(drv.on_scan))
			for k := range drv.on_scan {
				if time.Since(drv.on_scan[k].lastreq) > time.Minute {
					delete(drv.on_scan, k)
					continue
				}
				tags = append(tags, k)
				types = append(types, gologix.CIPTypeUnknown)
				elements = append(elements, 1)
			}
			vals, err := drv.client.ReadList(tags, types, elements)
			if err != nil {
				log.Printf("error reading tags in driver %s: %v", drv.DriverName, err)
				drv.status = fmt.Sprintf("error: %v", err)
				continue
			}
			for i := range tags {
				drv.on_scan[tags[i]].value = vals[i]
			}

			drv.m.Unlock()
		case <-ctx.Done():
			return

		}
	}
}

func (drv *LogixDriver) Change(a *admin.Admin, new_data any) {
	n, ok := new_data.(*LogixDriver)
	if !ok {
		log.Printf("shoudl have gotten a *LogixDriver???")
		return
	}
	drv.client.Disconnect()
	drv.client.IPAddress = n.IP
	p, err := gologix.ParsePath(n.Path)
	if err != nil {
		log.Panicf("bad path: %v", err)
	}
	drv.client.Path = p
	drv.client.Connect()
	drv.IP = n.IP
	drv.Path = n.Path
	drv.Rate = n.Rate
	if drv.DriverName != n.DriverName {
		a.UnRegisterStruct(drv.DriverName)
		drv.DriverName = n.DriverName
		a.RegisterStruct(drv.DriverName, drv)
	}
}
