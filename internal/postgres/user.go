package postgres

import (
	"log"

	_ "github.com/lib/pq"
)

// create a new user without a password via OAuth
func CreatePasswordlessUser(email, username string) (string, error) {
	var id string
	err := DB.QueryRow(`INSERT INTO users (email, username) VALUES ($1, $2) RETURNING id`, email, username).Scan(&id)
	return id, err
}

// create a new user without a password via OAuth
func CreateUser(email, passwordHash, username string) error {
	_, err := DB.Exec(`INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)`, email, username, passwordHash)
	return err
}

func UpdatePassword(email, passwordHash string) error {
	_, err := DB.Exec(`UPDATE users SET password_hash = $1 WHERE email = $2`, passwordHash, email)
	return err
}

// update a username given userID and new username
func UpdateUsername(id string, username string) error {
	_, err := DB.Exec(`UPDATE users SET username = $1 WHERE id = $2`, username, id)
	return err
}

// get password hash and activation status
func GetUserCredentials(email string) (string, string, bool, error) {
	var id string
	var passwordHash string
	var isActive bool
	err := DB.QueryRow(`SELECT id, password_hash, is_active FROM users WHERE email = $1`, email).Scan(&id, &passwordHash, &isActive)
	return id, passwordHash, isActive, err
}

// activate a user after email confirmation
func ActivateUser(email string) error {
	_, err := DB.Exec(`UPDATE users SET is_active = true WHERE email = $1`, email)
	return err
}

func GetUsernameById(id string) (string, error) {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE id = $1", id).Scan(&username)
	return username, err
}

// get a username and ID from email
func GetUserIdByEmail(email string) (string, error) {
	var id string
	err := DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	return id, err
}

// return true if the user has a password with the account
func PasswordExists(email string) (bool, error) {
	var exists bool
	err := DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND password_hash IS NOT NULL)`, email).Scan(&exists)
	log.Println(err)
	return exists, err
}

// return true if username already in DB
func UsernameExists(username string) (bool, error) {
	var exists bool
	err := DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`,
		username,
	).Scan(&exists)
	return exists, err
}

// return true if username already in DB
func EmailExists(email string) (bool, error) {
	var exists bool
	err := DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exists)
	return exists, err
}

// return true if account is activated
func IsActivated(email string) (bool, error) {
	var isActive bool
	err := DB.QueryRow(
		`SELECT is_active FROM users WHERE email = $1`,
		email,
	).Scan(&isActive)
	return isActive, err
}
