package http

import (
	"errors"
	"net/http"

	"github.com/theovidal/105chat/db"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := ParseUserFromURL(&w, r)
	if err != nil {
		return
	}

	Response(w, http.StatusOK, user)
}

func ParseUserFromURL(w *http.ResponseWriter, r *http.Request) (user *db.User, err error) {
	user, err = db.FindUserFromURL(r)
	if errors.Is(err, db.InvalidType) {
		Response(*w, http.StatusBadRequest, nil)
	} else if errors.Is(err, db.UnknownUser) {
		Response(*w, http.StatusNotFound, nil)
	}
	return
}
