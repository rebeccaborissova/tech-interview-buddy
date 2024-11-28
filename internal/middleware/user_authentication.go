package middleware

import (
	"context"
	"net/http"
	"time"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	"github.com/gofrs/uuid/v5"
	log "github.com/sirupsen/logrus"
)

func AuthenticateAndRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Get the session cookie from the HTTP request
		cookie, err := request.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		// Get the session UUID from the cookie
		sessionToken := cookie.Value
		sessionUUID, err := uuid.FromString(sessionToken)
		if err != nil {
			api.InternalErrorHandler(writer)
			return
		}

		// Get an instance of the user and session collection
		store, err := tools.NewPostgresStore()
		if err != nil {
			api.InternalErrorHandler(writer)
			return
		}
		usersCollection := tools.GetUserCollection(store.DB)
		sessionCollection := tools.GetSessionCollection(store.DB)

		// Get an instance of the user's session
		userSession := tools.GetSession(sessionUUID, sessionCollection)

		// Ensure that the session is valid
		err = tools.CheckSession(userSession, sessionCollection, usersCollection)
		if err != nil {
			log.Error(err)
			api.RequestErrorHandler(writer, err)
			return
		}

		// If the previous session is valid, create a new session token for the current user
		newSessionToken, err := uuid.NewV4()
		if err != nil {
			log.Error("Failed to generate UUID: %v", err)
			return
		}

		// Make the session expire after 10 minutes (periodic refresh required)
		expiresAt := time.Now().Add(600 * time.Second)

		// Add the new session to the database
		tools.AddSession(newSessionToken, userSession.Username, expiresAt, sessionCollection, usersCollection)

		// Delete the older session token
		tools.DeleteSession(userSession, sessionCollection)

		// Set the new token as the users `session_token` cookie
		http.SetCookie(writer, &http.Cookie{
			Name:    "session_token",
			Value:   newSessionToken.String(),
			Expires: expiresAt,
		})

		ctx := context.WithValue(request.Context(), "username", userSession.Username)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
