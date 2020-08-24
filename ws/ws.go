package ws

import (
	"sync"

	"golang.org/x/net/websocket"
)

// clients stores connections to the WebSocket server
// In the future, would be replaced by a concurrent-specific map : https://github.com/streamrail/concurrent-map
var station = struct {
	sync.RWMutex
	clients map[uint]*websocket.Conn
}{clients: make(map[uint]*websocket.Conn)}

// H is a shortcut to easily create a JSON map
type H map[string]interface{}
