package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"GO_PRACTICE_PROJECT/api"
	"GO_PRACTICE_PROJECT/internal/tools"

	log "github.com/sirupsen/logrus"
)

func getSignUpReponse(writer http.ResponseWriter, request *http.Request) {
	// Set logging settings
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetOutput(os.Stdout)

	var (
		// Declare errors
		MalformedRequestError  = errors.New("Malformed request body")
		IncompleteRequestError = errors.New("Form have empty fields.")

		params = api.SignUpParams{}
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
		username   = params.Username
		token      = params.Authorization
		firstName  = params.FirstName
		lastName   = params.LastName
		takenDSA   = params.DSA
		schoolYear = params.Year
	)

	// Ensure all required fields are assigned
	if username == "" || token == "" || firstName == "" || lastName == "" || schoolYear == 0 {
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

	// Create a new account
	usersCollection := tools.GetCollection(store.DB)
	// Change to := if you use the array of bools
	err, _ = tools.InsertAccount(username, token, firstName, lastName, takenDSA, schoolYear, usersCollection)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}   

	/// TODO: here down

	var response = api.LoginResponse{
		Code:     http.StatusOK,
		Username: params.Username,
		Message:  "Successfully registered " + params.Username,
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
