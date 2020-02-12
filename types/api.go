package types

import "encoding/json"

// API is the entrance point of all apis to connect to a client
type API interface {
	// Configure settings from raw json message
	Configure(body *json.RawMessage)

	// Main loop that faciliates interaction between outside world and the client widet
	Run()

	SetWidget(w Widget)
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
	APIName string `json:"Name"`
	Widget  Widget `json:"-"`
}

// Name gets name of the api
func (a *BaseAPI) Name() string {
	return a.APIName
}

// SetWidget sets the widget for the api
func (a *BaseAPI) SetWidget(w Widget) {
	a.Widget = w
}
