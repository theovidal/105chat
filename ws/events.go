package ws

import (
	"fmt"

	"github.com/theovidal/105chat/models"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type MessageCreateEvent struct {
	Author  *models.User `json:"author"`
	Room    int          `json:"room"`
	Content string       `json:"content"`
}

func (m *MessageCreateEvent) String() string {
	return fmt.Sprint(m.Author.Name, " : ", m.Content)
}
