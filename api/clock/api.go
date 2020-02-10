package clock

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pisign/pisign-backend/api"
)

// API yay
type API struct {
	api.BaseAPI
	Location *time.Location
	Format   string
}

// NewAPI creates a new clock api for a client
func NewAPI() *API {
	a := new(API)
	a.APIName = "clock"
	a.Location, _ = time.LoadLocation("Local")
	return a
}

type configurationArgs struct {
	Location string
}

// Configure for clock
func (a *API) Configure(j []byte) {
	fmt.Println("Configuring CLOCK!", j)
	var config configurationArgs
	err := json.Unmarshal(j, &config)
	if err != nil {
		log.Println("Error configuring Clock api:", err)
		return
	}
	loc, err := time.LoadLocation(config.Location)
	if err != nil {
		log.Printf("Could not load timezone %s: %s", config.Location, err)
		return
	}
	a.Location = loc
	log.Println("Clock configuration successful!", "a.Location = ", a.Location.String())
}

// Run main entry point to clock API
func (a *API) Run(w api.Widget) {
	fmt.Println("Running CLOCK")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		fmt.Println("STOPPING CLOCK")
	}()
	for {
		select {
		case <-w.Close():
			return
		default:
			t := <-ticker.C
			t = t.In(a.Location)
			w.Send(t)
		}
	}
}
