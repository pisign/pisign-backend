package pool

import (
	"errors"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/api"

	"github.com/pisign/pisign-backend/types"
)

// Pool holds multiple apis and handles registration and deletion
type Pool struct {
	registerChan   chan types.API
	unregisterChan chan types.Unregister
	Map            map[uuid.UUID]types.API
	name           string
}

// NewPool generates a new pool
func NewPool() *Pool {
	return &Pool{
		registerChan:   make(chan types.API),
		unregisterChan: make(chan types.Unregister),
		Map:            make(map[uuid.UUID]types.API),
		name:           "main",
	}
}

// List turns map into list
func (pool *Pool) List() []types.API {
	keys := make([]types.API, len(pool.Map))
	i := 0
	for _, a := range pool.Map {
		keys[i] = a
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

// Add a new API to the pool
func (pool *Pool) Add(apiName string, id uuid.UUID, sockets map[types.Socket]bool) (types.API, error) {
	if a := pool.containsUUID(id); a != nil {
		for existingSocket := range a.GetSockets() {
			log.Printf("Remote addr: %+v\n", existingSocket.RemoteAddr())
			for newSocket := range sockets {
				if equal, ip := equalAddresses(existingSocket.RemoteAddr(), newSocket.RemoteAddr()); equal {
					log.Printf("API (%v) already exists on client %v, cannot recreate!\n", id, ip)
					return nil, errors.New("API already exists on client")
				}
			}
		}
		// Add new socket to existing api
		log.Printf("API (%s) Already exists!\n", id)
		for socket := range sockets {
			a.AddSocket(socket)
		}
		// TODO Force data send to new socket

		return a, nil
	}

	// Else Spin up a new api
	log.Printf("API (%s) Being created now!\n", id)
	a, err := api.NewAPI(apiName, sockets, pool, id)
	if err != nil {
		return nil, err
	}
	pool.Register(a)

	go a.Run()
	return a, nil

}

// Register a new API
func (pool *Pool) Register(a types.API) {
	pool.Map[a.GetUUID()] = a
	log.Printf("New API: %s\n", a.GetName())
	log.Println("Size of Connection Pool: ", len(pool.Map))
	pool.Save()
}

// Unregister a new API
func (pool *Pool) Unregister(data types.Unregister) {
	delete(pool.Map, data.API.GetUUID())
	if data.Save {
		log.Printf("Saving API: %s\n", data.API.GetName())
		pool.Save()
	}
	log.Printf("Deleted API: %s\n", data.API.GetName())
	log.Println("Size of Connection Pool: ", len(pool.Map))
}

// Switch from one API type to another, while maintaining the same socket
func (pool *Pool) Switch(a types.API, name string) error {
	log.Printf("Switching API: %s -> %s\n", a.GetName(), name)
	sockets := a.GetSockets()
	id := a.GetUUID()
	pos := a.GetPosition()
	a.Stop()
	pool.Unregister(types.Unregister{API: a, Save: true})
	newAPI, err := pool.Add(name, id, sockets)
	if err != nil {
		return err
	}
	newAPI.SetPosition(pos)
	return nil
}

func (pool *Pool) containsUUID(targetUUID uuid.UUID) types.API {
	for id, a := range pool.Map {
		if id == targetUUID {
			return a
		}
	}
	return nil
}

func equalAddresses(a net.Addr, b net.Addr) (bool, net.IP) {
	// TODO figure out if we'll ever encounter non-tcp addresses?
	tcpA := a.(*net.TCPAddr)
	tcpB := b.(*net.TCPAddr)
	if tcpA.IP.Equal(tcpB.IP) {
		return true, tcpA.IP
	}

	return false, nil
}
