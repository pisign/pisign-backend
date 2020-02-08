package widget

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/api/manager"
)

// Widget struct for a single frontend widget
type Widget struct {
	ID   int             `json:"id"`
	API  api.API         `json:"api"`
	Conn *websocket.Conn `json:"-"`
	Pool *Pool           `json:"-"`
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
		ID:   len(pool.Widgets),
		API:  a,
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- widget
	go widget.Read()
	go widget.API.Run(conn)
	widget.Conn.WriteJSON(widget)
	return nil
}

// Message holds a messsage
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (w *Widget) Read() {
	defer func() {
		w.Pool.Unregister <- w
		w.Conn.Close()
	}()

	for {
		messageType, p, err := w.Conn.ReadMessage()
		if err != nil {
			log.Println("Read->w.Conn.ReadMessage", err)
			return
		}

		message := Message{Type: messageType, Body: string(p)}
		fmt.Printf("Message Received from %s: %+v\n", w, message)
	}
}

func (w *Widget) String() string {
	str, _ := json.Marshal(w)
	return string(str)
}
