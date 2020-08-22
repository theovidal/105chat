package db

type Room struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       uint   `json:"color"`
}

func FindRoom(id int) (room Room, err error) {
	err = Database.First(&room, id).Error
	return
}
