package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"rms/handlers"
	"rms/middleware"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(struct{ Message string }{Message: "server is running"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logrus.Errorf("Error encoding response: %v", err)
			return
		}
	}).Methods("GET")

	// Public Routes
	//r.HandleFunc("/signup", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/signin", handlers.LoginHandler).Methods("POST")

	// Session protected routes
	session := r.PathPrefix("/session").Subrouter()
	session.Use(middleware.SessionValidationMiddleware)
	session.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods("POST")
	session.HandleFunc("/signout", handlers.LogoutHandler).Methods("POST")

	//only auth require
	openRoutes := r.PathPrefix("/").Subrouter()
	openRoutes.Use(middleware.AuthMiddleware)
	openRoutes.HandleFunc("/restaurants", handlers.GetAllRestaurants).Methods("GET")
	openRoutes.HandleFunc("/restaurants/{restaurant_id}/dishes", handlers.GetDishesByRestaurant).Methods("GET")
	openRoutes.HandleFunc("/user-address", handlers.AddUserAddress).Methods("POST")
	openRoutes.HandleFunc("/distance", handlers.GetDistanceFromAddress).Methods("GET")

	//only for admin
	adminOnly := r.PathPrefix("/admin-only").Subrouter()
	adminOnly.Use(middleware.AuthMiddleware)
	adminOnly.Use(middleware.RequireRolesMiddleware("admin"))
	adminOnly.HandleFunc("/subadmins", handlers.CreateSubadmin).Methods("POST")
	adminOnly.HandleFunc("/subadmins", handlers.ListSubadmins).Methods("GET")

	//for admin or subadmin
	adminSubadmin := r.PathPrefix("/admin-subadmin").Subrouter()
	adminSubadmin.Use(middleware.AuthMiddleware)
	adminSubadmin.Use(middleware.RequireRolesMiddleware("admin", "subadmin"))
	adminSubadmin.HandleFunc("/users", handlers.CreateUserByAdminOrSubadmin).Methods("POST")
	adminSubadmin.HandleFunc("/restaurants", handlers.CreateRestaurant).Methods("POST")
	adminSubadmin.HandleFunc("/dishes", handlers.CreateDish).Methods("POST")
	adminSubadmin.HandleFunc("/users", handlers.ListUsers).Methods("GET")
	adminSubadmin.HandleFunc("/restaurants", handlers.ListRestaurants).Methods("GET")
	adminSubadmin.HandleFunc("/dishes", handlers.ListDishes).Methods("GET")

	return r
}
