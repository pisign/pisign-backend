package utils

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_GetAPIData(t *testing.T) {
	// Test normal execution
	normalurl := "https://api.openweathermap.org/data/2.5/weather?zip=98105,us&APPID=fec97c5a56a5d2966ad1c16f98b9e19a"

	_, err := GetAPIData(normalurl)
	if err != nil {
		t.Error("Returned error for valid URL get request")
	}

	_, err = GetAPIData("bad url ?")
	if err == nil {
		t.Error("Failed to return error for invalid URL get request")
	}
}

func Test_CreateDirectory(t *testing.T) {
	fakedir := "./fake"
	if CreateDirectory(fakedir) != nil {
		t.Error("error creating fake directory")
	}

	if CreateDirectory(fakedir) != nil {
		t.Error("error when processing existing directory")
	}

	os.Remove(fakedir)
}

func TestWrapError(t *testing.T) {
	var str bytes.Buffer
	errorMessage := "error message"
	log.SetOutput(&str)
	WrapError(errors.New(errorMessage))

	if !strings.Contains(string(str.Bytes()), errorMessage) {
		t.Error("error message not logged!")
		t.Error(string(str.Bytes()))
	}

	log.SetOutput(os.Stderr)
}
