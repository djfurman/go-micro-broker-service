package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify permission for who can connect
	mux.Use(cors.Handler(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Amzn-Trace-Id", "X-CSRF-Token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"https://*", "http://*"},
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	return mux
}
