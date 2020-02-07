package socket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client client
type Client struct {
	ID      int             `json:"id"`
	APIName string          `json:"api"`
	Conn    *websocket.Conn `json:"-"`
	Pool    *Pool           `json:"-"`
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
		fmt.Printf("Message Received from %s: %+v\n", c, message)
	}
}

func (c *Client) String() string {
	return fmt.Sprintf("Client{ID=%v, API:%s}", c.ID, c.APIName)
}
