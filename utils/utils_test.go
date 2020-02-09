package utils

import (
	"fmt"
	"testing"
)

func Test_GetAPIData(t *testing.T) {
	// Test normal execution
	normalurl := "https://api.openweathermap.org/data/2.5/weather?zip=98105,us&APPID=fec97c5a56a5d2966ad1c16f98b9e19a"

	GetAPIData(normalurl)
	// no panic = success

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		} else {
			t.Error("Failed to panic for invalid URL get request")
		}
	}()

	GetAPIData("bad url ?")
	// this should trigger panic, which is caught by the defer. If that is not called,
	// then test should fail!
}
