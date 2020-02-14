// Package socket creates and manages Sockets from the client
// Each 'Socket' represents a single Socket on the client, and has its own websocket connection
package socket

import (
	"encoding/json"
	"log"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/api/manager"
)

// Socket struct for a single frontend Socket
type Socket struct {
	API       types.API
	Position  *json.RawMessage
	Conn      *websocket.Conn `json:"-"`
	Pool      *Pool           `json:"-"`
	CloseChan chan bool       `json:"-"`
}

// Create creates a new Socket, with a valid api attached
func Create(a api.API, conn *websocket.Conn, pool *Pool) *Socket {
	socket := &Socket{
		API:       a,
		Conn:      conn,
		Pool:      pool,
		CloseChan: make(chan bool),
	}

	return socket
}

func (w *Socket) Read() {
	// The only time the socket is recieving data is when it is getting configutation data
	defer func() {
		w.Pool.Unregister <- w
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

		w.API.Configure(message.API)
		w.Position = message.Position

		log.Printf("Socket with new data: %+v\n", w)
		w.Pool.save()
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

// UnmarshalJSON for Socket
// Dynamically creates correct type of api for properly unmarshalling
func (w *Socket) UnmarshalJSON(body []byte) error {
	//TODO: find better way to Unmarshal Socket
	log.Printf("w.API: %v\n", w.API)
	var fields struct {
		API      *json.RawMessage
		Position *json.RawMessage
	}
	err := utils.ParseJSON(body, &fields)
	if err != nil {
		log.Println("Could not unmarshal Socket: ", err)
		return err
	}

	var APIFields struct {
		Name string
	}

	err = utils.ParseJSON(*fields.API, &APIFields)
	if err != nil {
		log.Printf("Could not unmarshal Socket: no `API.Name` field present: %v\n", err)
		return err
	}
	log.Printf("API: %s, Position: %s\n", fields.API, fields.Position)

	w.API, err = manager.NewAPI(APIFields.Name)
	if err != nil {
		log.Printf("Unknown API type: %s\n", APIFields.Name)
		return err
	}

	utils.ParseJSON(*fields.API, w.API)
	w.Position = fields.Position

	return nil
}
