package types

import (
	"log"

	"github.com/google/uuid"
)

// CloseType represents when a socket closes, or a forceful close if Socket is nil
type CloseType struct {
	Socket Socket
	Save   bool
}

// BaseAPI base for all APIs
type BaseAPI struct {
	Position
	Name       string
	UUID       uuid.UUID
	Sockets    map[Socket]bool    `json:"-"`
	Pool       Pool               `json:"-"`
	CloseChan  chan CloseType     `json:"-"`
	ConfigChan chan ClientMessage `json:"-"`
}

// Init Initialization
func (b *BaseAPI) Init(name string, sockets map[Socket]bool, pool Pool, id uuid.UUID) {
	b.Name = name
	b.Sockets = sockets
	b.Pool = pool
	b.UUID = id
	// TODO Remove arbitrary magic numbers?
	b.CloseChan = make(chan CloseType, 10)
	b.ConfigChan = make(chan ClientMessage, 10)
}

// GetName returns the name (or type) of the api
func (b *BaseAPI) GetName() string {
	return b.Name
}

// GetUUID returns the uuid
func (b *BaseAPI) GetUUID() uuid.UUID {
	return b.UUID
}

// GetSockets returns the apis socket connection
func (b *BaseAPI) GetSockets() map[Socket]bool {
	return b.Sockets
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
	b.CloseChan <- CloseType{Socket: nil, Save: true}
	log.Printf("Done stopping!\n")
}

// StopAll force stops the whole api
func (b *BaseAPI) StopAll(save bool) {
	for socket := range b.Sockets {
		socket.Close()
	}
}

// Send to websocket
func (b *BaseAPI) Send(data interface{}, err error) {
	if err != nil {
		for socket := range b.Sockets {
			socket.SendErrorMessage(err)
		}
	} else {
		for socket := range b.Sockets {
			socket.SendSuccess(data)
		}
	}
}

// AddSocket to the api's internal socket list
func (b *BaseAPI) AddSocket(socket Socket) {
	if _, ok := b.Sockets[socket]; !ok {
		b.Sockets[socket] = true
		go func(config chan ClientMessage, close chan bool) {
			defer func() {
				log.Printf("Exiting channel forwarding for socket!\n")
			}()
			for msg := range config {
				b.ConfigChan <- msg
			}
			for save := range close {
				b.CloseChan <- CloseType{Socket: socket, Save: save}
				return
			}
		}(socket.Config(), socket.Close())
	}
}
