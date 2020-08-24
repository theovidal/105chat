package controllers

import (
	"net/http"

	"github.com/theovidal/105chat/db"
)

// FindRoomFromURL parses request's URL and find the corresponding room thanks to the ID
func FindRoomFromURL(r *http.Request) (*db.Room, error) {
	roomID, err := FindIDFromURL(r, "room")
	if err != nil {
		return &db.Room{}, err
	}

	var room db.Room
	err = db.Database.First(&room, roomID).Error
	if err != nil {
		return &db.Room{}, UnknownRoom
	}

	return &room, nil
}
