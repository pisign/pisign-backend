package clock

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/api"
)

// API yay
type API struct {
	api.BaseAPI
}

// NewAPI creates a new clock api for a client
func NewAPI() *API {
	a := new(API)
	a.APIName = "clock"
	return a
}

// Configure for clock
func (a *API) Configure(json string) {
	fmt.Println("Configuring CLOCK!")
}

// Run main entry point to clock API
func (a *API) Run(conn *websocket.Conn) {
	fmt.Println("Running CLOCK")
	for {
		t := time.Now()
		conn.WriteJSON(t)
		time.Sleep(1 * time.Second)
	}
}
