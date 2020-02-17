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
	Config     types.WeatherConfig
	LastCalled time.Time `json:"-"`
	// This is the object we get from the backend API - we could possible remove this and just have the ResponseObject
	DataObject OpenWeatherResponse `json:"-"`
	// This is the object we are passing to the frontend - only need to rebuild it when its stale
	ResponseObject types.WeatherResponse `json:"-"`
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
	a.DataObject.Update(a.Config)
	response := a.DataObject.Transform()
	a.ResponseObject = *(response.(*types.WeatherResponse))
	a.LastCalled = time.Now()
	return &a.ResponseObject
}

// NewAPI creates a new weather api for a client
func NewAPI(configChan chan types.ConfigMessage, pool types.Pool) *API {
	a := new(API)
	a.BaseAPI.Init("weather", configChan, pool)
	if a.Pool != nil {
		a.Pool.Register(a)
	}
	return a
}

// Configure for weather
func (a *API) Configure(body types.ConfigMessage) {
	a.ConfigurePosition(body.Position)
	log.Println("Configuring WEATHER:", body)

	if err := json.Unmarshal(body.Config, &a.Config); err != nil {
		log.Println("Could not properly configure weather")
		return
	}
	log.Println("Weather configuration successfully:", a)
}

// Run main entry point to weather API
func (a *API) Run(w types.Socket) {
	log.Println("Running WEATHER")
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		log.Println("STOPPING WEATHER")
		a.Pool.Unregister(a)
	}()
	for {
		select {
		case data := <-a.ConfigChan:
			a.Configure(data)
		case <-w.Close():
			return
		default:
			<-ticker.C
			w.Send(a.Data())
		}
	}
}
