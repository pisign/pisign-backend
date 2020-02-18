# API Structure


This folder contains the apis that the frontend hits for info. The basic job of an API is to gather data from
an external source and return a parsed version for the frontend to later use. 

The types and apis that are avaliable to the frontend are summarized in the `./spec` folder. (still WIP)

Below are some files showing the barebones structure of an api called `NAME`

**api/NAME/api.go** contains most of the core logic
```go
package NAME

import (
    "github.com/pisign/pisign-backend/types"
    // import other necessary packages
)

// API struct to hold proper data types
type API struct {
	types.BaseAPI
	Config types.NAMEConfig
    ResponseObject Response

    // Any other necessary fields
}

// NewAPI creates a new API
func NewAPI(configChan chan types.ClientMessage, pool types.Pool) *API {
	a := new(API)
	a.BaseAPI.Init("<NAME>", configChan, pool)
	if a.Pool != nil {
		a.Pool.Register(a)
	}
	
    // Default config values, eg:
    // a.Location = "Local"
	return a
}

// Configure from json message sent from client
func (a *API) Configure(body types.ClientMessage) {
    // Also call parent Configure first
	a.ConfigurePosition(body.Position)

	if len(body.Config) == 0 {
		return
	}

    // Catch error if object can not be configured properly
	if err := json.Unmarshal(body.Config, &a.Config); err != nil {
		log.Println("Could not properly configure clock")
		a.Config = oldConfig
		return
	}
    // Add other checks or logic as necessary	
    
    // Update the connected pool object
	a.Pool.Save()
}

// Data performs the bulk work of retrieving the data from an external source
func (a *API) Data() interface{} {
    // Main logic here

	return Response{/*Fill in proper fields*/}
}

// Run main entry point to API
func (a *API) Run(w types.Socket) {
    log.Println("Running <NAME>")
    
    // Creates a timer that goes off every second
    ticker := time.NewTicker(1 * time.Second)

	defer func() {
		log.Println("STOPPING <NAME>")
        ticker.Stop()
        // Any additional cleanup work goes here
	}()
	for {
		select {
		case body := <-a.ConfigChan: // Configuration message received
			a.Configure(body)
		case save := <-w.Close(): // Socket connection was closed, stop running
			a.Pool.Unregister(types.Unregister{
				API:  a,
				Save: save,
			})			
            return
		case <- ticker.C:
			// Do any necessary pre-processing logic
			w.Send(a.Data()) // Send data to client through websocket
		}
	}
}
```

**api/NAME/types.go** holds all internal types not needed by the client
```go
package NAME

// Internal type definitions go here
// All types should all be unexported (lowercase)
```

**api/NAME/utils.go** holds utility functions particular to that API.
If you think a function can be made generic, consider putting it in the higher level `utils/` package.

```go
package NAME

// Utility functions go here
// All functions should be unexported (lowercase)
```

**types/NAME.go** holds all the public interfaces that need to be accesible both by the api and the client
```go
package types

// All types should be exported (capital first letter)

// NAMEResponse holds all of the data being sent to the client
type NAMEResponse struct {
	// Add any necessary fields
}

// NAMEConfig holds all necessary configuration parameters being sent from the client
type NAMEConfig struct {
	// Add any necessary fields
}

```
