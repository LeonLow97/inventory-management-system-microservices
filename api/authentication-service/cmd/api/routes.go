package main

import (
	"net/http"

	"github.com/LeonLow97/internal/authenticate"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func routes(db *sqlx.DB) http.Handler {
	r := mux.NewRouter()

	authRepo := authenticate.NewRepo(db)
	authService := authenticate.NewService(authRepo)
	authHandler := authenticate.NewAuthenticateHandler(authService)
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")

	return r
}
