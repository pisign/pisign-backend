package types

import (
	"encoding/json"
	"github.com/pisign/pisign-backend/utils"
)

// DataObject holds the data from the external API
type DataObject interface {
	// Build builds the data object
	Update(interface{})

	// Transform turns the data object into a front-end parsable object
	Transform() interface{}
}

// API is the entrance point of all apis to connect to a client
type API interface {
	// Configure settings from raw json message
	Configure(body *json.RawMessage)

	// Main loop that faciliates interaction between outside world and the client widet
	Run(w Socket)

	// Data gets the data to send
	Data() interface{}
}

// Socket interface, needed to avoid circular dependency with Socket package
// TODO: See if we can remove this interface without adding a circular dependency?
type Socket interface {
	json.Unmarshaler

	// Read information from the client
	Read()

	// Send information to the client
	Send(interface{})
	Close() chan bool
}

// BaseAPI base for all APIs
type BaseAPI struct {
	Position   Position
	Name       string
	ConfigChan chan *json.RawMessage `json:"-"`
}

func (b * BaseAPI) ConfigurePosition(body *json.RawMessage) {
	err := utils.ParseJSON(*body, &b.Position)
	if err != nil {
		panic("OH NO")
	}
}