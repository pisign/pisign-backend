// package utils provides utility functions to use elsewhere
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

// GetAPIData wraps the get request for a given url string
func GetAPIData(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	// TODO more elegant error handling
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ParseResponse transforms an entire http.Response object into raw bytes
func ParseResponse(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	// TODO more elegant error handling
	if err != nil {
		return nil, err
	}
	return body, nil
}

// ParseJSON parses byte slice body into struct i
func ParseJSON(body []byte, i interface{}) error {
	err := json.Unmarshal(body, i)
	return err
}

// WrapError wraps an error and prints it to the log
func WrapError(e error) {
	if e != nil {
		log.Println(e.Error())
	}
}

// DeleteEmpty deletes empty elements from a string slice, and returns the modified slice
//https://dabase.com/e/15006/
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// StringInSlice searches for a given string within a string slice, and returns if it is found or not
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
