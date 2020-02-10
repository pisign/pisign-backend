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
	Position  map[string]interface{}
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

		var data struct {
			API      map[string]interface{}
			Position map[string]interface{}
		}
		log.Printf("message: %v\n", p)
		err = json.Unmarshal(p, &data)
		if err != nil {
			log.Println("Json Unmarshal to data struct:", err)
			return
		}
		log.Printf("data with position: %+v\n", data)
		w.Position = data.Position
		w.API.Configure(data.API)
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
	var a struct {
		API api.BaseAPI
	}

	err := json.Unmarshal(b, &a)
	if err != nil {
		log.Println("Could not unmarshal widget error 1: ", err)
		return err
	}
	var t thing
	name := a.API.APIName
	switch name {
	case "weather":
		t.API = weather.NewAPI()
	case "clock":
		t.API = clock.NewAPI()
	default:
		msg := fmt.Sprintf("Unknown api type: %s", name)
		return errors.New(msg)
	}
	err = json.Unmarshal(b, &t)
	log.Printf("API data: %+v\n", t.API)
	w.API = t.API
	w.Position = t.Position
	return nil
}

type thing struct {
	API      api.API
	Position map[string]interface{}
}
