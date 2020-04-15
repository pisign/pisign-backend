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
	"path"
	"strconv"
	"strings"

	"github.com/pisign/pisign-backend/utils"

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

func uploadImageHandler(w http.ResponseWriter, r *http.Request, pool *pool.Pool) {
	db := pool.ImageDB
	log.Println("Upload Image endpoint hit!")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var maxMemoryToUse int64
	maxMemoryToUse = 200000
	switch r.Method {
	case http.MethodPost:
		// From https://socketloop.com/tutorials/upload-multiple-files-golang
		err := r.ParseMultipartForm(maxMemoryToUse) // grab the multipart form
		if err != nil {
			log.Println(err)
			return
		}

		formdata := r.MultipartForm // ok, no problem so far, read the Form data

		// Get the tag list from the form
		tags := r.Form.Get("tag")
		numfiles, _ := strconv.Atoi(r.Form.Get("length"))

		tagsList := strings.Split(tags, ",")

		// Update the unique tags list
		for _, tag := range tagsList {
			if !utils.StringInSlice(tag, db.UniqueTags) {
				db.UniqueTags = append(db.UniqueTags, tag)
			}
		}

		for n := 0; n < numfiles; n++ { // loop through the files one by one
			key := "files_" + strconv.Itoa(n)
			formfile := formdata.File[key][0]
			file, err := formfile.Open()
			defer file.Close()
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

			// Create a new UUID for the file, and get the file extention, and save to that file
			newFileName := uuid.New().String()
			extension := path.Ext(formfile.Filename)
			filepath := "./assets/images/" + newFileName + extension
			err = db.AddImage(filepath, file, tagsList)
			if err != nil {
				fmt.Fprint(w, err.Error())
			}

			fmt.Fprintf(w, "Files uploaded successfully : ")
		}
		pool.SaveImageDB()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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

func setupStaticFiles(directory string, path string) {
	fileServer := http.FileServer(http.Dir(directory))
	http.Handle(path, http.StripPrefix(path, fileServer))
}

func setupRoutes() *pool.Pool {
	p := pool.NewPool()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socketConnectionHandler(p, w, r)
	})
	http.HandleFunc("/layouts", serveLayouts)
	http.HandleFunc("/uploads", func(w http.ResponseWriter, r *http.Request) {
		uploadImageHandler(w, r, p)
	})

	http.HandleFunc("/assets/images/", func(w http.ResponseWriter, r *http.Request) {
		fileServer := http.FileServer(http.Dir("assets/images"))
		// Allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.StripPrefix("/assets/images/", fileServer).ServeHTTP(w, r)
	})
	setupStaticFiles("dist", "/")
	return p
}

// StartLocalServer creates a new server on localhost with the specified port
func StartLocalServer(port int) error {
	addr := fmt.Sprintf("0.0.0.0:%v", port)
	log.Printf("Running server at %v\n", addr)
	p := setupRoutes()
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// Make sure to save the ImageDB before shutdown
	defer func() {
		p.SaveImageDB()
	}()
	return nil
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
