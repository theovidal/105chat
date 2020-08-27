package http

import (
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
	"github.com/theovidal/105chat/utils"
	"github.com/theovidal/105chat/ws"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	if user := r.Context().Value("user").(db.User); !user.HasGlobalPermission(db.MANAGE_ROOM) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var payload RoomPayload
	if err := ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	payload.Name = govalidator.Trim(payload.Name, "")
	payload.Description = govalidator.Trim(payload.Description, "")
	room := db.Room{
		ID:          utils.GenerateSnowflake(),
		Name:        payload.Name,
		AvatarURL:   payload.AvatarURL,
		Description: payload.Description,
		Color:       payload.Color,
		Timestamp:   utils.Now(),
	}
	err := db.Database.Create(&room).Error
	if err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	ws.Pipeline <- ws.Event{
		Event: ws.ROOM_CREATE,
		Data:  &room,
	}
	Response(w, http.StatusCreated, &room)
}

// GetRooms returns all the rooms the user has access to
func GetRooms(w http.ResponseWriter, r *http.Request) {
	var rooms []db.Room

	if user := r.Context().Value("user").(db.User); user.HasGlobalPermission(db.READ_MESSAGES) {
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

	if user := r.Context().Value("user").(db.User); !user.HasAnyPermission(room.ID, db.READ_MESSAGES) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	Response(w, http.StatusOK, room)
}

func UpdateRoom(w http.ResponseWriter, r *http.Request) {
	room, err := ParseRoomFromURL(&w, r)
	if err != nil {
		return
	}

	if user := r.Context().Value("user").(db.User); !user.HasAnyPermission(room.ID, db.MANAGE_ROOM) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var payload RoomPayload
	if err := ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	payload.Name = govalidator.Trim(payload.Name, "")
	payload.Description = govalidator.Trim(payload.Description, "")
	if err = db.Database.Model(room).Updates(payload).Error; err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	ws.Pipeline <- ws.Event{
		Event: ws.ROOM_UPDATE,
		Data:  &room,
	}
	Response(w, http.StatusOK, &room)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	room, err := ParseRoomFromURL(&w, r)
	if err != nil {
		return
	}

	if user := r.Context().Value("user").(db.User); !user.HasAnyPermission(room.ID, db.MANAGE_ROOM) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	remaining := utils.H{
		"id": room.ID,
	}
	ws.Pipeline <- ws.Event{
		Event: ws.ROOM_DELETE,
		Data:  &remaining,
	}
	Response(w, http.StatusAccepted, &remaining)

	db.Database.Delete(&room)
	db.Database.Where("room_id = ?", room.ID).Delete(&db.Message{})
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
