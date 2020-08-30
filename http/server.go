package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/gorilla/mux"

	"github.com/theovidal/105chat/utils"
)

// Methods is a shortcut that stores a list of operations
type Methods map[string]Operation

// Operation is a shortcut that stores a handler function
type Operation func(w http.ResponseWriter, r *http.Request)

// operations contains a list of available endpoints on the API with their method
var operations = map[string]Methods{
	"auth": {
		"POST": Authenticate,
	},
	"users": {
		"GET": GetUsers,
	},
	"users/{user}": {
		"GET":   GetUser,
		"PATCH": UpdateUserProfile,
		"PUT":   UpdateUser,
	},
	"users/{user}/group": {
		"GET":   GetUserGroup,
		"PATCH": UpdateUserGroup,
	},
	"groups": {
		"GET":  GetGroups,
		"POST": CreateGroup,
	},
	"groups/{group}": {
		"GET":    GetGroup,
		"PATCH":  UpdateGroup,
		"DELETE": DeleteGroup,
	},
	"rooms": {
		"GET":  GetRooms,
		"POST": CreateRoom,
	},
	"rooms/{room}": {
		"GET":    GetRoom,
		"PATCH":  UpdateRoom,
		"DELETE": DeleteRoom,
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
			httpServer.HandleFunc(fullPath, handler).Methods(method, http.MethodOptions)
		}
	}

	httpServer.Use(AuthenticationMiddleware)
	httpServer.Use(mux.CORSMethodMiddleware(httpServer))

	addr := utils.GenerateAddress("HTTP")
	log.Println("▶ HTTP server listening on", color.CyanString("http://"+addr))
	if err := http.ListenAndServe(addr, httpServer); err != nil {
		log.Panicf("‼ HTTP server fatal error: %s", err.Error())
	}
}
