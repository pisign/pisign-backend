package api

import "github.com/gorilla/websocket"

// InternalAPI is the interface our internal API uses
type InternalAPI interface {
	// Serialize transforms the data structure into a byte slice to be sent
	// via the websocket - this should be done as a serialized JSON string
	Serialize() []byte
}

// ExternalAPI is the interface for all our APIs
type ExternalAPI interface {
	// Get takes in an arbitrary argument, and builds the struct
	Get(interface{})

	// Transform takes the API data and turns it the data we are going to
	// send to the frontend
	Transform(interface{}) InternalAPI
}

// API is the entrance point of all apis to connect to a client
type API interface {
	Configure(json string)
	Name() string
	Run(*websocket.Conn)
}

// BaseAPI base for all APIs
type BaseAPI struct {
	APIName string `json:"apiName"`
}

// Name gets name of the api
func (a *BaseAPI) Name() string {
	return a.APIName
}
