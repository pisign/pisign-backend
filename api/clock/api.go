// Package clock retrieves time data from the server and forwards it to the client
package clock

import (
	"encoding/json"
	"log"
	"time"

	"github.com/pisign/pisign-backend/types"
)

// API for clock
type API struct {
	types.BaseAPI
	config types.ClockConfig
	Time   time.Time `json:"-"`
}

// NewAPI creates a new clock api
func NewAPI(configChan chan types.ConfigMessage, pool types.Pool) *API {
	a := new(API)
	a.BaseAPI.Init("clock", configChan, pool)
	if a.Pool != nil {
		a.Pool.Register(a)
	}
	a.config.Location = "Local"
	return a
}

func (a *API) loc() *time.Location {
	t, err := time.LoadLocation(a.config.Location)
	if err != nil {
		t, _ = time.LoadLocation("Local")
	}
	return t
}

// Configure for clock
func (a *API) Configure(body types.ConfigMessage) {
	a.ConfigurePosition(body.Position)
	log.Println("Configuring CLOCK:", body)
	oldConfig := a.config

	if err := json.Unmarshal(body.Config, &a.config); err != nil {
		log.Println("Could not properly configure clock")
		a.config = oldConfig
		return
	}
	if _, err := time.LoadLocation(a.config.Location); err != nil {
		log.Printf("Could not load timezone %s: %s", a.config.Location, err)
		a.config.Location = oldConfig.Location // Revert to old location
	}
	log.Println("Clock configuration successful:", a)
	a.Pool.Save()
}

// Data gets the current time!
func (a *API) Data() interface{} {
	return types.ClockResponse{Time: a.Time.In(a.loc()).String()}
}

// Run main entry point to clock API
func (a *API) Run(w types.Socket) {
	log.Println("Running CLOCK")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		log.Println("STOPPING CLOCK")
	}()
	for {
		select {
		case body := <-a.ConfigChan:
			a.Configure(body)
		case <-w.Close():
			return
		default:
			t := <-ticker.C
			a.Time = t
			w.Send(a.Data())
		}
	}
}
