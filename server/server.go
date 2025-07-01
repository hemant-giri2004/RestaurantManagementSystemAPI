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
	r.HandleFunc("/sign-in", handlers.LoginHandler).Methods("POST")

	// Session protected routes
	session := r.PathPrefix("/session").Subrouter()
	session.Use(middleware.SessionValidationMiddleware)
	session.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods("POST")
	session.HandleFunc("/sign-out", handlers.LogoutHandler).Methods("POST")

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
	adminOnly.HandleFunc("/sub-admin", handlers.CreateSubAdmin).Methods("POST")
	adminOnly.HandleFunc("/sub-admin", handlers.ListOfSubAdmins).Methods("GET")

	//for admin or sub-admin
	adminSubAdmin := r.PathPrefix("/admin-sub-admin").Subrouter()
	adminSubAdmin.Use(middleware.AuthMiddleware)
	adminSubAdmin.Use(middleware.RequireRolesMiddleware("admin", "subadmin"))
	adminSubAdmin.HandleFunc("/user", handlers.CreateUserByAdminOrSubadmin).Methods("POST")
	adminSubAdmin.HandleFunc("/restaurant", handlers.CreateRestaurant).Methods("POST")
	adminSubAdmin.HandleFunc("/dish", handlers.CreateDish).Methods("POST")
	adminSubAdmin.HandleFunc("/users", handlers.ListUsers).Methods("GET")
	adminSubAdmin.HandleFunc("/restaurants", handlers.ListRestaurants).Methods("GET")
	adminSubAdmin.HandleFunc("/dishes", handlers.ListDishes).Methods("GET")

	return r
}
