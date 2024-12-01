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

	// Get an instance of the user and session collections
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	
	usersCollection := tools.GetUserCollection(store.DB)
	sessionCollection := tools.GetSessionCollection(store.DB)

	// Check if there is an existing session
	// Get the session cookie from the HTTP request
	cookie, err := request.Cookie("session_token")
	if err != nil {
		log.Info("No existing session cookie found, continuing login")
	} else {
		// Get the session UUID from the cookie
		sessionUUID, err := uuid.FromString(cookie.Value)
		if err != nil {
			api.InternalErrorHandler(writer)
			return
		}

		// Check if the existing session is valid
		userSession, err := tools.GetSession(sessionUUID, sessionCollection)
		err = tools.CheckSession(userSession, sessionCollection, usersCollection)
		if err != nil {
			// If the session did not pass the check then continue and issue a new session :3
			log.Error(err)
		} else {
			var response = api.LoginResponse{
				Code:    http.StatusOK,
				Session: sessionUUID.String(),
				Message: "User is already logged in, using existing session",
			}

			// Send the HTTP response
			err = json.NewEncoder(writer).Encode(response)
			if err != nil {
				log.Error(err)
				api.InternalErrorHandler(writer)
				return
			}
			return
		}
	}

	if username == "" || token == "" {
		log.Error(IncompleteRequestError)
		api.RequestErrorHandler(writer, IncompleteRequestError)
		return
	}

	// Delete any sessions that may have existed previously for this user
	err = tools.DeleteSessionByUsername(username, sessionCollection, usersCollection)

	// Get account with primary key "username"
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

	// Create a new session token (UUID v4)
	sessionToken, err = uuid.NewV4()
	if err != nil {
		log.Error("Failed to generate UUID: %v", err)
	}
	// Make the session expire after 2 hours (periodic refresh required)
	expiresAt := time.Now().Add(2 * time.Hour)

	// Delete old sessions
	tools.DeleteSessionByUsername(username, sessionCollection, usersCollection)

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
