package http

import (
	"net/http"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
	"github.com/theovidal/105chat/utils"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload AuthenticatePayload
	if err := ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	var user db.User
	if err := db.Database.Where("email = ?", payload.Email).Find(&user).Error; err != nil {
		Response(w, http.StatusUnauthorized, nil)
		return
	}

	if controllers.ComparePasswords(payload.Password, user.Password) {
		w.Header().Set("Cache-Control", "private")
		Response(w, http.StatusOK, utils.H{
			"token": user.Token,
		})
	} else {
		Response(w, http.StatusUnauthorized, nil)
	}
}
