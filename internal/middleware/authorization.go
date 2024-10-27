package middleware

import (
	"GO_PRACTICE_PROJECT/api"
	"GO_PRACTICE_PROJECT/internal/tools"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var UnauthorizedError = errors.New("Invalid username or token.")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, router *http.Request) {
		var username = router.Header.Get("Username")
		var token = router.Header.Get("Authorization")
		var err error

		if username == "" || token == "" {
			log.Error(UnauthorizedError)
			api.RequestErrorHandler(writer, UnauthorizedError)
			return
		}

		// Get an instance of the database
		store, err := tools.NewPostgresStore()
		if err != nil {
			api.InternalErrorHandler(writer)
			return
		}

		usersCollection := store.DB.Collection("users")
		account := tools.EmailInDatabase(username, usersCollection)
		isValidToken, err := tools.IsCorrectPassword(username, token, usersCollection)
		if err != nil {
			api.InternalErrorHandler(writer)
			return
		}

		// Check if provided token is valid
		if account == nil || !isValidToken {
			log.Error(UnauthorizedError)
			api.RequestErrorHandler(writer, UnauthorizedError)
			return
		}

		next.ServeHTTP(writer, router)
	})
}
