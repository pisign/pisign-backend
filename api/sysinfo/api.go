package sysinfo

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/types"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type API struct {
	types.BaseAPI
	Config         *types.SysInfoConfig
	ResponseObject types.SysInfoResponse
	LastCalled     time.Time `json:"-"`
	ValidCache     bool
	// Add other fields as necessary, (keep lowercase to avoid being stored in json)
}

// NewAPI creates a new API
func NewAPI(sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	a.BaseAPI.Init(types.APISysinfo, sockets, pool, id)
	a.ValidCache = false
	return a
}

// Configure based on client message
func (a *API) Configure(message types.ClientMessage) error {
	// TODO this requires no configuration
	// Make sure the client widget is immediately sent new data with new config options
	a.Config = &types.SysInfoConfig{}
	defer func() {
		if a.Pool != nil && a.Sockets != nil {
			a.Pool.Save()
			a.Send(a.Data())
		}
	}()

	if err := a.BaseAPI.Configure(message); err != nil {
		return err
	}

	switch message.Action {
	case types.ConfigureAPI, types.Initialize:
		log.Printf("Configuring %s: %v\n", a.Name, message)
		log.Printf("%s configuration successful: %s", a, a)
	}

	return nil
}

// Data performs a bulk of the computational logic
func (a *API) Data() (interface{}, error) {
	// Perform logic here (including call to external API)
	// 5 second cache
	if time.Now().Sub(a.LastCalled) < (time.Second*5) && a.ValidCache {
		// Send the old response object
		log.Println("using cached value")
		return a.ResponseObject, nil
	}

	v, err := mem.VirtualMemory()

	if err != nil {
		return nil, errors.New("Error getting memory info")
	}

	d, err := disk.Usage("/")

	if err != nil {
		return nil, errors.New("Error getting disk info")
	}

	a.ResponseObject = types.SysInfoResponse{
		MemTotal:        v.Total,
		MemUsed:         v.Used,
		MemUsedPercent:  v.UsedPercent,
		DiskUsedPercent: d.UsedPercent,
		DiskTotal:       d.Total,
		DiskUsed:        d.Used,
		DiskFree:        d.Free,
	}

	a.ValidCache = true

	return a.ResponseObject, nil
}

// Run main entry point to the API
func (a *API) Run() {
	// Start up the BaseAPI to handle core API stuff
	go a.BaseAPI.Run()

	log.Printf("Running %s\n", a)

	a.Send(a.Data())

	// check every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING %s\n", a)
	}()

	// Create a new channel to receive termination messages on
	stop := a.AddStopChan()
	for {
		select {
		case body := <-a.ConfigChan: // Configuration messages
			if err := a.Configure(body); err != nil {
				a.Send(nil, err)
			}
		case <-ticker.C: // Update timer tick
			a.Send(a.Data())
		case <-stop: // Terminate
			return
		}
	}
}
