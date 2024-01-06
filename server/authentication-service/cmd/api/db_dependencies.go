package main

import (
	"net/http"

	"github.com/LeonLow97/internal/authenticate"
	"github.com/LeonLow97/internal/users"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func (app *application) setupDBDependencies(db *sqlx.DB) http.Handler {
	r := mux.NewRouter()

	authRepo := authenticate.NewRepo(db)
	authService := authenticate.NewService(authRepo)
	authenticate.NewAuthenticateGRPCHandler(authService)

	userRepo := users.NewRepo(db)
	userService := users.NewService(userRepo)
	users.NewUsersGRPCHandler(userService)

	return r
}
