// Package socket creates and manages Sockets from the client
// Each 'Socket' represents has its own websocket connection to a single client
package socket

import (
	"errors"
	"log"
	"net"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"

	"github.com/gorilla/websocket"
)

// Socket struct for a single frontend Socket
type Socket struct {
	conn       *websocket.Conn
	closeChan  chan bool
	configChan chan types.ClientMessage
}

// Create creates a new Socket, wrapped around a primitive websocket.Conn connection
func Create(configChan chan types.ClientMessage, conn *websocket.Conn) *Socket {
	socket := &Socket{
		conn:       conn,
		closeChan:  make(chan bool),
		configChan: configChan,
	}

	return socket
}

// Read continuously waits for incoming information from the underlying socket, and forwards the data
// appropriately through the proper channel objects
// Note: this function should ALWAYS be run asynchronously, since it blocks while waiting for new data
func (w *Socket) Read() {
	// The only time the socket is recieving data is when it is getting configutation data
	defer func() {
		log.Println("CLOSING SOCKET")
	}()

	for {
		_, p, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("Read->w.conn.ReadMessage", err)
			w.closeChan <- false
			return
		}

		var message types.ClientMessage
		if err := utils.ParseJSON(p, &message); err != nil {
			log.Println("Socket Message Parsing:", err)
			continue
		}

		switch message.Action {
		case types.Initialize, types.ConfigurePosition, types.ConfigureAPI, types.ChangeAPI, types.Delete:
			w.configChan <- message
		default:
			log.Printf("Unknown client action: %s\n", message.Action)
		}
	}
}

// SendSuccess sends a success message
func (w *Socket) SendSuccess(msg interface{}, position types.Position) {
	w.Send(types.BaseMessage{
		Status:   types.StatusSuccess,
		Data:     msg,
		Position: position,
	})
}

// SendErrorMessage sends an error message
func (w *Socket) SendErrorMessage(err error) {
	w.Send(types.BaseMessage{
		Status: types.StatusFailure,
		Error:  err.Error(),
	})
}

// Send out to client through websocket
// The SendSuccess and SendErrorMessage should be used to construct the messages. Send should rarely be used directly
func (w *Socket) Send(msg interface{}) {
	if err := w.conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending JSON to client: %v\n", err)
	}
}

func (w *Socket) CloseChan() chan bool {
	return w.closeChan
}

func (w *Socket) ConfigChan() chan types.ClientMessage {
	return w.configChan
}

// Close the underlying websocket connection
func (w *Socket) Close() error {
	if w.conn == nil {
		return errors.New("type Socket has no conn")
	}
	return w.conn.Close()
}

// RemoteAddr retrieve respective remote address from underlying websocket connection
func (w *Socket) RemoteAddr() net.Addr {
	if w.conn == nil {
		return nil
	}
	return w.conn.RemoteAddr()
}
