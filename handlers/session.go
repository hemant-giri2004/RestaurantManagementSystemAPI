package handlers

import (
	"encoding/json"
	"net/http"
	"rms/database/dbHelper"
	"rms/models"
	"rms/utils"
	"strings"
	"time"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	//get refresh token
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	refreshToken := strings.TrimPrefix(tokenStr, "Bearer ")

	// Check if refresh token exists and is valid
	session, err := dbHelper.GetSessionByToken(refreshToken)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}
	if session.ExpiresAt.Before(time.Now()) {
		_ = dbHelper.DeleteSession(refreshToken)
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}

	// Get user's roles from DB
	roles, err := dbHelper.GetUserRoles(session.UserID)
	if err != nil {
		http.Error(w, "failed to get user roles", http.StatusInternalServerError)
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(session.UserID.String(), roles)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Generate new refresh token
	newRefreshToken, err := utils.CreateRefreshToken(session.UserID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Remove old session
	_ = dbHelper.DeleteSession(refreshToken)

	response := models.Response{
		Message:      "Refresh token generated",
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
