package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
	"github.com/theovidal/105chat/ws"
)

// CreateMessage sends a message from a user in a room
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	room, err := ParseRoomFromURL(&w, r)
	if err != nil {
		return
	}

	var payload MessageCreatePayload
	if err = ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	message := db.Message{
		ID:           uint(node.Generate()),
		RoomID:       room.ID,
		UserID:       user.ID,
		Content:      govalidator.Trim(payload.Content, ""),
		Announcement: payload.Announcement,
		Private:      false,
		Timestamp:    time.Now().Unix(),
	}

	if err = db.Database.Create(&message).Error; err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	ws.Pipeline <- ws.Event{
		Type: ws.MESSAGE_CREATE,
		Data: &message,
	}

	Response(w, http.StatusNoContent, nil)
}

// GetRoomMessages returns up to 25 messages in a specific room
func GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	room, err := ParseRoomFromURL(&w, r)
	if err != nil {
		return
	}

	var messages []db.Message
	query := fmt.Sprintf("room_id = %d", room.ID)
	if before := r.URL.Query().Get("before"); before != "" {
		query += fmt.Sprintf(" AND id < %s", before)
	}
	if after := r.URL.Query().Get("after"); after != "" {
		query += fmt.Sprintf(" AND id > %s", after)
	}
	db.Database.Order("id desc").Where(query).Limit(25).Find(&messages)

	Response(w, http.StatusOK, messages)
}

// GetRoomMessage returns data about a specific message in a room
func GetRoomMessage(w http.ResponseWriter, r *http.Request) {
	message, err := ParseMessageFromURL(&w, r)
	if err != nil {
		return
	}
	Response(w, http.StatusOK, message)
}

// UpdateRoomMessage is used by a user to edit one of their message
func UpdateRoomMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	message, err := ParseMessageFromURL(&w, r)
	if err != nil {
		return
	}

	if message.UserID != user.ID {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var payload MessageUpdatePayload
	if err = ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	payload.Content = govalidator.Trim(payload.Content, "")
	err = db.Database.Model(message).Updates(payload).Error
	if err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	ws.Pipeline <- ws.Event{
		Type: ws.MESSAGE_UPDATE,
		Data: &message,
	}
	Response(w, http.StatusNoContent, nil)
}

// UpdateRoomMessage is used by a user to delete one of their message
func DeleteRoomMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(db.User)
	message, err := ParseMessageFromURL(&w, r)
	if err != nil {
		return
	}

	if message.UserID != user.ID {
		Response(w, http.StatusForbidden, nil)
		return
	}

	db.Database.Delete(message)

	ws.Pipeline <- ws.Event{
		Type: ws.MESSAGE_DELETE,
		Data: ws.H{
			"id":      message.ID,
			"room_id": message.RoomID,
		},
	}
	Response(w, http.StatusNoContent, nil)
}

// ParseMessageFromURL checks for errors in the passed message ID inside request's URL
func ParseMessageFromURL(w *http.ResponseWriter, r *http.Request) (message *db.Message, err error) {
	message, err = controllers.FindMessageFromURL(r)
	if errors.Is(err, controllers.InvalidType) {
		Response(*w, http.StatusBadRequest, nil)
	} else if errors.Is(err, controllers.UnknownMessage) {
		Response(*w, http.StatusNotFound, nil)
	}
	return
}
