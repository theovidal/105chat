package db

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Room represents the place where users exchange messages between them
type Room struct {
	// Identifier of the room, Twitter snowflake
	ID uint `json:"id"          gorm:"primary_key"`
	// Name of the room
	Name string `json:"name"        gorm:"size:32"`
	// URL pointing to room's avatar
	AvatarURL string `json:"avatar_url"                     valid:"url"`
	// Description of the room, 0~512
	Description string `json:"description" gorm:"size:512"    valid:"length(0|512)"`
	// Color of the room (stored as an integer for less complexity)
	Color uint `json:"color"                          valid:"range(0|16777215)"`
	// When the room was created
	Timestamp int64 `json:"timestamp"`
}

// FindRoomFromURL parses request's URL and find the corresponding room thanks to the ID
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
