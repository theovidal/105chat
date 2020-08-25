package http

import (
	"errors"
	"net/http"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
)

// GetRooms returns all the rooms the user has access to
func GetRooms(w http.ResponseWriter, r *http.Request) {
	var rooms []db.Room

	user := r.Context().Value("user").(db.User)
	if user.HasGlobalPermission(db.READ_MESSAGES) {
		db.Database.Find(&rooms)
	} else {
		var accessibleRooms []uint
		for id := range user.Group.RoomPermissions {
			if user.HasRoomPermission(id, db.READ_MESSAGES) {
				accessibleRooms = append(accessibleRooms, id)
			}
		}
		db.Database.Where(accessibleRooms).Find(&rooms)
	}

	Response(w, http.StatusOK, rooms)
}

// GetRoom returns data about a specific room thanks to its ID
func GetRoom(w http.ResponseWriter, r *http.Request) {
	room, err := ParseRoomFromURL(&w, r)
	if err != nil {
		return
	}

	user := r.Context().Value("user").(db.User)
	if !user.HasAnyPermission(room.ID, db.READ_MESSAGES) {
		Response(w, http.StatusForbidden, nil)
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
