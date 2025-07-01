package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"rms/middleware"
	"strings"
	"time"

	"rms/database/dbHelper"
	"rms/models"
	"rms/utils"

	"github.com/sirupsen/logrus"
)

//func RegisterHandler(w http.ResponseWriter, r *http.Request) {
//	var req models.RegisterRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "invalid request body", http.StatusBadRequest)
//		return
//	}
//
//	// Check if user already exists
//	exists, err := dbHelper.IsEmailAlreadyRegistered(req.Email)
//	if err != nil {
//		http.Error(w, "internal server error", http.StatusInternalServerError)
//		return
//	}
//	if exists {
//		http.Error(w, "email already in use", http.StatusConflict)
//		return
//	}
//
//	// Hash password
//	hashedPassword, err := utils.HashPassword(req.Password)
//	if err != nil {
//		http.Error(w, "failed to hash password", http.StatusInternalServerError)
//		return
//	}
//
//	// Create user
//	userID := uuid.New()
//	err = dbHelper.CreateUser(userID, req.Username, req.Email, hashedPassword)
//	if err != nil {
//		http.Error(w, "failed to create user", http.StatusInternalServerError)
//		return
//	}
//
//	// Assign roles
//	if len(req.Roles) == 0 {
//		req.Roles = []string{"user"}
//	}
//	for _, roleName := range req.Roles {
//		err := dbHelper.AssignRoleToUser(userID, roleName)
//		if err != nil {
//			logrus.Warnf("role assignment failed: %s", err.Error())
//		}
//	}
//
//	// Add addresses
//	if len(req.Addresses) > 0 {
//		for _, addr := range req.Addresses {
//			err := dbHelper.InsertAddress(userID, addr)
//			if err != nil {
//				logrus.Warnf("Failed to insert address for user %s: %v", userID, err)
//			}
//		}
//	} else {
//		logrus.Infof("No address provided for user %s", userID)
//	}
//
//	//create access token
//	accessToken, err := utils.GenerateJWT(userID.String(), req.Roles)
//	if err != nil {
//		http.Error(w, "failed to generate access token", http.StatusInternalServerError)
//		return
//	}
//
//	// Create refresh token
//	refreshToken, err := utils.CreateRefreshToken(userID)
//	if err != nil {
//		http.Error(w, "failed to create session", http.StatusInternalServerError)
//		return
//	}
//
//	//create response
//	resp := models.Response{
//		Message:      "User registered successfully",
//		AccessToken:  accessToken,
//		RefreshToken: refreshToken,
//	}
//	//send response
//	w.WriteHeader(http.StatusCreated)
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(resp)
//
//}

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	// Parse body
	var req models.CreateSubadminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "username, email and password are required", http.StatusBadRequest)
		return
	}

	// Get creator (admin) userID from context (set by middleware)
	creatorIDRaw := r.Context().Value(middleware.UserIDKey)
	creatorID, ok := creatorIDRaw.(uuid.UUID)
	if !ok {
		//fmt.Println("CreateSubadminHandler")
		http.Error(w, "Unauthorized: userID missing in context", http.StatusUnauthorized)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Create user in DB
	userID, err := dbHelper.CreateUserWithRole(req.Username, req.Email, hashedPassword, "subadmin", creatorID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		logrus.Errorf("Error creating subadmin: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Success response
	resp := map[string]interface{}{
		"message": "Subadmin created successfully",
		"user_id": userID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func CreateUserByAdminOrSubadmin(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "username, email and password are required", http.StatusBadRequest)
		return
	}

	// Get creator (admin or subadmin) ID
	rawCreatorID := r.Context().Value(middleware.UserIDKey)
	creatorID, ok := rawCreatorID.(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Create user with role "user"
	userID, err := dbHelper.CreateUserWithRole(req.Username, req.Email, hashedPassword, "user", creatorID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		logrus.Errorf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message": "User created successfully",
		"user_id": userID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := dbHelper.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		//fmt.Println("LoginHandler :%w", err)
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Get roles
	roles, err := dbHelper.GetUserRoles(user.ID)
	if err != nil {
		http.Error(w, "failed to fetch roles", http.StatusInternalServerError)
		return
	}

	// Generate JWT
	accessToken, err := utils.GenerateJWT(user.ID.String(), roles)
	if err != nil {
		http.Error(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Refresh token
	refreshToken, err := utils.CreateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "failed to create refresh token", http.StatusInternalServerError)
		return
	}

	//  Respond
	res := models.Response{
		Message:      "User logged in Successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//get refresh token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
		return
	}

	refreshToken := tokenParts[1]

	// Check if session exists
	session, err := dbHelper.GetSessionByToken(refreshToken)
	if err != nil {
		logrus.Errorf("error fetching session for token %s: %v", refreshToken, err)
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		logrus.Warnf("token expired: %s", refreshToken)
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Delete session (logout)
	err = dbHelper.DeleteSession(refreshToken)
	if err != nil {
		logrus.Errorf("failed to delete session: %v", err)
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}

	logrus.Infof("logout successful for token: %s", refreshToken)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

func AddUserAddress(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var req models.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized: user ID not found", http.StatusUnauthorized)
		return
	}

	// Save to DB
	addressID, err := dbHelper.InsertUserAddress(userID, req.Label, req.Lat, req.Lng)
	if err != nil {
		logrus.Errorf("Failed to insert address: %v", err)
		http.Error(w, "Failed to insert address", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message":    "Address added successfully",
		"address_id": addressID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handlers/users.go

func ListOfSubAdmins(w http.ResponseWriter, r *http.Request) {
	subadmins, err := dbHelper.GetUsersByRole("subadmin")
	if err != nil {
		logrus.Errorf("Error fetching subadmins: %v", err)
		http.Error(w, "Failed to fetch subadmins", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subadmins)
}

// handlers/users.go

func ListUsers(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid user ID in context", http.StatusUnauthorized)
		return
	}

	isAdmin := false
	for _, role := range roles {
		if strings.ToLower(role) == "admin" {
			isAdmin = true
			break
		}
	}

	users, err := dbHelper.GetUsersVisibleTo(userID, isAdmin)
	if err != nil {
		logrus.Errorf("Failed to get users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
