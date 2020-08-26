package http

import (
	"net/http"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
	"github.com/theovidal/105chat/ws"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload AuthenticatePayload
	if err := ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	var user db.User
	err := db.Database.Where("email = ?", payload.Email).Find(&user).Error
	if err != nil {
		Response(w, http.StatusUnauthorized, nil)
		return
	}

	match := controllers.ComparePasswords(payload.Password, user.Password)
	if match {
		w.Header().Set("Cache-Control", "private")
		Response(w, http.StatusOK, ws.H{
			"token": user.Token,
		})
	} else {
		Response(w, http.StatusUnauthorized, nil)
	}
}
