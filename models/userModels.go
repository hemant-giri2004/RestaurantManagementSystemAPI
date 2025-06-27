package models

import (
	"github.com/google/uuid"
	"time"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // omit in JSON response
	CreatedBy uuid.UUID `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
type AddressRequest struct {
	Label string  `json:"label"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}
type RegisterRequest struct {
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Password  string           `json:"password"`
	Roles     []string         `json:"roles,omitempty"`
	Addresses []AddressRequest `json:"addresses"` // Optional
}

type Response struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateSubadminRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
