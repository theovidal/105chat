package db

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Room struct {
	ID          uint   `json:"id"          gorm:"primary_key"`
	Name        string `json:"name"        gorm:"size:32"`
	AvatarURL   string `json:"avatar_url"                     valid:"url"`
	Description string `json:"description" gorm:"size:512"`
	Color       uint   `json:"color"                          valid:"range(0|16777215)"`
	Timestamp   int64  `json:"timestamp"`
}

func FindRoomFromURL(r *http.Request) (*Room, error) {
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["room"])
	if err != nil {
		return &Room{}, InvalidType
	}

	var room Room
	err = Database.First(&room, roomID).Error
	if err != nil {
		return &Room{}, UnknownRoom
	}

	return &room, nil
}
