// Package clock retrieves time data from the server and forwards it to the client
package clock

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/pisign/pisign-backend/types"
)

// API for clock
type API struct {
	types.BaseAPI
	Location string    `json:"-"`
	Time     time.Time `json:"-"`
}

// NewAPI creates a new clock api
func NewAPI(configChan chan *json.RawMessage, pool types.Pool) *API {
	a := new(API)
	a.BaseAPI.Init("clock", configChan, pool)
	if a.Pool != nil {
		a.Pool.Register(a)
	}

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
	a.ConfigurePosition(body)
	log.Println("Configuring CLOCK!")

	var config types.ClockConfig
	if err := json.Unmarshal(*body, &config); err == nil {
		v := reflect.ValueOf(config)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			log.Printf("Field: %v, value: %v\n", t.Field(i).Name, v.Field(i))
		}
		_, err = time.LoadLocation(config.Location)
		if err != nil {
			log.Printf("Could not load timezone %s: %s", config.Location, err)
			return
		}
		a.Location = config.Location
		a.Config["Location"] = config.Location
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
