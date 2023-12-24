package main

import (
	"net/http"

	"github.com/LeonLow97/internal/authenticate"
	"github.com/LeonLow97/internal/users"
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
	authenticate.NewAuthenticateGRPCHandler(authService)

	userRepo := users.NewRepo(db)
	userService := users.NewService(userRepo)
	users.NewUserHandler(userService)

	return r
}
