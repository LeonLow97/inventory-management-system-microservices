package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func (app *application) setupDBDependencies(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	return r
}
