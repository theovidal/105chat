package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// Methods is a shortcut that stores a list of operations
type Methods map[string]Operation

// Operation is a shortcut that stores a handler function
type Operation func(w http.ResponseWriter, r *http.Request)

// operations contains a list of available endpoints on the API with their method
var operations = map[string]Methods{
	"users/{user}": {
		"GET":   GetUser,
		"PATCH": UpdateUserProfile,
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
	"rooms/{room}/messages/{message}": {
		"GET":    GetRoomMessage,
		"PATCH":  UpdateRoomMessage,
		"DELETE": DeleteRoomMessage,
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

	addr := os.Getenv("HTTP_ADDRESS") + ":" + os.Getenv("HTTP_PORT")
	log.Println("HTTP server listening on", color.CyanString(addr))
	err := http.ListenAndServe(addr, httpServer)
	if err != nil {
		log.Panicf("HTTP server fatal error: %s", err.Error())
	}
}
