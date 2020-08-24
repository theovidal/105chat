package ws

import "golang.org/x/net/websocket"

// Pipeline is the main channel where events are dispatched to connected users
var Pipeline = make(chan Event)

// HandlePipeline gets all the events sent in the pipeline and broadcasts them to users
func HandlePipeline() {
	for {
		event := <-Pipeline
		var toClose []uint

		station.RLock()
		for id, client := range station.clients {
			err := websocket.JSON.Send(client, event)
			if err != nil {
				client.Close()
				toClose = append(toClose, id)
			}
		}
		station.RUnlock()

		if len(toClose) != 0 {
			station.Lock()
			for _, id := range toClose {
				delete(station.clients, id)
			}
			station.Unlock()
		}
	}
}
