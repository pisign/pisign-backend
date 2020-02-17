package main

import (
	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

func main() {
	// Weather API
	utils.StructPrint(types.WeatherResponse{})
	utils.StructPrint(types.Coord{})
	utils.StructPrint(types.Main{})
	utils.StructPrint(types.Wind{})
	utils.StructPrint(types.Rain{})
	utils.StructPrint(types.Clouds{})
	utils.StructPrint(types.Sys{})
	utils.StructPrint(types.BaseMessage{})
	utils.StructPrint(types.ClockResponse{})
}
