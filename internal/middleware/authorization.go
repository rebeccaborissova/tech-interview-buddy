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
var ErrBlankHeaders = errors.New("username or password cannot be blank")
var ErrUserNotFound = errors.New("invalid username")
var ErrUnauthorized = errors.New("invalid password")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, router *http.Request) {
		var username = router.Header.Get("Username")
		var token = router.Header.Get("Authorization")
		var err error

		if username == "" || token == "" {
			log.Error(ErrBlankHeaders)
			api.RequestErrorHandler(writer, ErrBlankHeaders)
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
			log.Error(ErrUserNotFound)
			api.RequestErrorHandler(writer, ErrUserNotFound)
		} else {
			// Check if provided token is valid
			isValidToken, err := tools.IsCorrectPassword(username, token, usersCollection)
			if err != nil {
				api.InternalErrorHandler(writer)
				return
			}
			if !isValidToken {
				log.Error(ErrUnauthorized)
				api.RequestErrorHandler(writer, ErrUnauthorized)
				return
			}
		}

		next.ServeHTTP(writer, router)
	})
}
