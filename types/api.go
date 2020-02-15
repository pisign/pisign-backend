package types

import "encoding/json"

// API is the entrance point of all apis to connect to a client
type API interface {
	// Configure settings from raw json message
	Configure(body *json.RawMessage)

	// Main loop that faciliates interaction between outside world and the client widet
	Run()

	// Sets values needed to communicate with client and shut down api
	SetThings(send func(interface{}), close chan bool)
}

// InternalAPI is the interface our internal API uses
type InternalAPI interface {
	// Cache caches the internal API data
	Cache()
}

// ExternalAPI is the interface for all our APIs
type ExternalAPI interface {
	// Transform turns the ExternalAPI response into an InternalAPI response
	Transform() InternalAPI
}

// BaseAPI base for all APIs
type BaseAPI struct {
	APIName string            `json:"Name"`
	Send    func(interface{}) `json:"-"`
	Close   chan bool         `json:"-"`
}

// SetThings sets the function to send data to client
func (a *BaseAPI) SetThings(send func(interface{}), close chan bool) {
	a.Send = send
	a.Close = close
}
