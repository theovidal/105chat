package db

// Message is the entity that holds users' communications inside 105chat
type Message struct {
	// Identifier of the message, Twitter snowflake
	ID uint `json:"id"            gorm:"primary_key"`
	// Identifier of the room where the message was sent, Twitter snowflake
	RoomID uint `json:"room_id"`
	// Identifier of the user who sent the message, Twitter snowflake
	UserID uint `json:"user_id"`
	// Actual message content, 1~2000 chars
	Content string `json:"content"       gorm:"size:2000,not null" valid:"required,length(1|2000)"`
	// Whether the message is an announcement (is more visible in the room)
	Announcement bool `json:"announcement"`
	// Whether the message is private or not (direct message to another user)
	Private bool `json:"private"`
	// When the message was sent
	Timestamp int64 `json:"timestamp"`
}
