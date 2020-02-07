package socket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client client
type Client struct {
	APIName string
	Conn    *websocket.Conn
	Pool    *Pool
}

// Message holds a messsage
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read->c.Conn.ReadMessage", err)
			return
		}

		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
