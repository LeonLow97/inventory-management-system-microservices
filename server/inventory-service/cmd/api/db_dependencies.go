package main

import (
	"database/sql"
	"net/http"
	"time"

	inventory "github.com/LeonLow97/internal"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *application) setupDBDependencies(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	inventoryRepo := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepo)
	inventory.NewInventoryGRPCHandler(inventoryService)

	return r
}
