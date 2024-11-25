package handlers

import (
	"CODE_CONNECT_API/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(router *chi.Mux) {
	// Global midderware
	router.Use(chimiddle.StripSlashes)
	router.Use(middleware.CORSMiddleware)

	router.Route("/account", func(router chi.Router) {
		// Middleware for /account route
		// router.Use(middleware.Authorization)

		router.Post("/login", getLoginReponse)
		router.Post("/signup", getSignUpReponse)
		router.Post("/logout", logout)
	})
}
