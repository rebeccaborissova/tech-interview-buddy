package handlers

import (
	"encoding/json"
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	log "github.com/sirupsen/logrus"
)

func getActiveUsers(writer http.ResponseWriter, request *http.Request) {
	// Get an instance of the users collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetUserCollection(store.DB)
	// sessionCollection := tools.GetUserCollection(store.DB)

	// Get an array of active users
	activeAccounts, err := tools.GetOnlineAccounts(usersCollection)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}

	// Encode the array of active users as JSON
	activeUsersJSON, err := json.Marshal(activeAccounts)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}

	// Attempt to remove backslashes (did not work)
	// activeUsersString := string(activeUsersJSON)
	// strings.ReplaceAll(activeUsersString, `\`, "")
	// strings.ReplaceAll(activeUsersString, "[", "{")
	// strings.ReplaceAll(activeUsersString, "]", "}")

	var response = api.SimpleResponse{
		Code:    http.StatusOK,
		Message: string(activeUsersJSON),
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
