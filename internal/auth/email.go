package auth

import (
	"chatapp/internal/config"
	"fmt"
	"log"
	"net/smtp"
)

// SMTP util to send confirmation email
func SendConfirmationEmail(email, token string) error {
	// appPass was a second password used with 2FA in order to send emails. Not same as gmail password
	to := email
	link := fmt.Sprintf("%s?token=%s", config.App.Email.ConfirmEmailURL, token)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Confirm your ChatApp account\r\n" +
		"\r\n" +
		"Click the following link to confirm your account:\r\n" +
		link + "\r\n")

	log.Printf("Sent confirmation email to %s", to)
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", config.App.Email.FromAddress, config.App.Email.SMTPPass, "smtp.gmail.com"), config.App.Email.FromAddress, []string{to}, msg)
	return err
}

func SendPasswordResetEmail(email, token string) error {
	to := email

	resetLink := fmt.Sprintf("%s?token=%s", config.App.Email.ResetPasswordURL, token)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Reset your ChatApp password\r\n" +
		"\r\n" +
		"We received a request to reset your password.\r\n" +
		"Click the following link to reset it:\r\n" +
		resetLink + "\r\n" +
		"If you did not request this, you can safely ignore this email.\r\n")
	log.Printf("Sent password reset email to %s", to)
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", config.App.Email.FromAddress, config.App.Email.SMTPPass, "smtp.gmail.com"), config.App.Email.FromAddress, []string{to}, msg)
	return err
}
