package db

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User represents a user of 105chat, who communicates to others via messages in rooms.
type User struct {
	// Identifier of the user, Twitter snowflake
	ID uint `json:"id"          gorm:"primary_key"`
	// Name of the user, 2~32 characters
	Name string `json:"name"        gorm:"size:32"     valid:"length(2|32)"`
	// URL pointing to user's avatar
	AvatarURL string `json:"avatar_url"                     valid:"url"`
	// Description of the user, 0~512 characters
	Description string `json:"description" gorm:"size:512"    valid:"length(0|512)"`
	// Color of the user (stored as an integer for less complexity)
	Color uint `json:"color"                          valid:"range(0|16777215)"`
	// When the user was created (via registration or administrator action)
	Timestamp int64 `json:"timestamp"`

	// Email of the user, used to communicate and authenticate
	Email string `json:"-" gorm:"unique" valid:"email"`
	// Password of the user, used to authenticate
	Password string `json:"-"`
	// Token of the user, used to interact with the API (WS and HTTP)
	// Is obtained after user's login via email and password
	Token string `json:"-"`
}

// FindUserFromRequest look at request's headers for a user token
func FindUserFromRequest(r *http.Request) (user User, err error) {
	token := r.Header.Get("Authentication")
	user, err = FindUserByToken(token)
	return
}

// FindRoomFromURL parses request's URL and find the corresponding user thanks to the ID
func FindUserFromURL(r *http.Request) (*User, error) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user"])
	if err != nil {
		return &User{}, InvalidType
	}

	var user User
	err = Database.First(&user, userID).Error
	if err != nil {
		return &User{}, UnknownUser
	}

	return &user, nil
}

// FindUserByToken searches for a user with a specific token
func FindUserByToken(token string) (user User, err error) {
	err = Database.Where("token = ?", token).First(&user).Error
	return
}
