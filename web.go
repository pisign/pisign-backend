package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":1718", "http service address")

func web() {
	fmt.Println("Starting up server...")
	flag.Parse()
	http.Handle("/", http.HandlerFunc(root))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}

}

func root(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Responding...")
}
