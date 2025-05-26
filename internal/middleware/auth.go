package middleware

import (
	"chatapp/internal/auth"
	"net/http"
)

// authenticate short term session token on cookies
func AuthenticateSession(next http.Handler) http.Handler {
	return Authenticate(next, true)
}

// authenticate longer term tokens in query params for activation email or password resets
func AuthenticateParamToken(next http.Handler) http.Handler {
	return Authenticate(next, false)
}

// authenticate JWTs on requests
func Authenticate(next http.Handler, isSession bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.VerifyToken(r, isSession)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
