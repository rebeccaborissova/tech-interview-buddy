// Currently unused file

package middleware

import (
	"GO_PRACTICE_PROJECT/api"
	"GO_PRACTICE_PROJECT/internal/tools"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Declare errors
var BlankHeadersError = errors.New("Username or password cannot be blank")
var UserNotFoundError = errors.New("Invalid username")
var UnauthorizedError = errors.New("Invalid password")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, router *http.Request) {
		var username = router.Header.Get("Username")
		var token = router.Header.Get("Authorization")
		var err error

		if username == "" || token == "" {
			log.Error(BlankHeadersError)
			api.RequestErrorHandler(writer, BlankHeadersError)
			return
		}

		// Get an instance of the database
		store, err := tools.NewPostgresStore()
		if err != nil {
			api.InternalErrorHandler(writer)
			return
		}

		// Get account with primary key "username"
		usersCollection := store.DB.Collection("users")
		account := tools.EmailInDatabase(username, usersCollection)
		if account == nil {
			log.Error(UserNotFoundError)
			api.RequestErrorHandler(writer, UserNotFoundError)
		}

		// Check if provided token is valid
		isValidToken, err := tools.IsCorrectPassword(username, token, usersCollection)
		if err != nil {
			api.InternalErrorHandler(writer)
			return
		}

		if !isValidToken {
			log.Error(UnauthorizedError)
			api.RequestErrorHandler(writer, UnauthorizedError)
			return
		}

		next.ServeHTTP(writer, router)
	})
}
