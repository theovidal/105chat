package http

import (
	"net/http"

	"github.com/theovidal/105chat/db"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := db.FindUserFromRequest(r); err == nil {
			next.ServeHTTP(w, r)
		} else {
			Response(w, http.StatusUnauthorized, nil)
		}
	})
}
