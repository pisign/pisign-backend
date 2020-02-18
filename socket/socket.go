// Package socket creates and manages Sockets from the client
// Each 'Socket' represents a single Socket on the client, and has its own websocket connection
package socket

import (
	"encoding/json"
	"log"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"

	"github.com/gorilla/websocket"
)

// Socket struct for a single frontend Socket
type Socket struct {
	Conn       *websocket.Conn          `json:"-"`
	CloseChan  chan bool                `json:"-"`
	ConfigChan chan types.ClientMessage `json:"-"`
	DeleteChan chan bool                `json:"-"`
}

// Create creates a new Socket, with a valid api attached
func Create(configChan chan types.ClientMessage, conn *websocket.Conn) *Socket {
	socket := &Socket{
		Conn:       conn,
		CloseChan:  make(chan bool),
		ConfigChan: configChan,
	}

	return socket
}

func (w *Socket) Read() {
	// The only time the socket is recieving data is when it is getting configutation data
	defer func() {
		log.Println("CLOSING SOCKET")
	}()

	for {
		_, p, err := w.Conn.ReadMessage()
		if err != nil {
			log.Println("Read->w.Conn.ReadMessage", err)
			w.CloseChan <- false
			return
		}

		var message types.ClientMessage
		if err := utils.ParseJSON(p, &message); err != nil {
			log.Println("Socket Message Parsing:", err)
			continue
		}

		switch message.Action {
		case types.Initialize, types.ConfigurePosition, types.ConfigureAPI, types.ChangeAPI:
			w.ConfigChan <- message
		case types.Delete:
			w.CloseChan <- true
			return
		default:
			log.Printf("Unknown client action: %s\n", message.Action)
		}
	}
}

// SendSuccess sends a success message
func (w *Socket) SendSuccess(msg interface{}) {
	w.Send(types.BaseMessage{
		Status: types.StatusSuccess,
		Data:   msg,
	})
}

// SendErrorMessage sends the error message
func (w *Socket) SendErrorMessage(err error) {
	w.Send(types.BaseMessage{
		Status: types.StatusFailure,
		Error:  err.Error(),
	})
}

// SendDataOrError sends the appropriate data or error
func (w *Socket) SendDataOrError(data interface{}, err error) {
	if err != nil {
		w.SendErrorMessage(err)
	} else {
		w.SendSuccess(data)
	}
}

// Send out to client through websocket
func (w *Socket) Send(msg interface{}) {
	w.Conn.WriteJSON(msg)
}

// Close returns a channel that signifies the closing of the Socket
func (w *Socket) Close() chan bool {
	return w.CloseChan
}

// Config for api and position configuration
func (w *Socket) Config() chan types.ClientMessage {
	return w.ConfigChan
}

func (w *Socket) String() string {
	str, _ := json.Marshal(w)
	return string(str)
}
