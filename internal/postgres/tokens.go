package postgres

import (
	"time"
)

func CreateRefreshToken(id, token string, expiresAt time.Time) (err error) {
	_, err = DB.Exec(
		`INSERT INTO refresh_tokens (user_id, token, expires_at) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
    		token = EXCLUDED.token,
    		expires_at = EXCLUDED.expires_at,
    		created_at = CURRENT_TIMESTAMP;`, id, token, expiresAt)
	return
}

func DeleteRefreshToken(id string) (err error) {
	_, err = DB.Exec(`DELETE FROM refresh_tokens WHERE user_id = $1;`, id)
	return
}

func GetRefreshTokenInfo(token string) (id string, expiresAt time.Time, err error) {
	row := DB.QueryRow(`SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`, token)
	err = row.Scan(&id, &expiresAt)
	return
}
