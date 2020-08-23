package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Methods is a shortcut that stores a list of operations
type Methods map[string]Operation

// Operation is a shortcut that stores a handler function
type Operation func(w http.ResponseWriter, r *http.Request)

// operations contains a list of available endpoints on the API with their method
var operations = map[string]Methods{
	"users/{user}": {
		"GET": GetUser,
	},
	"rooms": {
		"GET": GetRooms,
	},
	"rooms/{room}": {
		"GET": GetRoom,
	},
	"rooms/{room}/messages": {
		"GET":  GetRoomMessages,
		"POST": CreateMessage,
	},
}

// Server runs the HTTP server
func Server() {
	httpServer := mux.NewRouter().StrictSlash(true)

	for path, methods := range operations {
		fullPath := fmt.Sprintf("/v1/http/%s", path)
		for method, handler := range methods {
			httpServer.HandleFunc(fullPath, handler).Methods(method)
		}
	}

	httpServer.Use(AuthenticationMiddleware)
	httpServer.Use(mux.CORSMethodMiddleware(httpServer))

	log.Println("HTTP server ready")
	err := http.ListenAndServe("localhost:1052", httpServer)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
