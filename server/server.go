package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pisign/pisign-backend/api"

	"github.com/pisign/pisign-backend/socket"
)

func serveWs(pool *socket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket endpoing hit!")
	conn, err := socket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	apiName := r.FormValue("api")
	if !api.Valid(apiName) {
		log.Printf("Invalid APIName in incoming websocket connection: %v\n", apiName)
		conn.Close()
		return
	}

	client := &socket.Client{
		ID:      len(pool.Clients),
		APIName: apiName,
		Conn:    conn,
		Pool:    pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := socket.NewPool()
	go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

// StartLocalServer creates a new server on localhost
func StartLocalServer(port int) {
	addr := fmt.Sprintf("localhost:%v", port)
	fmt.Printf("Running server at %v\n", addr)
	setupRoutes()
	http.ListenAndServe(addr, nil)
}
