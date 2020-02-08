package weather

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/api"
)

// APIFunc is the function that gets called to get the data from the weather API
// This should be the only function that gets called from an external service
func APIFunc(data chan string) {
	var openWeatherResponse OpenWeatherResponse
	zipcode := 98105
	apikey := os.Getenv("WEATHER_KEY")

	if apikey == "" {
		fmt.Println("no key found")
		data <- ""
		return
	}

	args := Args{Zipcode: zipcode, Apikey: apikey}
	openWeatherResponse.Get(args)
	var internalWeatherResponse api.InternalAPI
	internalWeatherResponse = openWeatherResponse.Transform()
	data <- string(internalWeatherResponse.Serialize())
}

// API yay
type API struct {
	api.BaseAPI
	zip    string
	apiKey string
}

// NewAPI creates a new weather api for a client
func NewAPI() *API {
	a := new(API)
	a.APIName = "weather"
	return a
}

// Configure for weather
func (a *API) Configure(json string) {
	fmt.Println("Configuring WEATHER!")
}

// Run main entry point to weather API
func (a *API) Run(conn *websocket.Conn) {
	fmt.Println("Running WEATHER")
}
