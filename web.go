package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pisign/backend/sockets"
)

func serveWs(pool *sockets.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket endpoing hit!")
	conn, err := sockets.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &sockets.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := sockets.NewPool()
	go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

var addr = flag.String("addr", "localhost:9000", "http service address")

func main() {
	fmt.Printf("Running server at %v\n", *addr)
	setupRoutes()
	http.ListenAndServe(*addr, nil)
}
