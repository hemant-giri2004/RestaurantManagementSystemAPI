package dbHelper

import (
	"database/sql"
	"github.com/google/uuid"
	"rms/database"
	"rms/models"
)

func CreateRestaurant(name string, lat, lng float64, createdBy uuid.UUID) (uuid.UUID, error) {
	query := `
		INSERT INTO restaurants (id, restaurantname, lat, lng, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	id := uuid.New()
	err := database.RMS.QueryRow(query, id, name, lat, lng, createdBy).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func DoesRestaurantExist(id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM restaurants WHERE id = $1 AND archived_at IS NULL)`
	var exists bool
	err := database.RMS.QueryRow(query, id).Scan(&exists)
	return exists, err
}

func CreateDish(dishName string, restaurantID, createdBy uuid.UUID, price float64) (uuid.UUID, error) {
	var dishID uuid.UUID
	query := `
		INSERT INTO dishes (id, dishname, restaurant_id, created_by, price)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
		RETURNING id`
	err := database.RMS.QueryRow(query, dishName, restaurantID, createdBy, price).Scan(&dishID)
	return dishID, err
}

func FetchAllRestaurants(page, limit int) ([]models.Restaurant, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, restaurantname, created_by, lat, lng
		FROM restaurants
		WHERE archived_at IS NULL
		ORDER BY created_by
		LIMIT $1 OFFSET $2
	`

	rows, err := database.RMS.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var r models.Restaurant
		if err := rows.Scan(&r.ID, &r.Name, &r.CreatedBy, &r.Lat, &r.Lng); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	return restaurants, nil
}

func FetchDishesByRestaurant(restaurantID uuid.UUID, page, limit int) ([]models.Dish, error) {
	offset := (page - 1) * limit
	query := `
        SELECT id, dishname, price
        FROM dishes
        WHERE restaurant_id = $1 AND archived_at IS NULL
        ORDER BY id 
        LIMIT $2 OFFSET $3
    `
	rows, err := database.RMS.Query(query, restaurantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []models.Dish
	for rows.Next() {
		var dish models.Dish
		err := rows.Scan(&dish.ID, &dish.DishName, &dish.Price)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return dishes, nil
}

// dbHelper/restaurant.go

func GetRestaurantsVisibleTo(userID uuid.UUID, isAdmin bool, page, limit int) ([]models.Restaurant, error) {
	offset := (page - 1) * limit //for pagination
	var rows *sql.Rows
	var err error

	if isAdmin {
		rows, err = database.RMS.Query(`SELECT id, restaurantname, lat, lng, created_by FROM restaurants ORDER BY created_by LIMIT $1 OFFSET $2`, limit, offset)
	} else {
		rows, err = database.RMS.Query(`SELECT id, restaurantname, lat, lng, created_by FROM restaurants WHERE created_by = $1 ORDER BY  created_by LIMIT $2 OFFSET $3`, userID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var r models.Restaurant
		if err := rows.Scan(&r.ID, &r.Name, &r.Lat, &r.Lng, &r.CreatedBy); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}

	return restaurants, nil
}

// dbHelper/dishes.go

func GetDishesVisibleTo(userID uuid.UUID, isAdmin bool, page, limit int) ([]models.Dishes, error) {
	offset := (page - 1) * limit
	var rows *sql.Rows
	var err error

	if isAdmin {
		rows, err = database.RMS.Query(`SELECT id, dishname, restaurant_id, created_by, price FROM dishes ORDER BY created_by  LIMIT $1 OFFSET $2`, limit, offset)
	} else {
		rows, err = database.RMS.Query(`SELECT id, dishname, restaurant_id, created_by, price FROM dishes WHERE created_by = $1 ORDER BY created_by LIMIT $2 OFFSET $3`, userID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []models.Dishes
	for rows.Next() {
		var d models.Dishes
		if err := rows.Scan(&d.ID, &d.Name, &d.RestaurantID, &d.CreatedBy, &d.Price); err != nil {
			return nil, err
		}
		dishes = append(dishes, d)
	}

	return dishes, nil
}
