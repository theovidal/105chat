package db

import (
	"net/http"
)

type User struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Timestamp   int64  `json:"timestamp"`
	Description string `json:"description"`

	Email    string `json:"-"`
	Password string `json:"-"`
	Token    string `json:"-"`
}

func FindUserFromRequest(r *http.Request) (user User, err error) {
	token := r.Header.Get("Authentication")
	user, err = FindUserByToken(token)
	return
}

func FindUserByToken(token string) (user User, err error) {
	err = Database.Where("token = ?", token).First(&user).Error
	return
}
