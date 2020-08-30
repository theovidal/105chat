package http

import (
	"context"
	"net/http"

	"github.com/theovidal/105chat/ws"
)

// AuthenticationMiddleware checks that user is authenticated before doing a request
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authentication")
		if r.Method == http.MethodOptions {
			Response(w, http.StatusOK, nil)
			return
		}

		// Ignoring for this endpoint - We need it to get our access token
		if r.URL.Path == "/v1/http/auth" {
			next.ServeHTTP(w, r)
		}

		token := r.Header.Get("Authentication")
		if user, found := ws.Station.GetUser(token); found {
			if user.Disabled {
				Response(w, http.StatusForbidden, nil)
			} else {
				userContext := context.WithValue(r.Context(), "user", user)
				next.ServeHTTP(w, r.WithContext(userContext))
			}
		} else {
			Response(w, http.StatusUnauthorized, nil)
		}
	})
}
