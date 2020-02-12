package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

func (a *API) get() api.ExternalAPI {
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
	var openWeatherResponse OpenWeatherResponse
	utils.ParseJSON(body, &openWeatherResponse)
	return &openWeatherResponse
}

func (a *API) data() api.InternalAPI {
	if time.Now().Sub(a.LastCalled) < (time.Second * 10) {
		log.Println("using cached value")
		return &a.CachedResponse
	}

	response := a.get().Transform()
	res := response.(*types.WeatherResponse)

	if res.COD > 300 {
		log.Println("API error")
		log.Println(res)
	} else {
		a.LastCalled = time.Now()
		a.CachedResponse = *res
	}

	return res
}

// API yay
type API struct {
	api.BaseAPI
	Zip        int
	APIKey     string
	LastCalled time.Time
	// TODO get rid of the cached response on the API struct?
	CachedResponse types.WeatherResponse
}

// NewAPI creates a new weather api for a client
func NewAPI() *API {
	a := new(API)
	a.APIName = "weather"
	return a
}

// Configure for weather
func (a *API) Configure(body *json.RawMessage) {
	log.Println("Configuring WEATHER!")
	err := json.Unmarshal(*body, a)
	if err != nil {
		log.Println("Error configuring weather api:", err)
		return
	}
}

// Run main entry point to weather API
func (a *API) Run(w api.Widget) {
	log.Println("Running WEATHER")
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		log.Println("STOPPING WEATHER")
	}()
	for {
		select {
		case <-w.Close():
			return
		default:
			<-ticker.C
			w.Send(a.data())
		}
	}
}
