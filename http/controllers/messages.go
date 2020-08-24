package controllers

import (
	"net/http"

	"github.com/theovidal/105chat/db"
)

// FindMessageFromURL parses request's URL and find the corresponding message thanks to the ID
func FindMessageFromURL(r *http.Request) (*db.Message, error) {
	roomID, err := FindIDFromURL(r, "room")
	if err != nil {
		return &db.Message{}, err
	}
	messageID, err := FindIDFromURL(r, "message")
	if err != nil {
		return &db.Message{}, err
	}

	var message db.Message
	err = db.Database.Where("id = ? AND room_id = ?", messageID, roomID).Find(&message).Error
	if err != nil {
		return &db.Message{}, UnknownMessage
	}

	return &message, nil
}
