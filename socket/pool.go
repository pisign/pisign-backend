package socket

import "fmt"

// Pool holds multiple clients
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

// NewPool generates a new pool
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// Start main entry point of a pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Printf("New Client: %s\n", client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			client.Conn.WriteJSON(client)
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Printf("Lost Client: %s\n", client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println("Error sending JSON:", err)
					return
				}
			}
		}
	}
}
