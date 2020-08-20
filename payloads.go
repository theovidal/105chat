package main

type Payload struct {
	Type string
	Data map[string]interface{}
}

type MessageCreatePayload struct {
	Room    int
	Content string
}

func ParseMessageCreatePayload(data map[string]interface{}) (MessageCreatePayload, error) {
	room := data["room"].(float64)
	content := data["content"].(string)
	return MessageCreatePayload{
		Room:    int(room),
		Content: content,
	}, nil
}
