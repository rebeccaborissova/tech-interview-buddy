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

		router.Post("/login", getLoginReponse)
		router.Post("/signup", getSignUpReponse)
		router.Post("/logout", logout)
	})

	router.Route("/app", func(router chi.Router) {
		// Middleware for /app route
		router.Use(middleware.AuthenticateUser)

		router.Post("/activeusers", getActiveUsers)
		router.Post("/userinfo", getUserInfo)
		router.Post("/useredit", updateUserInfo)
		router.Post("/userdelete", deleteUser)
		router.Post("/refresh", refreshUserSession)
	})
}
