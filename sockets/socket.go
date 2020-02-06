// Package sockets Code from https://tutorialedge.net/projects/chat-system-in-go-and-react/part-4-handling-multiple-clients/
package sockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Upgrade upgrades an http request to a websocket connection
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade->upgrader:", err)
		return ws, err
	}
	return ws, nil
}
