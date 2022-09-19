package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yinho999/go-bookings/pkg/config"
	"github.com/yinho999/go-bookings/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	// Chi
	mux := chi.NewRouter()

	// Use middleware
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	return mux
}
