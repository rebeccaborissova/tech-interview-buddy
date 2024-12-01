package handlers

import (
	"encoding/json"
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"
	log "github.com/sirupsen/logrus"
)

// getPushToken handles requests to retrieve a user's push token
func getPushToken(writer http.ResponseWriter, request *http.Request) {
    var params = api.GetPushTokenRequest{}
    decoder := json.NewDecoder(request.Body)
    err := decoder.Decode(&params)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(writer)
        return
    }

	// Ensure that the username is not blank
    if params.Username == "" {
        api.InternalErrorHandler(writer)
        return
    }

    store, err := tools.NewPostgresStore()
    if err != nil {
        api.InternalErrorHandler(writer)
        return
    }
    usersCollection := tools.GetUserCollection(store.DB)

    pushToken := tools.GetPushToken(params.Username, usersCollection)
    if pushToken == nil {
        api.InternalErrorHandler(writer)
        return
    }

    var response = api.SimpleResponse{
        Code:    http.StatusOK,
        Message: *pushToken,
    }

	// Send the HTTP response
    err = json.NewEncoder(writer).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(writer)
        return
    }
}

// setPushToken handles requests to update a user's push token
func setPushToken(writer http.ResponseWriter, request *http.Request) {
	var params = api.SimpleRequest{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}

	// Get an instance of the user collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetUserCollection(store.DB)

	// Set the user's push token
	username := request.Context().Value("username").(string)
	err = tools.UpdatePushToken(username, params.Token, usersCollection)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}

	var response = api.SimpleResponse{
		Code:    http.StatusOK,
		Message: "Successfully updated pushToken for user " + username,
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
