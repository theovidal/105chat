package db

type Message struct {
	ID           uint   `json:"id"            gorm:"primary_key"`
	RoomID       uint   `json:"room_id"`
	UserID       uint   `json:"user_id"`
	Content      string `json:"content"       gorm:"size:2000,not null" valid:"required,length(0|2000)"`
	Announcement bool   `json:"announcement"`
	Private      bool   `json:"private"`
	Timestamp    int64  `json:"timestamp"`
}
