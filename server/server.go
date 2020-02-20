// Package server much of the code comes from https://tutorialedge.net/projects/chat-system-in-go-and-react/part-4-handling-multiple-clients/
// Two main routes:
// - /ws creates a new websocket connection
// - /layouts?name=<name> retrieves a given layout, stored in a `json` file on the server
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pisign/pisign-backend/pool"
	"github.com/pisign/pisign-backend/types"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/socket"
)

func socketConnectionHandler(pool types.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket endpoint hit!")
	conn, err := upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	configChan := make(chan types.ClientMessage)

	apiName := r.FormValue("api")
	idString := r.FormValue("uuid")
	id := uuid.MustParse(idString)

	ws := socket.Create(configChan, conn)
	sockets := make(map[types.Socket]bool)
	sockets[ws] = true

	_, err = pool.Add(apiName, id, sockets)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.Close()
		return
	}

	// Socket connection handler should be the one to register, call the read method,
	go ws.Read()
}

func serveLayouts(w http.ResponseWriter, r *http.Request) {
	log.Println("Layouts endpoing hit!")
	layoutName := r.FormValue("name")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if layoutName == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, "Must supply `name` parameter")
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Printf("Retrieving layout data for %s...\n", layoutName)
		bytes, err := json.Marshal(pool.LoadLayout(layoutName))
		if err != nil {
			log.Printf("Error loading layout %s: %v\n", layoutName, err)
			return
		}
		fmt.Fprintf(w, "%s", string(bytes))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func setupRoutes() {
	p := pool.NewPool()
	go p.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socketConnectionHandler(p, w, r)
	})
	http.HandleFunc("/layouts", serveLayouts)
}

// StartLocalServer creates a new server on localhost
func StartLocalServer(port int) {
	addr := fmt.Sprintf("0.0.0.0:%v", port)
	log.Printf("Running server at %v\n", addr)
	setupRoutes()
	http.ListenAndServe(addr, nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Upgrades the HTTPS protocol to socket protocol
func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade->upgrader:", err)
		return ws, err
	}
	return ws, nil
}
