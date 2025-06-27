package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"rms/database/dbHelper"
	"rms/middleware"
)

func GetDistanceFromAddress(w http.ResponseWriter, r *http.Request) {
	// Get query params
	restaurantID := r.URL.Query().Get("restaurant_id")
	addressID := r.URL.Query().Get("address_id")

	if restaurantID == "" || addressID == "" {
		http.Error(w, "restaurant_id and address_id are required", http.StatusBadRequest)
		return
	}

	// Parse UUIDs
	restID, err := uuid.Parse(restaurantID)
	if err != nil {
		http.Error(w, "Invalid restaurant_id", http.StatusBadRequest)
		return
	}
	addrID, err := uuid.Parse(addressID)
	if err != nil {
		http.Error(w, "Invalid address_id", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call DB helper to compute distance
	distance, err := dbHelper.GetDistanceBetweenAddressAndRestaurant(userID, addrID, restID)
	if err != nil {
		logrus.Errorf("Error fetching distance: %v", err)
		http.Error(w, "Error fetching distance", http.StatusInternalServerError)
		return
	}

	// Respond
	resp := map[string]interface{}{
		"distance_km": distance,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
