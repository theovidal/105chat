package ws

import (
	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
)

// Server to client events
const (
	CONNECT = "CONNECT"
	PING    = "PING"
)

// ClientEvent is a shortcut for a handler function
type ClientEvent func(*websocket.Conn, *db.User, *Event)

// clientEvents is the list of available events the client can send to the server
var clientEvents = map[string]ClientEvent{
	PING: Ping,
}

// Ping is used to check the connection with the server
func Ping(ws *websocket.Conn, _ *db.User, _ *Event) {
	_ = websocket.JSON.Send(ws, Event{
		Event: PONG,
	})
}
