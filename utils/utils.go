package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*
 * If these need to go into individual packages, then that is also fine,
 * but right now they are generic
 */

// GetAPIData wraps the get request
func GetAPIData(url string) *http.Response {
	resp, err := http.Get(url)
	// TODO more elegant error handling
	if err != nil {
		panic("error in GET request")
	}
	return resp
}

// ParseResponse wraps ioutil.Readall
func ParseResponse(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	// TODO more elegant error handling
	if err != nil {
		panic("error in reading API response")
	}
	return body
}

// ParseJSON parses byte slice body into struct i
func ParseJSON(body []byte, i interface{}) {
	json.Unmarshal(body, i)
}
