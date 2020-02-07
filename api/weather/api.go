package weather

import (
	"fmt"
	"os"

	"github.com/pisign/pisign-backend/api"
)

// API is the function that gets called to get the data from the weather API
// This should be the only function that gets called from an external service
func API(data chan string) {
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
