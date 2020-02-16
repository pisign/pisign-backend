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
	Config types.ClockConfig
	Time   time.Time `json:"-"`
}

// NewAPI creates a new clock api
func NewAPI(configChan chan types.ConfigMessage, pool types.Pool) *API {
	a := new(API)
	a.BaseAPI.Init("clock", configChan, pool)
	if a.Pool != nil {
		a.Pool.Register(a)
	}
	a.Config.Location = "Local"
	return a
}

func (a *API) loc() *time.Location {
	t, err := time.LoadLocation(a.Config.Location)
	if err != nil {
		t, _ = time.LoadLocation("Local")
	}
	return t
}

// Configure for clock
func (a *API) Configure(body types.ConfigMessage) {
	a.ConfigurePosition(body.Position)
	log.Println("Configuring CLOCK!")
	log.Printf("Config: %+v\n", body)

	if err := json.Unmarshal(body.Config, &a.Config); err == nil {
		_, err = time.LoadLocation(a.Config.Location)
		if err != nil {
			log.Printf("Could not load timezone %s: %s", a.Config.Location, err)
			return
		}
		log.Println("Clock configuration successful!", "a = ", a)
	}
	a.Pool.Save()
}

// Data gets the current time!
func (a *API) Data() interface{} {
	return types.ClockOut{Time: a.Time.In(a.loc()).String()}
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
