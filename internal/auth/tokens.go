package auth

import (
	"chatapp/internal/config"
	"chatapp/internal/postgres"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// create short term access token
func CreateAccessToken(id string) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.App.Auth.AccessTokenKey)
}

// create longer term refresh token with random string instead of JWT, stored in postgres
func CreateRefreshToken(id string) (string, error) {
	token, err := GenerateRandomString()
	if err != nil {
		return "", err
	}
	if err := postgres.CreateRefreshToken(id, token, time.Now().Add(time.Hour*24*14)); err != nil {
		return "", err
	}
	return token, nil
}

// create activation token used for email verification and password reset
func CreateActivationToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 2).Unix(), // 2 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.App.Auth.ActivationTokenKey)
}

// creates and returns new access and refresh tokens
func CreateSessionTokens(id string) (string, string, error) {
	accessToken, err := CreateAccessToken(id)
	if err != nil {
		return "", "", errors.New("Failed to generate access token")
	}
	refreshToken, err := CreateRefreshToken(id)
	if err != nil {
		return "", "", errors.New("Failed to generate refresh token")
	}
	return accessToken, refreshToken, err
}

// delete a refresh token on logout
func DeleteRefreshToken(id string) error {
	if err := postgres.DeleteRefreshToken(id); err != nil {
		return err
	}
	return nil
}

// create a random string for tokens or keys
func GenerateRandomString() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
