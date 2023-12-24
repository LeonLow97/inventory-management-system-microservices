package main

import (
	"net/http"

	"github.com/LeonLow97/internal/authenticate"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func routes(db *sqlx.DB) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successful Reply From Authentication Service!"))
	}).Methods("GET")

	authRepo := authenticate.NewRepo(db)
	authService := authenticate.NewService(authRepo)
	_ = authenticate.NewAuthenticateGRPCHandler(authService)
	// r.HandleFunc("/login", authHandler.Login).Methods("POST")
	// r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")

	// userRepo := users.NewRepo(db)
	// userService := users.NewService(userRepo)
	// userHandler := users.NewUserHandler(userService)
	// r.HandleFunc("/user", userHandler.UpdateUser).Methods("PATCH")
	// r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")

	return r
}
