package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"rms/database/dbHelper"
	"rms/middleware"
	"rms/models"
	"rms/utils"
	"strings"
)

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	// Parse JSON body
	var req models.CreateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.RestaurantName = strings.TrimSpace(req.RestaurantName)
	if req.RestaurantName == "" || req.Lat == 0 || req.Lng == 0 {
		http.Error(w, "restaurant_name, lat, and lng are required", http.StatusBadRequest)
		return
	}

	// Get userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized: user ID missing in context", http.StatusUnauthorized)
		return
	}

	// Insert restaurant
	restaurantID, err := dbHelper.CreateRestaurant(req.RestaurantName, req.Lat, req.Lng, userID)
	if err != nil {
		logrus.Errorf("CreateRestaurant error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Restaurant created successfully",
		"restaurant_id": restaurantID,
	})
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.DishName = strings.TrimSpace(req.DishName)
	if req.DishName == "" || req.RestaurantID == "" {
		http.Error(w, "dish_name and restaurant_id are required", http.StatusBadRequest)
		return
	}

	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		http.Error(w, "Invalid restaurant_id format", http.StatusBadRequest)
		return
	}

	// Check if restaurant exists
	exists, err := dbHelper.DoesRestaurantExist(restaurantID)
	if err != nil {
		logrus.Errorf("Error checking restaurant existence: %v", err)
		http.Error(w, "Failed to verify restaurant", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	//fmt.Println("CreateDish", req)
	// Get creator user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized: user ID missing", http.StatusUnauthorized)
		return
	}

	if req.Price <= 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}
	// Insert the dish
	dishID, err := dbHelper.CreateDish(req.DishName, restaurantID, userID, req.Price)
	if err != nil {
		logrus.Errorf("CreateDish error: %v", err)
		http.Error(w, "Failed to create dish", http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Dish created successfully",
		"dish_id": dishID,
	})
}

func GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	page, limit := utils.ParsePageAndLimit(r)
	restaurants, err := dbHelper.FetchAllRestaurants(page, limit)
	if err != nil {
		logrus.Errorf("Failed to fetch restaurants: %v", err)
		http.Error(w, "Failed to fetch restaurants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func GetDishesByRestaurant(w http.ResponseWriter, r *http.Request) {
	page, limit := utils.ParsePageAndLimit(r)
	//fmt.Println("GetDishesByRestaurant")
	vars := mux.Vars(r)
	restaurantIDStr := vars["restaurant_id"]

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	dishes, err := dbHelper.FetchDishesByRestaurant(restaurantID, page, limit)
	if err != nil {
		logrus.Errorf("Failed to fetch dishes: %v", err)
		http.Error(w, "Failed to fetch dishes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dishes)
}

func ListRestaurants(w http.ResponseWriter, r *http.Request) {
	//parse
	page, limit := utils.ParsePageAndLimit(r)

	ctx := r.Context()
	rolesRaw := ctx.Value(middleware.RolesKey)
	userIDRaw := ctx.Value(middleware.UserIDKey)

	roles, ok := rolesRaw.([]string)
	if !ok || len(roles) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	isAdmin := false
	for _, role := range roles {
		if strings.ToLower(role) == "admin" {
			isAdmin = true
			break
		}
	}

	restaurants, err := dbHelper.GetRestaurantsVisibleTo(userID, isAdmin, page, limit)
	if err != nil {
		logrus.Errorf("Failed to fetch restaurants: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func ListDishes(w http.ResponseWriter, r *http.Request) {
	//parse
	page, limit := utils.ParsePageAndLimit(r)

	//use context
	ctx := r.Context()

	rolesRaw := ctx.Value(middleware.RolesKey)
	userIDRaw := ctx.Value(middleware.UserIDKey)

	roles, ok := rolesRaw.([]string)
	if !ok || len(roles) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	isAdmin := false
	for _, role := range roles {
		if strings.ToLower(role) == "admin" {
			isAdmin = true
			break
		}
	}

	dishes, err := dbHelper.GetDishesVisibleTo(userID, isAdmin, page, limit)
	if err != nil {
		logrus.Errorf("Failed to fetch dishes: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dishes)
}
