package auth

import (
	"chatapp/internal/config"
	"fmt"
	"log"
	"net/smtp"
)

// send email to activate a new account
func SendConfirmationEmail(email, token string) error {
	link := fmt.Sprintf("%s?token=%s", config.App.Email.ConfirmEmailURL, token)
	subject := "Confirm your ChatApp account"
	body := fmt.Sprintf("Click the following link to confirm your account:\r\n%s\r\n", link)
	log.Printf("Sending confirmation email to %s", email)
	return sendEmail(email, subject, body)
}

// send email to reset a users password
func SendPasswordResetEmail(email, token string) error {
	link := fmt.Sprintf("%s?token=%s", config.App.Email.ResetPasswordURL, token)
	subject := "Reset your ChatApp password"
	body := fmt.Sprintf("We received a request to reset your password.\r\nClick the following link to reset it:\r\n%s\r\nIf you did not request this, you can safely ignore this email.\r\n", link)
	log.Printf("Sending password reset email to %s", email)
	return sendEmail(email, subject, body)
}

func sendEmail(to, subject, body string) error {
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s",
		to, subject, body,
	))

	smtpAddr := fmt.Sprintf("%s:%s", config.App.Email.SMTPHost, config.App.Email.SMTPPort)
	auth := smtp.PlainAuth(
		"",
		config.App.Email.FromAddress,
		config.App.Email.SMTPPassword,
		config.App.Email.SMTPHost,
	)

	err := smtp.SendMail(smtpAddr, auth, config.App.Email.FromAddress, []string{to}, msg)
	return err
}
