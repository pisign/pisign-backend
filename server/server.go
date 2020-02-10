package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pisign/pisign-backend/socket"
	"github.com/pisign/pisign-backend/widget"
)

func serveWs(pool *widget.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket endpoing hit!")
	conn, err := socket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	apiName := r.FormValue("api")

	widget.Create(apiName, conn, pool)
}

func serveLayouts(w http.ResponseWriter, r *http.Request) {
	log.Println("Layouts endpoing hit!")
	layoutName := r.FormValue("name")
	if layoutName == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, "Must supply `name` parameter")
		return
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Printf("Retrieving layout data for %s...\n", layoutName)
		fmt.Fprintf(w, "%+v", widget.LoadLayout(layoutName))
	case http.MethodPost:
		fmt.Fprintf(w, "Updating layout data for %s...\n", layoutName)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func setupRoutes() {
	pool := widget.NewPool()
	go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	http.HandleFunc("/layouts", serveLayouts)
}

// StartLocalServer creates a new server on localhost
func StartLocalServer(port int) {
	addr := fmt.Sprintf("localhost:%v", port)
	fmt.Printf("Running server at %v\n", addr)
	setupRoutes()
	http.ListenAndServe(addr, nil)
}
