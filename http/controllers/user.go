package controllers

import (
	"net/http"

	"github.com/theovidal/105chat/db"
)

// FindRoomFromURL parses request's URL and find the corresponding user thanks to the ID
func FindUserFromURL(r *http.Request) (*db.User, error) {
	userID, err := FindIDFromURL(r, "user")
	if err != nil {
		return &db.User{}, err
	}

	var user db.User
	if err = db.Client.First(&user, userID).Error; err != nil {
		return &db.User{}, UnknownUser
	}

	return &user, nil
}
