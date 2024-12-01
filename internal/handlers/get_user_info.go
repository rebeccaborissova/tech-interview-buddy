package handlers

import (
	"encoding/json"
	"net/http"

	"CODE_CONNECT_API/api"
	"CODE_CONNECT_API/internal/tools"

	log "github.com/sirupsen/logrus"
)

func getUserInfo(writer http.ResponseWriter, request *http.Request) {
	// Get an instance of the user collection
	store, err := tools.NewPostgresStore()
	if err != nil {
		api.InternalErrorHandler(writer)
		return
	}
	usersCollection := tools.GetUserCollection(store.DB)

	// Get an instance of the user's account
	username := request.Context().Value("username").(string)

	userAccount := tools.EmailInDatabase(username, usersCollection)
	if err != nil || userAccount == nil {
		api.InternalErrorHandler(writer)
		return
	}

	// Debug hell
	log.WithFields(log.Fields{
		"username from request context": username,
		"Returned account from DB":      userAccount,
		"Error":                         err,
	}).Info("HTTP request received")

	var response = api.UserInfoResponse{
		Code:        http.StatusOK,
		Username:    userAccount.Email,
		FirstName:   userAccount.FirstName,
		LastName:    userAccount.LastName,
		InvitedBy:   userAccount.InvitedBy,
		TakenDSA:    userAccount.TakenDSA,
		Year:        userAccount.Year,
		Description: userAccount.Description,
	}

	// Send the HTTP response
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(writer)
		return
	}
}
