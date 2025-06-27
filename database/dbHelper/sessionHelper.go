package dbHelper

import (
	"rms/database"
	"rms/models"
)

func GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	query := `SELECT * FROM sessions WHERE refresh_token = $1`
	err := database.RMS.Get(&session, query, token)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteSession(refreshToken string) error {
	_, err := database.RMS.Exec(`DELETE FROM sessions WHERE refresh_token = $1`, refreshToken)
	return err
}
