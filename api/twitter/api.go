package twitter

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

// API for twitter
type API struct {
	types.BaseAPI
	Config     types.TwitterConfig
	LastCalled time.Time `json:"-"`
	// This is the object we get from the backend API - we could possible remove this and just have the ResponseObject
	DataObject TwitterResponse `json:"-"`
	// This is the object we are passing to the frontend - only need to rebuild it when its stale
	ResponseObject types.TwitterResponse `json:"-"`
	ValidCache     bool
}

// NewAPI creates a new API
func NewAPI(sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	// Twitter is the name of the api visible to the client
	a.BaseAPI.Init(types.APITwitter, sockets, pool, id)
	a.Config = types.TwitterConfig{}

	// Configure default values as necessary, for example:
	// a.Config.Variable = "Value"

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
		log.Println("Configuring Twitter:", message)
		oldConfig := a.Config
		if err := utils.ParseJSON(message.Config, &a.Config); err != nil {
			log.Println("Could not properly configure Twitter")
			a.Config = oldConfig
			return errors.New("could not properly configure Twitter")
		}

		// Add custom checks for config fields here (see the `time` api as an example)

		log.Println("Twitter configuration successful:", a)
	}

	return nil
}

// Data performs a bulk of the computational logic
func (a *API) Data() (interface{}, error) {
	// Perform logic here (including call to external API)

	// get cached response if still valid
	if time.Now().Sub(a.LastCalled) < (time.Second*60) && a.ValidCache {
		// Send the old response object
		log.Println("using cached value")
		return &a.ResponseObject, nil
	}

	// make twitter api request
	err := a.DataObject.Update(a.Config)
	if err != nil {
		a.ValidCache = false
		return nil, err
	}

	response := a.DataObject.Transform()
	a.ResponseObject = *(response.(*types.TwitterResponse))
	a.ResponseObject.TwitterConfig = a.Config
	a.LastCalled = time.Now()
	a.ValidCache = true
	return &a.ResponseObject, nil
}

// Run main entry point to the API
func (a *API) Run() {
	// Start up the BaseAPI to handle core API stuff
	go a.BaseAPI.Run()

	log.Println("Running Twitter")

	// Send data to client (using default config values)
	a.Send(a.Data())

	// Set up a new ticker (if you want to send info the client at set time intervals)
	// You can change the time length as necessary
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING Twitter: %s\n", a.UUID)
	}()

	// Create a new channel to recieve termination messages on
	stop := a.AddStopChan()
	for {
		select {
		case body := <-a.ConfigChan: // Configuration messages
			if err := a.Configure(body); err != nil {
				a.Send(nil, err)
			}
		case <-ticker.C: // Update timer tick
			a.Send(a.Data())
		case <-stop: // Terminate
			return
		}
	}
}
