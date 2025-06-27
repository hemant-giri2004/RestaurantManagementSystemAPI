package dbHelper

import (
	"github.com/google/uuid"
	"rms/database"
)

func GetDistanceBetweenAddressAndRestaurant(userID, addressID, restaurantID uuid.UUID) (float64, error) {
	query := `
		SELECT 
			111.111 * DEGREES(ACOS(LEAST(1.0, 
				COS(RADIANS(a.lat)) * COS(RADIANS(r.lat)) *
				COS(RADIANS(a.lng - r.lng)) +
				SIN(RADIANS(a.lat)) * SIN(RADIANS(r.lat))
			)))
		FROM addresses a
		JOIN restaurants r ON r.id = $3
		WHERE a.id = $2 AND a.user_id = $1
	`
	var distance float64
	err := database.RMS.QueryRow(query, userID, addressID, restaurantID).Scan(&distance)
	return distance, err
}
