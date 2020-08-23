package db

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID          uint   `json:"id"          gorm:"primary_key"`
	Name        string `json:"name"        gorm:"size:32"     valid:"length(0|32)"`
	AvatarURL   string `json:"avatar_url"                     valid:"url"`
	Description string `json:"description" gorm:"size:512"    valid:"length(0|512)"`
	Color       uint   `json:"color"                          valid:"range(0|16777215)"`
	Timestamp   int64  `json:"timestamp"`

	Email    string `json:"-" gorm:"unique" valid:"email"`
	Password string `json:"-"`
	Token    string `json:"-"`
}

func FindUserFromRequest(r *http.Request) (user User, err error) {
	token := r.Header.Get("Authentication")
	user, err = FindUserByToken(token)
	return
}

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

func FindUserByToken(token string) (user User, err error) {
	err = Database.Where("token = ?", token).First(&user).Error
	return
}
