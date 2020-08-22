package ws

import "golang.org/x/net/websocket"

var Pipeline = make(chan Event)

func HandlePipeline() {
	for {
		event := <-Pipeline
		for client := range clients {
			err := websocket.JSON.Send(client, event)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
