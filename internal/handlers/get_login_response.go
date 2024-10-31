package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"GO_PRACTICE_PROJECT/api"
	"GO_PRACTICE_PROJECT/internal/tools"

	log "github.com/sirupsen/logrus"
)

func getLoginReponse(writer http.ResponseWriter, router *http.Request) {
	// Enable logging
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	
	// Declare errors
	var (
		MalformedRequestError  = errors.New("Malformed request body")
		IncompleteRequestError = errors.New("Username or password cannot be blank")
		UserNotFoundError      = errors.New("User does not exist")
		UnauthorizedError      = errors.New("Invalid password")
	)

	// var username = router.Header.Get("Username")
	// var token = router.Header.Get("Authorization")
	var (
		params = api.LoginParams{}
		err    error
	)

	// Decode HTTP request body into params struct
	decoder := json.NewDecoder(router.Body)
	err = decoder.Decode(&params)

	log.WithFields(log.Fields{
		"Decoder error": err,
		"Login parameters": params,
		"Request Header": router.Header,
		"Request Body": router.Body,
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
		log.Error(UserNotFoundError)
		api.RequestErrorHandler(writer, UserNotFoundError)
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

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	writer.Header().Set("Access-Control-Max-Age", "15")
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
