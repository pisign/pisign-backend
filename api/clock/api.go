package clock

import (
	"fmt"
	"time"

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
func (a *API) Run(w api.Widget) {
	fmt.Println("Running CLOCK")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		fmt.Println("STOPPING CLOCK")
	}()
	for {
		w.Send()
	}
}
