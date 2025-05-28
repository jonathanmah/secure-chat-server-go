package auth

import (
	"chatapp/internal/config"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// used to Parse signature with JWT Parse(), returns the secret key used to sign JWT
func KeyFuncAccess(token *jwt.Token) (interface{}, error) {
	return config.App.Auth.AccessTokenKey, nil
}

func KeyFuncActivation(token *jwt.Token) (interface{}, error) {
	return config.App.Auth.ActivationTokenKey, nil
}

// Parses, validates JWT token, keyFunc tells Parse() which key it should use
func ParseToken(token string, keyFunc func(*jwt.Token) (interface{}, error)) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, keyFunc)
	if err != nil || !parsedToken.Valid {
		return nil, errors.New("Invalid or expired token.")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid token claims.")
	}
	return claims, nil
}

// Parse JWT in HTTP only cookie, used for access to protected routes
func ParseAccessCookie(r *http.Request) (jwt.MapClaims, error) {
	token, err := GetTokenFromCookie(config.AccessCookieName, r)
	if err != nil {
		return nil, err
	}
	return ParseToken(token, KeyFuncAccess)
}

// Parse JWT in query parameter, used for activation emails / password resets
func ParseQueryParams(r *http.Request) (jwt.MapClaims, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return nil, errors.New("Missing JWT in query parameters.")
	}
	return ParseToken(token, KeyFuncActivation)
}
