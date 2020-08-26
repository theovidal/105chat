package http

import (
	"context"
	"net/http"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
)

// AuthenticationMiddleware checks that user is authenticated before doing a request
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignoring for this endpoint - We need it to get our access token
		if r.URL.Path == "/v1/http/auth" {
			next.ServeHTTP(w, r)
		}

		if user, err := controllers.FindUserFromRequest(r); err == nil {
			if user.Disabled {
				Response(w, http.StatusForbidden, nil)
			} else {
				db.FetchPermissions(&user.Group, user.GroupID)
				userContext := context.WithValue(r.Context(), "user", user)
				next.ServeHTTP(w, r.WithContext(userContext))
			}
		} else {
			Response(w, http.StatusUnauthorized, nil)
		}
	})
}
