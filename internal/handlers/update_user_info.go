package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	log "github.com/sirupsen/logrus"
)

func updateUserInfo(writer http.ResponseWriter, request *http.Request) {
	var (
		MalformedRequestError = errors.New("Malformed request body")
		PasswordUpdateError   = errors.New("Unable to update user password")
		FNameUpdateError      = errors.New("Unable to update user first name")
		LNameUpdateError      = errors.New("Unable to update user last name")
		DSAUpdateError        = errors.New("Unable to update user DSA status")
		YearUpdateError       = errors.New("Unable to update user year")
		DescUpdateError       = errors.New("Unable to update user description")

		params = api.SignUpParams{}
	)

	// Decode HTTP request body into params struct
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(writer, MalformedRequestError)
		return
	}

	// Get an instance of the user collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetUserCollection(store.DB)

	username := request.Context().Value("username").(string)
	userAccount := tools.EmailInDatabase(username, usersCollection)

	// Check to see which fields need to be updated (TODO: email?)
	// if params.Username != userAccount.Email {}

	passwordsMatch, err := tools.IsCorrectPassword(username, params.Authorization, usersCollection)
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	if !passwordsMatch {
		err = tools.UpdatePassword(username, params.Authorization, usersCollection)
		if err != nil {
			api.RequestErrorHandler(writer, PasswordUpdateError)
			return
		}
	}
	if params.FirstName != userAccount.FirstName {
		err = tools.UpdateFirstName(username, params.FirstName, usersCollection)
		if err != nil {
			api.RequestErrorHandler(writer, FNameUpdateError)
			return
		}
	}
	if params.LastName != userAccount.LastName {
		err = tools.UpdateFirstName(username, params.FirstName, usersCollection)
		if err != nil {
			api.RequestErrorHandler(writer, LNameUpdateError)
			return
		}
	}
	if params.TakenDSA != userAccount.TakenDSA {
		err = tools.UpdateDSA(username, params.TakenDSA, usersCollection)
		if err != nil {
			api.RequestErrorHandler(writer, DSAUpdateError)
			return
		}
	}
	if params.Year != userAccount.Year {
		err = tools.UpdateYear(username, params.Year, usersCollection)
		if err != nil {
			api.RequestErrorHandler(writer, YearUpdateError)
			return
		}
	}
	if params.Description != userAccount.Description {
		err = tools.UpdateDescription(username, params.Description, usersCollection)
		if err != nil {
			api.RequestErrorHandler(writer, DescUpdateError)
			return
		}
	}

	var response = api.SimpleResponse{
		Code:    http.StatusOK,
		Message: "Successfully updated user fields",
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
