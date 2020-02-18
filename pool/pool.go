package pool

import (
	"log"

	"github.com/pisign/pisign-backend/types"
)

// Pool holds multiple apis and handles registration and deletion
type Pool struct {
	registerChan   chan types.API
	unregisterChan chan types.Unregister
	Map            map[types.API]bool
	name           string
}

// NewPool generates a new pool
func NewPool() *Pool {
	return &Pool{
		registerChan:   make(chan types.API),
		unregisterChan: make(chan types.Unregister),
		Map:            make(map[types.API]bool),
		name:           "main",
	}
}

// List turns map into list
func (pool *Pool) List() []types.API {
	keys := make([]types.API, len(pool.Map))
	i := 0
	for k := range pool.Map {
		keys[i] = k
		i++
	}
	return keys
}

// Save the current state of the pool
func (pool *Pool) Save() {
	list := pool.List()
	layout := Layout{List: list, Name: pool.name}
	err := SaveLayout(layout)
	if err != nil {
		panic("Error saving layout from pool!")
	}
}

// Start main entry point of a pool
func (pool *Pool) Start() {
	for {
		select {
		case api := <-pool.registerChan:
			pool.Map[api] = true
			log.Printf("New API: %s\n", api.GetName())
			log.Println("Size of Connection Pool: ", len(pool.Map))
			pool.Save()
		case data := <-pool.unregisterChan:
			delete(pool.Map, data.API)
			if data.Save {
				pool.Save()
			}
			log.Printf("Deleted API: %s\n", data.API.GetName())
			log.Println("Size of Connection Pool: ", len(pool.Map))
		}
	}
}

// Register a new API
func (pool *Pool) Register(a types.API) {
	pool.registerChan <- a
}

// Unregister a new API
func (pool *Pool) Unregister(data types.Unregister) {
	pool.unregisterChan <- data
}
