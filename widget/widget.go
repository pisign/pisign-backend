package widget

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/pisign/pisign-backend/api/clock"
	"github.com/pisign/pisign-backend/api/weather"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/api/manager"
)

// Widget struct for a single frontend widget
type Widget struct {
	API       api.API
	Position  *json.RawMessage
	Conn      *websocket.Conn `json:"-"`
	Pool      *Pool           `json:"-"`
	CloseChan chan bool       `json:"-"`
}

// Create creates a new widget, with a valid api attached
func Create(apiName string, conn *websocket.Conn, pool *Pool) error {
	a, err := manager.Connect(apiName)
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

	pool.Register <- widget
	go widget.Read()
	go widget.API.Run(widget)
	return nil
}

// Message holds a messsage
type Message struct {
	Type int
	Body string
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
		json.Unmarshal(p, &message)
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
func (w *Widget) UnmarshalJSON(body []byte) error {
	//TODO: find better way to Unmarshal widget
	log.Printf("w.API: %v\n", w.API)
	var fields struct {
		API      *json.RawMessage
		Position *json.RawMessage
	}
	err := json.Unmarshal(body, &fields)
	if err != nil {
		log.Println("Could not unmarshal widget: ", err)
		return err
	}

	var APIFields struct {
		Name string
	}

	err = json.Unmarshal(*fields.API, &APIFields)
	if err != nil {
		log.Printf("Could not unmarshal widget: no `API.Name` field present: %v\n", err)
		return err
	}
	log.Printf("fields: %v\n", fields)
	log.Printf("API: %s, Position: %s\n", fields.API, fields.Position)

	switch APIFields.Name {
	case "weather":
		w.API = weather.NewAPI()
	case "clock":
		w.API = clock.NewAPI()
	default:
		msg := fmt.Sprintf("Unknown api type: %s", APIFields.Name)
		return errors.New(msg)
	}

	json.Unmarshal(*fields.API, w.API)
	w.Position = fields.Position

	return nil
}
