package ws

import "golang.org/x/net/websocket"

// Pipeline is the main channel where events are dispatched to connected users
var Pipeline = make(chan Event)

// HandlePipeline gets all the events sent in the pipeline and broadcasts them to users
func HandlePipeline() {
	for {
		event := <-Pipeline
		for id, client := range clients {
			err := websocket.JSON.Send(client, event)
			if err != nil {
				client.Close()
				delete(clients, id)
			}
		}
	}
}
