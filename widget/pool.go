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
	name       string
}

// NewPool generates a new pool
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Widget),
		Unregister: make(chan *Widget),
		Widgets:    make(map[*Widget]bool),
		Broadcast:  make(chan Message),
		name:       "unknown",
	}
}

func (pool *Pool) widgetList() []*Widget {
	keys := make([]*Widget, len(pool.Widgets))
	i := 0
	for k := range pool.Widgets {
		keys[i] = k
		i++
	}
	return keys
}

func (pool *Pool) save() error {
	widgets := pool.widgetList()
	layout := Layout{Widgets: widgets}
	return SaveLayout(pool.name, layout)
}

// Start main entry point of a pool
func (pool *Pool) Start() {
	for {
		select {
		case widget := <-pool.Register:
			pool.Widgets[widget] = true
			fmt.Printf("New Widget: %s\n", widget)
			fmt.Println("Size of Connection Pool: ", len(pool.Widgets))
			pool.save()
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
