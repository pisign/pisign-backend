package weather

import (
	"encoding/json"
	"log"
	"time"

	"github.com/pisign/pisign-backend/types"
)

// API for weather
type API struct {
	types.BaseAPI
	APISettings APISettings
	LastCalled  time.Time
	// This is the object we get from the backend API - we could possible remove this and just have the ResponseObject
	DataObject OpenWeatherResponse
	// This is the object we are passing to the frontend - only need to rebuild it when its stale
	ResponseObject types.WeatherResponse
}

// APISettings are the config settings for the API
type APISettings struct {
	Zip    int
	APIKey string
}

// Data gets the data to send to the websocket
func (a *API) Data() interface{} {
	// Data should handle how to retrieve the data
	if time.Now().Sub(a.LastCalled) < (time.Second * 10) {
		// Send the old response object
		log.Println("using cached value")
		return &a.ResponseObject
	}

	// Otherwise, update the response object
	a.DataObject.Update(a.APISettings)
	response := a.DataObject.Transform()
	a.ResponseObject = *(response.(*types.WeatherResponse))
	a.LastCalled = time.Now()
	return &a.ResponseObject
}

// NewAPI creates a new weather api for a client
func NewAPI() *API {
	a := new(API)
	a.APIName = "weather"
	return a
}

// Configure for weather
func (a *API) Configure(body *json.RawMessage) {
	log.Println("Configuring WEATHER!")
	err := json.Unmarshal(*body, &a.APISettings)
	if err != nil {
		log.Println("Error configuring weather api:", err)
		return
	}
}

// Run main entry point to weather API
func (a *API) Run(w api.Socket) {
	log.Println("Running WEATHER")
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		log.Println("STOPPING WEATHER")
	}()
	for {
		select {
		case <-w.Close():
			return
		default:
			<-ticker.C
			w.Send(a.Data())
		}
	}
}
