package auth

import (
	"chatapp/internal/config"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// header is signing algorithm
// claims are just JWT payload, data to store
// signature is a hash of header+claims with a key
func CreateSessionToken(id string) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.App.Auth.SessionKey)
}

func CreateActivationToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 48).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.App.Auth.ActivationKey)
}

// used to verify signature with JWT Parse(), returns the secret key used to sign JWT
func KeyFuncSession(token *jwt.Token) (interface{}, error) {
	return config.App.Auth.SessionKey, nil
}

func KeyFuncActivation(token *jwt.Token) (interface{}, error) {
	return config.App.Auth.ActivationKey, nil
}

// verify JWT token, keyFunc tells Parse() which key it should use
func VerifyTokenString(tokenString string, keyFunc func(*jwt.Token) (interface{}, error)) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil || !token.Valid {
		return nil, errors.New("Invalid or expired token.")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid token claims.")
	}
	return claims, nil
}

// verify JWT in HTTP only cookie, used for access to protected routes
func VerifySessionToken(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Missing JWT cookie in request from address: %s", r.RemoteAddr)
		return nil, errors.New("Missing session cookie.")
	}
	tokenString := cookie.Value
	return VerifyTokenString(tokenString, KeyFuncSession)
}

// verify JWT in query parameter, used for activation emails / password resets
func VerifyQueryParamToken(r *http.Request) (jwt.MapClaims, error) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		return nil, errors.New("Missing JWT in query parameters.")
	}
	return VerifyTokenString(tokenString, KeyFuncActivation)
}

func VerifyToken(r *http.Request, isSession bool) (jwt.MapClaims, error) {
	if isSession {
		return VerifySessionToken(r)
	} else {
		return VerifyQueryParamToken(r)
	}
}
