package handlers

import (
	"chatapp/internal/auth"
	"chatapp/internal/chat"
	"chatapp/internal/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

// HTTP handler called when a client first logs on, gets the id and username for an active peer
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := auth.GetClaimFromAccessCookie("id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	username, err := postgres.GetUsernameById(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var userInfo = chat.UserItem{
		ID:       id,
		Username: username,
	}
	// write id and username to client
	json.NewEncoder(w).Encode(userInfo)
}

// HTTP handler, attempts to change usernames for a client
func UpdateUsernameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := auth.GetClaimFromAccessCookie("id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// checks if new username has valid syntax
	username, err := parseAndValidateUsername(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// check postgres if username already exists and will cause a collision error
	exists, err := postgres.UsernameExists(username)
	if err != nil {
		http.Error(w, "Database error checking for username existence", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}
	// update the username in postgres
	if err := postgres.UpdateUsername(id, username); err != nil {
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// used to verify a valid username syntatically (no DB collision check here)
func parseAndValidateUsername(r *http.Request) (string, error) {
	var payload struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return "", errors.New("Invalid request payload.")
	}
	// new username with leading/trailing whitespace removed
	username := strings.TrimSpace(payload.Username)
	// password length condition
	if len(username) < 3 || len(username) > 20 {
		return "", errors.New("Username must be between 3 and 20 characters.")
	}
	return username, nil
}

// generate random postfix and let users update usernames after account creation
func GenerateRandomUsername() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return "user_" + string(bytes)
}

// create a random username when creating an account, can change later
func CreateUniqueUsername() (string, error) {
	for i := 0; i < 10; i++ {
		username := GenerateRandomUsername()
		exists, err := postgres.UsernameExists(username)
		if err != nil { // failed to query db
			return "", err
		}
		if !exists { // if it's not used, then return
			return username, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique username")
}
