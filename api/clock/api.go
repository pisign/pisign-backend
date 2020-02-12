// Package clock retrieves time data from the server and forwards it to the client
package clock

import (
	"encoding/json"
	"log"
	"time"

	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/types"
)

// API for clock
type API struct {
	api.BaseAPI
	Location string
	Format   string
}

// NewAPI creates a new clock api
func NewAPI() *API {
	a := new(API)
	a.APIName = "clock"
	a.Location = "Local"
	return a
}

func (a *API) loc() *time.Location {
	t, err := time.LoadLocation(a.Location)
	if err != nil {
		t, _ = time.LoadLocation("Local")
	}
	return t
}

// Configure for clock
func (a *API) Configure(body *json.RawMessage) {
	log.Println("Configuring CLOCK!")

	var config types.ClockConfig
	if err := json.Unmarshal(*body, &config); err == nil {
		_, err = time.LoadLocation(config.Location)
		if err != nil {
			log.Printf("Could not load timezone %s: %s", config.Location, err)
			return
		}
		a.Location = config.Location
		log.Println("Clock configuration successful!", "a = ", a)
	}

}

// Run main entry point to clock API
func (a *API) Run(w api.Widget) {
	log.Println("Running CLOCK")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		log.Println("STOPPING CLOCK")
	}()
	for {
		select {
		case <-w.Close():
			return
		default:
			t := <-ticker.C
			t = t.In(a.loc())
			out := types.ClockOut{Time: t.String()}
			w.Send(out)
		}
	}
}
