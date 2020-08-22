package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/theovidal/105chat/models"
	"github.com/theovidal/105chat/ws"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := models.FindUserFromRequest(r); err == nil {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func Server() {
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
	user, _ := models.FindUserFromRequest(r)

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

	ws.Pipeline <- ws.Event{
		Type: "CREATE_MESSAGE",
		Data: ws.MessageCreateEvent{
			Author:  &user,
			Content: payload.Content,
			Room:    room,
		},
	}
}
