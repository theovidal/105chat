package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting 105chat...")

	go HandleBroadcasts()
	http.Handle("/ws", websocket.Handler(WebSocketServer))

	go HTTPServer()

	log.Println("WebSocket server ready")
	err := http.ListenAndServe("localhost:1051", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
