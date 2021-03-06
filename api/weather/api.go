package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/types"
)

// API for weather
type API struct {
	types.BaseAPI
	Config     types.WeatherConfig
	LastCalled time.Time `json:"-"`
	// This is the object we get from the backend API - we could possible remove this and just have the ResponseObject
	DataObject OpenWeatherResponse `json:"-"`
	// This is the object we are passing to the frontend - only need to rebuild it when its stale
	ResponseObject types.WeatherResponse `json:"-"`
	ValidCache     bool
}

// Data gets the data to send to the websocket
func (a *API) Data() (interface{}, error) {
	// Data should handle how to retrieve the data
	if time.Now().Sub(a.LastCalled) < (time.Second*10) && a.ValidCache {
		// Send the old response object
		log.Println("using cached value")
		return &a.ResponseObject, nil
	}

	// Otherwise, update the response object
	err := a.DataObject.Update(a.Config)

	// If there was an error updating the data object,
	// set response object to error'ed out and return it
	// TODO possibly delete the ResponseObject before doing this
	// so that the data is all 0'ed?
	if err != nil {
		a.ValidCache = false
		return nil, err
	}

	response := a.DataObject.Transform()
	a.ResponseObject = *(response.(*types.WeatherResponse))
	a.ResponseObject.WeatherConfig = a.Config
	a.LastCalled = time.Now()
	a.ValidCache = true
	return &a.ResponseObject, nil
}

// NewAPI creates a new weather api for a client
func NewAPI(sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	a.BaseAPI.Init(types.APIWeather, sockets, pool, id)
	a.ValidCache = false
	return a
}

// Configure for weather
func (a *API) Configure(message types.ClientMessage) error {
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
		log.Printf("Configuring %s: %+v", a, message)
		if err := json.Unmarshal(message.Config, &a.Config); err != nil {
			return errors.New(fmt.Sprintf("could not properly configure %s", a))
		}
		log.Printf("%s configuration successfully: %+v", a, a)
		a.ValidCache = false
	}
	return nil
}

// Run main entry point to weather API
func (a *API) Run() {
	log.Printf("Running %s\n", a)
	go a.BaseAPI.Run()
	a.Send(a.Data())
	ticker := time.NewTicker(10 * time.Second)
	stop := a.AddStopChan()
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING %s\n", a)
	}()
	for {
		select {
		case body := <-a.ConfigChan:
			err := a.Configure(body)
			if err != nil {
				a.Send(nil, err)
			}
		case <-ticker.C:
			a.Send(a.Data())
		case <-stop:
			return
		}
	}
}
