package models

import (
	"github.com/google/uuid"
)

type Restaurant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"restaurantname"`
	CreatedBy uuid.UUID `json:"created_by"`
	//CreatedAt time.Time `json:"created_at"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type CreateRestaurantRequest struct {
	RestaurantName string  `json:"restaurant_name"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
}

type CreateDishRequest struct {
	DishName     string  `json:"dish_name"`
	RestaurantID string  `json:"restaurant_id"`
	Price        float64 `json:"price"`
}

type Dish struct {
	ID       uuid.UUID `json:"id"`
	DishName string    `json:"dishname"`
	Price    float64   `json:"price"`
}

// models/dish.go

type Dishes struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	CreatedBy    uuid.UUID `json:"created_by"`
	Price        float64   `json:"price"`
}
