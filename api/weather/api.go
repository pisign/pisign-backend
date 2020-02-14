package weather

import (
	"encoding/json"
	"log"
	"time"

	"github.com/pisign/pisign-backend/types"
)

func (a *API) get() types.ExternalAPI {
	apikey := a.APIKey
	zipcode := a.Zip

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

// API yay
type API struct {
	types.BaseAPI
	Zip        int
	APIKey     string
	LastCalled time.Time
	// TODO get rid of the cached response on the API struct?
	CachedResponse types.WeatherResponse
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
func (a *API) Run() {
	log.Println("Running WEATHER")
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		log.Println("STOPPING WEATHER")
	}()
	for {
		select {
		case <-a.Widget.Close():
			return
		default:
			<-ticker.C
			a.Widget.Send(a.data())
		}
	}
}
