package ws

import "golang.org/x/net/websocket"

// clients stores connections to the WebSocket server
var clients = make(map[*websocket.Conn]bool)

// Event is the contained for all events that'll be sent inside the pipeline
type Event struct {
	// The type of the event, with format: "ENTITY_ACTION" (e.g: MESSAGE_CREATE)
	Type string `json:"type"`
	// The data related to the event (message, rooom, user...)
	Data interface{} `json:"data"`
}

// H is a shortcut to easily create a JSON map
type H map[string]interface{}
