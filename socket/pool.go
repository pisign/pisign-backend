package socket

import (
	"log"
)

// Pool holds multiple sockets and handles registration and deletion
type Pool struct {
	Register   chan *Socket
	Unregister chan *Socket
	Sockets    map[*Socket]bool
	name       string
}

// NewPool generates a new pool
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Socket),
		Unregister: make(chan *Socket),
		Sockets:    make(map[*Socket]bool),
		name:       "main",
	}
}

func (pool *Pool) SocketList() []*Socket {
	keys := make([]*Socket, len(pool.Sockets))
	i := 0
	for k := range pool.Sockets {
		keys[i] = k
		i++
	}
	return keys
}

func (pool *Pool) save() error {
	Sockets := pool.SocketList()
	layout := Layout{Sockets: Sockets, Name: pool.name}
	return SaveLayout(layout)
}

// Start main entry point of a pool
func (pool *Pool) Start() {
	for {
		select {
		case Socket := <-pool.Register:
			pool.Sockets[Socket] = true
			log.Printf("New Socket: %s\n", Socket)
			log.Println("Size of Connection Pool: ", len(pool.Sockets))
			pool.save()
		case Socket := <-pool.Unregister:
			delete(pool.Sockets, Socket)
			log.Printf("Lost Socket: %s\n", Socket)
			log.Println("Size of Connection Pool: ", len(pool.Sockets))
		}
	}
}
