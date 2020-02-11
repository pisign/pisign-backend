package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
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

// ParseJSONMap parses json byte slice into a generic key-value map
func ParseJSONMap(body []byte) (map[string]*json.RawMessage, error) {
	var data map[string]*json.RawMessage
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// StructPrint prints out the structure of a Struct
func StructPrint(v interface{}) {

	t := reflect.Indirect(reflect.ValueOf(v)).Type()

	fieldFmt := ""

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// If Name() returns empty string, means we have something like a slice
		name := field.Type.Name()
		if name == "" {
			name = field.Type.String()
		}

		line := fmt.Sprintf("%-13v", field.Name) + name
		fieldFmt += "    " + line + "\n"
	}

	log.Println("type " + t.Name() + " {\n" + fieldFmt + "}\n")
}
