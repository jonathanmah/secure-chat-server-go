package config

import (
	"net/http"
	"time"
)

// used with short term access tokens for sessions
func NewSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set to true later for https
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hr
	}
}

// used to nullify a session cookie when a user logs out
func NewExpiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

// add oauth random state on cookie
func NewOAuthStateCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
