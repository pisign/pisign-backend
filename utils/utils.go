package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	fmt.Println("type " + t.Name() + " {\n" + fieldFmt + "}\n")
}
