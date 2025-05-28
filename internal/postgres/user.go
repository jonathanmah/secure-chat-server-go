package postgres

import (
	_ "github.com/lib/pq"
)

// create a new user without a password via OAuth
func CreatePasswordlessUser(email, username string) (id string, err error) {
	err = DB.QueryRow(`INSERT INTO users (email, username) VALUES ($1, $2) RETURNING id`, email, username).Scan(&id)
	return
}

// create a new user without a password via OAuth
func CreateUser(email, passwordHash, username string) (err error) {
	_, err = DB.Exec(`INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)`, email, username, passwordHash)
	return
}

func UpdatePassword(email, passwordHash string) (err error) {
	_, err = DB.Exec(`UPDATE users SET password_hash = $1 WHERE email = $2`, passwordHash, email)
	return
}

// update a username given userID and new username
func UpdateUsername(id string, username string) (err error) {
	_, err = DB.Exec(`UPDATE users SET username = $1 WHERE id = $2`, username, id)
	return
}

// get password hash and activation status
func GetUserCredentials(email string) (id string, passwordHash string, isActive bool, err error) {
	err = DB.QueryRow(`SELECT id, password_hash, is_active FROM users WHERE email = $1`, email).Scan(&id, &passwordHash, &isActive)
	return
}

// activate a user after email confirmation
func ActivateUser(email string) (err error) {
	_, err = DB.Exec(`UPDATE users SET is_active = true WHERE email = $1`, email)
	return
}

func GetUsernameById(id string) (string, error) {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE id = $1", id).Scan(&username)
	return username, err
}

// get a username and ID from email
func GetUserIdByEmail(email string) (id string, err error) {
	err = DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	return
}

// return true if the user has a password with the account
func PasswordExists(email string) (exists bool, err error) {
	err = DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND password_hash IS NOT NULL)`, email).Scan(&exists)
	return
}

// return true if username already in DB
func UsernameExists(username string) (exists bool, err error) {
	err = DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`,
		username,
	).Scan(&exists)
	return
}

// return true if username already in DB
func EmailExists(email string) (exists bool, err error) {
	err = DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exists)
	return
}

// return true if account is activated
func IsActivated(email string) (isActive bool, err error) {
	err = DB.QueryRow(
		`SELECT is_active FROM users WHERE email = $1`,
		email,
	).Scan(&isActive)
	return
}
