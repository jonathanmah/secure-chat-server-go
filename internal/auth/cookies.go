package auth

import (
	"chatapp/internal/config"
	"fmt"
	"net/http"
)

func SetNewSessionCookies(id string, w http.ResponseWriter) error {
	accessToken, refreshToken, err := CreateSessionTokens(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	// Set access tokens and refresh tokens in HTTP-only cookie
	http.SetCookie(w, config.NewAccessCookie(accessToken))
	http.SetCookie(w, config.NewRefreshCookie(refreshToken))
	return nil
}

func ExpireSessionCookies(w http.ResponseWriter) {
	http.SetCookie(w, config.ExpiredAccessCookie())
	http.SetCookie(w, config.ExpiredRefreshCookie())
}

func GetTokenFromCookie(cookieName string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("Missing %s cookie.", cookieName)
	}
	return cookie.Value, nil
}
