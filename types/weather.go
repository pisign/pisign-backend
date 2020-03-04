package types

// weather API response data types, should implement api.InternalAPI

// Coord contains the lat and long
type Coord struct {
	Lon float64
	Lat float64
}

// Weather contains the current weather
type Weather struct {
	ID          float64
	Description string
	Main        string
	Icon        string
}

// Main contains a quick description of weather
type Main struct {
	Temp      float64
	FeelsLike float64
	TempMin   float64
	TempMax   float64
	Pressure  float64
	Humidity  float64
}

// Wind contains wind information
type Wind struct {
	Speed float64
	Deg   float64
}

// Clouds contains cloud info
type Clouds struct {
	All float64
}

// Rain contains rain info
type Rain struct {
	OneHR float64 `json:"1h"`
}

// Sys contains info about the system?
type Sys struct {
	Type    float64
	ID      float64
	Country string
	Sunrise float64
	Sunset  float64
}

// WeatherResponse is the struct that encodes the API data from our weather API
type WeatherResponse struct {
	WeatherConfig
	Coord      Coord
	Weather    []Weather
	Base       string
	Main       Main
	Visibility float64
	Wind       Wind
	Rain       Rain
	Clouds     Clouds
	DT         float64
	Sys        Sys
	Timezone   float64
	ID         float64
	Name       string
	COD        float64
	Zipcode    float64
}

// WeatherConfig are the config settings for the API
type WeatherConfig struct {
	Zip    int
	APIKey string
}
