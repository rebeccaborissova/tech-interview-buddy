package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	"github.com/gofrs/uuid/v5"
	log "github.com/sirupsen/logrus"
)

func logout(writer http.ResponseWriter, request *http.Request) {
	// Get the cookie from the HTTP request
	cookie, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		api.InternalErrorHandler(writer)
		return
	}

	sessionToken := cookie.Value

	// Remove the user's session from the database collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	sessionCollection := tools.GetSessionCollection(store.DB)

	sessionUUID, err := uuid.FromString(sessionToken)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	userSession := tools.GetSession(sessionUUID, sessionCollection)

	// Remove the session from the database
	tools.DeleteSessionByUsername(userSession.Username, sessionCollection)

	var response = api.SimpleResponse{
		Code:    http.StatusOK,
		Message: "Successfully logged out",
	}

	// Set the client cookie
	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
