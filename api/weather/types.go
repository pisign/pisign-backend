package weather

import (
	"fmt"
	"os"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

// openweathermap.org response data types

type coord struct {
	Lon float64
	Lat float64
}

type weather struct {
	ID          float64
	Description string
	Main        string
	Icon        string
}

type main struct {
	Temp      float64
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  float64
	Humidity  float64
}

type wind struct {
	Speed float64
	Deg   float64
}

type clouds struct {
	All float64
}

type rain struct {
	OneHR float64 `json:"1h"`
}

type sys struct {
	Type    float64
	ID      float64
	Country string
	Sunrise float64
	Sunset  float64
}

// OpenWeatherResponse is the struct that encodes the API data from openweathermap.org
type OpenWeatherResponse struct {
	Coord      coord
	Weather    []weather
	Base       string
	Main       main
	Visibility float64
	Wind       wind
	Rain       rain
	Clouds     clouds
	DT         float64
	Sys        sys
	Timezone   float64
	ID         float64
	Name       string
	COD        float64
}

// Update builds the data object
func (o *OpenWeatherResponse) Update(arguments interface{}) {
	a := (arguments).(types.WeatherConfig)
	apikey := a.APIKey
	zipcode := a.Zip

	if apikey == "" {
		// TODO better error handling
		fmt.Fprintf(os.Stderr, "No API key found for weather API")
		//panic("no api key found")
		return
	}

	url := buildurl(zipcode, apikey)
	resp := utils.GetAPIData(url)
	defer resp.Body.Close()

	body := utils.ParseResponse(resp)
	utils.ParseJSON(body, &o)
}

// Transform turns the OpenWeatherResponse into a WeatherResponse
func (o *OpenWeatherResponse) Transform() interface{} {
	weatherResponse := types.WeatherResponse{
		Name: o.Name,
		Main: types.Main{
			Temp: kelvinToF(o.Main.Temp),
		},
	}

	return &weatherResponse
}
