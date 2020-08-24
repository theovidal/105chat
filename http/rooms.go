package http

import (
	"errors"
	"net/http"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
)

// GetRooms returns all the rooms the user has access to
func GetRooms(w http.ResponseWriter, _ *http.Request) {
	var rooms []db.Room
	db.Database.Find(&rooms)
	Response(w, http.StatusOK, rooms)
}

// GetRoom returns data about a specific room thanks to its ID
func GetRoom(w http.ResponseWriter, r *http.Request) {
	room, err := ParseRoomFromURL(&w, r)
	if err != nil {
		return
	}

	Response(w, http.StatusOK, room)
}

// ParseRoomFromURL checks for errors in the passed room ID inside request's URL
func ParseRoomFromURL(w *http.ResponseWriter, r *http.Request) (room *db.Room, err error) {
	room, err = controllers.FindRoomFromURL(r)
	if errors.Is(err, controllers.InvalidType) {
		Response(*w, http.StatusBadRequest, nil)
	} else if errors.Is(err, controllers.UnknownRoom) {
		Response(*w, http.StatusNotFound, nil)
	}
	return
}
