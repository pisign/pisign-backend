package types

import (
	"log"

	"github.com/google/uuid"
)

// BaseAPI base for all APIs
type BaseAPI struct {
	Position
	Name     string
	UUID     uuid.UUID
	Socket   Socket    `json:"-"`
	Pool     Pool      `json:"-"`
	StopChan chan bool `json:"-"`
}

// Init Initialization
func (b *BaseAPI) Init(name string, socket Socket, pool Pool, id uuid.UUID) {
	b.Name = name
	b.Socket = socket
	b.Pool = pool
	b.UUID = id
	b.StopChan = make(chan bool, 1)
}

// GetName returns the name (or type) of the api
func (b *BaseAPI) GetName() string {
	return b.Name
}

// GetUUID returns the uuid
func (b *BaseAPI) GetUUID() uuid.UUID {
	return b.UUID
}

// GetSocket returns the apis socket connection
func (b *BaseAPI) GetSocket() Socket {
	return b.Socket
}

// GetPosition returns position
func (b *BaseAPI) GetPosition() Position {
	return b.Position
}

// Configure based on an incoming client message
func (b *BaseAPI) Configure(message ClientMessage) {
	switch message.Action {
	case ConfigurePosition, Initialize:
		b.SetPosition(message.Position)
	}
}

// SetPosition configures position
func (b *BaseAPI) SetPosition(pos Position) {
	b.Position = pos
}

// Stop the api
func (b *BaseAPI) Stop() {
	log.Printf("Stopping api %s (%s)\n", b.Name, b.UUID)
	b.StopChan <- true
	log.Printf("Done stopping!\n")
}

// Send to websocket
func (b *BaseAPI) Send(data interface{}, err error) {
	if err != nil {
		b.Socket.SendErrorMessage(err)
	} else {
		b.Socket.SendSuccess(data)
	}
}
