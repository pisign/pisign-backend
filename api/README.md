# API Structure


This folder contains the apis that the frontend hits for info. The basic job of an API is to gather data from
an external source and return a parsed version for the frontend to later use. 

The types and apis that are avaliable to the frontend are summarized in the `./spec` folder. (still WIP)

Below are some files showing the barebones structure of an api called `<NAME>`

<<<<<<< HEAD
**api/<NAME>/api.go** contains most of the core logic

```text
package <NAME>
=======
**api/NAME/api.go** contains most of the core logic

```text
package NAME
>>>>>>> bc5cfe521ea2d81d12abe8690931977fe40fd83e

import ...

type API struct {
	types.BaseAPI
	Config types.<NAME>Config
    // Add other fields as necessary, (keep lowercase to avoid being stored in json)
}

// NewAPI creates a new API
func NewAPI(sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	// <NAME> is the name of the api visible to the client
	a.BaseAPI.Init("<NAME>", sockets, pool, id)

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
		log.Println("Configuring <NAME>:", message)
		oldConfig := a.Config
		if err := utils.ParseJSON(message.Config, &a.Config); err != nil {
			log.Println("Could not properly configure <NAME>")
			a.Config = oldConfig
			return errors.New("could not properly configure <NAME>")
		}

		// Add custom checks for config fields here (see the `time` api as an example)

		log.Println("<NAME> configuration successful:", a)
	}

	return nil
}

// Data performs a bulk of the computational logic
func (a *API) Data() (interface{}, error) {
	// Perform logic here (including call to external API)
    
    // Successful:
    // return types.<NAME>Response{/* Fields go here */}, nil

    // Error:
    // return nil, errors.New("ERROR GETTING DATA FOR <NAME>")
}

// Run main entry point to the API
func (a *API) Run() {
    // Start up the BaseAPI to handle core API stuff
	go a.BaseAPI.Run()

	log.Println("Running <NAME>")

    // Send data to client (using default config values)
	a.Send(a.Data())
    
    // Set up a new ticker (if you want to send info the client at set time intervals)
    // You can change the time length as necessary
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING <NAME>: %s\n", a.UUID)
	}()

    // Create a new channel to recieve termination messages on
	stop := a.AddStopChan()
	for {
		select {
		case body := <-a.ConfigChan: // Configuration messages
			if err := a.Configure(body); err != nil {
				a.Send(nil, err)
			}
		case t := <-ticker.C: // Update timer tick
			a.Send(a.Data())
		case <-stop: // Terminate
			return
		}
	}
}
```

**api/<NAME>/types.go** holds all internal types not needed by the client
```go
package <NAME>

// Internal type definitions go here
// All types should all be unexported (lowercase)
```

**api/<NAME>/utils.go** holds utility functions particular to that API.
If you think a function can be made generic, consider putting it in the higher level `utils/` package.

```go
package <NAME>

// Utility functions go here
// All functions should be unexported (lowercase)
```

**types/<NAME>.go** holds all the public interfaces that need to be accesible both by the api and the client
```go
package types

// All types should be exported (capital first letter)

// <NAME>Response holds all of the data being sent to the client
type <NAME>Response struct {
	// Add any necessary fields
}

// <NAME>Config holds all necessary configuration parameters being sent from the client
type <NAME>Config struct {
	// Add any necessary fields
}

```
