package http

import (
	"errors"
	"net/http"

	"github.com/theovidal/105chat/db"
)

func GetRooms(w http.ResponseWriter, _ *http.Request) {
	var rooms []db.Room
	db.Database.Find(&rooms)
	Response(w, http.StatusOK, rooms)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	room, err := ParseRoomFromURL(&w, r)
	if err != nil {
		return
	}

	Response(w, http.StatusOK, room)
}

func ParseRoomFromURL(w *http.ResponseWriter, r *http.Request) (room *db.Room, err error) {
	room, err = db.FindRoomFromURL(r)
	if errors.Is(err, db.InvalidType) {
		Response(*w, http.StatusBadRequest, nil)
	} else if errors.Is(err, db.UnknownRoom) {
		Response(*w, http.StatusNotFound, nil)
	}
	return
}
