package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/ws"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	user, _ := db.FindUserFromRequest(r)

	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["room"])
	if err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	_, err = db.FindRoom(roomID)
	if err != nil {
		Response(w, http.StatusNotFound, nil)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var payload MessageCreatePayload
	err = json.Unmarshal(body, &payload)
	if err != nil || payload.Content == "" {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	message := db.Message{
		ID:        uint(node.Generate()),
		RoomID:    uint(roomID),
		UserID:    user.ID,
		Timestamp: time.Now().Unix(),
		Content:   payload.Content,
	}
	db.Database.Create(&message)

	ws.Pipeline <- ws.Event{
		Type: "CREATE_MESSAGE",
		Data: message,
	}

	Response(w, http.StatusNoContent, nil)
}

func ReadMessages(w http.ResponseWriter, _ *http.Request) {
	var messages []db.Message
	db.Database.Order("id desc").Limit(25).Find(&messages)
	Response(w, http.StatusOK, messages)
}
