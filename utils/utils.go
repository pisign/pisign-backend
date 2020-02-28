package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

/*
 * If these need to go into individual packages, then that is also fine,
 * but right now they are generic
 */

// GetAPIData wraps the get request
func GetAPIData(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	// TODO more elegant error handling
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ParseResponse wraps ioutil.Readall
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

// WrapError wraps an error
func WrapError(e error) {
	if e != nil {
		log.Println(e.Error())
	}
}

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

// CreateDirectory checks if a directory exists and if not creates it
// taken from https://www.socketloop.com/tutorials/golang-check-if-directory-exist-and-create-if-does-not-exist
func CreateDirectory(dirName string) error {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			return err
		}
		log.Printf("Created directory: '%s'\n", dirName)
		return nil
	}

	if err != nil {
		return err
	}

	if src.Mode().IsRegular() {
		return fmt.Errorf("CreateDirectory: '%s' already exists as a file, not a directory", dirName)
	}

	return nil
}
