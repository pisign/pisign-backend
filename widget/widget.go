package widget

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/pisign/pisign-backend/utils"

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
		log.Printf("message: %v\n", p)
		message, err := utils.ParseJSONMap(p)
		if err != nil {
			log.Println("Widget Message Parsing:", err)
			continue
		}

		if API, ok := message["API"]; ok {
			w.API.Configure(API)
		}

		if pos, ok := message["Position"]; ok && pos != nil {
			log.Printf("Setting position: %s\n", *pos)
			if err = json.Unmarshal(*pos, &w.Position); err != nil {
				log.Printf("Pos err: %v\n", err)
			}
		}

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
func (w *Widget) UnmarshalJSON(b []byte) error {
	//TODO: find better way to Unmarshal widget
	fields, err := utils.ParseJSONMap(b)
	if err != nil {
		log.Println("Could not unmarshal widget: ", err)
		return err
	}

	API, err := utils.ParseJSONMap(*fields["API"])
	if err != nil {
		log.Printf("Could not unmarshal widget: no `API` field present: %v\n", err)
		return err
	}

	var name string
	err = json.Unmarshal(*API["Name"], &name)
	if err != nil {
		log.Printf("Could not unmarshal widget: no `API.Name` field present: %v\n", err)
		return err
	}
	switch string(name) {
	case "weather":
		w.API = weather.NewAPI()
	case "clock":
		w.API = clock.NewAPI()
	default:
		msg := fmt.Sprintf("Unknown api type: %s", name)
		return errors.New(msg)
	}

	err = json.Unmarshal(*fields["API"], w.API)
	if err != nil {
		log.Println("Could not unmarshal widget.API:", err)
		return err
	}
	log.Printf("API data: %+v\n", w.API)

	position, ok := fields["Position"]
	if ok && position != nil {
		log.Printf("ok: %v, Position: %v\n", ok, position)
		err = json.Unmarshal(*position, &w.Position)
		if err != nil {
			log.Println("Could not unmarshal widget.Position:", err)
			return err
		}
	}
	return nil
}
