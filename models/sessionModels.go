package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}
type CustomClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}
