package config

import (
	"net/http"
	"time"
)

const (
	AccessCookieName     string = "access_token"
	RefreshCookieName    string = "refresh_token"
	OAuthStateCookieName string = "oauth_state"
)

// used with short term access tokens for sessions
func NewAccessCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     AccessCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set to true later for https
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60, // 15 mins
	}
}

// used to nullify a session cookie when a user logs out
func ExpiredAccessCookie() *http.Cookie {
	return &http.Cookie{
		Name:     AccessCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func NewRefreshCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshCookieName,
		Value:    token,
		Path:     "/auth/refresh", // limit to only used when refresh
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   14 * 24 * 60 * 60, // 	14 days
	}
}

func ExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     "/auth/refresh", // limit to only used when refresh
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
		Name:     OAuthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}
