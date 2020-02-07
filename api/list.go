package api

import (
	"fmt"
)

// InvalidAPIError error for missing API
type InvalidAPIError struct {
	APIName string
}

func (e InvalidAPIError) Error() string {
	return fmt.Sprintf("Invalid API: %s", e.APIName)
}

// Connect connects a client to a new API
func Connect(name string) (API, error) {
	switch name {
	case "weather":
		return NewWeatherAPI(), nil
	case "clock":
		return NewClockAPI(), nil
	default:
		return nil, InvalidAPIError{name}
	}
}

type baseAPI struct {
	APIName string `json:"apiName"`
}

func (a *baseAPI) Name() string {
	return a.APIName
}

// WeatherAPI yay
type WeatherAPI struct {
	*baseAPI
	zip    string
	apiKey string
}

// NewWeatherAPI creates a new weather api for a client
func NewWeatherAPI() *WeatherAPI {
	a := new(WeatherAPI)
	a.APIName = "weather"
	return a
}

// Configure for weather
func (a *WeatherAPI) Configure(json string) {
	fmt.Println("Configuring WEATHER!")
}

// ClockAPI yay
type ClockAPI struct {
	baseAPI
}

// NewClockAPI creates a new clock api for a client
func NewClockAPI() *ClockAPI {
	a := new(ClockAPI)
	a.APIName = "clock"
	return a
}

// Configure for clock
func (a *ClockAPI) Configure(json string) {
	fmt.Println("Configuring CLOCK!")
}
