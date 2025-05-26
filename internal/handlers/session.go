package handlers

import (
	"chatapp/internal/auth"
	"chatapp/internal/config"
	"chatapp/internal/postgres"
	"net/http"
	"net/mail"

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
	// create JWT
	sessionToken, err := auth.CreateSessionToken(id)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	// Set token in HTTP-only cookie
	http.SetCookie(w, config.NewSessionCookie(sessionToken))
	w.WriteHeader(http.StatusOK)
}

// set JWT on cookie to be expired/invalid after a logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, config.NewExpiredSessionCookie())
	w.WriteHeader(http.StatusOK)
}
