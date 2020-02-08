package clock

import (
	"fmt"

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
