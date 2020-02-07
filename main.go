package main

import (
	"fmt"

	"github.com/pisign/pisign-backend/api/weather"
	"github.com/pisign/pisign-backend/server"
)

func main() {
	fmt.Print(string(weather.API()))
	server.StartLocalServer(9000)
}
