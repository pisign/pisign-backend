package api

import "encoding/json"

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

// API is the entrance point of all apis to connect to a client
type API interface {
	// Configure settings from raw json message
	Configure(body *json.RawMessage)

	// Main loop that faciliates interaction between outside world and the client widet
	Run(w Widget)
}

// Widget interface, needed to avoid circular dependency with widget package
// TODO: See if we can remove this interface without adding a circular dependency?
type Widget interface {
	json.Unmarshaler

	// Read information from the client
	Read()

	// Send information to the client
	Send(interface{})
	Close() chan bool
}

// BaseAPI base for all APIs
type BaseAPI struct {
	APIName string `json:"Name"`
}

// Name gets name of the api
func (a *BaseAPI) Name() string {
	return a.APIName
}
