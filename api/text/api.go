// Package text for displaying a static text message
package text

import (
	"errors"
	"log"

	"github.com/pisign/pisign-backend/utils"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/types"
)

type API struct {
	types.BaseAPI
	Config types.TextConfig

	// Add other fields as necessary, (keep lowercase to avoid being stored in json)
}

// NewAPI creates a new API
func NewAPI(sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	a.BaseAPI.Init("text", sockets, pool, id)

	a.Config.Text = ""
	a.Config.Title = ""

	return a
}

// Configure based on client message
func (a *API) Configure(message types.ClientMessage) error {
	// Make sure the client widget is immediately sent new data with new config options
	defer func() {
		if a.Pool != nil && a.Sockets != nil {
			a.Pool.Save()
			a.Send(a.Data())
		}
	}()

	if err := a.BaseAPI.Configure(message); err != nil {
		return err
	}

	switch message.Action {
	case types.ConfigureAPI, types.Initialize:
		log.Println("Configuring TEXT:", message)
		oldConfig := a.Config
		if err := utils.ParseJSON(message.Config, &a.Config); err != nil {
			log.Println("Could not properly configure TEXT")
			a.Config = oldConfig
			return errors.New("could not properly configure TEXT")
		}

		// Add custom checks for config fields here (see the `time` api as an example)

		log.Println("TEXT configuration successful:", a)
	}

	return nil
}

// Data performs a bulk of the computational logic
func (a *API) Data() (interface{}, error) {
	return types.TextResponse{TextConfig: a.Config}, nil
}

// Run main entry point to the API
func (a *API) Run() {
	// Start up the BaseAPI to handle core API stuff
	go a.BaseAPI.Run()

	log.Println("Running TEXT")

	// Send data to client (using default config values)
	a.Send(a.Data())

	defer func() {
		log.Printf("STOPPING TEXT: %s\n", a.UUID)
	}()

	// Create a new channel to recieve termination messages on
	stop := a.AddStopChan()
	for {
		select {
		case body := <-a.ConfigChan: // Configuration messages
			if err := a.Configure(body); err != nil {
				a.Send(nil, err)
			}
		case <-stop: // Terminate
			return
		}
	}
}
