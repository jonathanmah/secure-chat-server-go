package handlers

import (
	"chatapp/internal/auth"
	"chatapp/internal/config"
	"chatapp/internal/postgres"
	"net/http"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// handler for login portal
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	// check email format
	_, err = mail.ParseAddress(email)
	if err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	// get the users hashed password from DB and isactive status
	id, storedHash, isActive, err := postgres.GetUserCredentials(email)

	if err != nil {
		http.Error(w, "Account not created with this email yet.", http.StatusUnauthorized)
		return
	}
	if !isActive {
		http.Error(w, "Please check your email to confirm and activate account.", http.StatusForbidden)
		return
	}
	// compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}
	if err := auth.SetNewSessionCookies(id, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// set JWT on cookie to be expired/invalid after a logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := auth.GetClaimFromAccessCookie("id", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = auth.DeleteRefreshToken(id)
	if err != nil {
		http.Error(w, "Failed to delete refresh token on logout", http.StatusInternalServerError)
		return
	}
	auth.ExpireSessionCookies(w)
	w.WriteHeader(http.StatusOK)
}

// creates new access and refresh token when access token expires
func RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	currRefreshToken, err := auth.GetTokenFromCookie(config.RefreshCookieName, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, expiresAt, err := postgres.GetRefreshTokenInfo(currRefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
	}
	if time.Now().After(expiresAt) {
		postgres.DeleteRefreshToken(id) // if its already expired just delete it
		http.Error(w, "Refresh token has expired", http.StatusUnauthorized)
		return
	}
	if err := auth.SetNewSessionCookies(id, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
