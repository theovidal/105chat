package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := FindUserFromRequest(r); err == nil {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func HTTPServer() {
	httpServer := mux.NewRouter().StrictSlash(true)
	httpServer.HandleFunc("/http/room/{room}/messages", CreateMessage)

	httpServer.Use(AuthenticationMiddleware)
	httpServer.Use(mux.CORSMethodMiddleware(httpServer))

	log.Println("HTTP server ready")
	err := http.ListenAndServe("localhost:1052", httpServer)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	user, _ := FindUserFromRequest(r)

	vars := mux.Vars(r)
	room, err := strconv.Atoi(vars["room"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var payload MessageCreatePayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	broadcast <- Event{
		Type: "CREATE_MESSAGE",
		Data: MessageCreateEvent{
			Author:  &user,
			Content: payload.Content,
			Room:    room,
		},
	}
}
