package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/utils"
)

// Data is the function that gets called to get the data from the weather API
// This should be the only function that gets called from an external service
func (a *API) Data() string {
	var openWeatherResponse OpenWeatherResponse
	if a.APIKey == "" {
		fmt.Println("no key found")
		return ""
	}

	fmt.Println("Getting weather with args: ", a)
	apikey := a.APIKey
	zipcode := a.Zip

	if apikey == "" {
		// TODO better error handling
		fmt.Fprintf(os.Stderr, "No API key found for weather API")
		panic("no api key found")
	}

	url := buildurl(zipcode, apikey)
	resp := utils.GetAPIData(url)
	defer resp.Body.Close()

	body := utils.ParseResponse(resp)

	fmt.Println("Weather Response: ", string(body))

	utils.ParseJSON(body, openWeatherResponse)
	fmt.Printf("Weather returned: %+v", openWeatherResponse)
	var internalWeatherResponse api.InternalAPI
	internalWeatherResponse = openWeatherResponse.Transform()
	return string(internalWeatherResponse.Serialize())
}

// API yay
type API struct {
	api.BaseAPI
	Zip    int    `json:"zip"`
	APIKey string `json:"apiKey"`
}

// NewAPI creates a new weather api for a client
func NewAPI() *API {
	a := new(API)
	a.APIName = "weather"
	return a
}

// UnmarshalJSON for weather
func (a *API) UnmarshalJSON(b []byte) error {
	a.APIName = "WEATHER"
	return nil
}

// Configure for weather
func (a *API) Configure(j []byte) {
	fmt.Println("Configuring WEATHER!")
	err := json.Unmarshal(j, &a)
	if err != nil {
		log.Println("Error configuring weather api:", err)
		return
	}
}

// Run main entry point to weather API
func (a *API) Run(w api.Widget) {
	fmt.Println("Running WEATHER")
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		fmt.Println("STOPPING WEATHER")
	}()
	for {
		select {
		case <-w.Close():
			return
		default:
			<-ticker.C
			w.Send(a.Data())
		}
	}
}
