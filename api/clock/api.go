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
	Location string
	Format   string
	Time     time.Time
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

// Data gets the current time!
func (a *API) Data() interface{} {
	return types.ClockOut{Time: a.Time.In(a.loc()).String()}
}

// Run main entry point to clock API
func (a *API) Run() {
	log.Println("Running CLOCK")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		log.Println("STOPPING CLOCK")
	}()
	for {
		select {
		case <-a.Close:
			return
		default:
			t := <-ticker.C
			t = t.In(a.loc())
			out := types.ClockOut{Time: t.String()}
			a.Send(out)
		}
	}
}
