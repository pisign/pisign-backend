package main

import (
	"fmt"

	"github.com/pisign/pisign-backend/api/weather"
)

func main() {
	fmt.Print(string(weather.API()))
}
