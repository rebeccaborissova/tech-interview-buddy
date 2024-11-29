package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	"github.com/gofrs/uuid/v5"
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

		params       = api.LoginParams{}
		sessionToken uuid.UUID
		err          error
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
	usersCollection := tools.GetUserCollection(store.DB)
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

	sessionCollection := tools.GetSessionCollection(store.DB)

	// Check if there is an existing session

	// Create a new session token (UUID v4)
	sessionToken, err = uuid.NewV4()
	if err != nil {
		log.Error("Failed to generate UUID: %v", err)
	}
	// Make the session expire after 10 minutes (periodic refresh required)
	expiresAt := time.Now().Add(time.Hour)

	// Delete old sessions 
	tools.DeleteSessionByUsername(username, sessionCollection)

	// Add the new session to the database
	tools.AddSession(sessionToken, username, expiresAt, sessionCollection, usersCollection)

	var response = api.LoginResponse{
		Code:    http.StatusOK,
		Session: sessionToken.String(),
		Message: "Successfully authenticated " + params.Username,
	}

	// Set the client cookie
	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: expiresAt,
	})

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
