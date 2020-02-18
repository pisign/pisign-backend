package weather

import (
	"encoding/json"
	"errors"
	"log"
	"time"

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
func (a *API) Data() interface{} {
	// Data should handle how to retrieve the data
	if time.Now().Sub(a.LastCalled) < (time.Second*10) && a.ValidCache {
		// Send the old response object
		log.Println("using cached value")
		return &a.ResponseObject
	}

	// Otherwise, update the response object
	err := a.DataObject.Update(a.Config)

	// If there was an error updating the data object,
	// set response object to error'ed out and return it
	// TODO possibly delete the ResponseObject before doing this
	// so that the data is all 0'ed?
	if err != nil {
		a.ResponseObject.Status = types.StatusFailure
		a.ResponseObject.ErrorMessage = err.Error()
		a.ValidCache = false
		return &a.ResponseObject
	}

	a.ResponseObject.Status = types.StatusSuccess

	response := a.DataObject.Transform()
	a.ResponseObject = *(response.(*types.WeatherResponse))
	a.LastCalled = time.Now()
	a.ValidCache = true
	return &a.ResponseObject
}

// NewAPI creates a new weather api for a client
func NewAPI(configChan chan types.ClientMessage, pool types.Pool) *API {
	a := new(API)
	a.BaseAPI.Init("weather", configChan, pool)
	a.ValidCache = false
	return a
}

// Configure for weather
func (a *API) Configure(body types.ClientMessage) error {
	a.ConfigurePosition(body.Position)
	log.Println("Configuring WEATHER:", body)

	if len(body.Config) == 0 {
		return nil
	}
	if err := json.Unmarshal(body.Config, &a.Config); err != nil {
		return errors.New("could not properly configure weather")
	}
	log.Println("Weather configuration successfully:", a)
	return nil
}

// Run main entry point to weather API
func (a *API) Run(w types.Socket) {
	log.Println("Running WEATHER")
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING WEATHER: %s\n", a.UUID)
	}()
	for {
		select {
		case body := <-a.ConfigChan:
			if err := a.Configure(body); err != nil {
				w.SendErrorMessage(err.Error())
			} else {
				w.Send(a.Data())
			}
		case save := <-w.Close():
			a.Pool.Unregister(types.Unregister{
				API:  a,
				Save: save,
			})
			return
		case <-ticker.C:
			w.Send(a.Data())
		}
	}
}
