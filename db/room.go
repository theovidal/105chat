package db

// Room represents the place where users exchange messages between them
type Room struct {
	// Identifier of the room, Twitter snowflake
	ID uint `json:"id" gorm:"primary_key"`
	// Name of the room
	Name string `json:"name" gorm:"size:32,not null" valid:"required"`
	// URL pointing to room's avatar
	AvatarURL string `json:"avatar_url" valid:"url"`
	// Description of the room, 0~512
	Description string `json:"description" gorm:"size:512" valid:"length(0|512)"`
	// Color of the room (stored as an integer for less complexity)
	Color uint `json:"color" valid:"range(0|16777215)"`
	// When the room was created
	Timestamp int64 `json:"timestamp"`
}
