package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"GO_PRACTICE_PROJECT/api"
	"GO_PRACTICE_PROJECT/internal/tools"

	log "github.com/sirupsen/logrus"
)

func getLoginReponse(writer http.ResponseWriter, request *http.Request) {
	// Set logging settings
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetOutput(os.Stdout)

	var (
		// Declare errors
		MalformedRequestError  = errors.New("Malformed request body")
		IncompleteRequestError = errors.New("Username or password cannot be blank")
		UnauthorizedError      = errors.New("Invalid username or password")

		params = api.LoginParams{}
		err    error
	)

	// Decode HTTP request body into params struct
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&params)

	log.WithFields(log.Fields{
		"Login parameters": params,
		"Request Header":   request.Header,
		"Request Body":     request.Body,
	}).Info("HTTP request received")

	if err != nil {
		log.Error(MalformedRequestError)
		api.RequestErrorHandler(writer, MalformedRequestError)
		return
	}

	var (
		username = params.Username
		token    = params.Authorization
	)

	if username == "" || token == "" {
		log.Error(IncompleteRequestError)
		api.RequestErrorHandler(writer, IncompleteRequestError)
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
		log.Error(UnauthorizedError)
		api.RequestErrorHandler(writer, UnauthorizedError)
		return
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

	var response = api.LoginResponse{
		Code:     http.StatusOK,
		Username: params.Username,
		Message:  "Successfully authenticated " + params.Username,
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
