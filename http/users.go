package http

import (
	"errors"
	"net/http"

	"github.com/theovidal/105chat/db"
)

// GetUser returns information about a specific user thanks to their ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := ParseUserFromURL(&w, r)
	if err != nil {
		return
	}

	Response(w, http.StatusOK, user)
}

// ParseUserFromURL checks for errors in the passed user ID inside request's URL
func ParseUserFromURL(w *http.ResponseWriter, r *http.Request) (user *db.User, err error) {
	user, err = db.FindUserFromURL(r)
	if errors.Is(err, db.InvalidType) {
		Response(*w, http.StatusBadRequest, nil)
	} else if errors.Is(err, db.UnknownUser) {
		Response(*w, http.StatusNotFound, nil)
	}
	return
}
