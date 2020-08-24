package ws

import "golang.org/x/net/websocket"

// clients stores connections to the WebSocket server
var clients = make(map[uint]*websocket.Conn)

// H is a shortcut to easily create a JSON map
type H map[string]interface{}
