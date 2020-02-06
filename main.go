package main

import (
	"fmt"

	"github.com/pisign/pisign-backend/api/weather"
)

func main() {
	datachan := make(chan string)
	go weather.API(datachan)
	output := <-datachan
	fmt.Println(output)
}
