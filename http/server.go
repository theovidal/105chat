package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Methods map[string]Operation
type Operation func(w http.ResponseWriter, r *http.Request)

var operations = map[string]Methods{
	"/http/room/{room}/messages": {
		"GET":  ReadMessages,
		"POST": CreateMessage,
	},
}

func Server() {
	httpServer := mux.NewRouter().StrictSlash(true)

	for path, methods := range operations {
		for method, handler := range methods {
			httpServer.HandleFunc(path, handler).Methods(method)
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

func Response(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
