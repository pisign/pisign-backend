// Package socket creates and manages Sockets from the client
// Each 'Socket' represents a single Socket on the client, and has its own websocket connection
package socket

import (
	"encoding/json"
	"log"

	"github.com/pisign/pisign-backend/utils"

	"github.com/gorilla/websocket"
)

// Socket struct for a single frontend Socket
type Socket struct {
	Conn       *websocket.Conn       `json:"-"`
	CloseChan  chan bool             `json:"-"`
	ConfigChan chan *json.RawMessage `json:"-"`
}

// Create creates a new Socket, with a valid api attached
func Create(configChan chan *json.RawMessage, conn *websocket.Conn) *Socket {
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
		w.CloseChan <- true
		w.Conn.Close()
	}()

	for {
		_, p, err := w.Conn.ReadMessage()
		if err != nil {
			log.Println("Read->w.Conn.ReadMessage", err)
			return
		}
		log.Printf("message: %v\n", p)
		var message struct {
			API      *json.RawMessage
			Position *json.RawMessage
		}
		utils.ParseJSON(p, &message)
		if err != nil {
			log.Println("Socket Message Parsing:", err)
			continue
		}

		w.ConfigChan <- message.API
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

func (w *Socket) String() string {
	str, _ := json.Marshal(w)
	return string(str)
}
