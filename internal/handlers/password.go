package handlers

import (
	"chatapp/internal/auth"
	"chatapp/internal/postgres"
	"errors"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

// HTTP handler when a user submits a request to reset password
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	// verify email / password
	email := r.FormValue("email")
	if err := ValidateEmail(email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// send confirmation email required to activate account
	if err := sendResetLink(email); err != nil {
		http.Error(w, "Failed to send password reset email.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Password reset link has been sent."))
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("Invalid email format.")
	}
	emailExists, err := postgres.EmailExists(email)
	if err != nil {
		return errors.New("Failed to verify if email exists.")
	}
	if !emailExists {
		return errors.New("No account has been created with this email.")
	}
	return nil
}

// send confirmation email required to activate account
func sendResetLink(email string) error {
	token, err := auth.CreateActivationToken(email)
	if err != nil {
		return err
	}
	return auth.SendPasswordResetEmail(email, token)
}

// HTTP handler when a user submits a request to reset password
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	token := r.FormValue("token")
	newPassword := r.FormValue("password")

	email, err := auth.GetClaimFromActivationToken("email", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := GetHashedPassword(newPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = postgres.UpdatePassword(email, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to update password.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Password successfully updated."))
}

func GetHashedPassword(password string) (string, error) {
	// hash the plaintext password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("Failed to hash password.")
	}
	return string(hashedPassword), err
}
