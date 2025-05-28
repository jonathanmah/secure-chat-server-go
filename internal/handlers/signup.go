package handlers

import (
	"chatapp/internal/auth"
	"chatapp/internal/postgres"
	"errors"
	"log"
	"net/http"
	"net/mail"
)

// HTTP handler for creating an account
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	// verify email / password
	email := r.FormValue("email")
	password := r.FormValue("password")
	if err := validateCredentials(email, password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// create new user or update users password in DB (may have been created with OAuth earlier without a password)
	if err := createOrUpdateUser(email, password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// send confirmation email required to activate account
	if err := sendConfirmation(email); err != nil {
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Signup successful. Please check your email to confirm your account."))
}

// verify email and password user gave
func validateCredentials(email, password string) error {
	if email == "" || password == "" {
		return errors.New("Email and password are required.")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("Invalid email format.")
	}
	return nil
}

// create a new user, or create a password for an existing one
func createOrUpdateUser(email, password string) error {
	emailExists, err := postgres.EmailExists(email)
	if err != nil {
		return errors.New("Failed to check for existing email in database.")
	}
	// // hash the plaintext password
	hashedPassword, err := GetHashedPassword(password)
	if err != nil {
		return err
	}
	// if email doesn't exist, create new user
	if !emailExists {
		if err := createNewUser(email, hashedPassword); err != nil {
			return err
		}
	} else {
		if err := createUserPassword(email, hashedPassword); err != nil {
			return err
		}
	}
	return nil
}

// create a password for user, or error if user already registered with a password
func createUserPassword(email, hashedPassword string) error {
	passwordExists, err := postgres.PasswordExists(email)
	if err != nil {
		return errors.New("Failed to check for existing password in database.")
	}
	if passwordExists { // if email and password are both already created, give error
		isActivated, err := postgres.IsActivated(email)
		if err != nil {
			return errors.New("Failed to check if account is activated.")
		}
		if !isActivated {
			return errors.New("Pending account activation. Please confirm email to activate account.")
		}
		return errors.New("Email is already registered. Please sign in or reset password.")
	} else {
		// edge case where a user has created an account with OAuth but not registered with password
		if err := postgres.UpdatePassword(email, hashedPassword); err != nil {
			return errors.New("Failed to update password.")
		}
	}
	return nil
}

// create a new user in database
func createNewUser(email, hashedPassword string) error {
	username, err := CreateUniqueUsername()
	if err != nil {
		return errors.New("Failed to create unique username.")
	}
	if err := postgres.CreateUser(email, hashedPassword, username); err != nil {
		return errors.New("Failed to create user.")
	}
	return nil
}

// send confirmation email required to activate account
func sendConfirmation(email string) error {
	token, err := auth.CreateActivationToken(email)
	if err != nil {
		return err
	}
	return auth.SendConfirmationEmail(email, token)
}

// verify token in query params from email and activate a new user
func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	email, err := auth.GetClaimFromQueryParams("email", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = postgres.ActivateUser(email)
	if err != nil {
		http.Error(w, "Failed to confirm user", http.StatusInternalServerError)
	}
	log.Println("Account activated for: ", email)
	w.Header().Set("Cache-Control", "no-store") // prevent browser caching
	http.Redirect(w, r, "/?msg=Email+confirmed!+Account is activated.", http.StatusSeeOther)
}
