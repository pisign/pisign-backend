package types

import (
	"log"

	"github.com/google/uuid"
)

// CloseType represents when a socket closes, or a forceful close if Socket is nil
type CloseType struct {
	Sockets map[Socket]bool
	Save    bool
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
	running    bool
}

// Init Initialization
func (b *BaseAPI) Init(name string, sockets map[Socket]bool, pool Pool, id uuid.UUID) {
	b.Name = name
	b.Sockets = make(map[Socket]bool)
	for socket := range sockets {
		b.AddSocket(socket)
	}
	b.Pool = pool
	b.UUID = id
	// TODO Remove arbitrary magic numbers?
	b.CloseChan = make(chan CloseType, 10)
	b.ConfigChan = make(chan ClientMessage, 10)
	b.running = true
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
func (b *BaseAPI) Configure(message ClientMessage) error {
	switch message.Action {
	case ConfigurePosition, Initialize:
		b.SetPosition(message.Position)
	case ChangeAPI:
		return b.Pool.Switch(b, message.Name)
	case Delete:
		b.CloseChan <- CloseType{Sockets: b.Sockets, Save: true}
	}
	return nil
}

func (b *BaseAPI) Run() {
	defer func() {
		log.Printf("STOPPING BASE API: %s\n", b.UUID)
	}()
	for b.Running() {
		select {
		case msg := <-b.CloseChan:
			log.Printf("Closing sockets: %v\n", msg.Sockets)
			for socket := range msg.Sockets {
				socket.Close()
				delete(b.Sockets, socket)
			}
			if len(b.Sockets) == 0 {
				b.Stop()
				b.Pool.Unregister(Unregister{
					API:  b,
					Save: msg.Save,
				})
			}
		default:
			continue
		}
	}
}

func (b *BaseAPI) Data() (interface{}, error) {
	return nil, nil
}

// SetPosition configures position
func (b *BaseAPI) SetPosition(pos Position) {
	b.Position = pos
}

// Stop the api
func (b *BaseAPI) Stop() {
	log.Printf("Stopping api %s (%s)\n", b.Name, b.UUID)
	b.running = false
}

// Send to websocket
func (b *BaseAPI) Send(data interface{}, err error) {
	// If API has already been closed
	if !b.Running() {
		return
	}
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
		log.Printf("Adding new socket: %v!\n", socket)
		b.Sockets[socket] = true
		go func(config chan ClientMessage, close chan bool) {
			defer func() {
				log.Printf("Exiting channel forwarding for socket!\n")
			}()
			for b.Running() {
				select {
				case msg := <-config:
					b.ConfigChan <- msg
				case save := <-close:
					sockets := make(map[Socket]bool)
					sockets[socket] = true
					b.CloseChan <- CloseType{Sockets: sockets, Save: save}
				default:
					continue
				}
			}
		}(socket.ConfigChan(), socket.CloseChan())
	}
}

func (b *BaseAPI) Running() bool {
	return b.running
}
