package controllers

import (
	"net/http"

	"github.com/theovidal/105chat/db"
)

// FindUserFromRequest look at request's headers for a user token
func FindUserFromRequest(r *http.Request) (user db.User, err error) {
	token := r.Header.Get("Authentication")
	user, err = db.FindUserByToken(token)
	return
}

// FindRoomFromURL parses request's URL and find the corresponding user thanks to the ID
func FindUserFromURL(r *http.Request) (*db.User, error) {
	userID, err := FindIDFromURL(r, "user")
	if err != nil {
		return &db.User{}, err
	}

	var user db.User
	if err = db.Database.First(&user, userID).Error; err != nil {
		return &db.User{}, UnknownUser
	}

	return &user, nil
}
