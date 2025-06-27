package dbHelper

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"rms/database"
	"rms/models"
)

func CreateUserWithRole(username, email, hashedPassword, roleName string, createdBy uuid.UUID) (string, error) {
	userID := uuid.New().String()

	tx, err := database.RMS.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Insert user
	_, err = tx.Exec(`
		INSERT INTO users (id, username, email, password, created_by)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, username, email, hashedPassword, createdBy)
	if err != nil {
		return "", err
	}

	// Get role ID
	var roleID string
	err = tx.QueryRow(`SELECT id FROM roles WHERE role_name = $1`, roleName).Scan(&roleID)
	if err != nil {
		return "", err
	}

	// Insert into user_roles
	_, err = tx.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
	`, userID, roleID)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func IsEmailAlreadyRegistered(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 AND archived_at IS NULL)`
	err := database.RMS.Get(&exists, query, email)
	return exists, err
}

func CreateUser(id uuid.UUID, username, email, hashedPassword string) error {
	query := `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err := database.RMS.Exec(query, id, username, email, hashedPassword)
	return err
}

func AssignRoleToUser(userID uuid.UUID, roleName string) error {
	var roleID uuid.UUID
	err := database.RMS.Get(&roleID, `SELECT id FROM roles WHERE role_name = $1`, roleName)
	if err != nil {
		return fmt.Errorf("role not found: %s", roleName)
	}
	_, err = database.RMS.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
	return err
}

func InsertAddress(userID uuid.UUID, addr models.AddressRequest) error {
	query := `
	INSERT INTO addresses (id, user_id, label, lat, lng)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := database.RMS.Exec(query, uuid.New(), userID, addr.Label, addr.Lat, addr.Lng)
	return err
}

func GetUserRoles(userID uuid.UUID) ([]string, error) {
	var roles []string
	query := `
        SELECT r.role_name 
        FROM roles r
        JOIN user_roles ur ON ur.role_id = r.id 
        WHERE ur.user_id = $1
    `
	err := database.RMS.Select(&roles, query, userID)
	return roles, err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, username, email, password 
        FROM users 
        WHERE email = $1 AND archived_at IS NULL
    `
	err := database.RMS.Get(&user, query, email)
	if err != nil {
		//fmt.Println("GetUserByEmail : %w ", err)
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func InsertUserAddress(userID uuid.UUID, label string, lat, lng float64) (uuid.UUID, error) {
	id := uuid.New()
	query := `
		INSERT INTO addresses (id, user_id, label, lat, lng)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := database.RMS.Exec(query, id, userID, label, lat, lng)
	return id, err
}

func GetUsersByRole(roleName string) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		JOIN roles r ON ur.role_id = r.id
		WHERE LOWER(r.role_name) = LOWER($1)
	`

	rows, err := database.RMS.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// dbHelper/user.go

func GetUsersVisibleTo(requesterID uuid.UUID, isAdmin bool) ([]models.User, error) {
	var rows *sql.Rows
	var err error

	if isAdmin {
		rows, err = database.RMS.Query(`SELECT id, username, email, created_at FROM users`)
	} else {
		rows, err = database.RMS.Query(`SELECT id, username, email, created_at FROM users WHERE created_by = $1`, requesterID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
