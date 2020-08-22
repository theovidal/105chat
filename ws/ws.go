package ws

import "golang.org/x/net/websocket"

var clients = make(map[*websocket.Conn]bool)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type H map[string]interface{}
