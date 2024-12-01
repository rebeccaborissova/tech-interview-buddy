package handlers

import (
	"net/http"
	"time"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	"github.com/gofrs/uuid/v5"
	log "github.com/sirupsen/logrus"
)

func refreshUserSession(writer http.ResponseWriter, request *http.Request) {
	// Get an instance of the user and session collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetUserCollection(store.DB)
	sessionCollection := tools.GetSessionCollection(store.DB)

	newSessionToken, err := uuid.NewV4()
	if err != nil {
		log.Error("Failed to generate UUID: %v", err)
		return
	}

	// Make the session expire after 1 hour (periodic refresh required)
	expiresAt := time.Now().Add(time.Hour)

	// Delete the older session token
	username := request.Context().Value("username").(string)
	tools.DeleteSessionByUsername(username, sessionCollection, usersCollection)

	// Add the new session to the database
	tools.AddSession(newSessionToken, username, expiresAt, sessionCollection, usersCollection)

	// Set the new token as the users `session_token` cookie
	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken.String(),
		Expires: expiresAt,
	})
}
