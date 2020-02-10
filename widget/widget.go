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
	ID        int
	APIName   string
	API       api.API
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
		ID:        len(pool.Widgets),
		API:       a,
		APIName:   apiName,
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
		messageType, p, err := w.Conn.ReadMessage()
		if err != nil {
			log.Println("Read->w.Conn.ReadMessage", err)
			return
		}

		message := Message{Type: messageType, Body: string(p)}
		log.Printf("Message Received from %s: %+v\n", w, message)
		w.API.Configure(p)
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
func (w *Widget) UnmarshalJSON(b []byte) error {
	//TODO: find better way to Unmarshal widget
	var name struct {
		APIName string
	}

	err := json.Unmarshal(b, &name)
	if err != nil {
		log.Println("Could not unmarshal widget error 1: ", err)
		return err
	}
	var t thing
	switch name.APIName {
	case "weather":
		t.API = new(weather.API)
	case "clock":
		t.API = new(clock.API)
	default:
		msg := fmt.Sprintf("Unknown api type: %s", name.APIName)
		return errors.New(msg)
	}
	err = json.Unmarshal(b, &t)
	w.API = t.API
	w.APIName = t.APIName
	w.ID = t.ID
	return nil
}

type thing struct {
	ID      int     `json:"id"`
	API     api.API `json:"api"`
	APIName string  `json:"apiName"`
}
