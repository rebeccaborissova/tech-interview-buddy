package handlers

import (
	_ "GO_PRACTICE_PROJECT/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(router *chi.Mux) {
	// Global midderware
	router.Use(chimiddle.StripSlashes)

	router.Route("/account", func(router chi.Router) {
		// Middleware for /account route
		// router.Use(middleware.Authorization)

		router.Post("/", getLoginReponse)
	})
}
