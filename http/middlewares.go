package http

import (
	"context"
	"net/http"

	"github.com/theovidal/105chat/db"
)

// AuthenticationMiddleware checks that user is authenticated before doing a request
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, err := db.FindUserFromRequest(r); err == nil {
			userContext := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(userContext))
		} else {
			Response(w, http.StatusUnauthorized, nil)
		}
	})
}
