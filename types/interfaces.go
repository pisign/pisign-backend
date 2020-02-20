package types

import "github.com/google/uuid"

// DataObject holds the data from the external API
type DataObject interface {
	// Build builds the data object
	Update(interface{}) error

	// Transform turns the data object into a front-end parsable object
	Transform() interface{}
}

// API is the entrance point of all apis to connect to a client
type API interface {
	// Configure settings from raw json message
	Configure(message ClientMessage) error

	// Main loop that faciliates interaction between outside world and the client widet
	Run()

	// Data gets the data to send
	Data() (interface{}, error)

	GetName() string
	GetUUID() uuid.UUID
	GetSockets() map[Socket]bool
	GetPosition() Position
	SetPosition(Position)
	Stop()
	Send(interface{}, error)

	AddSocket(Socket)
}

// Socket interface, needed to avoid circular dependency with Socket package
type Socket interface {
	// Read information from the client
	Read()

	// Send information to the client
	Send(interface{})

	// Close the socket
	Close() chan bool

	// Configuration data
	Config() chan ClientMessage

	// SendErrorMessage sends error message
	SendErrorMessage(error)
	SendSuccess(interface{})
}

// Unregister stores info about which api to unregister, and weather the pool should be saved
type Unregister struct {
	API  API
	Save bool
}

// Pool pool
type Pool interface {
	Start()
	Register(API)
	Unregister(Unregister)
	Switch(API, string) error
	Save()
	Add(string, uuid.UUID, map[Socket]bool) (API, error)
}
