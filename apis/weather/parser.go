package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Parse the openweathermap.org API to get weather data
func Parse(zipcode int, apikey string) OpenWeatherResponse {
	if apikey == "" {
		// TODO error handling could be here, or not, idk
		fmt.Fprintf(os.Stderr, "No API key found for weather API")
		panic("no api key found")
	}

	url := buildurl(zipcode, apikey)
	resp := getAPIData(url)
	defer resp.Body.Close()

	body := parseResponse(resp)

	return parseJSON(body)
}

const endpoint = "https://api.openweathermap.org/data/2.5/weather?zip="

// Builds the URL from the endpoint, adding required args
func buildurl(zipcode int, apikey string) string {
	return fmt.Sprintf("%s%d,us&APPID=%s", endpoint, zipcode, apikey)
}

// Sends a get request to the url and returns the response object
// TODO consider taking this out into a util function
func getAPIData(url string) *http.Response {
	resp, err := http.Get(url)
	// TODO more elegant error handling
	if err != nil {
		panic("error in GET request to weather api")
	}
	return resp
}

// Parse the http response object and return the body
// TODO consider taking this out into a util function
func parseResponse(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	// TODO more elegant error handling
	if err != nil {
		panic("error in reading weather API response")
	}
	return body
}

// Parses out the JSON from the byte array and returns it as an OpenWeatherResponse
// TOCO consider taking this out into a util function (assuming you can pass a type to it)
func parseJSON(body []byte) OpenWeatherResponse {
	var openWeatherResponse OpenWeatherResponse
	json.Unmarshal(body, &openWeatherResponse)
	return openWeatherResponse
}
