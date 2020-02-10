package weather

import (
	"fmt"
	"log"
	"os"

	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

// Args for the Get method
type Args struct {
	Zipcode int
	Apikey  string
}

// Get hits the openweathermap.org API to get weather data
func (o *OpenWeatherResponse) Get(a interface{}) {
	args, ok := a.(API)
	log.Println("Getting weather with args: ", args)
	apikey := args.APIKey
	zipcode := args.Zip

	if !ok {
		// TODO better error handling
		panic("error in parsing arg struct, most likely called with bad type")
	}

	if apikey == "" {
		// TODO better error handling
		fmt.Fprintf(os.Stderr, "No API key found for weather API")
		panic("no api key found")
	}

	url := buildurl(zipcode, apikey)
	resp := utils.GetAPIData(url)
	defer resp.Body.Close()

	body := utils.ParseResponse(resp)

	utils.ParseJSON(body, o)
}

// Transform method converts to the InternalAPI type
// All business logic for converting from the external API to the internal one should be here
func (o *OpenWeatherResponse) Transform() api.InternalAPI {
	weatherResponse := types.WeatherResponse{
		Name: o.Name,
	}

	return &weatherResponse
}
