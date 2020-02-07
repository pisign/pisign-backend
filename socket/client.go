package socket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/api"
)

// Client client
type Client struct {
	ID   int             `json:"id"`
	API  api.API         `json:"api"`
	Conn *websocket.Conn `json:"-"`
	Pool *Pool           `json:"-"`
}

// CreateClient creates a new client, with a valid api attached
func CreateClient(apiName string, conn *websocket.Conn, pool *Pool) error {
	a, err := api.Connect(apiName)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.Close()
		return err
	}

	client := &Client{
		ID:   len(pool.Clients),
		API:  a,
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
	return nil
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
	str, _ := json.Marshal(c)
	return string(str)
}
