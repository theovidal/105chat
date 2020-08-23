package http

// MessageCreatePayload is sent to create a message in a room
type MessageCreatePayload struct {
	Content      string
	Announcement bool
}
