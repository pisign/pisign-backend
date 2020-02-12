// Package widget creates and manages widgets from the client
// Each 'widget' represents a single widget on the client, and has its own websocket connection
package widget

import (
	"encoding/json"
	"log"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/api/manager"
)

// Widget struct for a single frontend widget
type Widget struct {
	API       types.API
	Position  *json.RawMessage
	Conn      *websocket.Conn `json:"-"`
	Pool      *Pool           `json:"-"`
	CloseChan chan bool       `json:"-"`
}

// Create creates a new widget, with a valid api attached
func Create(apiName string, conn *websocket.Conn, pool *Pool) error {
	a, err := manager.NewAPI(apiName)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.Close()
		return err
	}

	widget := &Widget{
		API:       a,
		Conn:      conn,
		Pool:      pool,
		CloseChan: make(chan bool),
	}
	a.SetWidget(widget)

	pool.Register <- widget
	go widget.Read()
	go widget.API.Run()
	return nil
}

func (w *Widget) Read() {
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
			log.Println("Widget Message Parsing:", err)
			continue
		}

		w.API.Configure(message.API)
		w.Position = message.Position

		log.Printf("widget with new data: %+v\n", w)
		w.Pool.save()
	}
}

// Send out to client through websocket
func (w *Widget) Send(msg interface{}) {
	w.Conn.WriteJSON(msg)
}

// Close returns a channel that signifies the closing of the widget
func (w *Widget) Close() chan bool {
	return w.CloseChan
}

func (w *Widget) String() string {
	str, _ := json.Marshal(w)
	return string(str)
}

// UnmarshalJSON for widget
// Dynamically creates correct type of api for properly unmarshalling
func (w *Widget) UnmarshalJSON(body []byte) error {
	//TODO: find better way to Unmarshal widget
	log.Printf("w.API: %v\n", w.API)
	var fields struct {
		API      *json.RawMessage
		Position *json.RawMessage
	}
	err := utils.ParseJSON(body, &fields)
	if err != nil {
		log.Println("Could not unmarshal widget: ", err)
		return err
	}

	var APIFields struct {
		Name string
	}

	err = utils.ParseJSON(*fields.API, &APIFields)
	if err != nil {
		log.Printf("Could not unmarshal widget: no `API.Name` field present: %v\n", err)
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
