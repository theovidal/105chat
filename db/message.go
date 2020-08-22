package db

type Message struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	RoomID    uint   `json:"room_id"`
	UserID    uint   `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content" gorm:"primary_key"`
}
