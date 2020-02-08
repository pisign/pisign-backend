package widget

import (
	"fmt"
)

// Pool holds multiple widgets
type Pool struct {
	Register   chan *Widget
	Unregister chan *Widget
	Widgets    map[*Widget]bool
	Broadcast  chan Message
}

// NewPool generates a new pool
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Widget),
		Unregister: make(chan *Widget),
		Widgets:    make(map[*Widget]bool),
		Broadcast:  make(chan Message),
	}
}

// Start main entry point of a pool
func (pool *Pool) Start() {
	for {
		select {
		case widget := <-pool.Register:
			pool.Widgets[widget] = true
			fmt.Printf("New Widget: %s\n", widget)
			fmt.Println("Size of Connection Pool: ", len(pool.Widgets))
		case widget := <-pool.Unregister:
			delete(pool.Widgets, widget)
			fmt.Printf("Lost Widget: %s\n", widget)
			fmt.Println("Size of Connection Pool: ", len(pool.Widgets))
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all widgets in Pool")
			for widget := range pool.Widgets {
				if err := widget.Conn.WriteJSON(message); err != nil {
					fmt.Println("Error sending JSON:", err)
					return
				}
			}
		}
	}
}
