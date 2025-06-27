package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"rms/database"
	"rms/models"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenStr string) (*models.CustomClaims, error) {
	claims := &models.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Generate a secure random string
func generateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateRefreshToken creates and saves a refresh token in the DB
func CreateRefreshToken(userID uuid.UUID) (string, error) {
	token, err := generateSecureToken(32)
	if err != nil {
		return "", err
	}

	expiryDays := 7
	if envDays := os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS"); envDays != "" {
		if parsed, err := time.ParseDuration(envDays + "d"); err == nil {
			expiryDays = int(parsed.Hours() / 24)
		}
	}

	session := models.Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: token,
		ExpiresAt:    time.Now().Add(time.Hour * 24 * time.Duration(expiryDays)),
		CreatedAt:    time.Now(),
	}
	//start session
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, expires_at, created_at)
		VALUES (:id, :user_id, :refresh_token, :expires_at, :created_at)
	`

	if _, err := database.RMS.NamedExec(query, &session); err != nil {
		return "", errors.New("could not save refresh token")
	}

	return token, nil
}
