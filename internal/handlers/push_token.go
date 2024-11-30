package handlers

import (
	"encoding/json"
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	log "github.com/sirupsen/logrus"
)

func getPushToken(writer http.ResponseWriter, request *http.Request) {
	// Get an instance of the user collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetUserCollection(store.DB)

	// Obtain the user's push token
	username := request.Context().Value("username").(string)
	pushToken := tools.GetPushToken(username, usersCollection)

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

func setPushToken(writer http.ResponseWriter, request *http.Request) {
	// Decode HTTP request body into params struct
	var params = api.SimpleRequest{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)

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
