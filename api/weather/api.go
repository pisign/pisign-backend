package weather

import (
	"fmt"
	"os"
	"time"

	"github.com/pisign/pisign-backend/api"
)

// APIFunc is the function that gets called to get the data from the weather API
// This should be the only function that gets called from an external service
func APIFunc() string {
	var openWeatherResponse OpenWeatherResponse
	zipcode := 98105
	apikey := os.Getenv("WEATHER_KEY")

	if apikey == "" {
		fmt.Println("no key found")
		return ""
	}

	args := Args{Zipcode: zipcode, Apikey: apikey}
	openWeatherResponse.Get(args)
	var internalWeatherResponse api.InternalAPI
	internalWeatherResponse = openWeatherResponse.Transform()
	return string(internalWeatherResponse.Serialize())
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
func (a *API) Configure(j []byte) {
	fmt.Println("Configuring WEATHER!")
}

// Run main entry point to weather API
func (a *API) Run(w api.Widget) {
	fmt.Println("Running WEATHER")
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
			<-ticker.C
			w.Send(APIFunc())
		}
	}
}
