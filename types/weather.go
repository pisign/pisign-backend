package types

import (
	"encoding/json"
	"fmt"
)

// weather API response data types, should implement api.InternalAPI

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

// WeatherResponse is the struct that encodes the API data from our weather API
type WeatherResponse struct {
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

// Serialize serializes the WeatherResponse
func (w *WeatherResponse) Serialize() []byte {
	bytes, err := json.Marshal(w)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return bytes
}
