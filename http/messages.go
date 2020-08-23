package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/theovidal/105chat/db"
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
		Content:      payload.Content,
		Announcement: payload.Announcement,
		Private:      false,
		Timestamp:    time.Now().Unix(),
	}

	if err = db.Database.Create(&message).Error; err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	ws.Pipeline <- ws.Event{
		Type: "MESSAGE_CREATE",
		Data: message,
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
