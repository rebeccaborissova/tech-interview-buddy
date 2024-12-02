package handlers

import (
	"encoding/json"
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	log "github.com/sirupsen/logrus"
)

func deleteUser(writer http.ResponseWriter, request *http.Request) {
	// Get an instance of the user and session collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetUserCollection(store.DB)
	sessionCollection := tools.GetSessionCollection(store.DB)

	username := request.Context().Value("username").(string)

	// Purge all active sessions for the given account
	err = tools.DeleteSessionByUsername(username, usersCollection, sessionCollection)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}

	// Delete the user's account
	err = tools.DeleteAccount(username, usersCollection)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}

	var response = api.SimpleResponse{
		Code: http.StatusOK,
		Message: "Successfully deleted user: " + username,
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
