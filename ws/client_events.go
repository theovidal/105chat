package ws

import (
	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
)

// Client to API events
const (
	CONNECT = "CONNECT"
	PING    = "PING"
)

type ClientEvent func(*websocket.Conn, *db.User, *Event)

var clientEvents = map[string]ClientEvent{
	PING: Ping,
}

func Ping(ws *websocket.Conn, _ *db.User, _ *Event) {
	websocket.JSON.Send(ws, Event{
		Event: PONG,
	})
}
