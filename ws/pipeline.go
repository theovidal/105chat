package ws

// Pipeline is the main channel where events are dispatched to connected users
var Pipeline = make(chan Event)

// HandlePipeline gets all the events sent in the pipeline and broadcasts them to users
func HandlePipeline() {
	for {
		event := <-Pipeline
		for _, client := range Station.clients {
			client.Pipeline <- &event
		}
	}
}
